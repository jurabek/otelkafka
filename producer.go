package otelkafka

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// Producer supports only tracing mechanism for Produce method over deprecated ProduceChannel method
type Producer struct {
	*kafka.Producer
	cfg                         config
	msgCounter                  metric.Int64Counter
	clientOperationDurHistogram metric.Float64Histogram
}

func NewProducer(conf *kafka.ConfigMap, opts ...Option) (*Producer, error) {
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, err
	}

	opts = append(opts, withConfig(conf))
	cfg := newConfig(opts...)
	meter := cfg.MeterProvider.Meter("kafka_producer")

	// Stability: experimental https://opentelemetry.io/docs/specs/semconv/messaging/messaging-metrics/#metric-messagingclientsentmessages
	msgCounter, err := meter.Int64Counter(
		"messaging.client.sent.messages",
		metric.WithDescription("Number of messages producer attempted to send to the broker"),
		metric.WithUnit("{message}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create produced message counter metric: %w", err)
	}

	// Stability: experimental https://opentelemetry.io/docs/specs/semconv/messaging/messaging-metrics/#metric-messagingclientoperationduration
	clientOperationDurHistogram, err := meter.Float64Histogram(
		"messaging.client.operation.duration",
		metric.WithDescription("Duration of messaging operation initiated by a producer or consumer client."),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client.operation.duration: %w", err)
	}

	return &Producer{Producer: p, cfg: cfg, msgCounter: msgCounter, clientOperationDurHistogram: clientOperationDurHistogram}, nil
}

// Produce calls the underlying Producer.Produce and traces the request.
func (p *Producer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	return p.ProduceWithContext(context.Background(), msg, deliveryChan)
}

func (p *Producer) ProduceWithContext(context context.Context, msg *kafka.Message, deliveryChan chan kafka.Event) error {
	// If there's a span context in the message, use that as the parent context.
	carrier := NewMessageCarrier(msg)
	ctx := p.cfg.Propagators.Extract(context, carrier)
	start := time.Now()

	lowCardinalityAttrs := []attribute.KeyValue{
		semconv.MessagingOperationName("produce"),
		semconv.MessagingOperationTypePublish,
		semconv.MessagingSystemKafka,
		semconv.ServerAddress(p.cfg.bootstrapServers),
		semconv.MessagingDestinationName(*msg.TopicPartition.Topic),
	}

	highCardinalityAttrs := []attribute.KeyValue{
		semconv.MessagingKafkaMessageKey(string(msg.Key)),
		semconv.MessagingMessageBodySize(getMsgSize(msg)),
	}
	attrs := append(lowCardinalityAttrs, highCardinalityAttrs...)

	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindProducer),
	}

	ctx, span := p.cfg.Tracer.Start(ctx, fmt.Sprintf("%s publish", *msg.TopicPartition.Topic), opts...)
	p.cfg.Propagators.Inject(ctx, carrier)

	// if the user has selected a delivery channel, we will wrap it and
	// wait for the delivery event to finish the span
	if deliveryChan != nil {
		oldDeliveryChan := deliveryChan
		deliveryChan = make(chan kafka.Event)
		go func() {
			evt := <-deliveryChan
			if resMsg, ok := evt.(*kafka.Message); ok {
				if err := resMsg.TopicPartition.Error; err != nil {
					span.RecordError(resMsg.TopicPartition.Error)
					span.SetStatus(codes.Error, err.Error())
					span.SetAttributes(semconv.ErrorTypeKey.String("publish_error"))
					lowCardinalityAttrs = append(lowCardinalityAttrs, semconv.ErrorTypeKey.String("publish_error"))
				} else {
					if resMsg.TopicPartition.Partition >= 0 {
						partitionIDAttr := semconv.MessagingDestinationPartitionID(fmt.Sprintf("%d", resMsg.TopicPartition.Partition))
						lowCardinalityAttrs = append(lowCardinalityAttrs, partitionIDAttr)

						// update the span attribute with the partition ID and offset
						span.SetAttributes(partitionIDAttr)
						span.SetAttributes(semconv.MessagingKafkaMessageOffset(int(resMsg.TopicPartition.Offset)))
					}
				}
			}
			p.msgCounter.Add(ctx, 1, metric.WithAttributes(lowCardinalityAttrs...))
			p.clientOperationDurHistogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(lowCardinalityAttrs...))
			span.End()
			oldDeliveryChan <- evt
		}()
	}

	err := p.Producer.Produce(msg, deliveryChan)
	if err != nil {
		span.RecordError(err)
	}

	// with no delivery channel or enqueue error, finish immediately
	if deliveryChan == nil {
		p.msgCounter.Add(ctx, 1, metric.WithAttributes(lowCardinalityAttrs...))
		p.clientOperationDurHistogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(lowCardinalityAttrs...))
		span.End()
	}

	return err
}

// Close calls the underlying Producer.Close and also closes the internal
// wrapping producer channel.
func (p *Producer) Close() {
	p.Producer.Close()
}
