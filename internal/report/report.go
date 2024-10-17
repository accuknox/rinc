package report

import (
	"context"
	"time"
)

// Reporter defines an interface for reporting data. Implementations of this
// interface should collect and write metrics to the database returning any
// errors encountered during the process.
type Reporter interface {
	Report(ctx context.Context, now time.Time) error
}
