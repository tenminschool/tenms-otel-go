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
1. Add the following line in the `main()` function at main.go file
```go
 cleanup := otel.Boot(os.Getenv("APP_NAME"),
    os.Getenv("INSECURE_MODE"),
    os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")).
	 Init(artifact.Router, nil)
defer cleanup(context.TODO())
```
2. add the following line before `routes.Register()`
```go
artifact.Router.Use(otel.GetTenMsOtel().TraceMiddleware())
```

full `main` func will look like this:
```go
func main() {
    // Initialize the application
    artifact.New()
    config.Register() // will load the config file
    newrelic.Setup()
	
    rabbit_mq.RabbitMQConnection()

    cleanup := otel.Boot(os.Getenv("APP_NAME"), 
		os.Getenv("INSECURE_MODE"), 
		os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")).
        Init(artifact.Router, nil)
    
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

## Get current span from gin context
1. Use the following line for getting current span context
```go
c.Request.Context()
```

## MongoDB Tracing
1. Add the following line while init mongo connection
```go
 artifact.NoSqlConnectionWithOtelMonitoring()
```
2. Pass the current span context whenever you are doing any operation with mongo
example
```go
result, err := models.VideoCollection.Collection.InsertOne(ctx, video)
````
here `ctx` is the current span context

## RabbitMQ Tracing
### Publisher
1. Add the following line before publishing to rabbitmq
```go
span, _ := tenmsOtel.TraceRabbitMqPublisher(exchange, routingKey, ctx)
defer span.End()
```
### Consumer
1. Add the following line before consuming from rabbitmq
```go
span, spanCtx := rabbitMqTrace.TraceRabbitMqConsumer(
				"queue name", //read it from env
				"ExamSessionConsumer", //consumer name
				context.TODO(), // pass the context if you have
			)
defer span.End()
```
## Tracing Http Client Calls
### Resty Client(Recommended)
1. Add the following line during resty client setup
```go
client := resty.New()
client.OnBeforeRequest(tenmsOtelhttpClientTrace.UseOtelBeforeRequestMiddleware(ctx)) // pass current span context
client.OnAfterResponse(tenmsOtelhttpClientTrace.UseOtelAfterResponseMiddleware())
client.OnError(tenmsOtelhttpClientTrace.UseOtelOnErrorHook())
```
### Nahid Http Client
This library does not support hooks. We fork the repo from and added the hooks. You can use that library for now.
1. Install the library
```go
go get github.com/tenminschool/gohttp
```
2. Replace all the `github.com/nahid/gohttp` import statements with `github.com/tenminschool/gohttp`
3. Add the following line during nahid http req setup
```go
req := gohttp.NewRequest()
req.OnBeforeRequest(nahid.UseOtelBeforeRequestHook(ctx)) // pass current span context
req.OnAfterResponse(nahid.UseOtelAfterResponseHook())
req.OnError(nahid.UseOtelOnErrorHook())
```
