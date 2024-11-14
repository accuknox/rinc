package connectivity

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/db"
	"github.com/accuknox/rinc/internal/report"
	types "github.com/accuknox/rinc/types/connectivity"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"k8s.io/client-go/kubernetes"
)

// Reporter is the connectivity status reporter.
type Reporter struct {
	kubeClient *kubernetes.Clientset
	conf       conf.Connectivity
	mongo      *mongo.Client
}

// NewReporter creates a new connectivity status reporter.
func NewReporter(c conf.Connectivity, k *kubernetes.Clientset, mongo *mongo.Client) Reporter {
	return Reporter{
		conf:       c,
		kubeClient: k,
		mongo:      mongo,
	}
}

// Report satisfies the report.Reporter interface by writing the connectivity
// status to the database.
func (r Reporter) Report(ctx context.Context, now time.Time) error {
	metrics := types.Metrics{Timestamp: now}

	if r.conf.Vault.Enable {
		vault, err := r.vaultReport(ctx)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"fetching vault report",
				slog.String("error", err.Error()),
			)
		}
		if vault != nil {
			metrics.Vault = *vault
		}
	}

	if r.conf.Mongodb.Enable {
		mongodb, err := r.mongodbReport(ctx)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"fetching mongodb report",
				slog.String("error", err.Error()),
			)
		}
		if mongodb != nil {
			metrics.Mongodb = *mongodb
		}
	}

	if r.conf.Neo4j.Enable {
		neo4j, err := r.neo4jReport(ctx)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"fetching neo4j report",
				slog.String("error", err.Error()),
			)
		}
		if neo4j != nil {
			metrics.Neo4j = *neo4j
		}
	}

	if r.conf.Postgres.Enable {
		postgres, err := r.postgresReport(ctx)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"fetching postgres report",
				slog.String("error", err.Error()),
			)
		}
		if postgres != nil {
			metrics.Postgres = *postgres
		}
	}

	if r.conf.Redis.Enable {
		redis, err := r.redisReport(ctx)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"fetching redis report",
				slog.String("error", err.Error()),
			)
		}
		if redis != nil {
			metrics.Redis = *redis
		}
	}

	if r.conf.Metabase.Enable {
		metabase, err := r.metabaseReport(ctx)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"fetching metabase report",
				slog.String("error", err.Error()),
			)
		}
		if metabase != nil {
			metrics.Metabase = *metabase
		}
	}

	result, err := db.
		Database(r.mongo).
		Collection(db.CollectionConnectivity).
		InsertOne(ctx, metrics)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"connectivity: inserting metrics into mongodb",
			slog.Time("timestamp", now),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("inserting metrics into mongodb: %w", err)
	}
	slog.LogAttrs(
		ctx,
		slog.LevelDebug,
		"connectivity: inserted metrics into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	alerts := report.SoftEvaluateAlerts(ctx, r.conf.Alerts, metrics)
	result, err = db.
		Database(r.mongo).
		Collection(db.CollectionAlerts).
		InsertOne(ctx, bson.M{
			"timestamp": now,
			"from":      db.CollectionConnectivity,
			"alerts":    alerts,
		})
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"connectivity: inserting alerts into mongodb",
			slog.Time("timestamp", now),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("inserting alerts into mongodb: %w", err)
	}
	slog.LogAttrs(
		ctx,
		slog.LevelDebug,
		"connectivity: inserted alerts into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	return nil
}
