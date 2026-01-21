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

	auth "auth/api/auth"
	"auth/internal/biz"
	"auth/internal/conf"
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

// Permission validates JWT tokens and enforces RBAC permissions.
//
// Matching rules:
// - HTTP: METHOD + URI path
// - gRPC: operation ("/package.Service/Method")
//
// The Auth.Permission endpoint is treated specially:
// - public (/pub/) requests are always allowed (for auth_request)
// - when the Whitelist module is enabled, whitelist rules can short-circuit it
// - otherwise the service handler performs the actual permission check.
func Permission(c *conf.Bootstrap, rds redis.UniversalClient, permission *biz.PermissionUseCase, user *biz.UserUseCase, whitelist *biz.WhitelistUseCase) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			tr := otel.Tracer("middleware")
			ctx, span := tr.Start(ctx, "Permission")
			defer span.End()

			trans, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			operation := trans.Operation()
			// Always allow standard gRPC health checks.
			if strings.HasPrefix(operation, "/grpc.health.v1.Health/") {
				return handler(ctx, req)
			}

			// HTTP request info (if applicable).
			var (
				method string
				uri    string
			)
			if trans.Kind() == transport.KindHTTP {
				if ht, ok := trans.(kratosHttp.Transporter); ok && ht.Request() != nil && ht.Request().URL != nil {
					method = ht.Request().Method
					uri = ht.Request().URL.Path
				}
				// Public HTTP endpoints are always allowed.
				if strings.Contains(uri, pubURIPrefix) {
					return handler(ctx, req)
				}
				// Skip preflight/probe methods.
				if method == http.MethodOptions || method == http.MethodHead {
					return handler(ctx, req)
				}
				// Common health endpoint; keep it open.
				if strings.HasPrefix(uri, "/healthz") {
					return handler(ctx, req)
				}
			}

			// Special-case: /permission endpoint (used by nginx auth_request).
			if operation == auth.OperationAuthPermission {
				checkReq := permissionRequest(ctx, req)
				// Public resources are always allowed for permission checks.
				if strings.Contains(checkReq.URI, pubURIPrefix) ||
					checkReq.Method == http.MethodOptions ||
					checkReq.Method == http.MethodHead {
					return &emptypb.Empty{}, nil
				}

				// Permission whitelist can short-circuit the permission check endpoint.
				if permissionWhitelist(ctx, whitelist, checkReq) {
					return &emptypb.Empty{}, nil
				}
				// Optional: allow skipping JWT parsing via whitelist rules.
				if jwtWhitelist(ctx, whitelist, operation) {
					return handler(ctx, req)
				}

				var u *jwt.User
				u, err = parseJwt(ctx, c, rds, c.Server.Jwt.Key)
				if err != nil {
					return nil, err
				}
				// Pass user info into ctx for the service handler.
				ctx = jwt.NewServerContextByUser(ctx, *u)
				return handler(ctx, req)
			}

			// JWT whitelist is checked first (skip authentication & authorization).
			if jwtWhitelist(ctx, whitelist, operation) {
				return handler(ctx, req)
			}

			var u *jwt.User
			u, err = parseJwt(ctx, c, rds, c.Server.Jwt.Key)
			if err != nil {
				return nil, err
			}
			ctx = jwt.NewServerContextByUser(ctx, *u)

			// Permission whitelist skips permission checks (JWT still required).
			matchResource := operation
			if method != "" && uri != "" {
				// Prefer HTTP matching; keep grpc op as an optional 3rd segment.
				matchResource = strings.Join([]string{method, uri, operation}, "|")
			}
			ok2, _ := whitelist.Match(ctx, biz.WhitelistPermissionCategory, matchResource)
			if ok2 {
				return handler(ctx, req)
			}

			// Resolve userID from JWT attrs (code/platform) and check permissions.
			info := user.InfoFromCtx(ctx)
			if info == nil || info.Id == 0 {
				return nil, biz.ErrNoPermission(ctx)
			}

			resource := operation
			checkMethod := ""
			if trans.Kind() == transport.KindHTTP {
				resource = uri
				checkMethod = method
			}

			pass, perr := permission.CheckPermission(ctx, info.Id, resource, checkMethod)
			if perr != nil || !pass {
				return nil, biz.ErrNoPermission(ctx)
			}

			return handler(ctx, req)
		}
	}
}

func permissionRequest(ctx context.Context, req interface{}) (r biz.CheckPermission) {
	tr := otel.Tracer("middleware")
	ctx, span := tr.Start(ctx, "permissionRequest")
	defer span.End()

	copierx.Copy(&r, req)
	trans, _ := transport.FromServerContext(ctx)

	// Prefer nginx-provided headers when present.
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
		// Prefer HTTP matching; include grpc resource as an optional 3rd segment when present.
		matchResource = strings.Join([]string{r.Method, r.URI}, "|")
		if strings.TrimSpace(r.Resource) != "" {
			matchResource = strings.Join([]string{r.Method, r.URI, r.Resource}, "|")
		}
	}

	log.
		WithContext(ctx).
		Info("method: %s, uri: %s, resource: %s", r.Method, r.URI, r.Resource)

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

	if client != nil {
		key := strings.Join([]string{c.Name, jwtTokenCachePrefix, utils.StructMd5(user.Token)}, ".")
		res, _ := client.Get(ctx, key).Result()
		if res != "" {
			utils.JSON2Struct(user, res)
			return user, nil
		}

		// Parse token and cache derived user info.
		var info *jwtV4.Token
		info, err = parseToken(ctx, jwtKey, user.Token)
		if err != nil {
			return nil, err
		}
		ctx = jwt.NewServerContext(ctx, info.Claims, "code", "platform")
		user = jwt.FromServerContext(ctx)
		client.Set(ctx, key, utils.Struct2JSON(user), jwtTokenCacheExpire).Err()
		return user, nil
	}

	// No cache client; just parse the token.
	var info *jwtV4.Token
	info, err = parseToken(ctx, jwtKey, user.Token)
	if err != nil {
		return nil, err
	}
	ctx = jwt.NewServerContext(ctx, info.Claims, "code", "platform")
	user = jwt.FromServerContext(ctx)
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
