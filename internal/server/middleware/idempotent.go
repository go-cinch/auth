package middleware

import (
	"auth/api/reason"
	"auth/internal/biz"
	"auth/internal/pkg/idempotent"
	"context"
	"github.com/go-cinch/common/middleware/i18n"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
)

func Idempotent(idt *idempotent.Idempotent) middleware.Middleware {
	return selector.Server(
		func(handler middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
				if tr, ok := transport.FromServerContext(ctx); ok {
					token := tr.RequestHeader().Get("x-idempotent")
					if token != "" {
						if idt.Check(ctx, token) {
							return handler(ctx, req)
						}
						err = reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(biz.IdempotentTokenExpired))
						return
					}
				}
				err = reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(biz.IdempotentMissingToken))
				return
			}
		},
	).Match(idempotentBlacklist()).Build()
}
