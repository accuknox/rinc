package pv

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/db"
	"github.com/accuknox/rinc/internal/report"
	types "github.com/accuknox/rinc/types/pv"

	"github.com/prometheus/client_golang/api"
	promV1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"k8s.io/client-go/kubernetes"
)

// Reporter is the PV utilization reporter.
type Reporter struct {
	kubeClient *kubernetes.Clientset
	conf       conf.PVUtilization
	mongo      *mongo.Client
}

// NewReporter creates a new PV utilization reporter.
func NewReporter(c conf.PVUtilization, k *kubernetes.Clientset, mongo *mongo.Client) Reporter {
	return Reporter{
		conf:       c,
		kubeClient: k,
		mongo:      mongo,
	}
}

// Report satisfies the report.Reporter interface by fetching the PV
// utilizations by querying prometheus, and writes the report to the database.
func (r Reporter) Report(ctx context.Context, now time.Time) error {
	client, err := api.NewClient(api.Config{
		Address: r.conf.PrometheusURL,
	})
	if err != nil {
		return fmt.Errorf("creating prometheus client: %w", err)
	}

	api := promV1.NewAPI(client)
	pvs := make(types.PVs, 0)

	for metric, q := range queries {
		vector, err := query(ctx, api, q)
		if err != nil {
			return err
		}
		for _, sample := range vector {
			var ns, pvc string
			if label := sample.Metric["namespace"]; label.IsValid() {
				ns = string(label)
			}
			if label := sample.Metric["persistentvolumeclaim"]; label.IsValid() {
				pvc = string(label)
			}
			slog.LogAttrs(
				ctx,
				slog.LevelDebug,
				"sample",
				slog.Int("metric", metric),
				slog.String("namespace", ns),
				slog.String("pvc", pvc),
				slog.Float64("value", float64(sample.Value)),
			)
			switch metric {
			case metricCapacity:
				pvs = pvs.AppendCapacity(pvc, ns, float64(sample.Value))
			case metricUsed:
				pvs = pvs.AppendUsed(pvc, ns, float64(sample.Value))
			case metricAvailable:
				pvs = pvs.AppendAvailable(pvc, ns, float64(sample.Value))
			case metricUtilization:
				pvs = pvs.AppendUtilization(pvc, ns, float64(sample.Value))
			}
		}
	}

	metrics := types.Metrics{
		Timestamp: now,
		PVs:       pvs,
	}

	result, err := db.
		Database(r.mongo).
		Collection(db.CollectionPVUtilizaton).
		InsertOne(ctx, metrics)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"inserting into mongodb",
			slog.Time("timestamp", now),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("inserting into mongodb: %w", err)
	}
	slog.LogAttrs(
		ctx,
		slog.LevelDebug,
		"pv: inserted document into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	alerts := report.SoftEvaluateAlerts(ctx, r.conf.Alerts, metrics)
	result, err = db.
		Database(r.mongo).
		Collection(db.CollectionAlerts).
		InsertOne(ctx, bson.M{
			"timestamp": now,
			"from":      db.CollectionPVUtilizaton,
			"alerts":    alerts,
		})
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"pv: inserting alerts into mongodb",
			slog.Time("timestamp", now),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("inserting alerts into mongodb: %w", err)
	}
	slog.LogAttrs(
		ctx,
		slog.LevelDebug,
		"pv: inserted alerts into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	return nil
}
