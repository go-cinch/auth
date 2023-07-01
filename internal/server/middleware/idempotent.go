package middleware

import (
	"context"

	"auth/api/reason"
	"auth/internal/biz"
	"auth/internal/pkg/idempotent"
	"github.com/go-cinch/common/middleware/i18n"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
)

func Idempotent(idt *idempotent.Idempotent, whitelist *biz.WhitelistUseCase) middleware.Middleware {
	return selector.Server(
		func(handler middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
				tr, ok := transport.FromServerContext(ctx)
				if !ok {
					err = reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(biz.IdempotentMissingToken))
					return
				}
				token := tr.RequestHeader().Get("x-idempotent")
				if token == "" {
					err = reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(biz.IdempotentMissingToken))
					return
				}
				if !idt.Check(ctx, token) {
					err = reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(biz.IdempotentTokenExpired))
					return
				}
				return handler(ctx, req)
			}
		},
	).Match(idempotentBlacklist(whitelist)).Build()
}
