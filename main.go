package tenms_otel_go

import (
	"context"
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
func (tenmsOtel *TenMsOtel) Init(Router *gin.Engine) {
	cleanup := tracer.InitTracer(tenmsOtel.tenMsOtelConfig)
	defer cleanup(context.Background())

	provider := metrics.InitMeter(tenmsOtel.tenMsOtelConfig)
	defer provider.Shutdown(context.Background())

	meter := provider.Meter(tenmsOtel.tenMsOtelConfig.ServiceName)
	metrics.GenerateMetrics(meter)

	Router.Use(otelgin.Middleware(tenmsOtel.tenMsOtelConfig.ServiceName))
}
