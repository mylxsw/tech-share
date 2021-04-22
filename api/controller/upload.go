package controller

import (
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/tech-share/config"
)

type UploadController struct {
	cc   infra.Resolver
	conf *config.Config
}

func NewUploadController(cc infra.Resolver, conf *config.Config) web.Controller {
	return &UploadController{cc: cc, conf: conf}
}

func (ctl UploadController) Register(router web.Router) {
	router.Group("upload/", func(router web.Router) {
		router.Get("/", ctl.Upload)
	})
}

// Upload upload a file
func (ctl UploadController) Upload(ctx web.Context, req web.Request) web.Response {
	return ctx.JSON(web.M{})
}
