package resource

import "time"

type Metrics struct {
	Timestamp time.Time
	Nodes     []Node
	Pods      []Pod
}

type Node struct {
	Name string
	Usage
}

type Pod struct {
	Name      string
	Namespace string
	Usage
}

type Usage struct {
	CPU float64
	Mem float64
}
