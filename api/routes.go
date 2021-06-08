package api

import (
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/tech-share/config"

	"github.com/mylxsw/tech-share/api/controller"
)

func authedControllers(cc container.Resolver, conf *config.Config) []web.Controller {
	return []web.Controller{
		controller.NewShareController(cc, conf),
		controller.NewUploadController(cc, conf),
		controller.NewUserController(cc, conf),
		controller.NewAuthController(cc, conf),
		controller.NewInspectController(cc, conf),
		controller.NewCreditController(cc, conf),
	}
}
