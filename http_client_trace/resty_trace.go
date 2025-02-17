package http_client_trace

import (
	"context"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	"go.opentelemetry.io/otel/trace"
)

func UseOtelBeforeRequestMiddleware(ctx context.Context) resty.RequestMiddleware {
	return func(client *resty.Client, request *resty.Request) error {
		traceContext, _ := otel.Tracer("http-client").
			Start(ctx, "http-client", trace.WithSpanKind(trace.SpanKindClient))

		otel.GetTextMapPropagator().Inject(traceContext, propagation.HeaderCarrier(request.Header))
		request.SetContext(traceContext)
		return nil
	}
}

func UseOtelAfterResponseMiddleware() resty.ResponseMiddleware {
	return func(client *resty.Client, response *resty.Response) error {
		host := response.Request.RawRequest.Host

		span := trace.SpanFromContext(response.Request.Context())
		span.SetAttributes(httpconv.ClientResponse(response.RawResponse)...)
		span.SetAttributes(attribute.String("http.host", host))
		span.SetAttributes(attribute.String("http.url", response.Request.RawRequest.URL.String()))
		// Setting request attributes here since res.Request.RawRequest is nil
		// in onBeforeRequest.
		// span.SetAttributes(httpconv.ClientRequest(response.ro)...)
		// span.SetAttributes(attribute.String("http.url", response.Request.RawRequest.URL.String()))
		span.SetAttributes(httpconv.ClientRequest(response.Request.RawRequest)...)
		span.SetAttributes(attribute.String("http.path", response.Request.RawRequest.URL.Path))

		span.End()
		return nil
	}
}

func UseOtelOnErrorHook() resty.ErrorHook {
	return func(request *resty.Request, err error) {
		span := trace.SpanFromContext(request.Context())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		span.End()
	}
}
