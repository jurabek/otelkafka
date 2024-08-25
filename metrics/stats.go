package metrics

const (
	VALUE = "value"
	HELP  = "help"
	TYPE  = "type"
)

func getMappings() []map[string]interface{} {

	metrics := []map[string]interface{}{
		{
			"value": "msg_cnt",
			"help":  "Current number of messages in all queues.",
			"type":  "gauge",
		},
		{
			"value": "msg_size",
			"help":  "Current total size of messages in all queues.",
			"type":  "gauge",
		},
		{
			"value": "tx",
			"help":  "Total number of requests sent to brokers.",
			"type":  "counter",
		},
		{
			"value": "tx_bytes",
			"help":  "Total number of bytes sent to brokers.",
			"type":  "counter",
		},
		{
			"value": "rx",
			"help":  "Total number of responses received from brokers.",
			"type":  "counter",
		},
		{
			"value": "rx_bytes",
			"help":  "Total number of bytes received from brokers.",
			"type":  "counter",
		},
		{
			"value": "metadata_cache_cnt",
			"help":  "Number of topics in the metadata cache.",
			"type":  "gauge",
		},
		{
			"value": "txmsgs",
			"help":  "Total number of messages transmitted (produced) to Kafka brokers",
			"type":  "counter",
		},
		{
			"value": "txmsg_bytes",
			"help":  "Total number of message bytes (including framing, such as per-Message framing and MessageSet/batch framing) transmitted to Kafka brokers",
			"type":  "counter",
		},
		{
			"value": "rxmsgs",
			"help":  "Total number of messages consumed, not including ignored messages (due to offset, etc), from Kafka brokers.",
			"type":  "counter",
		},
		{
			"value": "rxmsg_bytes",
			"help":  "Total number of message bytes (including framing) received from Kafka brokers",
			"type":  "counter",
		},
	}
	return metrics
}

func getWindowsStats() map[string]string {

	metrics := make(map[string]string)

	metrics["min"] = "Smallest value"
	metrics["max"] = "Largest value"
	metrics["avg"] = "Average value"
	metrics["sum"] = "Sum of values"
	metrics["cnt"] = "Number of values sampled"
	metrics["stddev"] = "Standard deviation (based on histogram)"
	metrics["hdrsize"] = "Memory size of Hdr Histogram"
	metrics["p50"] = "50th percentile"
	metrics["p75"] = "75th percentile"
	metrics["p90"] = "90th percentile"
	metrics["p95"] = "95th percentile"
	metrics["p99"] = "99th percentile"
	metrics["p99_99"] = "99.99th percentile"
	metrics["outofrange"] = "Values skipped due to out of histogram range"

	return metrics
}

func getTopicsMappings() []map[string]interface{} {

	metrics := []map[string]interface{}{
		{
			"value": "age",
			"help":  "Age of client's topic object (milliseconds)",
			"type":  "gauge",
		},
		{
			"value": "metadata_age",
			"help":  "Age of metadata from broker for this topic (milliseconds)",
			"type":  "gauge",
		},
		{
			"help":  "Batch sizes in bytes. See Window stats",
			"type":  "windowStats",
			"value": "batchsize",
		},
		{
			"help":  "Batch message counts. See Window stats",
			"type":  "windowStats",
			"value": "batchcnt",
		},
		{
			"help":   "Partitions dict, key is partition id.",
			"type":   "object",
			"value":  "partitions",
			"labels": []string{"partition", "broker", "leader"},
			"metrics": []map[string]interface{}{
				{
					"help":  "Number of messages waiting to be produced in first-level queue",
					"value": "msgq_cnt",
					"type":  "gauge",
				},
				{
					"help":  "Number of bytes in msgq_cnt",
					"value": "msgq_bytes",
					"type":  "gauge",
				},
				{
					"help":  "Number of messages ready to be produced in transmit queue",
					"value": "xmit_msgq_cnt",
					"type":  "gauge",
				},
				{
					"help":  "Number of bytes in xmit_msgq",
					"value": "xmit_msgq_bytes",
					"type":  "gauge",
				},
				{
					"help":  "Number of pre-fetched messages in fetch queue",
					"value": "fetchq_cnt",
					"type":  "gauge",
				},
				{
					"help":  "Bytes in fetchq",
					"value": "fetchq_size",
					"type":  "gauge",
				},
				{
					"help":  "Current/Last logical offset query",
					"value": "query_offset",
					"type":  "gauge",
				},
				{
					"help":  "Next offset to fetch",
					"value": "next_offset",
					"type":  "gauge",
				},
				{
					"help":  "Offset of last message passed to application + 1",
					"value": "app_offset",
					"type":  "gauge",
				},
				{
					"help":  "Offset to be committed",
					"value": "stored_offset",
					"type":  "gauge",
				},
				{
					"help":  "Partition leader epoch of stored offset",
					"value": "stored_leader_epoch",
					"type":  "counter",
				},
				{
					"help":  "Last committed offset",
					"value": "committed_offset",
					"type":  "gauge",
				},
				{
					"help":  "Partition leader epoch of committed offset",
					"value": "committed_leader_epoch",
					"type":  "counter",
				},
				{
					"help":  "Last PARTITION_EOF signaled offset",
					"value": "eof_offset",
					"type":  "gauge",
				},
				{
					"help":  "Partition's low watermark offset on broker",
					"value": "lo_offset",
					"type":  "gauge",
				},
				{
					"help":  "Partition's high watermark offset on broker",
					"value": "hi_offset",
					"type":  "gauge",
				},
				{
					"help":  "Partition's last stable offset on broker, or same as hi_offset is broker version is less than 0.11.0",
					"value": "ls_offset",
					"type":  "gauge",
				},
				{
					"help":  "Difference between (hi_offset or ls_offset) and committed_offset). hi_offset is used when isolation.level=read_uncommitted, otherwise ls_offset.",
					"value": "consumer_lag",
					"type":  "gauge",
				},
				{
					"help":  "Difference between (hi_offset or ls_offset) and stored_offset. See consumer_lag and stored_offset.",
					"value": "consumer_lag_stored",
					"type":  "gauge",
				},
				{
					"help":  "Last known partition leader epoch, or -1 if unknown.",
					"value": "leader_epoch",
					"type":  "counter",
				},
				{
					"help":  "Total number of messages transmitted (produced)",
					"value": "txmsgs",
					"type":  "counter",
				},
				{
					"help":  "Total number of bytes transmitted for txmsgs",
					"value": "txbytes",
					"type":  "counter",
				},
				{
					"help":  "Total number of messages consumed, not including ignored messages (due to offset, etc).",
					"value": "rxmsgs",
					"type":  "counter",
				},
				{
					"help":  "Total number of bytes received for rxmsgs",
					"value": "rxbytes",
					"type":  "counter",
				},
				{
					"help":  "Total number of messages received (consumer, same as rxmsgs), or total number of messages produced (possibly not yet transmitted) (producer).",
					"value": "msgs",
					"type":  "counter",
				},
				{
					"help":  "Dropped outdated messages",
					"value": "rx_ver_drops",
					"type":  "counter",
				},
				{
					"help":  "Current number of messages in-flight to/from broker",
					"value": "msgs_inflight",
					"type":  "gauge",
				},
				{
					"help":  "Next expected acked sequence (idempotent producer)",
					"value": "next_ack_seq",
					"type":  "gauge",
				},
				{
					"help":  "Next expected errored sequence (idempotent producer)",
					"value": "next_err_seq",
					"type":  "gauge",
				},
				{
					"help":  "Last acked internal message id (idempotent producer)",
					"value": "acked_msgid",
					"type":  "counter",
				},
			},
		},
	}
	return metrics
}

func getBrokerMappings() []map[string]interface{} {
	metrics := []map[string]interface{}{
		{
			"help":  "Time since last broker state change (microseconds)",
			"value": "stateage",
			"type":  "gauge",
		},
		{
			"help":  "Number of requests awaiting transmission to broker",
			"value": "outbuf_cnt",
			"type":  "gauge",
		},
		{
			"help":  "Number of messages awaiting transmission to broker",
			"value": "outbuf_msg_cnt",
			"type":  "gauge",
		},
		{
			"help":  "Number of requests in-flight to broker awaiting response",
			"value": "waitresp_cnt",
			"type":  "gauge",
		},
		{
			"help":  "Number of messages in-flight to broker awaiting response",
			"value": "waitresp_msg_cnt",
			"type":  "gauge",
		},
		{
			"help":  "Total number of requests sent",
			"value": "tx",
			"type":  "gauge",
		},
		{
			"help":  "Total number of bytes sent",
			"value": "txbytes",
			"type":  "gauge",
		},
		{
			"help":  "Total number of request retries",
			"value": "txretries",
			"type":  "gauge",
		},
		{
			"help":  "Total number of transmission errors",
			"value": "txerrs",
			"type":  "gauge",
		},
		{
			"value": "txidle",
			"help":  "Total number of transmission attempts during idle state.",
			"type":  "gauge",
		},
		{
			"value": "req_timeouts",
			"help":  "Total number of request timeouts.",
			"type":  "gauge",
		},
		{
			"value": "rx",
			"help":  "Total number of responses received.",
			"type":  "gauge",
		},
		{
			"value": "rxbytes",
			"help":  "Total number of bytes received.",
			"type":  "gauge",
		},
		{
			"value": "rxerrs",
			"help":  "Total number of reception errors.",
			"type":  "gauge",
		},
		{
			"value": "int_latency",
			"help":  "Internal producer queue latency in microseconds",
			"type":  "windowStats",
		},
		{
			"value": "outbuf_latency",
			"help":  "Internal request queue latency in microseconds. This is the time between a request is enqueued on the transmit (outbuf) queue and the time the request is written to the TCP socket. Additional buffering and latency may be incurred by the TCP stack and network",
			"type":  "windowStats",
		},
		{
			"value": "rtt",
			"help":  "Broker RTT histogram.",
			"type":  "windowStats",
		},
		{
			"value": "throttle",
			"help":  "Broker throttle time histogram.",
			"type":  "windowStats",
		},
	}

	return metrics
}

func getConsumerGroupMappings() []map[string]interface{} {
	metrics := []map[string]interface{}{
		{
			"help":  "Time elapsed since last state change (milliseconds).",
			"value": "stateage",
			"type":  "gauge",
		},
		{
			"help":  "Time elapsed since last rebalance (assign or revoke) (milliseconds).",
			"value": "rebalance_age",
			"type":  "gauge",
		},
		{
			"help":  "Total number of rebalances (assign or revoke).",
			"value": "rebalance_cnt",
			"type":  "counter",
		},
		{
			"help":  "Current assignment's partition count.",
			"value": "assignment_size",
			"type":  "gauge",
		},
	}
	return metrics
}

func getEOSMappings() []map[string]interface{} {
	metrics := []map[string]interface{}{
		{
			"help":  "Time elapsed since last idemp_state change (milliseconds).",
			"value": "idemp_stateage",
			"type":  "gauge",
		},
		{
			"help":  "Time elapsed since last txn_state change (milliseconds).",
			"value": "txn_stateage",
			"type":  "gauge",
		},
		{
			"help":  "The currently assigned Producer ID (or -1).",
			"value": "producer_id",
			"type":  "gauge",
		},
		{
			"help":  "The current epoch (or -1).",
			"value": "producer_epoch",
			"type":  "gauge",
		},
		{
			"help":  "The number of Producer ID assignments since start.",
			"value": "epoch_cnt",
			"type":  "counter",
		},
	}
	return metrics
}
