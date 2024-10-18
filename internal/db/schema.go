package db

import "github.com/accuknox/rinc/internal/conf"

// Alert is the alert object schema that should be stored in the `alerts`
// collection.
type Alert struct {
	Message  string        `bson:"message"`
	Severity conf.Severity `bson:"severity"`
}

const (
	CollectionAlerts   = "alerts"
	CollectionRabbitmq = "rabbitmq"
	CollectionCeph     = "ceph"
	CollectionImageTag = "imagetag"
	CollectionDass     = "dass"
	CollectionLongJobs = "longjobs"
)

// Collections is a list of MongoDB collection names, excluding the alerts
// collection.
var Collections = []string{
	CollectionRabbitmq,
	CollectionCeph,
	CollectionImageTag,
	CollectionDass,
	CollectionLongJobs,
}
