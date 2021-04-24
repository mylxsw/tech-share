package api

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/listener"
	"github.com/mylxsw/glacier/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mylxsw/tech-share/config"
)

type Provider struct{}

func (s Provider) Aggregates() []infra.Provider {
	return []infra.Provider{
		web.Provider(
			listener.FlagContext("listen"),
			web.SetRouteHandlerOption(s.routes),
			web.SetMuxRouteHandlerOption(s.muxRoutes),
			web.SetExceptionHandlerOption(s.exceptionHandler),
			web.SetIgnoreLastSlashOption(true),
		),
	}
}

func (s Provider) Register(app infra.Binder) {}
func (s Provider) Boot(app infra.Resolver)   {}

func (s Provider) exceptionHandler(ctx web.Context, err interface{}) web.Response {
	if err2, ok := err.(error); ok {
		if errors.Is(err2, query.ErrNoResult) {
			return ctx.JSONWithCode(web.M{
				"error": fmt.Sprintf("%v", err2),
			}, http.StatusNotFound)
		}
	}
	log.Errorf("error: %v, call stack: %s", err, debug.Stack())
	return ctx.JSONWithCode(web.M{
		"error": fmt.Sprintf("%v", err),
	}, http.StatusInternalServerError)
}

func (s Provider) muxRoutes(cc infra.Resolver, router *mux.Router) {
	cc.MustResolve(func(conf *config.Config) {
		// prometheus metrics
		router.PathPrefix("/metrics").Handler(promhttp.Handler())
		// health check
		router.PathPrefix("/health").Handler(HealthCheck{})
		// file storage
		router.PathPrefix("/storage").Handler(http.FileServer(http.Dir(conf.StoragePath)))
	})
}

func (s Provider) routes(cc infra.Resolver, router web.Router, mw web.RequestMiddleware) {
	conf := config.Get(cc)

	mws := make([]web.HandlerDecorator, 0)
	mws = append(mws, mw.AccessLog(log.Module("api")), mw.CORS("*"))

	router.WithMiddleware(mws...).Controllers(
		"/api",
		controllers(cc, conf)...,
	)
}

type HealthCheck struct{}

func (h HealthCheck) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(`{"status": "UP"}`))
}