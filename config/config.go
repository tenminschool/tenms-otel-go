package config

type TenMsOtelConfig struct {
	ServiceName              string
	InsecureMode             string
	OtelExporterOtlpEndpoint string
	SamplingRatio            float64
}
