package nahid

import (
	"context"
	"github.com/nahid/gohttp"
	"github.com/tenminschool/tenms-otel-go/trace/httpclient/intrumentation"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func UseOtelBeforeRequestMiddleware(ctx context.Context) func(request *gohttp.Request) error {
	return func(request *gohttp.Request) error {
		traceContext, _ := otel.Tracer("http-client").
			Start(ctx, "http-client", trace.WithSpanKind(trace.SpanKindClient))

		otel.GetTextMapPropagator().Inject(traceContext, nil)
		ctx = context.WithValue(ctx, "traceContext", traceContext)
		return nil
	}
}

func UseOtelAfterResponseMiddleware(ctx context.Context) func(response *gohttp.Response) error {
	return func(response *gohttp.Response) error {
		span := trace.SpanFromContext(ctx.Value("traceContext").(context.Context))
		intrumentation.InstrumentResponse(span, response.GetResp())

		span.End()
		return nil
	}
}

func UseOtelOnErrorHook(ctx context.Context) func(request *gohttp.Request, err error) {
	return func(request *gohttp.Request, err error) {
		span := trace.SpanFromContext(ctx.Value("traceContext").(context.Context))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		span.End()
	}
}
