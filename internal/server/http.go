package server

import (
	"auth/api/auth"
	"auth/internal/biz"
	"auth/internal/conf"
	localMiddleware "auth/internal/server/middleware"
	"auth/internal/service"
	"github.com/go-cinch/common/i18n"
	"github.com/go-cinch/common/idempotent"
	i18nMiddleware "github.com/go-cinch/common/middleware/i18n"
	"github.com/go-cinch/common/middleware/logging"
	tenantMiddleware "github.com/go-cinch/common/middleware/tenant"
	traceMiddleware "github.com/go-cinch/common/middleware/trace"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/pprof"
	"github.com/redis/go-redis/v9"
	"golang.org/x/text/language"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(
	c *conf.Bootstrap,
	client redis.UniversalClient,
	idt *idempotent.Idempotent,
	svc *service.AuthService,
	whitelist *biz.WhitelistUseCase,
) *http.Server {
	var middlewares []middleware.Middleware
	if c.Tracer.Enable {
		middlewares = append(middlewares, tracing.Server(), traceMiddleware.Id())
	}
	middlewares = append(
		middlewares,
		recovery.Recovery(),
		tenantMiddleware.Tenant(),
		ratelimit.Server(),
		localMiddleware.Header(),
		logging.Server(),
		i18nMiddleware.Translator(i18n.WithLanguage(language.Make(c.Server.Language)), i18n.WithFs(locales)),
		metadata.Server(),
	)
	if c.Server.Jwt.Enable {
		middlewares = append(middlewares, localMiddleware.Permission(c, client, whitelist))
	}
	if c.Server.Idempotent {
		middlewares = append(middlewares, localMiddleware.Idempotent(idt))
	}
	if c.Server.Validate {
		middlewares = append(middlewares, validate.Validator())
	}
	var opts = []http.ServerOption{
		// already set cors header in nginx
		// http.Filter(handlers.CORS(
		// 	handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Idempotent"}),
		// 	handlers.AllowedMethods([]string{"OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"}),
		// 	handlers.AllowedOrigins([]string{"*"}),
		// 	handlers.AllowCredentials(),
		// )),
		http.Middleware(middlewares...),
	}
	if c.Server.Http.Network != "" {
		opts = append(opts, http.Network(c.Server.Http.Network))
	}
	if c.Server.Http.Addr != "" {
		opts = append(opts, http.Address(c.Server.Http.Addr))
	}
	if c.Server.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Server.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	auth.RegisterAuthHTTPServer(srv, svc)
	srv.HandlePrefix("/debug/pprof", pprof.NewHandler())
	srv.HandlePrefix("/pub/healthcheck", HealthHandler(svc))
	return srv
}
