package dass

import "time"

type Metrics struct {
	Timestamp    time.Time
	Deployments  []Resource
	Statefulsets []Resource
}

type Resource struct {
	Name              string
	Namespace         string
	Age               time.Duration
	DesiredReplicas   int32
	ReadyReplicas     int32
	AvailableReplicas int32
	UpdatedReplicas   int32
	Events            []Event
	IsReplicaFailure  bool
	IsAvailable       bool
}

type Event struct {
	Type    string
	Reason  string
	Message string
}
