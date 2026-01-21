package server

import (
	"github.com/go-cinch/common/i18n"
	i18nMiddleware "github.com/go-cinch/common/middleware/i18n"
	"github.com/go-cinch/common/middleware/logging"
	tenantMiddleware "github.com/go-cinch/common/middleware/tenant/v2"
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

	v1 "auth/api/auth"
	"auth/internal/biz"
	"auth/internal/conf"
	localMiddleware "auth/internal/server/middleware"
	"auth/internal/service"
)

// NewHTTPServer creates an HTTP server.
func NewHTTPServer(
	c *conf.Bootstrap,
	svc *service.AuthService,
	rds redis.UniversalClient,
	permission *biz.PermissionUseCase,
	user *biz.UserUseCase,
	whitelist *biz.WhitelistUseCase,
) *http.Server {
	middlewares := []middleware.Middleware{
		recovery.Recovery(),
		tenantMiddleware.Tenant(), // Default required middleware for multi-tenancy
		i18nMiddleware.Translator(i18n.WithLanguage(language.Make(c.Server.Language)), i18n.WithFs(locales)),
		ratelimit.Server(),
		localMiddleware.Header(),
	}
	if c.Tracer.Enable {
		middlewares = append(middlewares, tracing.Server(), traceMiddleware.ID())
	}
	middlewares = append(middlewares, logging.Server(), metadata.Server())
	if c.Server.Validate {
		middlewares = append(middlewares, validate.Validator())
	}
	// Add Permission middleware for JWT parsing when enabled
	if c.Server.Jwt.Enable {
		middlewares = append(middlewares, localMiddleware.Permission(c, rds, permission, user, whitelist))
	}
	middlewares = append(middlewares, localMiddleware.Idempotent(rds))

	var opts = []http.ServerOption{
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
	v1.RegisterAuthHTTPServer(srv, svc)
	srv.HandlePrefix("/healthz", HealthHandler(svc))
	if c.Server.Http.Docs {
		srv.HandlePrefix("/docs/", DocsHandler())
	}
	if c.Server.EnablePprof {
		srv.HandlePrefix("/debug/pprof", pprof.NewHandler())
	}

	return srv
}
