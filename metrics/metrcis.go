package metrics

import (
	"fmt"

	"go.opentelemetry.io/otel/metric"
)

type TopLevelMetrics struct {
	ClientAgeGauge              metric.Int64Gauge
	ReplyQueueGauge             metric.Int64Gauge
	RequestsSentTotal           metric.Int64Gauge
	RequestSentBytesTotal       metric.Int64Gauge
	ResponseReceievedTotal      metric.Int64Gauge
	ResponseReceievedBytesTotal metric.Int64Gauge
}

type ConsumerMetrics struct {
	TotalNumberOfMessagesConsumed      metric.Int64Gauge
	TotalNumberOfMessagesConsumedBytes metric.Int64Gauge
}

type ConsumerGroupMetrics struct {
	Up              metric.Int64Gauge
	JoinState       metric.Int64Gauge
	StateAge        metric.Int64Gauge
	RebalanceAge    metric.Int64Gauge
	RebalanceCount  metric.Int64Gauge
	RebalanceReason metric.Int64Gauge
	AssignmentSize  metric.Int64Gauge
}

type ProducerMetrics struct {
	ProducerMsgQueueCountGauge         metric.Int64Gauge
	ProducerMsgQueueSizeGauge          metric.Int64Gauge
	TotalNumberOfMessagesProduced      metric.Int64Gauge
	TotalNumberOfMessagesProducedBytes metric.Int64Gauge
}

type ConsumerClientMetrics struct {
	TopLevel             *TopLevelMetrics
	Consumer             *ConsumerMetrics
	ConsumerGroupMetrics *ConsumerGroupMetrics
}

func NewConsumerClientMetrics(meter metric.Meter) (*ConsumerClientMetrics, error) {
	var consumerMetrics ConsumerMetrics
	var err error

	consumerMetrics.TotalNumberOfMessagesConsumed, err = meter.Int64Gauge(
		"kafka.messages.consumed.total",
		metric.WithDescription("Total number of messages consumed, not including ignored messages (due to offset, etc), from Kafka brokers."),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.messages.consumed.total failed: %w", err)
	}

	consumerMetrics.TotalNumberOfMessagesConsumedBytes, err = meter.Int64Gauge(
		"kafka.messages.consumed.bytes",
		metric.WithDescription("Total number of message bytes (including framing) received from Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.messages.consumed.bytes failed: %w", err)
	}

	cgrpMetris, err := newConsumerGroupMetrics(meter)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group metrics: %w", err)
	}

	topLevel, err := NewTopLevelMetrics(meter)
	if err != nil {
		return nil, fmt.Errorf("failed to create top level metrics: %w", err)
	}
	result := &ConsumerClientMetrics{
		TopLevel:             topLevel,
		Consumer:             &consumerMetrics,
		ConsumerGroupMetrics: cgrpMetris,
	}

	return result, nil
}

func newConsumerGroupMetrics(meter metric.Meter) (*ConsumerGroupMetrics, error) {
	var consumerGroupMetrics ConsumerGroupMetrics
	var err error

	consumerGroupMetrics.Up, err = meter.Int64Gauge(
		"kafka.consumer.group.up",
		metric.WithDescription("Consumer group up"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.consumer.group.up: %w", err)
	}

	consumerGroupMetrics.JoinState, err = meter.Int64Gauge(
		"kafka.consumer.group.join_state",
		metric.WithDescription("Consumer group join state"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.consumer.group.join_state: %w", err)
	}

	consumerGroupMetrics.StateAge, err = meter.Int64Gauge(
		"kafka.consumer.group.state_age",
		metric.WithDescription("Consumer group state age"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.consumer.group.state_age: %w", err)
	}

	consumerGroupMetrics.RebalanceAge, err = meter.Int64Gauge(
		"kafka.consumer.group.rebalance_age",
		metric.WithDescription("Consumer group rebalance age"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.consumer.group.rebalance_age: %w", err)
	}
	consumerGroupMetrics.RebalanceCount, err = meter.Int64Gauge(
		"kafka.consumer.group.rebalance_count",
		metric.WithDescription("Consumer group rebalance count"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.consumer.group.rebalance_count: %w", err)
	}
	consumerGroupMetrics.RebalanceReason, err = meter.Int64Gauge(
		"kafka.consumer.group.rebalance_reason",
		metric.WithDescription("Consumer group rebalance reason"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.consumer.group.rebalance_reason: %w", err)
	}
	consumerGroupMetrics.AssignmentSize, err = meter.Int64Gauge(
		"kafka.consumer.group.assignment_size",
		metric.WithDescription("Consumer group assignment size"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.consumer.group.assignment_size: %w", err)
	}

	return &consumerGroupMetrics, nil
}

func NewTopLevelMetrics(meter metric.Meter) (*TopLevelMetrics, error) {
	var topLevel TopLevelMetrics
	var err error

	topLevel.ClientAgeGauge, err = meter.Int64Gauge("kafka.client.age", metric.WithDescription("Time since the client instance was created (microseconds)."), metric.WithUnit("µs"))
	if err != nil {
		return nil, fmt.Errorf("kafka.client.age failed: %w", err)
	}

	topLevel.ReplyQueueGauge, err = meter.Int64Gauge(
		"kafka.client.reply_queue.size",
		metric.WithDescription("Number of ops (callbacks, events, etc) waiting in queue for application to serve with rd_kafka_poll()"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.client.reply_queue.size failed: %w", err)
	}

	topLevel.RequestsSentTotal, err = meter.Int64Gauge(
		"kafka.requests.sent.total",
		metric.WithDescription("Total number of requests sent to Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.requests.sent.total failed: %w", err)
	}

	topLevel.RequestSentBytesTotal, err = meter.Int64Gauge(
		"kafka.request.sent.bytes.total",
		metric.WithDescription("Total number of bytes transmitted to Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.request.sent.bytes.total failed: %w", err)
	}

	topLevel.ResponseReceievedTotal, err = meter.Int64Gauge(
		"kafka.response.recieved.total",
		metric.WithDescription("Total number of responses received from Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.response.recieved.total failed: %w", err)
	}

	topLevel.ResponseReceievedBytesTotal, err = meter.Int64Gauge(
		"kafka.response.recieved.bytes.total",
		metric.WithDescription("Total number of bytes received from Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.response.recieved.bytes.total failed: %w", err)
	}
	return &topLevel, nil
}

type ProducerClientMetrics struct {
	TopLevel *TopLevelMetrics
	Producer *ProducerMetrics
}

func NewProducerClientMetrics(meter metric.Meter) (*ProducerMetrics, error) {
	var producerMetrics ProducerMetrics
	var err error
	producerMetrics.ProducerMsgQueueCountGauge, err = meter.Int64Gauge(
		"kafka.producer.queue.msg_count",
		metric.WithDescription("Current number of messages in producer queues"),
		metric.WithUnit("{message}"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.producer.queue.msg_count failed: %w", err)
	}

	producerMetrics.ProducerMsgQueueSizeGauge, err = meter.Int64Gauge(
		"kafka.producer.queue.msg_size",
		metric.WithDescription("Current total size of messages in producer queues"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.producer.queue.msg_size failed: %w", err)
	}

	producerMetrics.TotalNumberOfMessagesProduced, err = meter.Int64Gauge(
		"kafka.messages.produced.total",
		metric.WithDescription("Total number of messages transmitted (produced) to Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.messages.produced.total failed: %w", err)
	}

	producerMetrics.TotalNumberOfMessagesProducedBytes, err = meter.Int64Gauge(
		"kafka.messages.produced.bytes",
		metric.WithDescription("Total number of message bytes (including framing, such as per-Message framing and MessageSet/batch framing) transmitted to Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.messages.produced.bytes failed: %w", err)
	}

	return &producerMetrics, nil
}
