package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func Header() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			tr, _ := transport.FromServerContext(ctx)
			switch tr.Kind() {
			case transport.KindHTTP:
				if tr.ReplyHeader() != nil {
					tr.ReplyHeader().Set("x-content-type-options", "nosniff")
					tr.ReplyHeader().Set("x-xss-protection", "1; mode=block")
					tr.ReplyHeader().Set("x-frame-options", "deny")
					tr.ReplyHeader().Set("cache-control", "private")
				}
			}
			return handler(ctx, req)
		}
	}
}
