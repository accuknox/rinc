package job

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/accuknox/rinc/internal/report/dass"
	"github.com/accuknox/rinc/internal/util"
)

// GenerateDaSSReport generates a status report for deployments and
// statefulsets.
func (j Job) GenerateDaSSReport(ctx context.Context, now time.Time) error {
	stamp := now.Format(util.IsosecLayout)
	err := j.initStamp(ctx, stamp)
	if err != nil {
		return fmt.Errorf("initializing %q stamp: %w", stamp, err)
	}

	path := filepath.Join(j.conf.Output, stamp, "deployment-statefulset-status.html")
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

	r := dass.NewReporter(j.conf.DaSS, j.kubeClient)
	err = r.Report(ctx, f, now)
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
