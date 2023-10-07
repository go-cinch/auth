package middleware

import (
	"context"
	"net/http"
	"strings"

	"auth/api/auth"
	"auth/internal/biz"
	"auth/internal/conf"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Whitelist(c *conf.Bootstrap, whitelist *biz.WhitelistUseCase) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				err = biz.ErrNoPermission(ctx)
				return
			}
			operation := tr.Operation()
			switch tr.Kind() {
			case transport.KindHTTP:
				var pass bool
				if operation == auth.OperationAuthPermission {
					// nginx auth_request
					v, ok2 := req.(*auth.PermissionRequest)
					if !ok2 {
						err = biz.ErrNoPermission(ctx)
						return
					}
					// check whitelist
					pass, err = hasPermissionWhitelist(ctx, whitelist, v)
					if err != nil {
						return
					}
					if pass {
						rp = &emptypb.Empty{}
						return
					}
					// not in whitelist
					return handler(ctx, req)
				}
				if !c.Server.Permission.Direct {
					// called OperationAuthPermission means has permission
					return handler(ctx, req)
				}
				// direct call other http api
				var v auth.PermissionRequest
				var method, path string
				if ht, ok3 := tr.(kratosHttp.Transporter); ok3 {
					method = ht.Request().Method
					path = strings.Join([]string{c.Server.Http.Path, ht.Request().URL.Path}, "")
					v.Method = &method
					v.Uri = &path
				}
				// check whitelist
				pass, err = hasPermissionWhitelist(ctx, whitelist, &v)
				if err != nil {
					return
				}
				if pass {
					return handler(ctx, req)
				}
			case transport.KindGRPC:
				var pass bool
				if operation == auth.OperationAuthPermission {
					// direct call /Permission
					v, ok2 := req.(*auth.PermissionRequest)
					if !ok2 {
						err = biz.ErrNoPermission(ctx)
						return
					}
					// check whitelist
					pass, err = hasPermissionWhitelist(ctx, whitelist, v)
					if err != nil {
						return
					}
					if pass {
						rp = &emptypb.Empty{}
						return
					}
					// not in whitelist
					return handler(ctx, req)
				}
				if !c.Server.Permission.Direct {
					// called OperationAuthPermission means has permission
					return handler(ctx, req)
				}
				// direct call other grpc api
				var v auth.PermissionRequest
				v.Resource = &operation
				// check whitelist
				pass, err = hasPermissionWhitelist(ctx, whitelist, &v)
				if err != nil {
					return
				}
				if pass {
					return handler(ctx, req)
				}
			}
			return handler(ctx, req)
		}
	}
}

func hasPermissionWhitelist(ctx context.Context, whitelist *biz.WhitelistUseCase, req *auth.PermissionRequest) (ok bool, err error) {
	var r biz.CheckPermission
	if req.Method != nil {
		r.Method = *req.Method
	}
	if req.Uri != nil {
		r.URI = *req.Uri
	}
	if req.Uri != nil {
		r.URI = *req.Uri
	}
	if req.Resource != nil {
		r.Resource = *req.Resource
	}
	// skip options
	if r.Method == http.MethodOptions {
		ok = true
		return
	}
	// check if it is on the whitelist
	ok, err = whitelist.Has(ctx, &biz.HasWhitelist{
		Category:   biz.WhitelistPermissionCategory,
		Permission: r,
	})
	return
}

func jwtWhitelist(whitelist *biz.WhitelistUseCase) selector.MatchFunc {
	return func(ctx context.Context, operation string) bool {
		pass, _ := whitelist.Has(ctx, &biz.HasWhitelist{
			Category: biz.WhitelistJwtCategory,
			Permission: biz.CheckPermission{
				Resource: operation,
			},
		})
		// has data means no need check jwt
		return !pass
	}
}

func idempotentBlacklist(whitelist *biz.WhitelistUseCase) selector.MatchFunc {
	return func(ctx context.Context, operation string) bool {
		pass, _ := whitelist.Has(ctx, &biz.HasWhitelist{
			Category: biz.WhitelistIdempotentCategory,
			Permission: biz.CheckPermission{
				Resource: operation,
			},
		})
		// has data means need check idempotent
		return pass
	}
}
