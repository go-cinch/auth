package data

import (
	"context"
	"github.com/go-cinch/common/middleware/i18n"
	"go.opentelemetry.io/otel/trace"
	"time"
)

// get default timeout ctx
func getDefaultTimeoutCtx(ctx context.Context) context.Context {
	translator := i18n.FromContext(ctx)
	span := trace.SpanFromContext(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
	ctx = i18n.NewContext(ctx, translator)
	ctx = trace.ContextWithSpan(ctx, span)
	return ctx
}
