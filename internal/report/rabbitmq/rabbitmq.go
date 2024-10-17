package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/db"
	types "github.com/accuknox/rinc/types/rabbitmq"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"k8s.io/client-go/kubernetes"
)

// Reporter is the rabbitmq health metrics reporter.
type Reporter struct {
	kubeClient *kubernetes.Clientset
	conf       conf.RabbitMQ
	mongo      *mongo.Client
}

// NewReporter creates a new of the rabbitmq reporter.
func NewReporter(c conf.RabbitMQ, k *kubernetes.Clientset, mongo *mongo.Client) Reporter {
	return Reporter{
		conf:       c,
		kubeClient: k,
		mongo:      mongo,
	}
}

// Report satisfies the report.Reporter interface by writing the RabbitMQ
// cluster status and fetched metrics to the mongodb database.
func (r Reporter) Report(ctx context.Context, now time.Time) error {
	up, err := r.IsClusterUp(ctx)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"fetching rabbitmq health status",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("fetching rabbitmq health status: %w", err)
	}
	if !up {
		slog.LogAttrs(
			ctx,
			slog.LevelInfo,
			"rabbitmq cluster is down",
		)
		result, err := db.
			Database(r.mongo).
			Collection(db.CollectionRabbitmq).
			InsertOne(ctx, types.Metrics{
				Timestamp:   now,
				IsClusterUp: false,
			})
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
			"rabbitmq: inserted document into mongodb",
			slog.Any("insertedId", result.InsertedID),
		)
		return nil
	}

	metrics, err := r.GetMetrics(ctx)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to fetch rabbitmq metrics",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to fetch rabbitmq metrics: %w", err)
	}
	metrics.Timestamp = now

	result, err := db.Database(r.mongo).
		Collection(db.CollectionRabbitmq).
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
		"rabbitmq: inserted document into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	return nil
}
