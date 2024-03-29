package api

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/listener"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/go-utils/str"
	"github.com/mylxsw/tech-share/internal/auth"
	"github.com/mylxsw/tech-share/internal/auth/database"
	"github.com/mylxsw/tech-share/internal/auth/ldap"
	"github.com/mylxsw/tech-share/internal/service"
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

func (s Provider) Register(app infra.Binder) {
	// 注册鉴权实现
	app.MustSingletonOverride(func(conf *config.Config, userSrv service.UserService) auth.Auth {
		switch conf.AuthProvider {
		case "database":
			return database.New(conf, userSrv)
		case "ldap":
			return ldap.New(conf)
		default:
			panic(fmt.Errorf("invalid auth provider: %s, only support database/ldap", conf.AuthProvider))
		}
	})
}
func (s Provider) Boot(app infra.Resolver) {}

func (s Provider) exceptionHandler(ctx web.Context, err interface{}) web.Response {
	if err2, ok := err.(error); ok {
		if errors.Is(err2, query.ErrNoResult) {
			return ctx.JSONWithCode(web.M{
				"error": fmt.Sprintf("%v", err2),
			}, http.StatusNotFound)
		}
		if verr, ok := err2.(*service.ValidateError); ok {
			return ctx.JSONWithCode(web.M{
				"error": fmt.Sprintf("%v", verr),
			}, http.StatusUnprocessableEntity)
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
		router.PathPrefix("/storage").Handler(http.StripPrefix("/storage", http.FileServer(http.Dir(conf.StoragePath))))
	})
}

func (s Provider) routes(cc infra.Resolver, router web.Router, mw web.RequestMiddleware) {
	conf := config.Get(cc)

	mws := make([]web.HandlerDecorator, 0)
	mws = append(mws,
		mw.AccessLog(log.Module("api")),
		mw.CORS("*"),
		mw.Session(sessions.NewCookieStore([]byte(conf.SessionKey)), "tech-share", &sessions.Options{MaxAge: 86400 * 30 * 12, Path: "/"}),
	)

	// 存储在 session 中的对象必须在这里注册，否则无法序列化
	gob.Register(service.UserInfo{})

	authMW := mw.BeforeInterceptor(func(ctx web.Context) web.Response {
		if str.HasPrefixes(ctx.CurrentRoute().GetName(), []string{"auth:login", "inspect:"}) || str.HasSuffixes(ctx.CurrentRoute().GetName(), []string{"@public"}) {
			return nil
		}

		_, ok := ctx.Session().Values["user_login"]
		if !ok {
			return ctx.JSONWithCode(web.M{"error": "access denied"}, http.StatusUnauthorized)
		}

		return nil
	})

	router.WithMiddleware(append(mws, authMW)...).Controllers(
		"/api",
		authedControllers(cc, conf)...,
	)
}

type HealthCheck struct{}

func (h HealthCheck) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(`{"status": "UP"}`))
}
