package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/report/pod"
)

// GeneratePodStatusReport generates pod status report.
func (j Job) GeneratePodStatusReport(ctx context.Context, now time.Time) error {
	r := pod.NewReporter(j.conf.PodStatus, j.kubeClient, j.mongo)
	err := r.Report(ctx, now)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"generating pod status report",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("generating pod status report: %w", err)
	}
	return nil
}
