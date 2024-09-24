package job

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/accuknox/rinc/internal/report/ceph"
	"github.com/accuknox/rinc/internal/util"
)

// GenerateCEPHReport generates ceph status report.
func (j Job) GenerateCEPHReport(ctx context.Context, now time.Time) error {
	stamp := now.Format(util.IsosecLayout)
	err := j.initStamp(ctx, stamp)
	if err != nil {
		return fmt.Errorf("initializing %q stamp: %w", stamp, err)
	}

	path := filepath.Join(j.conf.Output, stamp, "ceph.html")
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

	rmqR := ceph.NewReporter(j.conf.Ceph, j.kubeClient)
	err = rmqR.Report(ctx, f, now)
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
