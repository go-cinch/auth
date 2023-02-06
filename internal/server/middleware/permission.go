package middleware

import (
	"auth/api/auth"
	"auth/api/reason"
	"auth/internal/biz"
	"auth/internal/service"
	"context"
	"github.com/go-cinch/common/middleware/i18n"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
)

func Permission(svc *service.AuthService) middleware.Middleware {
	return selector.Server(
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
						err = reason.ErrorForbidden(i18n.FromContext(ctx).T(biz.NoPermission))
						return
					}
				}
				return handler(ctx, req)
			}
		},
	).Match(permissionWhitelist()).Build()
}
