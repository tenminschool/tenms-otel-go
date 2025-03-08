package intrumentation

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func InstrumentResponse(span trace.Span, response *http.Response) {
	methodName := response.Request.Method
	span.SetName("http-client:" + methodName)

	span.SetAttributes(httpconv.ClientResponse(response)...)
	span.SetAttributes(attribute.String("http.host", response.Request.URL.Host))
	span.SetAttributes(attribute.String("http.url", response.Request.URL.String()))
	span.SetAttributes(httpconv.ClientRequest(response.Request)...)
	span.SetAttributes(attribute.String("http.path", response.Request.URL.Path))
}

func InstrumentError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
