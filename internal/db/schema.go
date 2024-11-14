package db

import (
	"time"

	"github.com/accuknox/rinc/internal/conf"
)

// AlertDocument defines the schema that should be stored in the
// `alerts` collection.
type AlertDocument struct {
	Timestamp time.Time `bson:"timestamp"`
	From      string    `bson:"from"`
	Alerts    []Alert   `bson:"alerts"`
}

// Alert defines the schema that should be stored within the
// AlertDocument in the `alerts` collection.
type Alert struct {
	Message  string        `bson:"message"`
	Severity conf.Severity `bson:"severity"`
}

const (
	CollectionAlerts       = "alerts"
	CollectionRabbitmq     = "rabbitmq"
	CollectionCeph         = "ceph"
	CollectionImageTag     = "imagetag"
	CollectionDass         = "dass"
	CollectionLongJobs     = "longjobs"
	CollectionPVUtilizaton = "pv_utilization"
)

// Collections is a list of MongoDB collection names, excluding the alerts
// collection.
var Collections = []string{
	CollectionRabbitmq,
	CollectionCeph,
	CollectionImageTag,
	CollectionDass,
	CollectionLongJobs,
	CollectionPVUtilizaton,
}
