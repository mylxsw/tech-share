package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/tech-share/config"
	"github.com/mylxsw/tech-share/internal/service"
)

// currentUser extract current user from request
func currentUser(req web.Request) service.UserInfo {
	userLogin, ok := req.Session().Values["user_login"]
	if !ok {
		return service.UserInfo{}
	}

	return userLogin.(service.UserInfo)
}

type AuthController struct {
	cc   infra.Resolver
	conf *config.Config
}

func NewAuthController(cc infra.Resolver, conf *config.Config) web.Controller {
	return &AuthController{cc: cc, conf: conf}
}

func (ctl AuthController) Register(router web.Router) {
	router.Group("auth/", func(router web.Router) {
		router.Post("/logout/", ctl.Logout).Name("auth:logout")
		router.Post("/login-ldap/", ctl.LdapLogin).Name("auth:login-ldap")
		router.Get("/current", ctl.CurrentUser).Name("auth:current")
	})
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

func (ctl AuthController) Logout(req web.Request) error {
	delete(req.Session().Values, "user_login")
	return nil
}

// CurrentUser return current user info
func (ctl AuthController) CurrentUser(ctx web.Context, req web.Request) *User {
	if userLogin, ok := ctx.Session().Values["user_login"]; ok {
		user := userLogin.(service.UserInfo)
		return &User{ID: user.Id, Name: user.Name, UUID: user.Uuid}
	}

	return nil
}

// LdapLogin let user login to the system
func (ctl AuthController) LdapLogin(ctx web.Context, req web.Request, userSrv service.UserService) (*User, error) {
	username := req.Input("username")
	password := req.Input("password")

	if userLogin, ok := req.Session().Values["user_login"]; ok {
		return nil, service.NewValidateError(fmt.Errorf("%s, 你已经登录过了", userLogin.(service.UserInfo).Name))
	}

	if username == "" || password == "" {
		return nil, service.NewValidateError(fmt.Errorf("用户名或密码不能为空"))
	}

	l, err := ldap.DialURL(ctl.conf.LDAP.URL)
	if err != nil {
		return nil, fmt.Errorf("无法连接 LDAP 服务器: %w", err)
	}

	defer l.Close()

	if err := l.Bind(ctl.conf.LDAP.Username, ctl.conf.LDAP.Password); err != nil {
		return nil, fmt.Errorf("LDAP 服务器鉴权失败: %w", err)
	}

	searchReq := ldap.NewSearchRequest(
		ctl.conf.LDAP.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(%s=%s))", ctl.conf.LDAP.UID, ldap.EscapeFilter(username)),
		[]string{"objectguid", ctl.conf.LDAP.UID, ctl.conf.LDAP.DisplayName, "userAccountControl"},
		nil,
	)

	sr, err := l.Search(searchReq)
	if err != nil {
		return nil, fmt.Errorf("LDAP 用户查询失败: %w", err)
	}

	if len(sr.Entries) != 1 {
		return nil, service.NewValidateError(fmt.Errorf("用户不存在"))
	}

	// 514-禁用 512-启用
	if sr.Entries[0].GetAttributeValue("userAccountControl") == "514" {
		return nil, service.NewValidateError(errors.New("LDAP 用户账户已禁用"))
	}

	if ctl.conf.WeakPasswordMode {
		// TODO 弱密码模式
	} else {
		if err := l.Bind(sr.Entries[0].DN, password); err != nil {
			return nil, service.NewValidateError(fmt.Errorf("用户密码错误: %w", err))
		}
	}

	user, err := userSrv.LoadUser(
		context.TODO(),
		uuid.Must(uuid.FromBytes(sr.Entries[0].GetRawAttributeValue("objectGUID"))).String(),
		service.UserInfo{
			Name:    sr.Entries[0].GetAttributeValue(ctl.conf.LDAP.DisplayName),
			Account: sr.Entries[0].GetAttributeValue(ctl.conf.LDAP.UID),
			Status:  service.UserStatusEnabled,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("加载用户失败: %w", err)
	}

	req.Session().Values["user_login"] = *user
	return &User{ID: user.Id, Name: user.Name, UUID: user.Uuid}, nil
}
