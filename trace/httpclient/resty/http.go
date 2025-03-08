package resty

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/tenminschool/tenms-otel-go/trace/httpclient/intrumentation"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
		span := trace.SpanFromContext(response.Request.Context())
		intrumentation.InstrumentResponse(span, response.RawResponse)
		span.End()
		return nil
	}
}

func UseOtelOnErrorHook() resty.ErrorHook {
	return func(request *resty.Request, err error) {
		span := trace.SpanFromContext(request.Context())
		intrumentation.InstrumentError(span, err)
		span.End()
	}
}
