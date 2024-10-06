package metrics

type Stats struct {
	Name             string  `json:"name"`
	ClientID         string  `json:"client_id"`
	Type             string  `json:"type"`
	Ts               int64   `json:"ts"`
	Time             int64   `json:"time"`
	Age              int64   `json:"age"`
	Replyq           int64   `json:"replyq"`
	MsgCnt           int64   `json:"msg_cnt"`
	MsgSize          int64   `json:"msg_size"`
	MsgMax           int64   `json:"msg_max"`
	MsgSizeMax       int64   `json:"msg_size_max"`
	SimpleCnt        int64   `json:"simple_cnt"`
	MetadataCacheCnt int64   `json:"metadata_cache_cnt"`
	Brokers          Brokers `json:"brokers"`
	Topics           Topics  `json:"topics"`
	Cgrp             Cgrp    `json:"cgrp"`
	Tx               int64   `json:"tx"`
	TxBytes          int64   `json:"tx_bytes"`
	Rx               int64   `json:"rx"`
	RxBytes          int64   `json:"rx_bytes"`
	Txmsgs           int64   `json:"txmsgs"`
	TxmsgBytes       int64   `json:"txmsg_bytes"`
	Rxmsgs           int64   `json:"rxmsgs"`
	RxmsgBytes       int64   `json:"rxmsg_bytes"`
}

type Brokers struct {
	Broker290921     Broker           `json:"broker:29092/1"`
	GroupCoordinator GroupCoordinator `json:"GroupCoordinator"`
}

type Broker struct {
	Name           string               `json:"name"`
	Nodeid         int64                `json:"nodeid"`
	Nodename       string               `json:"nodename"`
	Source         string               `json:"source"`
	State          string               `json:"state"`
	Stateage       int64                `json:"stateage"`
	OutbufCnt      int64                `json:"outbuf_cnt"`
	OutbufMsgCnt   int64                `json:"outbuf_msg_cnt"`
	WaitrespCnt    int64                `json:"waitresp_cnt"`
	WaitrespMsgCnt int64                `json:"waitresp_msg_cnt"`
	Tx             int64                `json:"tx"`
	Txbytes        int64                `json:"txbytes"`
	Txerrs         int64                `json:"txerrs"`
	Txretries      int64                `json:"txretries"`
	Txidle         int64                `json:"txidle"`
	ReqTimeouts    int64                `json:"req_timeouts"`
	Rx             int64                `json:"rx"`
	Rxbytes        int64                `json:"rxbytes"`
	Rxerrs         int64                `json:"rxerrs"`
	Rxcorriderrs   int64                `json:"rxcorriderrs"`
	Rxpartial      int64                `json:"rxpartial"`
	Rxidle         int64                `json:"rxidle"`
	ZbufGrow       int64                `json:"zbuf_grow"`
	BufGrow        int64                `json:"buf_grow"`
	Wakeups        int64                `json:"wakeups"`
	Connects       int64                `json:"connects"`
	Disconnects    int64                `json:"disconnects"`
	IntLatency     map[string]int64     `json:"int_latency"`
	OutbufLatency  map[string]int64     `json:"outbuf_latency"`
	Rtt            map[string]int64     `json:"rtt"`
	Throttle       map[string]int64     `json:"throttle"`
	Req            map[string]int64     `json:"req"`
	Toppars        Broker290921_Toppars `json:"toppars"`
}

type GroupCoordinator struct {
	Name           string                  `json:"name"`
	Nodeid         int64                   `json:"nodeid"`
	Nodename       string                  `json:"nodename"`
	Source         string                  `json:"source"`
	State          string                  `json:"state"`
	Stateage       int64                   `json:"stateage"`
	OutbufCnt      int64                   `json:"outbuf_cnt"`
	OutbufMsgCnt   int64                   `json:"outbuf_msg_cnt"`
	WaitrespCnt    int64                   `json:"waitresp_cnt"`
	WaitrespMsgCnt int64                   `json:"waitresp_msg_cnt"`
	Tx             int64                   `json:"tx"`
	Txbytes        int64                   `json:"txbytes"`
	Txerrs         int64                   `json:"txerrs"`
	Txretries      int64                   `json:"txretries"`
	Txidle         int64                   `json:"txidle"`
	ReqTimeouts    int64                   `json:"req_timeouts"`
	Rx             int64                   `json:"rx"`
	Rxbytes        int64                   `json:"rxbytes"`
	Rxerrs         int64                   `json:"rxerrs"`
	Rxcorriderrs   int64                   `json:"rxcorriderrs"`
	Rxpartial      int64                   `json:"rxpartial"`
	Rxidle         int64                   `json:"rxidle"`
	ZbufGrow       int64                   `json:"zbuf_grow"`
	BufGrow        int64                   `json:"buf_grow"`
	Wakeups        int64                   `json:"wakeups"`
	Connects       int64                   `json:"connects"`
	Disconnects    int64                   `json:"disconnects"`
	IntLatency     map[string]int64        `json:"int_latency"`
	OutbufLatency  map[string]int64        `json:"outbuf_latency"`
	Rtt            map[string]int64        `json:"rtt"`
	Throttle       map[string]int64        `json:"throttle"`
	Req            map[string]int64        `json:"req"`
	Toppars        GroupCoordinatorToppars `json:"toppars"`
}

type Broker290921_Toppars struct {
	TestTopic0 TestTopic0 `json:"test-topic-0"`
}

type TestTopic0 struct {
	Topic     string `json:"topic"`
	Partition int64  `json:"partition"`
}

type GroupCoordinatorToppars struct {
}

type Cgrp struct {
	State           string `json:"state"`
	Stateage        int64  `json:"stateage"`
	JoinState       string `json:"join_state"`
	RebalanceAge    int64  `json:"rebalance_age"`
	RebalanceCnt    int64  `json:"rebalance_cnt"`
	RebalanceReason string `json:"rebalance_reason"`
	AssignmentSize  int64  `json:"assignment_size"`
}

type Topics struct {
	TestTopic TestTopic `json:"test-topic"`
}

type TestTopic struct {
	Topic       string           `json:"topic"`
	Age         int64            `json:"age"`
	MetadataAge int64            `json:"metadata_age"`
	Batchsize   map[string]int64 `json:"batchsize"`
	Batchcnt    map[string]int64 `json:"batchcnt"`
	Partitions  Partitions       `json:"partitions"`
}

type Partitions struct {
	The0 The1 `json:"0"`
	The1 The1 `json:"-1"`
}

type The1 struct {
	Partition            int64  `json:"partition"`
	Broker               int64  `json:"broker"`
	Leader               int64  `json:"leader"`
	Desired              bool   `json:"desired"`
	Unknown              bool   `json:"unknown"`
	MsgqCnt              int64  `json:"msgq_cnt"`
	MsgqBytes            int64  `json:"msgq_bytes"`
	XmitMsgqCnt          int64  `json:"xmit_msgq_cnt"`
	XmitMsgqBytes        int64  `json:"xmit_msgq_bytes"`
	FetchqCnt            int64  `json:"fetchq_cnt"`
	FetchqSize           int64  `json:"fetchq_size"`
	FetchState           string `json:"fetch_state"`
	QueryOffset          int64  `json:"query_offset"`
	NextOffset           int64  `json:"next_offset"`
	AppOffset            int64  `json:"app_offset"`
	StoredOffset         int64  `json:"stored_offset"`
	StoredLeaderEpoch    int64  `json:"stored_leader_epoch"`
	CommitedOffset       int64  `json:"commited_offset"`
	CommittedOffset      int64  `json:"committed_offset"`
	CommittedLeaderEpoch int64  `json:"committed_leader_epoch"`
	EOFOffset            int64  `json:"eof_offset"`
	LoOffset             int64  `json:"lo_offset"`
	HiOffset             int64  `json:"hi_offset"`
	LsOffset             int64  `json:"ls_offset"`
	ConsumerLag          int64  `json:"consumer_lag"`
	ConsumerLagStored    int64  `json:"consumer_lag_stored"`
	LeaderEpoch          int64  `json:"leader_epoch"`
	Txmsgs               int64  `json:"txmsgs"`
	Txbytes              int64  `json:"txbytes"`
	Rxmsgs               int64  `json:"rxmsgs"`
	Rxbytes              int64  `json:"rxbytes"`
	Msgs                 int64  `json:"msgs"`
	RxVerDrops           int64  `json:"rx_ver_drops"`
	MsgsInflight         int64  `json:"msgs_inflight"`
	NextACKSeq           int64  `json:"next_ack_seq"`
	NextErrSeq           int64  `json:"next_err_seq"`
	AckedMsgid           int64  `json:"acked_msgid"`
}
