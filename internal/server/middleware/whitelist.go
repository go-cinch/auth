package middleware

import (
	"context"
	"net/http"
	"strings"

	"auth/api/auth"
	"auth/api/reason"
	"auth/internal/biz"
	"auth/internal/conf"
	"github.com/go-cinch/common/middleware/i18n"
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
				err = reason.ErrorForbidden(i18n.FromContext(ctx).T(biz.NoPermission))
				return
			}
			operation := tr.Operation()
			switch tr.Kind() {
			case transport.KindHTTP:
				if operation == auth.OperationAuthPermission {
					// nginx auth_request
					v, ok2 := req.(*auth.PermissionRequest)
					if !ok2 {
						err = reason.ErrorForbidden(i18n.FromContext(ctx).T(biz.NoPermission))
						return
					}
					// check whitelist
					if hasPermissionWhitelist(ctx, whitelist, v) {
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
				if hasPermissionWhitelist(ctx, whitelist, &v) {
					return handler(ctx, req)
				}
			case transport.KindGRPC:
				if operation == auth.OperationAuthPermission {
					// direct call /Permission
					v, ok2 := req.(*auth.PermissionRequest)
					if !ok2 {
						err = reason.ErrorForbidden(i18n.FromContext(ctx).T(biz.NoPermission))
						return
					}
					// check whitelist
					if hasPermissionWhitelist(ctx, whitelist, v) {
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
				if hasPermissionWhitelist(ctx, whitelist, &v) {
					return handler(ctx, req)
				}
			}
			return handler(ctx, req)
		}
	}
}

func hasPermissionWhitelist(ctx context.Context, whitelist *biz.WhitelistUseCase, req *auth.PermissionRequest) (ok bool) {
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
	ok = whitelist.Has(ctx, &biz.HasWhitelist{
		Category:   biz.WhitelistPermissionCategory,
		Permission: r,
	})
	return
}

func jwtWhitelist(whitelist *biz.WhitelistUseCase) selector.MatchFunc {
	return func(ctx context.Context, operation string) bool {
		// has data means no need check jwt
		return !whitelist.Has(ctx, &biz.HasWhitelist{
			Category: biz.WhitelistJwtCategory,
			Permission: biz.CheckPermission{
				Resource: operation,
			},
		})
	}
}

func idempotentBlacklist(whitelist *biz.WhitelistUseCase) selector.MatchFunc {
	return func(ctx context.Context, operation string) bool {
		// has data means need check idempotent
		return whitelist.Has(ctx, &biz.HasWhitelist{
			Category: biz.WhitelistIdempotentCategory,
			Permission: biz.CheckPermission{
				Resource: operation,
			},
		})
	}
}
