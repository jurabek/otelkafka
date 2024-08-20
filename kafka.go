package otelkafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"strconv"
	"time"
)

type Consumer struct {
	*kafka.Consumer
	cfg  config
	prev trace.Span
}

func WrapConsumer(c *kafka.Consumer, cfg config) *Consumer {
	return &Consumer{Consumer: c, cfg: cfg}
}

func (c *Consumer) Pool(timeout int) kafka.Event {
	if c.prev != nil {
		c.prev.End()
	}
	event := c.Consumer.Poll(timeout)
	switch e := event.(type) {
	case *kafka.Message:
		span := c.startSpan(e)
		c.prev = span
	}

	return event
}

// ReadMessage polls the consumer for a message. Message will be traced.
func (c *Consumer) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	if c.prev != nil {
		c.prev.End()
		c.prev = nil
	}
	msg, err := c.Consumer.ReadMessage(timeout)
	if err != nil {
		return nil, err
	}
	c.prev = c.startSpan(msg)
	return msg, nil
}

// Close calls the underlying Consumer.Close and if polling is enabled, finishes
// any remaining span.
func (c *Consumer) Close() error {
	err := c.Consumer.Close()
	// we only close the previous span if consuming via the events channel is
	// not enabled, because otherwise there would be a data race from the
	// consuming goroutine.
	if c.prev != nil {
		c.prev.End()
		c.prev = nil
	}
	return err
}

func (c *Consumer) startSpan(msg *kafka.Message) trace.Span {
	carrier := NewMessageCarrier(msg)
	parentSpanContext := c.cfg.Propagators.Extract(context.Background(), carrier)

	// Create a span.
	attrs := []attribute.KeyValue{
		semconv.MessagingOperationTypeReceive,
		semconv.MessagingSystemKafka,
		semconv.MessagingKafkaMessageOffset(int(msg.TopicPartition.Offset)),
		semconv.MessagingKafkaConsumerGroup(c.cfg.consumerGroupID),
		semconv.ServerAddress(c.cfg.bootstrapServers),
		semconv.MessagingDestinationName(*msg.TopicPartition.Topic),
		semconv.MessagingMessageID(strconv.FormatInt(int64(msg.TopicPartition.Offset), 10)),
		semconv.MessagingDestinationPartitionID(strconv.Itoa(int(msg.TopicPartition.Partition))),
		semconv.MessagingKafkaMessageKey(string(msg.Key)),
		semconv.MessagingMessageBodySize(getMsgSize(msg)),
	}
	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	}
	newCtx, span := c.cfg.Tracer.Start(parentSpanContext, fmt.Sprintf("%v receive", msg.TopicPartition.Topic), opts...)

	// Inject current span context, so consumers can use it to propagate span.
	c.cfg.Propagators.Inject(newCtx, carrier)
	return span
}

func getMsgSize(msg *kafka.Message) (size int) {
	for _, header := range msg.Headers {
		size += len(header.Key) + len(header.Value)
	}
	return size + len(msg.Value) + len(msg.Key)
}
