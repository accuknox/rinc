package longjobs

import "time"

type Metrics struct {
	Timestamp time.Time
	OlderThan time.Duration
	Jobs      []Job
}

type Job struct {
	Name       string
	Namespace  string
	Suspended  bool
	ActivePods int32
	FailedPods int32
	ReadyPods  int32
	Age        time.Duration
}
