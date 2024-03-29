package event

import (
	"github.com/mylxsw/asteria/log"
	eloquentEvt "github.com/mylxsw/eloquent/event"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
)

type Provider struct{}

func (s Provider) Aggregates() []infra.Provider {
	return []infra.Provider{
		event.Provider(s.eventHandler, event.SetStoreOption(func(cc infra.Resolver) event.Store {
			return event.NewMemoryEventStore(true, 100)
		})),
	}
}

func (s Provider) Register(cc infra.Binder) {}
func (s Provider) Boot(cc infra.Resolver)   {}

func (s Provider) eventHandler(cc infra.Resolver, listener event.Listener) {
	listener.Listen(func(evt SystemUpDownEvent) {
		log.Debugf("new event received: %v", evt)
	})
	listener.Listen(func(evt eloquentEvt.QueryExecutedEvent) {
		log.WithFields(log.Fields{
			"sql":    evt.SQL,
			"params": evt.Bindings,
			"elapse": evt.Time,
		}).Debugf("database_sql")
	})
}
