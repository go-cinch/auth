package middleware

import (
	"context"

	"auth/internal/biz"
	"auth/internal/pkg/idempotent"
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
					err = biz.ErrIdempotentMissingToken(ctx)
					return
				}
				token := tr.RequestHeader().Get("x-idempotent")
				if token == "" {
					err = biz.ErrIdempotentMissingToken(ctx)
					return
				}
				if !idt.Check(ctx, token) {
					err = biz.ErrIdempotentTokenExpired(ctx)
					return
				}
				return handler(ctx, req)
			}
		},
	).Match(idempotentBlacklist(whitelist)).Build()
}
