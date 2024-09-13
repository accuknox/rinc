package job

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/accuknox/rinc/internal/report/imagetag"
	"github.com/accuknox/rinc/internal/util"
)

// GenerateImageTagReport generates an image tag report for deployments and
// statefulsets.
func (j Job) GenerateImageTagReport(ctx context.Context, now time.Time) error {
	stamp := now.Format(util.IsosecLayout)
	err := j.initStamp(ctx, stamp)
	if err != nil {
		return fmt.Errorf("initializing %q stamp: %w", stamp, err)
	}

	path := filepath.Join(j.conf.Output, stamp, "imagetag.html")
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

	r := imagetag.NewReporter(j.conf.ImageTag, j.kubeClient)
	err = r.Report(ctx, f, now)
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
