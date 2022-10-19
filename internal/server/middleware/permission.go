package middleware

import (
	"auth/api/auth"
	"auth/internal/biz"
	"auth/internal/conf"
	jwtLocal "auth/internal/pkg/jwt"
	"auth/internal/service"
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	jwtV4 "github.com/golang-jwt/jwt/v4"
	"strings"
)

func Permission(c *conf.Bootstrap, svc *service.AuthService) middleware.Middleware {
	return selector.Server(
		jwt(c),
		func(handler middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
				var resource string
				if tr, ok := transport.FromServerContext(ctx); ok {
					resource = tr.Operation()
				}
				user := jwtLocal.FromContext(ctx)
				res, err := svc.Permission(ctx, &auth.PermissionRequest{
					UserCode: user.Code,
					Resource: resource,
				})
				if err != nil {
					return
				}
				if !res.Pass {
					err = biz.NoPermission
					return

				}
				return handler(ctx, req)
			}
		},
	).Match(permissionWhitelist()).Build()
}

func jwt(c *conf.Bootstrap) func(handler middleware.Handler) middleware.Handler {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				switch tr.Kind() {
				case transport.KindGRPC:
					ctx = jwtLocal.NewContextByCode(ctx, tr.RequestHeader().Get("X-USER-CODE"))
				case transport.KindHTTP:
					// http need check Authorization header
					auths := strings.SplitN(tr.RequestHeader().Get("Authorization"), " ", 2)
					if len(auths) != 2 || !strings.EqualFold(auths[0], "Bearer") {
						err = MissingToken
						return
					}
					jwtToken := auths[1]
					var tokenInfo *jwtV4.Token
					tokenInfo, err = jwtV4.Parse(jwtToken, func(token *jwtV4.Token) (rp interface{}, err error) {
						rp = []byte(c.Auth.Jwt.Key)
						return
					})
					if err != nil {
						ve, ok := err.(*jwtV4.ValidationError)
						if !ok {
							return
						}
						if ve.Errors&jwtV4.ValidationErrorMalformed != 0 {
							err = TokenInvalid
							return
						}
						if ve.Errors&(jwtV4.ValidationErrorExpired|jwtV4.ValidationErrorNotValidYet) != 0 {
							err = TokenExpired
							return
						}
						err = TokenParseFail
						return
					}
					if !tokenInfo.Valid {
						err = TokenParseFail
						return
					}
					if tokenInfo.Method != jwtV4.SigningMethodHS512 {
						err = UnSupportSigningMethod
						return
					}
					ctx = jwtLocal.NewContext(ctx, tokenInfo.Claims)
				}
				return handler(ctx, req)
			}
			err = WrongContext
			return
		}
	}
}
