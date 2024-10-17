package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/report/ceph"
)

// GenerateCEPHReport generates ceph status report.
func (j Job) GenerateCEPHReport(ctx context.Context, now time.Time) error {
	r := ceph.NewReporter(j.conf.Ceph, j.kubeClient, j.mongo)
	err := r.Report(ctx, now)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"generating ceph report",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("generating ceph report: %w", err)
	}
	return nil
}
