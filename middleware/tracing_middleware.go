package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TraceMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := c.Request
		ctx := request.Context()

		// ref: https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/labstack/echo/otelecho/echo.go
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(request.Header))
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", request)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(request)...),
			oteltrace.WithAttributes(
				semconv.HTTPServerAttributesFromHTTPRequest(serviceName, c.FullPath(), request)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}
		spanName := c.FullPath()
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", request.Method)
		}

		ctx, span := otel.Tracer("gin-http").Start(ctx, spanName, opts...)
		defer span.End()

		// pass the span through the request context
		c.Request = request.WithContext(ctx)

		c.Next()

		// Record response details
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.Int("http.response_size", c.Writer.Size()),
		)

		// Record error if any
		if len(c.Errors) > 0 {
			span.RecordError(c.Errors.Last().Err)
		}
	}
}
