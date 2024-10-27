package metrics

import (
	"fmt"

	"go.opentelemetry.io/otel/metric"
)

type TopLevelMetrics struct {
	ClientAgeGauge                     metric.Int64ObservableGauge
	ReplyQueueGauge                    metric.Int64ObservableGauge
	MsgCountGauge                      metric.Int64ObservableGauge
	MsgSizeGauge                       metric.Int64ObservableGauge
	RequestsSentTotal                  metric.Int64ObservableCounter
	RequestSentBytesTotal              metric.Int64ObservableCounter
	ResponseReceievedTotal             metric.Int64ObservableCounter
	ResponseReceievedBytesTotal        metric.Int64ObservableCounter
	TotalNumberMessagesProduced        metric.Int64ObservableCounter
	TotalNumberOfMessagesProducedBytes metric.Int64ObservableCounter
	TotalNumberOfMessagesConsumed      metric.Int64ObservableCounter
	TotalNumberOfMessagesConsumedBytes metric.Int64ObservableCounter
}

func GetTopLevelMetrics(meter metric.Meter) (*TopLevelMetrics, error) {
	var topLevel TopLevelMetrics
	var err error

	topLevel.ClientAgeGauge, err = meter.Int64ObservableGauge(
		"kafka.client.age",
		metric.WithDescription("Time since the client instance was created (microseconds)."),
		metric.WithUnit("µs"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.client.age failed: %w", err)
	}

	topLevel.ReplyQueueGauge, err = meter.Int64ObservableGauge(
		"kafka.client.reply_queue.size",
		metric.WithDescription("Number of ops (callbacks, events, etc) waiting in queue for application to serve with rd_kafka_poll()"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.client.reply_queue.size failed: %w", err)
	}

	topLevel.MsgCountGauge, err = meter.Int64ObservableGauge(
		"kafka.producer.msg_count",
		metric.WithDescription("Current number of messages in producer queues"),
		metric.WithUnit("{message}"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.producer.msg_count failed: %w", err)
	}

	topLevel.MsgSizeGauge, err = meter.Int64ObservableGauge(
		"kafka.producer.msg_size",
		metric.WithDescription("Current total size of messages in producer queues"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.producer.msg_size failed: %w", err)
	}

	topLevel.RequestsSentTotal, err = meter.Int64ObservableCounter(
		"kafka.requests.sent.total",
		metric.WithDescription("Total number of requests sent to Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.requests.sent.total failed: %w", err)
	}

	topLevel.RequestSentBytesTotal, err = meter.Int64ObservableCounter(
		"kafka.request.sent.bytes.total",
		metric.WithDescription("Total number of bytes transmitted to Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.request.sent.bytes.total failed: %w", err)
	}

	topLevel.ResponseReceievedTotal, err = meter.Int64ObservableCounter(
		"kafka.response.recieved.total",
		metric.WithDescription("Total number of responses received from Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.response.recieved.total failed: %w", err)
	}

	topLevel.ResponseReceievedBytesTotal, err = meter.Int64ObservableCounter(
		"kafka.response.recieved.bytes.total",
		metric.WithDescription("Total number of bytes received from Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.response.recieved.bytes.total failed: %w", err)
	}

	topLevel.TotalNumberMessagesProduced, err = meter.Int64ObservableCounter(
		"kafka.messages.produced.total",
		metric.WithDescription("Total number of messages transmitted (produced) to Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.messages.produced.total failed: %w", err)
	}

	topLevel.TotalNumberOfMessagesProducedBytes, err = meter.Int64ObservableCounter(
		"kafka.messages.produced.bytes",
		metric.WithDescription("Total number of message bytes (including framing, such as per-Message framing and MessageSet/batch framing) transmitted to Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.messages.produced.bytes failed: %w", err)
	}

	topLevel.TotalNumberOfMessagesConsumed, err = meter.Int64ObservableCounter(
		"kafka.messages.consumed.total",
		metric.WithDescription("Total number of messages consumed, not including ignored messages (due to offset, etc), from Kafka brokers."),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.messages.consumed.total failed: %w", err)
	}

	topLevel.TotalNumberOfMessagesConsumedBytes, err = meter.Int64ObservableCounter(
		"kafka.messages.consumed.bytes",
		metric.WithDescription("Total number of message bytes (including framing) received from Kafka brokers"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka.messages.consumed.bytes failed: %w", err)
	}

	return &topLevel, nil
}
