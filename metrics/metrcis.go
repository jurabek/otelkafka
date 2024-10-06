package metrics

import "go.opentelemetry.io/otel/metric"

type Metrics struct {
	meter    metric.Meter
	topLevel TopLevelMetrics
}

type TopLevelMetrics struct {
	clientAgeGauge  metric.Int64Gauge
	replyQueueGauge metric.Int64Gauge
	msgCountGauge   metric.Int64Gauge
	msgSizeGauge    metric.Int64Gauge
}

func topMetrics() {
	m := Metrics{}
	meter := m.meter

	clientAgeGauge, err := meter.Int64Gauge(
		"kafka.client.age",
		metric.WithDescription("Time since the client instance was created (microseconds)."),
		metric.WithUnit("µs"),
	)

	m.topLevel.clientAgeGauge = clientAgeGauge

	replyQueueGauge, err := meter.Int64Gauge(
		"kafka.client.reply_queue.size",
		metric.WithDescription("Number of ops (callbacks, events, etc) waiting in queue for application to serve with rd_kafka_poll()"),
		metric.WithUnit("{ops}"),
	)
	m.topLevel.replyQueueGauge = replyQueueGauge

	msgCountGauge, err := meter.Int64Gauge(
		"kafka.producer.msg_count",
		metric.WithDescription("Current number of messages in producer queues"),
		metric.WithUnit("{message}"),
	)
	m.topLevel.msgCountGauge = msgCountGauge

	msgSizeGauge, err := meter.Int64Gauge(
		"kafka.producer.msg_size",
		metric.WithDescription("Current total size of messages in producer queues"),
		metric.WithUnit("By"),
	)
	m.topLevel.msgSizeGauge = msgSizeGauge

	txCounter, err := meter.Int64Counter(
		"kafka.tx.total",
		metric.WithDescription("Total number of requests sent to Kafka brokers"),
		metric.WithUnit("{request}"),
	)
	_ = txCounter

	txBytesCounter, err := meter.Int64Counter(
		"kafka.tx.bytes.total",
		metric.WithDescription("Total number of bytes transmitted to Kafka brokers"),
		metric.WithUnit("By"),
	)
	_ = txBytesCounter

	rxCounter, err := meter.Int64Counter(
		"kafka.rx.total",
		metric.WithDescription("Total number of responses received from Kafka brokers"),
		metric.WithUnit("{response}"),
	)
	_ = rxCounter

	rxBytesCounter, err := meter.Int64Counter(
		"kafka.rx.bytes.total",
		metric.WithDescription("Total number of bytes received from Kafka brokers"),
		metric.WithUnit("By"),
	)
	_ = rxBytesCounter

	// missing
	/*
		msg_max	int		Threshold: maximum number of messages allowed allowed on the producer queues
		msg_size_max	int		Threshold: maximum total size of messages allowed on the producer queues
		txmsgs	int		Total number of messages transmitted (produced) to Kafka brokers
		txmsg_bytes	int		Total number of message bytes (including framing, such as per-Message framing and MessageSet/batch framing) transmitted to Kafka brokers
		rxmsgs	int		Total number of messages consumed, not including ignored messages (due to offset, etc), from Kafka brokers.
		rxmsg_bytes	int		Total number of message bytes (including framing) received from Kafka brokers
		simple_cnt	int gauge		Internal tracking of legacy vs new consumer API state
		metadata_cache_cnt	int gauge		Number of topics in the metadata cache.
	*/

	if err != nil {
		panic(err)
	}
}
