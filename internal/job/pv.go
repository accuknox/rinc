package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/report/pv"
)

// GeneratePVUtilizationReport generates a PV utilization status report.
func (j Job) GeneratePVUtilizationReport(ctx context.Context, now time.Time) error {
	r := pv.NewReporter(j.conf.PVUtilization, j.kubeClient, j.mongo)
	err := r.Report(ctx, now)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"generating PV utilization report",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("generating PV utilization report: %w", err)
	}
	return nil
}
