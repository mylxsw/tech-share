package controller

import (
	"context"
	"fmt"

	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/tech-share/config"
	"github.com/mylxsw/tech-share/internal/auth"
	"github.com/mylxsw/tech-share/internal/service"
	"github.com/mylxsw/tech-share/pkg/bcrypt"
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
		router.Post("/login/", ctl.Login).Name("auth:login")
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

// Login let user login to the system
func (ctl AuthController) Login(ctx web.Context, req web.Request, userSrv service.UserService, authProvider auth.Auth) (*User, error) {
	username := req.Input("username")
	password := req.Input("password")

	if userLogin, ok := req.Session().Values["user_login"]; ok {
		return nil, service.NewValidateError(fmt.Errorf("%s, 你已经登录过了", userLogin.(service.UserInfo).Name))
	}

	if username == "" || password == "" {
		return nil, service.NewValidateError(fmt.Errorf("用户名或密码不能为空"))
	}

	authedUser, err := authProvider.Login(username, password)
	if err != nil {
		return nil, err
	}

	passwordHash, _ := bcrypt.Hash(password)
	user, err := userSrv.LoadUser(
		context.TODO(),
		authedUser.Account,
		service.UserInfo{
			Name:     authedUser.Name,
			Account:  authedUser.Account,
			Uuid:     authedUser.UUID,
			Status:   authedUser.Status,
			Password: passwordHash,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("加载用户失败: %w", err)
	}

	req.Session().Values["user_login"] = *user
	return &User{ID: user.Id, Name: user.Name, UUID: user.Uuid}, nil
}
