package tracing

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const errAttributeKey = "error"

func SetSpanError(span trace.Span, err error) {
	if err != nil {
		span.SetAttributes(attribute.String(errAttributeKey, err.Error()))
	}
}
