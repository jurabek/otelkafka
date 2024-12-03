package otelkafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jurabek/otelkafka/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// Producer supports only tracing mechanism for Produce method over deprecated ProduceChannel method
type Producer struct {
	*kafka.Producer
	cfg config

	statsEnabled bool
	metrics      *metrics.ProducerClientMetrics
}

func NewProducer(conf *kafka.ConfigMap, opts ...Option) (*Producer, error) {
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, err
	}
	opts = append(opts, withConfig(conf))
	cfg := newConfig("producer", opts...)

	if si, err := conf.Get("statistics.interval.ms", 0); err == nil && si != 0 {
		statsMetrics, err := metrics.NewProducerClientMetrics(cfg.Meter)
		if err != nil {
			return nil, fmt.Errorf("failed to get top level metrics: %w", err)
		}
		return &Producer{Producer: p, cfg: cfg, statsEnabled: true, metrics: statsMetrics}, nil
	}

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
			evt := <-deliveryChan
			switch e := evt.(type) {
			case *kafka.Message:
				if err := e.TopicPartition.Error; err != nil {
					span.RecordError(e.TopicPartition.Error)
					span.SetStatus(codes.Error, err.Error())
				}
			case *kafka.Stats:
				if !p.statsEnabled {
					break
				}

				var stats metrics.Stats
				err := json.Unmarshal([]byte(e.String()), &stats)
				if err != nil {
					fmt.Printf("Failed to unmarshal stats: %v\n", err)
				} else {
					metrics.ProducerStatsToMetrics(context.Background(), stats, p.metrics, metrics.Cfg{})
				}
			}
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

	var topicName string
	if msg.TopicPartition.Topic != nil {
		topicName = *msg.TopicPartition.Topic
	}

	attr := []attribute.KeyValue{
		semconv.MessagingOperationTypePublish,
		semconv.MessagingSystemKafka,
		semconv.ServerAddress(p.cfg.bootstrapServers),
		semconv.MessagingDestinationName(topicName),
		semconv.MessagingKafkaMessageKey(string(msg.Key)),
		semconv.MessagingMessageBodySize(getMsgSize(msg)),
	}

	opts := []trace.SpanStartOption{
		trace.WithAttributes(attr...),
		trace.WithSpanKind(trace.SpanKindProducer),
	}

	ctx, span := p.cfg.Tracer.Start(ctx, fmt.Sprintf("%s publish", topicName), opts...)
	p.cfg.Propagators.Inject(ctx, carrier)
	return span

}
