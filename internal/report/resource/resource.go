package resource

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/db"
	"github.com/accuknox/rinc/internal/report"
	types "github.com/accuknox/rinc/types/resource"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

// Reporter is the resource utilization reporter.
type Reporter struct {
	Config
}

type Config struct {
	ResourceUtilizationConfig conf.ResourceUtilization
	KubeClient                *kubernetes.Clientset
	MetricsClient             *metrics.Clientset
	MongoClient               *mongo.Client
}

// NewReporter creates a new resource utilization reporter.
func NewReporter(c Config) Reporter {
	return Reporter{Config: c}
}

// Report satisfies the report.Reporter interface by fetching the resource
// utilizations of nodes & pods from the Kubernetes metrics API server, and
// writes the report to the database.
func (r Reporter) Report(ctx context.Context, now time.Time) error {
	nodes, err := r.nodeUsage(ctx)
	if err != nil {
		return fmt.Errorf("fetching node usage: %w", err)
	}

	metrics := types.Metrics{
		Timestamp: now,
		Nodes:     nodes,
	}

	result, err := db.
		Database(r.MongoClient).
		Collection(db.CollectionResourceUtilization).
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
		"resource: inserted document into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	alerts := report.SoftEvaluateAlerts(
		ctx,
		r.ResourceUtilizationConfig.Alerts,
		metrics,
	)
	result, err = db.
		Database(r.MongoClient).
		Collection(db.CollectionAlerts).
		InsertOne(ctx, bson.M{
			"timestamp": now,
			"from":      db.CollectionResourceUtilization,
			"alerts":    alerts,
		})
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"resource: inserting alerts into mongodb",
			slog.Time("timestamp", now),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("inserting alerts into mongodb: %w", err)
	}
	slog.LogAttrs(
		ctx,
		slog.LevelDebug,
		"resource: inserted alerts into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	return nil
}

func (r Reporter) nodeUsage(ctx context.Context) ([]types.Node, error) {
	var (
		nodes   []types.Node
		usage   []metric
		cntinue string
	)
	for {
		metrics, err := r.MetricsClient.
			MetricsV1beta1().
			NodeMetricses().
			List(ctx, metav1.ListOptions{
				Limit:    30,
				Continue: cntinue,
			})
		if err != nil {
			return nil, fmt.Errorf("fetching node metrics: %w", err)
		}
		for _, m := range metrics.Items {
			usage = append(usage, metric{
				name: m.Name,
				cpu:  m.Usage.Cpu(),
				mem:  m.Usage.Memory(),
			})
		}
		cntinue = metrics.Continue
		if cntinue == "" {
			slog.LogAttrs(
				ctx,
				slog.LevelInfo,
				"received usage of all nodes",
			)
			break
		}
	}

	cntinue = ""
	for {
		nodeList, err := r.KubeClient.
			CoreV1().
			Nodes().
			List(ctx, metav1.ListOptions{
				Limit:    30,
				Continue: cntinue,
			})
		if err != nil {
			return nil, fmt.Errorf("fetching nodes: %w", err)
		}
		for _, n := range nodeList.Items {
			var u *metric = nil
			for _, m := range usage {
				if m.name == n.Name {
					u = &m
				}
			}
			if u == nil {
				continue
			}
			cpu := percentage(*u.cpu, *n.Status.Capacity.Cpu())
			mem := percentage(*u.mem, *n.Status.Capacity.Memory())
			slog.LogAttrs(
				ctx,
				slog.LevelDebug,
				"NODE UTILIZATION %",
				slog.String("name", n.Name),
				slog.Float64("cpu", cpu),
				slog.Float64("mem", mem),
			)
			nodes = append(nodes, types.Node{
				Name: n.Name,
				Usage: types.Usage{
					CPU: cpu,
					Mem: mem,
				},
			})
		}
		cntinue = nodeList.Continue
		if cntinue == "" {
			slog.LogAttrs(
				ctx,
				slog.LevelInfo,
				"received capacity of all nodes",
			)
			break
		}
	}

	return nodes, nil
}

type metric struct {
	name      string
	namespace string
	cpu       *resource.Quantity
	mem       *resource.Quantity
}

func percentage(used, total resource.Quantity) float64 {
	usedFloat := used.AsApproximateFloat64()
	totalFloat := total.AsApproximateFloat64()
	if totalFloat == 0 {
		return 0
	}
	return usedFloat * 100 / totalFloat
}
