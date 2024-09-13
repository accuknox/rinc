package longjobs

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/util"
	"github.com/accuknox/rinc/view/layout"
	tmpl "github.com/accuknox/rinc/view/longjobs"
	"github.com/accuknox/rinc/view/partial"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Reporter is the long-running jobs reporter.
type Reporter struct {
	kubeClient *kubernetes.Clientset
	conf       conf.LongJobs
}

// NewReporter creates a new long-running jobs reporter.
func NewReporter(c conf.LongJobs, kubeClient *kubernetes.Clientset) Reporter {
	return Reporter{
		conf:       c,
		kubeClient: kubeClient,
	}
}

// Report satisfies the report.Reporter interface by fetching the long-running
// jobs from the Kubernetes API server and writing it to the provided
// io.Writer.
func (r Reporter) Report(ctx context.Context, to io.Writer, now time.Time) error {
	threshold := now.Add(-r.conf.OlderThan)
	var longJobs []tmpl.Job
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

	Outer:
		for _, job := range jobs.Items {
			if isFinished(job.Status.Conditions) {
				continue Outer
			}
			old := job.CreationTimestamp.Time.Before(threshold)
			if !old {
				continue Outer
			}
			var readyPods int32
			if job.Status.Ready != nil {
				readyPods = *job.Status.Ready
			}
			longJobs = append(longJobs, tmpl.Job{
				Name:       job.GetName(),
				Namespace:  job.GetNamespace(),
				Suspended:  isSuspended(job.Status.Conditions),
				ActivePods: job.Status.Active,
				FailedPods: job.Status.Failed,
				ReadyPods:  readyPods,
			})
			slog.LogAttrs(
				ctx,
				slog.LevelDebug,
				"long running job found",
				slog.String("name", job.GetName()),
				slog.String("namespace", job.GetNamespace()),
				slog.Bool("suspended", isSuspended(job.Status.Conditions)),
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

	stamp := now.Format(util.IsosecLayout)
	c := layout.Base(
		fmt.Sprintf("Long Running Jobs - %s | AccuKnox Reports", stamp),
		partial.Navbar(false, false),
		tmpl.Report(tmpl.Data{
			Timestamp: now,
			Jobs:      longJobs,
		}),
	)
	err := c.Render(ctx, to)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"rendering long running jobs template",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("rendering long running jobs template: %w", err)
	}

	return nil
}
