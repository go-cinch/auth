package middleware

import (
	"context"

	"auth/api/auth"
	"auth/internal/biz"
	"github.com/go-cinch/common/idempotent"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func Idempotent(idt *idempotent.Idempotent) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			tr, _ := transport.FromServerContext(ctx)
			token := tr.RequestHeader().Get("x-idempotent")
			if token == "" || tr.Operation() == auth.OperationAuthPermission {
				// not have token, no need check
				return handler(ctx, req)
			}
			if !idt.Check(ctx, token) {
				err = biz.ErrIdempotentTokenExpired(ctx)
				return
			}
			return handler(ctx, req)
		}
	}
}
