package api

import (
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/tech-share/config"

	"github.com/mylxsw/tech-share/api/controller"
)

func controllers(cc container.Resolver, conf *config.Config) []web.Controller {
	return []web.Controller{
		controller.NewInspectController(cc, conf),
		controller.NewShareController(cc, conf),
		controller.NewAuthController(cc, conf),
		controller.NewUploadController(cc, conf),
	}
}
