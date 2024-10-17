package rabbitmq

import "time"

// Metrics contains a set of important RabbitMQ metrics that need to be
// included in the report.
type Metrics struct {
	Timestamp   time.Time `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	IsClusterUp bool      `json:"isClusterUp,omitempty" bson:"isClusterUp,omitempty"`
	Overview    Overview  `json:"overview,omitempty" bson:"overview,omitempty"`
	Nodes       Nodes     `json:"nodes,omitempty" bson:"nodes,omitempty"`
	Queues      Queues    `json:"queues,omitempty" bson:"queues,omitempty"`
	Consumers   Consumers `json:"consumers,omitempty" bson:"consumers,omitempty"`
	Exchanges   Exchanges `json:"exchanges,omitempty" bson:"exchanges,omitempty"`
}

type Overview struct {
	Version      string       `json:"rabbitmq_version,omitempty" bson:"rabbitmq_version,omitempty"`
	QueueTotals  queueTotals  `json:"queue_totals,omitempty" bson:"queue_totals,omitempty"`
	ObjectTotals objectTotals `json:"object_totals,omitempty" bson:"object_totals,omitempty"`
}

type Nodes []Node

type Node struct {
	Name           string   `json:"name,omitempty" bson:"name,omitempty"`
	Running        bool     `json:"running,omitempty" bson:"running,omitempty"`
	CPUCount       uint     `json:"processors,omitempty" bson:"processors,omitempty"`
	MemUsed        float64  `json:"mem_used,omitempty" bson:"mem_used,omitempty"`
	FreeDisk       float64  `json:"disk_free,omitempty" bson:"disk_free,omitempty"`
	ProcUsed       uint     `json:"proc_used,omitempty" bson:"proc_used,omitempty"`
	SocketsUsed    uint     `json:"sockets_used,omitempty" bson:"sockets_used,omitempty"`
	FDUsed         uint     `json:"fd_used,omitempty" bson:"fd_used,omitempty"`
	Uptime         uint64   `json:"uptime,omitempty" bson:"uptime,omitempty"`
	EnabledPlugins []string `json:"enabled_plugins,omitempty" bson:"enabled_plugins,omitempty"`
}

type Queues []Queue

type Queue struct {
	Durable                bool   `json:"durable,omitempty" bson:"durable,omitempty"`
	Messages               uint   `json:"messages,omitempty" bson:"messages,omitempty"`
	UnacknowledgedMessages uint   `json:"messages_unacknowledged,omitempty" bson:"messages_unacknowledged,omitempty"`
	ReadyMessages          uint   `json:"messages_ready,omitempty" bson:"messages_ready,omitempty"`
	Name                   string `json:"name,omitempty" bson:"name,omitempty"`
	State                  string `json:"state,omitempty" bson:"state,omitempty"`
}

type Consumers []Consumer

type Consumer struct {
	Active        bool              `json:"active,omitempty" bson:"active,omitempty"`
	Tag           string            `json:"consumer_tag,omitempty" bson:"consumer_tag,omitempty"`
	PrefetchCount uint              `json:"prefetch_count,omitempty" bson:"prefetch_count,omitempty"`
	Queue         consumerQueueInfo `json:"queue,omitempty" bson:"queue,omitempty"`
}

type Exchanges []Exchange

type Exchange struct {
	Durable bool   `json:"durable,omitempty" bson:"durable,omitempty"`
	Name    string `json:"name,omitempty" bson:"name,omitempty"`
	Typ     string `json:"type,omitempty" bson:"type,omitempty"`
}

type queueTotals struct {
	Messages               uint `json:"messages,omitempty" bson:"messages,omitempty"`
	ReadyMessages          uint `json:"messages_ready,omitempty" bson:"messages_ready,omitempty"`
	UnacknowledgedMessages uint `json:"messages_unacknowledged,omitempty" bson:"messages_unacknowledged,omitempty"`
}

type objectTotals struct {
	Channels    uint `json:"channels,omitempty" bson:"channels,omitempty"`
	Connections uint `json:"connections,omitempty" bson:"connections,omitempty"`
	Consumers   uint `json:"consumers,omitempty" bson:"consumers,omitempty"`
	Exchanges   uint `json:"exchanges,omitempty" bson:"exchanges,omitempty"`
	Queues      uint `json:"queues,omitempty" bson:"queues,omitempty"`
}

type consumerQueueInfo struct {
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}
