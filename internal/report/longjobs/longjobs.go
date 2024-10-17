package longjobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/db"
	types "github.com/accuknox/rinc/types/longjobs"

	"go.mongodb.org/mongo-driver/v2/mongo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Reporter is the long-running jobs reporter.
type Reporter struct {
	kubeClient *kubernetes.Clientset
	conf       conf.LongJobs
	mongo      *mongo.Client
}

// NewReporter creates a new long-running jobs reporter.
func NewReporter(c conf.LongJobs, k *kubernetes.Clientset, mongo *mongo.Client) Reporter {
	return Reporter{
		conf:       c,
		kubeClient: k,
		mongo:      mongo,
	}
}

// Report satisfies the report.Reporter interface by fetching the long-running
// jobs from the Kubernetes API server and writing it to the database.
func (r Reporter) Report(ctx context.Context, now time.Time) error {
	threshold := now.Add(-r.conf.OlderThan)
	var longJobs []types.Job
	var cntinue string

	for {
		jobs, err := r.kubeClient.
			BatchV1().
			Jobs(r.conf.Namespace).
			List(ctx, metav1.ListOptions{
				Continue: cntinue,
				Limit:    30,
			})
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"listing jobs",
				slog.String("namespace", r.conf.Namespace),
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("listing jobs in ns %q: %w", r.conf.Namespace, err)
		}

		for _, job := range jobs.Items {
			if isFinished(job.Status.Conditions) {
				continue
			}
			isSuspended := isSuspended(job.Status.Conditions)
			if isSuspended && !r.conf.IncludeSuspended {
				continue
			}
			old := job.CreationTimestamp.Time.Before(threshold)
			if !old {
				continue
			}
			var readyPods int32
			if job.Status.Ready != nil {
				readyPods = *job.Status.Ready
			}
			longJobs = append(longJobs, types.Job{
				Name:       job.GetName(),
				Namespace:  job.GetNamespace(),
				Suspended:  isSuspended,
				ActivePods: job.Status.Active,
				FailedPods: job.Status.Failed,
				ReadyPods:  readyPods,
				Age:        time.Now().UTC().Sub(job.CreationTimestamp.UTC()),
			})
			slog.LogAttrs(
				ctx,
				slog.LevelDebug,
				"long running job found",
				slog.String("name", job.GetName()),
				slog.String("namespace", job.GetNamespace()),
				slog.Bool("suspended", isSuspended),
				slog.Int("activePods", int(job.Status.Active)),
				slog.Int("failedPods", int(job.Status.Failed)),
				slog.Int("readyPods", int(readyPods)),
			)
		}

		cntinue = jobs.Continue
		if cntinue == "" {
			slog.LogAttrs(
				ctx,
				slog.LevelInfo,
				"all jobs diagnosed successfully",
			)
			break
		}
	}

	result, err := db.Database(r.mongo).
		Collection(db.CollectionLongJobs).
		InsertOne(ctx, types.Metrics{
			Timestamp: now,
			OlderThan: r.conf.OlderThan,
			Jobs:      longJobs,
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
		"longjobs: inserted document into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	return nil
}
