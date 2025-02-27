package http_client_trace

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func OtelHttpSpanStart(ctx context.Context, method string) trace.Span {
	var tracer = otel.Tracer(method)
	_, span := tracer.Start(ctx, fmt.Sprintf("http-client: %s", method))
	return span
}

func OtelHttpSpanEnd(span trace.Span, url string, status int, message string) {
	span.SetAttributes(
		attribute.String("http.url", url),
		attribute.Int("http.status_code", status),
		attribute.String("http.request.message", message),
	)
}
