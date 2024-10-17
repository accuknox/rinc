package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/report/imagetag"
)

// GenerateImageTagReport generates an image tag report for deployments and
// statefulsets.
func (j Job) GenerateImageTagReport(ctx context.Context, now time.Time) error {
	r := imagetag.NewReporter(j.conf.ImageTag, j.kubeClient, j.mongo)
	err := r.Report(ctx, now)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"generating image tag report",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("generating image tag report: %w", err)
	}
	return nil
}
