package ceph

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
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
	err := r.call(ctx, healthEndpoint, mediaTypeV10, status)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"fetching ceph health status",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("fetching ceph health status: %w", err)
	}

	var hosts []types.Host
	for offset := 0; ; offset += 30 {
		q := url.Values{}
		q.Set("limit", "30")
		q.Set("offset", fmt.Sprintf("%d", offset))

		var h []types.Host
		err = r.call(ctx, hostListEndpoint, mediaTypeV13, &h, q)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"fetching ceph hosts",
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("fetching ceph hosts: %w", err)
		}
		if len(h) == 0 {
			break
		}
		hosts = append(hosts, h...)
	}

	var devices []types.Device
	for _, h := range hosts {
		var d []types.Device
		endp := fmt.Sprintf(hostDevicesEndpoint, h.Hostname)
		err = r.call(ctx, endp, mediaTypeV10, &d)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"fetching ceph host devices",
				slog.String("error", err.Error()),
				slog.String("host", h.Hostname),
			)
			return fmt.Errorf("fetching ceph host devices: %w", err)
		}
		devices = append(devices, d...)
	}

	var inventories []types.Inventory
	err = r.call(ctx, hostInventoryEndpoint, mediaTypeV10, &inventories)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"fetching ceph host inventories",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("fetching ceph host inventories: %w", err)
	}

	stamp := now.Format(util.IsosecLayout)
	c := layout.Base(
		fmt.Sprintf("CEPH - %s | AccuKnox Reports", stamp),
		partial.Navbar(false, false),
		tmpl.Report(tmpl.Data{
			Timestamp:   now,
			Status:      *status,
			Hosts:       hosts,
			Devices:     devices,
			Inventories: inventories,
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
