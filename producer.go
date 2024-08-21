package otelkafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type Producer struct {
	*kafka.Producer
	cfg config
}

func NewProducer(conf *kafka.ConfigMap, opts ...Option) (*Producer, error) {
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, err
	}
	cfg := newConfig(opts...)
	return &Producer{Producer: p, cfg: cfg}, nil
}

// Produce calls the underlying Producer.Produce and traces the request.
func (p *Producer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	span := p.startSpan(msg)

	// if the user has selected a delivery channel, we will wrap it and
	// wait for the delivery event to finish the span
	if deliveryChan != nil {
		oldDeliveryChan := deliveryChan
		deliveryChan = make(chan kafka.Event)
		go func() {
			var err error
			evt := <-deliveryChan
			if msg, ok := evt.(*kafka.Message); ok {
				// delivery errors are returned via TopicPartition.Error
				err = msg.TopicPartition.Error
			}
			span.RecordError(err)
			span.End()
			oldDeliveryChan <- evt
		}()
	}

	err := p.Producer.Produce(msg, deliveryChan)
	// with no delivery channel or enqueue error, finish immediately
	if err != nil || deliveryChan == nil {
		span.RecordError(err)
		span.End()
	}

	return err
}

// Close calls the underlying Producer.Close and also closes the internal
// wrapping producer channel.
func (p *Producer) Close() {
	p.Producer.Close()
}

func (p *Producer) startSpan(msg *kafka.Message) trace.Span {
	// If there's a span context in the message, use that as the parent context.
	carrier := NewMessageCarrier(msg)
	ctx := p.cfg.Propagators.Extract(context.Background(), carrier)

	attr := []attribute.KeyValue{
		semconv.MessagingOperationTypePublish,
		semconv.MessagingSystemKafka,
		semconv.ServerAddress(p.cfg.bootstrapServers),
		semconv.MessagingDestinationName(*msg.TopicPartition.Topic),
		semconv.MessagingKafkaMessageKey(string(msg.Key)),
		semconv.MessagingMessageBodySize(getMsgSize(msg)),
	}

	opts := []trace.SpanStartOption{
		trace.WithAttributes(attr...),
		trace.WithSpanKind(trace.SpanKindProducer),
	}

	ctx, span := p.cfg.Tracer.Start(ctx, fmt.Sprintf("%s publish", *msg.TopicPartition.Topic), opts...)
	p.cfg.Propagators.Inject(ctx, carrier)
	return span

}
