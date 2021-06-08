package controller

import (
	"context"
	"time"

	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/tech-share/config"
	"github.com/mylxsw/tech-share/internal/service"
)

type CreditController struct {
	cc   infra.Resolver
	conf *config.Config
}

func NewCreditController(cc infra.Resolver, conf *config.Config) web.Controller {
	return &CreditController{cc: cc, conf: conf}
}

func (wel CreditController) Register(router web.Router) {
	router.Group("/credits", func(router web.Router) {
		router.Any("/rank", wel.Ranks).Name("credits:rank@public")
	})
}

func (wel CreditController) Ranks(ctx web.Context, credSrv service.CreditService) (service.CreditRanks, error) {
	startAt, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006")+"-01-01 00:00:00")
	return credSrv.CreditRanks(context.TODO(), startAt)
}
