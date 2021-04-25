package controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/tech-share/config"
	"github.com/mylxsw/tech-share/internal/service"
)

// currentUserID extract current user id from request
func currentUserID(req web.Request) int64 {
	userID, err := strconv.Atoi(req.Header("mock-user-id"))
	if err != nil {
		return 101
	}

	return int64(userID)
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
		router.Post("login/", ctl.Login)
	})
}

type LoginRes struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

// Login let user login to the system
func (ctl AuthController) Login(ctx web.Context, req web.Request, userSrv service.UserService) (*LoginRes, error) {
	username := req.Input("username")
	password := req.Input("password")

	l, err := ldap.DialURL(ctl.conf.LDAP.URL)
	if err != nil {
		return nil, err
	}

	defer l.Close()

	if err := l.Bind(ctl.conf.LDAP.Username, ctl.conf.LDAP.Password); err != nil {
		return nil, err
	}

	searchReq := ldap.NewSearchRequest(
		ctl.conf.LDAP.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(%s=%s))", ctl.conf.LDAP.UID, ldap.EscapeFilter(username)),
		[]string{"objectguid", ctl.conf.LDAP.UID, ctl.conf.LDAP.DisplayName},
		nil,
	)

	sr, err := l.Search(searchReq)
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) != 1 {
		return nil, fmt.Errorf("user does not exist or too many entries returned")
	}

	if err := l.Bind(sr.Entries[0].DN, password); err != nil {
		return nil, err
	}

	user, err := userSrv.LoadUser(
		context.TODO(),
		uuid.Must(uuid.FromBytes(sr.Entries[0].GetRawAttributeValue("objectGUID"))).String(),
		sr.Entries[0].GetAttributeValue(ctl.conf.LDAP.DisplayName),
	)
	if err != nil {
		return nil, err
	}

	return &LoginRes{ID: user.Id, Name: user.Name, UUID: user.Uuid}, nil
}
