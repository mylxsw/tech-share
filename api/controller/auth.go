package controller

import (
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/tech-share/config"
)

type AuthController struct {
	cc   infra.Resolver
	conf *config.Config
}

func NewAuthController(cc infra.Resolver, conf *config.Config) web.Controller {
	return &AuthController{cc: cc, conf: conf}
}

func (ctl AuthController) Register(router web.Router) {
	router.Group("auth/", func(router web.Router) {
		router.Get("login/", ctl.Login)
	})
}

// Login let user login to the system
func (ctl AuthController) Login(ctx web.Context, req web.Request) web.Response {
	return ctx.JSON(web.M{})
}
