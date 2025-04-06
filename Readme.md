# TenMs Otel Go
This is a Go wrapper for OpenTelemetry that provides support for the OpenTelemetry API
## Installation
```bash
go get github.com/tenminschool/tenms-otel-go
```
## Environment Variables
| Name                       | Description                            | Default | Example        |
|----------------------------|----------------------------------------|---------|----------------|
| OTEL_EXPORTER_OTLP_ENDPOINT | The endpoint to send the trace data to |         | localhost:4317 |
| INSECURE_MODE | https or http                          |         | true           |
| DISABLE_OTEL        | To disable the tracing                 | `false` | true/false     |
| APP_NAME                   | The name of the application            |         | demo-service   |


## Basic Usage (Gin)
1. Add the following line in the `main()` function at main.go file before `routes.Register()` 
```go
 cleanup := otel.Boot(Config.GetString("App.Name"),
    Config.GetString("App.InsecureMode"),
    Config.GetString("App.OtlpExporterOtlpEndpoint")).
	 Init(artifact.Router, nil)
defer cleanup(context.TODO())
```
full `main` func will look like this:
```go
func main() {
    // Initialize the application
    artifact.New()
    config.Register() // will load the config file
    newrelic.Setup()

    artifact.NoSqlConnectionWithOtelMonitoring()
    rabbit_mq.RabbitMQConnection()

    cleanup := otel.Boot(Config.GetString("App.Name"),
        Config.GetString("App.InsecureMode"),
        Config.GetString("App.OtlpExporterOtlpEndpoint")).Init(
        artifact.Router,
        nil,
    )
    defer cleanup(context.TODO())

    artifact.Router.Use(otel.GetTenMsOtel().TraceMiddleware())

    config.Boot()     // if you need any initialization
    routes.Register() // will register all the routes
	
    artifact.Run()

}

```
3. Add the following env's in the .env file.
```go
APP_NAME=demo-service
INSECURE_MODE=true
OTEL_EXPORTER_OTLP_ENDPOINT=
DISABLE_OTEL=false
```