package middleware

import (
	"context"
	"errors"
	"strings"
	"time"

	"auth/internal/biz"
	"auth/internal/conf"
	jwtLocal "github.com/go-cinch/common/jwt"
	"github.com/go-cinch/common/utils"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	jwtV4 "github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

func Jwt(c *conf.Bootstrap, client redis.UniversalClient, whitelist *biz.WhitelistUseCase) middleware.Middleware {
	return selector.Server(
		jwtHandler(c, client),
	).Match(jwtWhitelist(whitelist)).Build()
}

const jwtTokenUserKey = "jwt_token_"

func jwtHandler(c *conf.Bootstrap, client redis.UniversalClient) func(handler middleware.Handler) middleware.Handler {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			user := jwtLocal.FromServerContext(ctx)
			if user.Token == "" && user.Code == "" {
				err = biz.ErrJwtMissingToken(ctx)
				return
			}
			if user.Code != "" {
				return handler(ctx, req)
			}
			token := user.Token
			if token != "" {
				key := strings.Join([]string{jwtTokenUserKey, utils.StructMd5(token)}, "")
				// read user info from cache
				res, e := client.Get(ctx, key).Result()
				if e == nil {
					utils.Json2Struct(user, res)
				} else {
					// parse Authorization jwt token to get user info
					var info *jwtV4.Token
					info, err = parseToken(ctx, c.Server.Jwt.Key, token)
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

func parseToken(ctx context.Context, key, jwtToken string) (info *jwtV4.Token, err error) {
	info, err = jwtV4.Parse(jwtToken, func(token *jwtV4.Token) (rp interface{}, err error) {
		rp = []byte(key)
		return
	})
	if err != nil {
		var ve *jwtV4.ValidationError
		ok := errors.As(err, &ve)
		if !ok {
			return
		}
		if ve.Errors&jwtV4.ValidationErrorMalformed != 0 {
			err = biz.ErrJwtTokenInvalid(ctx)
			return
		}
		if ve.Errors&(jwtV4.ValidationErrorExpired|jwtV4.ValidationErrorNotValidYet) != 0 {
			err = biz.ErrJwtTokenExpired(ctx)
			return
		}
		err = biz.ErrJwtTokenParseFail(ctx)
		return
	}
	if !info.Valid {
		err = biz.ErrJwtTokenParseFail(ctx)
		return
	}
	if info.Method != jwtV4.SigningMethodHS512 {
		err = biz.ErrJwtUnSupportSigningMethod(ctx)
		return
	}
	return
}
