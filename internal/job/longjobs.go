package job

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/accuknox/rinc/internal/report/longjobs"
	"github.com/accuknox/rinc/internal/util"
)

// GenerateLongRunningJobsReport generates a RabbitMQ status and metrics
// report.
func (j Job) GenerateLongRunningJobsReport(ctx context.Context, now time.Time) error {
	stamp := now.Format(util.IsosecLayout)
	err := j.initStamp(ctx, stamp)
	if err != nil {
		return fmt.Errorf("initializing %q stamp: %w", stamp, err)
	}

	path := filepath.Join(j.conf.Output, stamp, "longrunningjobs.html")
	f, err := os.Create(path)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			fmt.Sprintf("creating %q file", path),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("creating %q file: %w", path, err)
	}
	defer f.Close()

	r := longjobs.NewReporter(j.conf.LongJobs, j.kubeClient)
	err = r.Report(ctx, f, now)
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
