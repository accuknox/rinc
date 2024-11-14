package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/report/connectivity"
)

// GenerateConnectivityReport generates connectivity status report.
func (j Job) GenerateConnectivityReport(ctx context.Context, now time.Time) error {
	r := connectivity.NewReporter(j.conf.Connectivity, j.kubeClient, j.mongo)
	err := r.Report(ctx, now)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"generating connectivity report",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("generating connectivity report: %w", err)
	}
	return nil
}
