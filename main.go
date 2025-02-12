package tenms_otel_go

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tenminschool/tenms-otel-go/config"
	"github.com/tenminschool/tenms-otel-go/metrics"
	"github.com/tenminschool/tenms-otel-go/tracer"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
func (tenmsOtel *TenMsOtel) Init(Router *gin.Engine) func() {
	shutDownTracer := tracer.InitTracer(tenmsOtel.tenMsOtelConfig)

	meterProvider := metrics.InitMeter(tenmsOtel.tenMsOtelConfig)

	meter := meterProvider.Meter(tenmsOtel.tenMsOtelConfig.ServiceName)
	metrics.GenerateMetrics(meter)

	Router.Use(otelgin.Middleware(tenmsOtel.tenMsOtelConfig.ServiceName))
	return func() {
		if err := shutDownTracer(context.Background()); err != nil {
			fmt.Println("error in shut down tracer")
		}

		if err := meterProvider.Shutdown(context.Background()); err != nil {
			fmt.Println("error in shut down meterProvider")
		}
	}
}
