package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/report/dass"
)

// GenerateDaSSReport generates a status report for deployments and
// statefulsets.
func (j Job) GenerateDaSSReport(ctx context.Context, now time.Time) error {
	r := dass.NewReporter(j.conf.DaSS, j.kubeClient, j.mongo)
	err := r.Report(ctx, now)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"generating DaSS report",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("generating DaSS report: %w", err)
	}
	return nil
}
