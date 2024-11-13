package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/report/resource"
)

// GenerateResourceUtilizationReport generates resource utilizaton report.
func (j Job) GenerateResourceUtilizationReport(ctx context.Context, now time.Time) error {
	r := resource.NewReporter(resource.Config{
		ResourceUtilizationConfig: j.conf.ResourceUtilization,
		KubeClient:                j.kubeClient,
		MetricsClient:             j.metricsClient,
		MongoClient:               j.mongo,
	})
	err := r.Report(ctx, now)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"generating resource utilization report",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("generating resource utilization report: %w", err)
	}
	return nil
}
