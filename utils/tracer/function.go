package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func GetTraceParent(ctx context.Context) string {
	propagator := otel.GetTextMapPropagator()
	carr := propagation.MapCarrier{}
	propagator.Inject(ctx, &carr)
	return carr.Get("traceparent")
}
