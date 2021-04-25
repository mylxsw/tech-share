package controller

import (
	"fmt"

	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/tech-share/config"
	"github.com/mylxsw/tech-share/internal/service"
)

type UserController struct {
	cc   infra.Resolver
	conf *config.Config
}

func NewUserController(cc infra.Resolver, conf *config.Config) web.Controller {
	return &UserController{cc: cc, conf: conf}
}

func (ctl UserController) Register(router web.Router) {
	router.Group("/users", func(router web.Router) {
		router.Get("/me", ctl.CurrentUser)
	})
}

func (ctl UserController) CurrentUser(req web.Request) (*User, error) {
	userLogin, ok := req.Session().Values["user_login"]
	if !ok {
		return nil, fmt.Errorf("no user logon")
	}

	user := userLogin.(service.UserInfo)
	return &User{ID: user.Id, Name: user.Name, UUID: user.Uuid}, nil
}
