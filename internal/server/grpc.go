package server

import (
	"auth/api/auth"
	"auth/internal/conf"
	"auth/internal/pkg/idempotent"
	localMiddleware "auth/internal/server/middleware"
	"auth/internal/service"
	"github.com/go-cinch/common/i18n"
	"github.com/go-cinch/common/log"
	i18nMiddleware "github.com/go-cinch/common/middleware/i18n"
	traceMiddleware "github.com/go-cinch/common/middleware/trace"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/redis/go-redis/v9"
	"golang.org/x/text/language"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Bootstrap, client redis.UniversalClient, idt *idempotent.Idempotent, svc *service.AuthService) *grpc.Server {
	middlewares := []middleware.Middleware{
		recovery.Recovery(),
		ratelimit.Server(),
	}
	if c.Tracer.Enable {
		middlewares = append(middlewares, tracing.Server(), traceMiddleware.Id())
	}
	middlewares = append(
		middlewares,
		logging.Server(log.DefaultWrapper.Options().Logger()),
		i18nMiddleware.Translator(i18n.WithLanguage(language.Make(c.Server.Language)), i18n.WithFs(locales)),
		metadata.Server(),
	)
	if c.Server.Jwt {
		middlewares = append(middlewares, localMiddleware.Jwt(c, client))
	}
	if c.Server.Permission {
		middlewares = append(middlewares, localMiddleware.Permission(svc))
	}
	if c.Server.Idempotent {
		middlewares = append(middlewares, localMiddleware.Idempotent(idt))
	}
	if c.Server.Validate {
		middlewares = append(middlewares, validate.Validator())
	}
	var opts = []grpc.ServerOption{grpc.Middleware(middlewares...)}
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
	auth.RegisterAuthServer(srv, svc)
	return srv
}
