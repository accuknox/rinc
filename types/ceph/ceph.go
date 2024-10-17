package ceph

import "time"

// Metrics contains a set of important CEPH metrics that need to be included in
// the report.
type Metrics struct {
	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	Summary   Summary   `json:"summary,omitempty" bson:"summary,omitempty"`
	Status    Status    `json:"status,omitempty" bson:"status,omitempty"`
	Devices   []Device  `json:"devices,omitempty" bson:"devices,omitempty"`
	Buckets   []Bucket  `json:"buckets,omitempty" bson:"buckets,omitempty"`
}
