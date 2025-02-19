package rabbit_mq_trace

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

func TraceRabbitMqConsumer(queueName string, consumerName string, ctx context.Context) (trace.Span, context.Context) {
	spanName := "RabbitMq Consumer"
	if consumerName != "" {
		spanName = consumerName
	}
	msgCtx, span := otel.Tracer("rabbitmq").Start(ctx, spanName,
		trace.WithAttributes(
			semconv.MessagingSystemKey.String("rabbitmq"),
			semconv.MessagingOperationKey.String("receive"),
			attribute.String("messaging.rabbitmq.queue", queueName),
			//attribute.String("messaging.message_id", d.MessageId),
			//attribute.Int("messaging.message_size", len(d.Body)),
		),
	)
	return span, msgCtx
}
