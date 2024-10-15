package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"auth/api/auth"
	"auth/internal/biz"
	"auth/internal/conf"
	"github.com/go-cinch/common/copierx"
	jwtLocal "github.com/go-cinch/common/jwt"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	jwtV4 "github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	pubURIPrefix           = "/pub/"
	jwtTokenCachePrefix    = "jwt.token"
	jwtTokenCacheExpire    = 10 * time.Minute
	permissionHeaderMethod = "x-original-method"
	permissionHeaderURI    = "x-permission-uri"
)

func Permission(c *conf.Bootstrap, client redis.UniversalClient, whitelist *biz.WhitelistUseCase) func(handler middleware.Handler) middleware.Handler {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			tr, _ := transport.FromServerContext(ctx)
			if tr.Kind() == transport.KindGRPC {
				// 1. grpc api for internal no need check
				// tip: u can add ur logic if need check grpc
				return handler(ctx, req)
			}
			// check http api
			var uri string
			operation := tr.Operation()
			if ht, ok := tr.(kratosHttp.Transporter); ok {
				uri = ht.Request().URL.Path
			}
			// 2. public api no need check
			if strings.Contains(uri, pubURIPrefix) {
				return handler(ctx, req)
			}
			if operation == auth.OperationAuthPermission && permissionWhitelist(ctx, whitelist, req) {
				// for nginx auth_request
				// 3. permission whitelist api no need check
				rp = &emptypb.Empty{}
				return
			} else if jwtWhitelist(ctx, whitelist) {
				// 4. jwt whitelist api no need check
				return handler(ctx, req)
			}
			user, err := parseJwt(ctx, c, client, c.Server.Jwt.Key)
			if err != nil {
				return
			}
			// pass user info into ctx
			ctx = jwtLocal.NewServerContextByUser(ctx, *user)
			return handler(ctx, req)
		}
	}
}

func permissionWhitelist(ctx context.Context, whitelist *biz.WhitelistUseCase, req interface{}) (ok bool) {
	tr, _ := transport.FromServerContext(ctx)
	var r biz.CheckPermission
	copierx.Copy(&r, req)
	// get from header if exist
	method := tr.RequestHeader().Get(permissionHeaderMethod)
	if method != "" {
		r.Method = method
	}
	uri := tr.RequestHeader().Get(permissionHeaderURI)
	if uri != "" {
		r.URI = uri
	}
	// public api no need check
	if strings.Contains(r.URI, pubURIPrefix) {
		return true
	}
	log.
		WithContext(ctx).
		Info("method: %s, uri: %s, resource: %s", r.Method, r.URI, r.Resource)
	// skip options
	if r.Method == http.MethodOptions || r.Method == http.MethodHead {
		return
	}
	// check if it is on the whitelist
	ok = whitelist.Has(ctx, &biz.HasWhitelist{
		Category:   biz.WhitelistPermissionCategory,
		Permission: &r,
	})
	// override params
	v, ok2 := req.(*auth.PermissionRequest)
	if ok2 {
		v.Method = &r.Method
		v.Uri = &r.URI
		req = v
		return
	}
	return
}

func jwtWhitelist(ctx context.Context, whitelist *biz.WhitelistUseCase) bool {
	tr, _ := transport.FromServerContext(ctx)
	return whitelist.Has(ctx, &biz.HasWhitelist{
		Category: biz.WhitelistJwtCategory,
		Permission: &biz.CheckPermission{
			Resource: tr.Operation(),
		},
	})
}

func parseJwt(ctx context.Context, c *conf.Bootstrap, client redis.UniversalClient, jwtKey string) (user *jwtLocal.User, err error) {
	user = jwtLocal.FromServerContext(ctx)
	if user.Token == "" {
		err = biz.ErrJwtMissingToken(ctx)
		return
	}
	key := strings.Join([]string{c.Name, jwtTokenCachePrefix, utils.StructMd5(user.Token)}, ".")
	res, _ := client.Get(ctx, key).Result()
	if res != "" {
		utils.Json2Struct(user, res)
		return
	}

	// parse Authorization jwt token to get user info
	var info *jwtV4.Token
	info, err = parseToken(ctx, jwtKey, user.Token)
	if err != nil {
		return
	}
	ctx = jwtLocal.NewServerContext(ctx, info.Claims, "code", "platform")
	user = jwtLocal.FromServerContext(ctx)
	client.Set(ctx, key, utils.Struct2Json(user), jwtTokenCacheExpire)
	return
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
