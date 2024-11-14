package resource

import "time"

type Metrics struct {
	Timestamp  time.Time
	Nodes      []Node
	Containers []Container
}

type Node struct {
	Name           string
	CPUUsedPercent float64
	MemUsedPercent float64
}

type Container struct {
	PodName        string
	Namespace      string
	Name           string
	CPULimit       float64
	MemLimit       float64
	CPUUsed        float64
	MemUsed        float64
	CPUUsedPercent float64
	MemUsedPercent float64
}
