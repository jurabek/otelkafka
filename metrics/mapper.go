package metrics

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func ConsumerStatsToMetrics(ctx context.Context, stats Stats, m *ConsumerClientMetrics, cfg Cfg) {
	attributes := []attribute.KeyValue{
		attribute.String("consumer_client_id", stats.ClientID),
		attribute.String("consumer_name", stats.Name),
		attribute.String("type", stats.Type),
	}
	attributes = append(attributes, cfg.Attributes...)

	m.TopLevel.ClientAgeGauge.Record(ctx, stats.Age, metric.WithAttributes(attributes...))
	m.TopLevel.ReplyQueueGauge.Record(ctx, stats.Replyq, metric.WithAttributes(attributes...))

	m.TopLevel.RequestsSentTotal.Record(ctx, stats.Tx, metric.WithAttributes(attributes...))
	m.TopLevel.RequestSentBytesTotal.Record(ctx, stats.TxBytes, metric.WithAttributes(attributes...))
	m.TopLevel.ResponseReceievedTotal.Record(ctx, stats.Rx, metric.WithAttributes(attributes...))
	m.TopLevel.ResponseReceievedBytesTotal.Record(ctx, stats.RxBytes, metric.WithAttributes(attributes...))

	m.Consumer.TotalNumberOfMessagesConsumed.Record(ctx, stats.Rxmsgs, metric.WithAttributes(attributes...))
	m.Consumer.TotalNumberOfMessagesConsumedBytes.Record(ctx, stats.RxmsgBytes, metric.WithAttributes(attributes...))

	recordConsumerGroupMetrics(ctx, stats.Cgrp, m.ConsumerGroupMetrics, attributes)
}

func ProducerStatsToMetrics(ctx context.Context, stats Stats, m *ProducerClientMetrics, cfg Cfg) {
	attributes := []attribute.KeyValue{
		attribute.String("producer_client_id", stats.ClientID),
		attribute.String("producer_name", stats.Name),
		attribute.String("type", stats.Type),
	}
	attributes = append(attributes, cfg.Attributes...)

	m.TopLevel.ClientAgeGauge.Record(ctx, stats.Age, metric.WithAttributes(attributes...))
	m.TopLevel.ReplyQueueGauge.Record(ctx, stats.Replyq, metric.WithAttributes(attributes...))

	m.TopLevel.RequestsSentTotal.Record(ctx, stats.Tx, metric.WithAttributes(attributes...))
	m.TopLevel.RequestSentBytesTotal.Record(ctx, stats.TxBytes, metric.WithAttributes(attributes...))
	m.TopLevel.ResponseReceievedTotal.Record(ctx, stats.Rx, metric.WithAttributes(attributes...))
	m.TopLevel.ResponseReceievedBytesTotal.Record(ctx, stats.RxBytes, metric.WithAttributes(attributes...))

	m.Producer.ProducerMsgQueueCountGauge.Record(ctx, stats.MsgCnt, metric.WithAttributes(attributes...))
	m.Producer.ProducerMsgQueueSizeGauge.Record(ctx, stats.MsgSize, metric.WithAttributes(attributes...))
	m.Producer.TotalNumberOfMessagesProduced.Record(ctx, stats.Txmsgs, metric.WithAttributes(attributes...))
	m.Producer.TotalNumberOfMessagesProducedBytes.Record(ctx, stats.TxmsgBytes, metric.WithAttributes(attributes...))
}

func recordConsumerGroupMetrics(ctx context.Context, s Cgrp, cm *ConsumerGroupMetrics, attr []attribute.KeyValue) {
	if s.State == "UP" {
		cm.Up.Record(ctx, 1, metric.WithAttributes(attr...))
	}
	cm.StateAge.Record(ctx, s.Stateage, metric.WithAttributes(attr...))

	if s.JoinState != "" {
		attr = append(attr, attribute.String("join_state", s.JoinState))
		cm.JoinState.Record(ctx, 1, metric.WithAttributes(attr...))
	}

	cm.RebalanceAge.Record(ctx, s.RebalanceAge, metric.WithAttributes(attr...))
	cm.RebalanceCount.Record(ctx, s.RebalanceCnt, metric.WithAttributes(attr...))

	if s.RebalanceReason != "" {
		attr = append(attr, attribute.String("rebalance_reason", s.RebalanceReason))
		cm.RebalanceReason.Record(ctx, 1, metric.WithAttributes(attr...))
	}
	cm.AssignmentSize.Record(ctx, s.AssignmentSize, metric.WithAttributes(attr...))
}
