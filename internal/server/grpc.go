package server

import (
	tenantMiddleware "github.com/go-cinch/common/middleware/tenant/v2"

	"github.com/go-cinch/common/i18n"
	i18nMiddleware "github.com/go-cinch/common/middleware/i18n"
	"golang.org/x/text/language"

	"github.com/go-cinch/common/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"

	traceMiddleware "github.com/go-cinch/common/middleware/trace"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/redis/go-redis/v9"

	v1 "auth/api/auth"
	"auth/internal/conf"
	localMiddleware "auth/internal/server/middleware"
	"auth/internal/service"
)

// NewGRPCServer creates a gRPC server.
func NewGRPCServer(c *conf.Bootstrap, svc *service.AuthService, rds redis.UniversalClient) *grpc.Server {
	middlewares := []middleware.Middleware{
		recovery.Recovery(),
		tenantMiddleware.Tenant(), // Default required middleware for multi-tenancy

		i18nMiddleware.Translator(i18n.WithLanguage(language.Make(c.Server.Language)), i18n.WithFs(locales)),

		ratelimit.Server(),
		logging.Server(),
		metadata.Server(),
	}

	if c.Tracer.Enable {
		middlewares = append(middlewares, tracing.Server(), traceMiddleware.ID())
	}

	if c.Server.Validate {
		middlewares = append(middlewares, validate.Validator())
	}
	middlewares = append(middlewares, localMiddleware.Idempotent(rds))

	var opts = []grpc.ServerOption{
		grpc.Middleware(middlewares...),
	}
	if c.Server.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Server.Grpc.Network))
	}
	if c.Server.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Server.Grpc.Addr))
	}
	if c.Server.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Server.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)
	v1.RegisterAuthServer(srv, svc)

	return srv
}
