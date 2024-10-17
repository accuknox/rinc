package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/report/longjobs"
)

// GenerateLongRunningJobsReport generates a report for Kubernetes jobs running
// older than the given provided threshold.
func (j Job) GenerateLongRunningJobsReport(ctx context.Context, now time.Time) error {
	r := longjobs.NewReporter(j.conf.LongJobs, j.kubeClient, j.mongo)
	err := r.Report(ctx, now)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"generating long running jobs report",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("generating long running jobs report: %w", err)
	}
	return nil
}
