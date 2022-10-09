package server

import (
	v1 "auth/api/auth/v1"
	"auth/internal/conf"
	"auth/internal/service"
	"context"
	"github.com/go-cinch/common/log"
	commonMiddleware "github.com/go-cinch/common/middleware"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtV4 "github.com/golang-jwt/jwt/v4"
)

func whitelist() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList[v1.OperationAuthLogin] = struct{}{}
	whiteList[v1.OperationAuthStatus] = struct{}{}
	whiteList[v1.OperationAuthLogout] = struct{}{}
	whiteList[v1.OperationAuthRegister] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Bootstrap, svc *service.AuthService) *http.Server {
	middlewares := []middleware.Middleware{
		recovery.Recovery(),
		ratelimit.Server(),
	}
	if c.Tracer.Enable {
		middlewares = append(middlewares, tracing.Server(), commonMiddleware.TraceId())
	}
	middlewares = append(
		middlewares,
		logging.Server(log.DefaultWrapper.Options().Logger()),
		validate.Validator(),
		selector.Server(
			jwt.Server(
				func(token *jwtV4.Token) (rp interface{}, err error) {
					rp = []byte(c.Auth.Jwt.Key)
					return
				},
				jwt.WithSigningMethod(jwtV4.SigningMethodHS512),
			),
		).Match(whitelist()).Build(),
	)
	var opts = []http.ServerOption{http.Middleware(middlewares...)}
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
	return srv
}
