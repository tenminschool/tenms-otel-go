package nahid

import (
	"context"
	"fmt"
	"github.com/tenminschool/gohttp"
	"github.com/tenminschool/tenms-otel-go/trace/httpclient/intrumentation"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func UseOtelBeforeRequestHook(ctx context.Context) gohttp.BeforeRequestHook {
	return func(request *gohttp.Request) error {
		traceContext, _ := otel.Tracer("http-client").
			Start(ctx, "http-client", trace.WithSpanKind(trace.SpanKindClient))

		otel.GetTextMapPropagator().Inject(traceContext, nil)
		request.SetContext(traceContext)
		fmt.Println("Span Created")
		return nil
	}
}

func UseOtelAfterResponseHook() gohttp.AfterResponseHook {
	return func(request *gohttp.Request, response *gohttp.Response) error {
		span := trace.SpanFromContext(request.Context())
		intrumentation.InstrumentResponse(span, response.GetResp())
		span.End()
		return nil
	}
}

func UseOtelOnErrorHook() gohttp.ErrorHook {
	return func(request *gohttp.Request, err error) {
		span := trace.SpanFromContext(request.Context())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		span.End()
	}
}
