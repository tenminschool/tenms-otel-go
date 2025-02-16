package otel

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tenminschool/tenms-otel-go/config"
	"github.com/tenminschool/tenms-otel-go/metrics"
	"github.com/tenminschool/tenms-otel-go/tracer"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type TenMsOtel struct {
	tenMsOtelConfig *config.TenMsOtelConfig
}

func NewTenMsOtel(serviceName string, insecureMode string, OtelExporterOtlpEndpoint string) *TenMsOtel {
	return &TenMsOtel{tenMsOtelConfig: &config.TenMsOtelConfig{
		ServiceName:              serviceName,
		InsecureMode:             insecureMode,
		OtelExporterOtlpEndpoint: OtelExporterOtlpEndpoint,
	}}
}

func (tenmsOtel *TenMsOtel) Init(
	Router *gin.Engine,
	db *gorm.DB,
) func(ctx context.Context) {
	shutDownTracer := tracer.InitTracer(tenmsOtel.tenMsOtelConfig)

	meterProvider := metrics.InitMeter(tenmsOtel.tenMsOtelConfig)

	meter := meterProvider.Meter(tenmsOtel.tenMsOtelConfig.ServiceName)
	metrics.GenerateMetrics(meter)

	RegisterMiddleware(Router, tenmsOtel.tenMsOtelConfig.ServiceName)

	if db != nil {
		if err := db.Use(tracing.NewPlugin()); err != nil {
			fmt.Println("error while connecting to db ", err.Error())
		}
	}

	fmt.Printf(
		"Tenms otel initialized with service_name %s, insecure mode %s, OtelExporterOtlpEndpoint %s\n",
		tenmsOtel.tenMsOtelConfig.ServiceName,
		tenmsOtel.tenMsOtelConfig.InsecureMode,
		tenmsOtel.tenMsOtelConfig.OtelExporterOtlpEndpoint,
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

func RegisterMiddleware(Router *gin.Engine, serviceName string) {
	Router.Use(otelgin.Middleware(serviceName))
}
