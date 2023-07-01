package middleware

import (
	"context"
	"strings"

	"auth/api/auth"
	"auth/api/reason"
	"auth/internal/biz"
	"auth/internal/conf"
	"github.com/go-cinch/common/jwt"
	"github.com/go-cinch/common/middleware/i18n"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
)

func Permission(c *conf.Bootstrap, permission *biz.PermissionUseCase) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				err = reason.ErrorForbidden(i18n.FromContext(ctx).T(biz.NoPermission))
				return
			}
			user := jwt.FromServerContext(ctx)
			operation := tr.Operation()
			if operation == auth.OperationAuthPermission {
				return handler(ctx, req)
			}
			if !c.Server.Permission.Direct {
				// called OperationAuthPermission means has permission
				return handler(ctx, req)
			}
			switch tr.Kind() {
			case transport.KindHTTP:
				// direct call other http api
				var method, path string
				if ht, ok3 := tr.(kratosHttp.Transporter); ok3 {
					method = ht.Request().Method
					path = strings.Join([]string{c.Server.Http.Path, ht.Request().URL.Path}, "")
				}
				// user code is empty
				if user.Code == "" {
					err = reason.ErrorForbidden(i18n.FromContext(ctx).T(biz.NoPermission))
					return
				}
				// has user, check permission
				pass := permission.Check(ctx, biz.CheckPermission{
					UserCode: user.Code,
					Method:   method,
					URI:      path,
				})
				if !pass {
					err = reason.ErrorForbidden(i18n.FromContext(ctx).T(biz.NoPermission))
					return
				}
				// has permission
			case transport.KindGRPC:
				// direct call other grpc api
				if user.Code == "" {
					err = reason.ErrorForbidden(i18n.FromContext(ctx).T(biz.NoPermission))
					return
				}
				// has user, check permission
				pass := permission.Check(ctx, biz.CheckPermission{
					UserCode: user.Code,
					Resource: operation,
				})
				if !pass {
					err = reason.ErrorForbidden(i18n.FromContext(ctx).T(biz.NoPermission))
					return
				}
				// has permission
			}
			return handler(ctx, req)
		}
	}
}
