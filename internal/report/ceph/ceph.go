package ceph

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/util"
	types "github.com/accuknox/rinc/types/ceph"
	tmpl "github.com/accuknox/rinc/view/ceph"
	"github.com/accuknox/rinc/view/layout"
	"github.com/accuknox/rinc/view/partial"

	"k8s.io/client-go/kubernetes"
)

// Reporter is the ceph status reporter.
type Reporter struct {
	kubeClient *kubernetes.Clientset
	conf       conf.Ceph
	token      *token
}

// NewReporter creates a new ceph status reporter.
func NewReporter(c conf.Ceph, kubeClient *kubernetes.Clientset) Reporter {
	return Reporter{
		conf:       c,
		kubeClient: kubeClient,
		token:      nil,
	}
}

// Report satisfies the report.Reporter interface by writing the CEPH status
// and fetched metrics to the provided io.Writer.
func (r Reporter) Report(ctx context.Context, to io.Writer, now time.Time) error {
	status := new(types.Status)
	err := r.call(ctx, healthEndpoint, mediaTypeV1, status)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"fetching ceph health status",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("fetching ceph health status: %w", err)
	}
	stamp := now.Format(util.IsosecLayout)
	c := layout.Base(
		fmt.Sprintf("CEPH - %s | AccuKnox Reports", stamp),
		partial.Navbar(false, false),
		tmpl.Report(tmpl.Data{
			Timestamp: now,
			Status:    *status,
		}),
	)
	err = c.Render(ctx, to)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"rendering ceph template",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("rendering ceph template: %w", err)
	}
	return nil
}
