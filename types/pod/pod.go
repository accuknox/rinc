package pod

import "time"

type Metrics struct {
	Timestamp    time.Time
	Deployments  []Resource
	Statefulsets []Resource
}

type Resource struct {
	Name      string
	Namespace string
	Pods      []Pod
}

type Pod struct {
	Name       string
	Status     string
	QOSClass   string
	StartTime  time.Time
	Containers []Container
}

type Container struct {
	Name                 string
	IsInit               bool
	Ready                bool
	State                string
	RestartCount         int32
	LastTerminationState string
}
