package middleware

import (
	"auth/api/auth"
	"auth/internal/biz"
	"auth/internal/conf"
	"auth/internal/service"
	"context"
	jwtLocal "github.com/go-cinch/common/jwt"
	"github.com/go-cinch/common/utils"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-redis/redis/v8"
	jwtV4 "github.com/golang-jwt/jwt/v4"
	"time"
)

func Permission(c *conf.Bootstrap, client redis.UniversalClient, svc *service.AuthService) middleware.Middleware {
	return selector.Server(
		jwt(c, client),
		func(handler middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
				// ok: call permission, no need check permission
				// not ok: not call permission, need check permission
				if _, ok := req.(*auth.PermissionRequest); !ok {
					var resource string
					if tr, ok2 := transport.FromServerContext(ctx); ok2 {
						resource = tr.Operation()
					}
					var res *auth.PermissionReply
					res, err = svc.Permission(ctx, &auth.PermissionRequest{
						Resource: resource,
					})
					if err != nil {
						return
					}
					if !res.Pass {
						err = biz.NoPermission
						return
					}
				}
				return handler(ctx, req)
			}
		},
	).Match(permissionWhitelist()).Build()
}

const jwtTokenUserKey = "jwt_token_"

func jwt(c *conf.Bootstrap, client redis.UniversalClient) func(handler middleware.Handler) middleware.Handler {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			user := jwtLocal.FromServerContext(ctx)
			if user.Token == "" && user.Code == "" {
				err = MissingToken
				return
			}
			token := user.Token
			if token != "" {
				key := jwtTokenUserKey + utils.StructMd5(token)
				// read user info from cache
				res, e := client.Get(ctx, key).Result()
				if e == nil {
					utils.Json2Struct(user, res)
				} else {
					// parse Authorization jwt token to get user info
					var info *jwtV4.Token
					info, err = parseToken(c.Auth.Jwt.Key, token)
					if err != nil {
						return
					}
					ctx = jwtLocal.NewServerContext(ctx, info.Claims)
					user = jwtLocal.FromServerContext(ctx)
					client.Set(ctx, key, utils.Struct2Json(user), time.Hour)
				}
			}
			ctx = jwtLocal.NewServerContextByUser(ctx, *user)
			return handler(ctx, req)
		}
	}
}

func parseToken(key, jwtToken string) (info *jwtV4.Token, err error) {
	info, err = jwtV4.Parse(jwtToken, func(token *jwtV4.Token) (rp interface{}, err error) {
		rp = []byte(key)
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
	if !info.Valid {
		err = TokenParseFail
		return
	}
	if info.Method != jwtV4.SigningMethodHS512 {
		err = UnSupportSigningMethod
		return
	}
	return
}
