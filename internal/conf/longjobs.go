package conf

import "time"

// LongJobs contains configuration related to the long-running job
// reporter.
type LongJobs struct {
	// Enable specifies whether the long-running job reporter should be
	// enabled.
	Enable bool `koanf:"enable"`
	// Namespace is the namespace in which the long-running jobs will be
	// reported.
	Namespace string `koanf:"namespace"`
	// OlderThan defines the duration threshold; jobs older than this
	// value will be reported.
	OlderThan time.Duration `koanf:"olderThan"`
}
