package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/util"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"k8s.io/client-go/kubernetes"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

// Job runs inside the Kubernetes cluster and generates status and metrics
// reports.
type Job struct {
	conf          conf.C
	kubeClient    *kubernetes.Clientset
	metricsClient *metrics.Clientset
	mongo         *mongo.Client
}

// New returns a new reporting Job object.
func New(c conf.C, k *kubernetes.Clientset, m *metrics.Clientset, mongo *mongo.Client) Job {
	slog.SetDefault(util.NewLogger(c.Log))
	return Job{
		conf:          c,
		kubeClient:    k,
		metricsClient: m,
		mongo:         mongo,
	}
}

// GenerateAll generates reports for all the configured tasks.
func (j Job) GenerateAll(ctx context.Context) error {
	now := time.Now().UTC().Round(time.Second)

	if j.conf.RabbitMQ.Enable {
		err := j.GenerateRMQReport(ctx, now)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"generating RMQ report",
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("generating RMQ report: %w", err)
		}
	}

	if j.conf.LongJobs.Enable {
		err := j.GenerateLongRunningJobsReport(ctx, now)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"generating long running jobs report",
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("generating long running jobs report: %w", err)
		}
	}

	if j.conf.ImageTag.Enable {
		err := j.GenerateImageTagReport(ctx, now)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"generating image tag report",
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("generating image tag report: %w", err)
		}
	}

	if j.conf.DaSS.Enable {
		err := j.GenerateDaSSReport(ctx, now)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"generating DaSS report",
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("generating DaSS report: %w", err)
		}
	}

	if j.conf.Ceph.Enable {
		err := j.GenerateCEPHReport(ctx, now)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"generating ceph status report",
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("generating ceph status report: %w", err)
		}
	}

	if j.conf.ResourceUtilization.Enable {
		err := j.GenerateResourceUtilizationReport(ctx, now)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"generating resource utilization report",
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("generating resource utilization report: %w", err)
		}
	}

	return nil
}
