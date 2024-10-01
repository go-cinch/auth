package middleware

import (
	"context"
	"net/http"

	"auth/api/auth"
	"auth/internal/biz"
	"github.com/go-cinch/common/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Whitelist(whitelist *biz.WhitelistUseCase) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				err = biz.ErrNoPermission(ctx)
				return
			}
			var pass bool
			operation := tr.Operation()
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
			}
			// not in whitelist or not auth.OperationAuthPermission
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
	log.WithContext(ctx).Info("method: %s, uri: %s, resource: %s", r.Method, r.URI, r.Resource)
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
