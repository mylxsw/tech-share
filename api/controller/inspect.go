package controller

import (
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/tech-share/pkg/bcrypt"

	"github.com/mylxsw/tech-share/config"
)

type InspectController struct {
	cc   infra.Resolver
	conf *config.Config
}

func NewInspectController(cc infra.Resolver, conf *config.Config) web.Controller {
	return &InspectController{cc: cc, conf: conf}
}

func (wel InspectController) Register(router web.Router) {
	router.Group("/inspect", func(router web.Router) {
		router.Any("/version", wel.Version).Name("inspect:version")
		router.Post("/helper/password-generator", wel.PasswordGenerator).Name("inspect:helper:password-generator")
	})
}

func (wel InspectController) Version(ctx web.Context) web.Response {
	return ctx.JSON(web.M{
		"version": wel.conf.Version,
		"git":     wel.conf.GitCommit,
	})
}

func (wel InspectController) PasswordGenerator(ctx web.Context) web.Response {
	password := ctx.Input("password")
	hash, _ := bcrypt.Hash(password)

	return ctx.JSON(web.M{"hash": hash})
}
