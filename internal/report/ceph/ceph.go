package ceph

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/db"
	"github.com/accuknox/rinc/internal/report"
	types "github.com/accuknox/rinc/types/ceph"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"k8s.io/client-go/kubernetes"
)

// Reporter is the ceph status reporter.
type Reporter struct {
	kubeClient *kubernetes.Clientset
	conf       conf.Ceph
	mongo      *mongo.Client
	token      *token
}

// NewReporter creates a new ceph status reporter.
func NewReporter(c conf.Ceph, k *kubernetes.Clientset, mongo *mongo.Client) Reporter {
	return Reporter{
		conf:       c,
		kubeClient: k,
		mongo:      mongo,
		token:      nil,
	}
}

// Report satisfies the report.Reporter interface by writing the CEPH status
// and fetched metrics to the provided io.Writer.
func (r Reporter) Report(ctx context.Context, now time.Time) error {
	summary := new(types.Summary)
	err := r.call(ctx, summaryEndpoint, mediaTypeV10, summary)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"fetching ceph summary",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("fetching ceph summary: %w", err)
	}

	status := new(types.Status)
	err = r.call(ctx, healthEndpoint, mediaTypeV10, status)
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

	var buckets []types.Bucket
	q := url.Values{}
	q.Set("stats", "true")
	err = r.call(ctx, bucketEndpoint, mediaTypeV11, &buckets, q)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"fetching ceph RGW buckets",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("fetching ceph RGW buckets: %w", err)
	}

	metrics := types.Metrics{
		Timestamp:   now,
		Summary:     *summary,
		Status:      *status,
		Buckets:     buckets,
		Devices:     devices,
		Hosts:       hosts,
		Inventories: inventories,
	}

	result, err := db.
		Database(r.mongo).
		Collection(db.CollectionCeph).
		InsertOne(ctx, metrics)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"ceph: inserting metrics into mongodb",
			slog.Time("timestamp", now),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("inserting metrics into mongodb: %w", err)
	}
	slog.LogAttrs(
		ctx,
		slog.LevelDebug,
		"ceph: inserted metrics into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	alerts := report.SoftEvaluateAlerts(ctx, r.conf.Alerts, metrics)
	result, err = db.
		Database(r.mongo).
		Collection(db.CollectionAlerts).
		InsertOne(ctx, bson.M{
			"timestamp": now,
			"from":      db.CollectionCeph,
			"alerts":    alerts,
		})
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"ceph: inserting alerts into mongodb",
			slog.Time("timestamp", now),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("inserting alerts into mongodb: %w", err)
	}
	slog.LogAttrs(
		ctx,
		slog.LevelDebug,
		"ceph: inserted alerts into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	return nil
}
