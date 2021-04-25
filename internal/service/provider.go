package service

import (
	eloquentEvt "github.com/mylxsw/eloquent/event"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
)

type Provider struct{}

func (p Provider) Register(cc infra.Binder) {
	cc.MustSingletonOverride(NewShareService)
	cc.MustSingletonOverride(NewAttachmentService)
	cc.MustSingletonOverride(NewUserService)
}

func (p Provider) Boot(cc infra.Resolver) {
	cc.MustResolve(func(publisher event.Publisher) {
		eloquentEvt.SetDispatcher(publisher)
	})
}
