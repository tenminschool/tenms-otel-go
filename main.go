package otel

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tenminschool/tenms-otel-go/config"
	"github.com/tenminschool/tenms-otel-go/metrics"
	"github.com/tenminschool/tenms-otel-go/tracer"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"os"
)

var tenMsOpenTelemetry *TenMsOtel

func GetTenMsOtel() *TenMsOtel {
	return tenMsOpenTelemetry
}

type TenMsOtel struct {
	config *config.TenMsOtelConfig
}

func Boot(serviceName string, insecureMode string, OtelExporterOtlpEndpoint string) *TenMsOtel {
	tenMsOpenTelemetry = &TenMsOtel{config: &config.TenMsOtelConfig{
		ServiceName:              serviceName,
		InsecureMode:             insecureMode,
		OtelExporterOtlpEndpoint: OtelExporterOtlpEndpoint,
	}}
	return tenMsOpenTelemetry
}

func (tenmsOtel *TenMsOtel) Init(
	Router *gin.Engine,
	db *gorm.DB,
) func(ctx context.Context) {
	shutDownTracer := tracer.InitTracer(tenmsOtel.config)

	meterProvider := metrics.InitMeter(tenmsOtel.config)

	meter := meterProvider.Meter(tenmsOtel.config.ServiceName)
	metrics.GenerateMetrics(meter)

	RegisterMiddleware(Router, tenmsOtel.config.ServiceName)

	if db != nil {
		if err := db.Use(tracing.NewPlugin()); err != nil {
			fmt.Println("error while connecting to db ", err.Error())
		}
	}

	fmt.Printf(
		"Tenms otel initialized with service_name %s, insecure mode %s, OtelExporterOtlpEndpoint %s\n",
		tenmsOtel.config.ServiceName,
		tenmsOtel.config.InsecureMode,
		tenmsOtel.config.OtelExporterOtlpEndpoint,
	)
	return func(ctx context.Context) {
		if err := shutDownTracer(ctx); err != nil {
			fmt.Println("error in shut down tracer")
		} else {
			fmt.Println("tracer shut down properly")
		}

		if err := meterProvider.Shutdown(ctx); err != nil {
			fmt.Println("error in shut down meterProvider")
		} else {
			fmt.Println("meter provider shut down properly")
		}
	}
}

func (tenmsOtel *TenMsOtel) TraceMiddleware() gin.HandlerFunc {
	disableOtel := os.Getenv("DISABLE_OTEL")
	if disableOtel == "true" {
		return func(c *gin.Context) {}
	}
	return func(c *gin.Context) {
		request := c.Request
		ctx := request.Context()

		// ref: https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/labstack/echo/otelecho/echo.go
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(request.Header))
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", request)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(request)...),
			oteltrace.WithAttributes(
				semconv.HTTPServerAttributesFromHTTPRequest(tenmsOtel.config.ServiceName, c.FullPath(), request)...),
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

func RegisterMiddleware(Router *gin.Engine, serviceName string) {
	Router.Use(otelgin.Middleware(serviceName))
}
