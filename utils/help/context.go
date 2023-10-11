package help

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// NewCtxFromTraceCtx 重新生成带trace的context，防止父context的timeout影响
func NewCtxFromTraceCtx(ctx context.Context) context.Context {
	return trace.ContextWithSpanContext(context.Background(), trace.SpanContextFromContext(ctx))
}
