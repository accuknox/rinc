package imagetag

import "time"

type Metrics struct {
	Timestamp    time.Time
	Deployments  []Resource
	Statefulsets []Resource
}

type Resource struct {
	Name      string
	Namespace string
	Images    []string
}
