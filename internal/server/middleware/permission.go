package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/jwt"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"

	auth "auth/api/auth"
	"auth/internal/biz"
	"auth/internal/conf"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
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

// Permission middleware with gateway-trust model.
//
// Design:
//  1. gRPC and other HTTP requests: trust gateway, allow directly
//     (user info already in x-md-global-* headers, metadata.Server() handles it)
//  2. /permission endpoint: gateway auth_request, parse JWT and check permission
func Permission(c *conf.Bootstrap, rds redis.UniversalClient, whitelist *biz.WhitelistUseCase) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			tr := otel.Tracer("middleware")
			ctx, span := tr.Start(ctx, "Permission")
			defer span.End()

			trans, _ := transport.FromServerContext(ctx)
			operation := trans.Operation()

			// 1. gRPC and other HTTP requests: trust gateway, allow directly
			if trans.Kind() == transport.KindGRPC || operation != auth.OperationAuthPermission {
				return handler(ctx, req)
			}

			// 2. /permission endpoint: gateway auth_request
			return handlePermissionEndpoint(ctx, req, c, rds, whitelist, handler)
		}
	}
}

// handlePermissionEndpoint processes /permission requests from gateway auth_request.
func handlePermissionEndpoint(
	ctx context.Context,
	req interface{},
	c *conf.Bootstrap,
	rds redis.UniversalClient,
	whitelist *biz.WhitelistUseCase,
	handler middleware.Handler,
) (rp interface{}, err error) {
	tr := otel.Tracer("middleware")
	ctx, span := tr.Start(ctx, "handlePermissionEndpoint")
	defer span.End()

	checkReq := permissionRequest(ctx, req)

	// Public resources are always allowed
	if strings.Contains(checkReq.URI, pubURIPrefix) ||
		checkReq.Method == http.MethodOptions ||
		checkReq.Method == http.MethodHead {
		return &emptypb.Empty{}, nil
	}

	// Permission whitelist can short-circuit
	if permissionWhitelist(ctx, whitelist, checkReq) {
		return &emptypb.Empty{}, nil
	}

	// JWT whitelist: no need JWT verification (check the requested resource, not /permission itself)
	if jwtWhitelist(ctx, whitelist, checkReq.Resource) {
		return &emptypb.Empty{}, nil
	}

	// Parse JWT and pass to handler for actual permission check
	var u *jwt.User
	u, err = parseJwt(ctx, c, rds, c.Server.Jwt.Key)
	if err != nil {
		return nil, err
	}
	ctx = jwt.NewServerContextByUser(ctx, *u)
	return handler(ctx, req)
}

func permissionRequest(ctx context.Context, req interface{}) (r biz.CheckPermission) {
	tr := otel.Tracer("middleware")
	ctx, span := tr.Start(ctx, "permissionRequest")
	defer span.End()

	copierx.Copy(&r, req)
	trans, _ := transport.FromServerContext(ctx)

	// Prefer gateway-provided headers when present.
	if method := strings.TrimSpace(trans.RequestHeader().Get(permissionHeaderMethod)); method != "" {
		r.Method = method
	}
	if uri := strings.TrimSpace(trans.RequestHeader().Get(permissionHeaderURI)); uri != "" {
		r.URI = uri
	}

	// Mutate request for downstream handler (auth.Permission).
	if v, ok := req.(*auth.PermissionRequest); ok {
		if r.Method != "" {
			v.Method = &r.Method
		}
		if r.URI != "" {
			v.Uri = &r.URI
		}
	}
	return r
}

func permissionWhitelist(ctx context.Context, whitelist *biz.WhitelistUseCase, r biz.CheckPermission) (ok bool) {
	tr := otel.Tracer("middleware")
	ctx, span := tr.Start(ctx, "permissionWhitelist")
	defer span.End()

	matchResource := strings.TrimSpace(r.Resource)
	if strings.TrimSpace(r.Method) != "" && strings.TrimSpace(r.URI) != "" {
		matchResource = strings.Join([]string{r.Method, r.URI}, "|")
		if strings.TrimSpace(r.Resource) != "" {
			matchResource = strings.Join([]string{r.Method, r.URI, r.Resource}, "|")
		}
	}

	log.WithContext(ctx).Info("permission check: method=%s, uri=%s, resource=%s", r.Method, r.URI, r.Resource)

	ok, _ = whitelist.Match(ctx, biz.WhitelistPermissionCategory, matchResource)
	return ok
}

func jwtWhitelist(ctx context.Context, whitelist *biz.WhitelistUseCase, operation string) bool {
	tr := otel.Tracer("middleware")
	ctx, span := tr.Start(ctx, "jwtWhitelist")
	defer span.End()

	ok, _ := whitelist.Match(ctx, biz.WhitelistJwtCategory, operation)
	return ok
}

func parseJwt(ctx context.Context, c *conf.Bootstrap, client redis.UniversalClient, jwtKey string) (user *jwt.User, err error) {
	tr := otel.Tracer("middleware")
	ctx, span := tr.Start(ctx, "parseJwt")
	defer span.End()

	user = jwt.FromServerContext(ctx)
	if user.Token == "" {
		return nil, biz.ErrJwtMissingToken(ctx)
	}

	// Try cache first
	var cacheKey string
	if client != nil {
		cacheKey = strings.Join([]string{c.Name, jwtTokenCachePrefix, utils.StructMd5(user.Token)}, ".")
		if res, _ := client.Get(ctx, cacheKey).Result(); res != "" {
			utils.JSON2Struct(user, res)
			return user, nil
		}
	}

	// Parse token
	info, err := parseToken(ctx, jwtKey, user.Token)
	if err != nil {
		return nil, err
	}
	ctx = jwt.NewServerContext(ctx, info.Claims, "code", "platform")
	user = jwt.FromServerContext(ctx)

	// Cache result
	if client != nil {
		client.Set(ctx, cacheKey, utils.Struct2JSON(user), jwtTokenCacheExpire).Err()
	}
	return user, nil
}

func parseToken(ctx context.Context, key, jwtToken string) (info *jwtV4.Token, err error) {
	tr := otel.Tracer("middleware")
	ctx, span := tr.Start(ctx, "parseToken")
	defer span.End()

	info, err = jwtV4.Parse(jwtToken, func(_ *jwtV4.Token) (rp interface{}, err error) {
		return []byte(key), nil
	})
	if err != nil {
		var ve *jwtV4.ValidationError
		if !errors.As(err, &ve) {
			return nil, err
		}
		switch {
		case ve.Errors&jwtV4.ValidationErrorMalformed != 0:
			return nil, biz.ErrJwtTokenInvalid(ctx)
		case ve.Errors&(jwtV4.ValidationErrorExpired|jwtV4.ValidationErrorNotValidYet) != 0:
			return nil, biz.ErrJwtTokenExpired(ctx)
		default:
			return nil, biz.ErrJwtTokenParseFail(ctx)
		}
	}
	if !info.Valid {
		return nil, biz.ErrJwtTokenParseFail(ctx)
	}
	if info.Method != jwtV4.SigningMethodHS512 {
		return nil, biz.ErrJwtUnSupportSigningMethod(ctx)
	}
	return info, nil
}
