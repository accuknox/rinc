package imagetag

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/util"
	tmpl "github.com/accuknox/rinc/view/imagetag"
	"github.com/accuknox/rinc/view/layout"
	"github.com/accuknox/rinc/view/partial"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Reporter is the image tag reporter.
type Reporter struct {
	kubeClient *kubernetes.Clientset
	conf       conf.ImageTag
}

// NewReporter creates a new image tag reporter.
func NewReporter(c conf.ImageTag, kubeClient *kubernetes.Clientset) Reporter {
	return Reporter{
		conf:       c,
		kubeClient: kubeClient,
	}
}

// Report satisfies the report.Reporter interface by fetching the image tags of
// deployments and statefulsets from the Kubernetes API server, and writes the
// report to the provided io.Writer.
func (r Reporter) Report(ctx context.Context, to io.Writer, now time.Time) error {
	depls, err := r.deployments(ctx)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"fetching deployment resources",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("fetching deployments: %w", err)
	}

	statefulsets, err := r.statefulsets(ctx)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"fetching statefulset resources",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("fetching statefulsets: %w", err)
	}

	stamp := now.Format(util.IsosecLayout)
	c := layout.Base(
		fmt.Sprintf("Image Tags - %s | AccuKnox Reports", stamp),
		partial.Navbar(false, false),
		tmpl.Report(tmpl.Data{
			Timestamp:    now,
			Deployments:  depls,
			Statefulsets: statefulsets,
		}),
	)
	err = c.Render(ctx, to)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"rendering image tag report template",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("rendering image tag report template: %w", err)
	}

	return nil
}

func (r Reporter) deployments(ctx context.Context) ([]tmpl.Resource, error) {
	var resources []tmpl.Resource
	var cntinue string

	for {
		depls, err := r.kubeClient.
			AppsV1().
			Deployments(r.conf.Namespace).
			List(ctx, metav1.ListOptions{
				Continue: cntinue,
				Limit:    30,
			})
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"listing deployments",
				slog.String("namespace", r.conf.Namespace),
				slog.String("error", err.Error()),
			)
			return nil, fmt.Errorf("listing deployments in ns %q: %w",
				r.conf.Namespace, err)
		}

		for _, d := range depls.Items {
			containers := d.Spec.Template.Spec.Containers
			images := make([]string, len(containers))
			for idx, c := range containers {
				images[idx] = c.Image
			}
			resources = append(resources, tmpl.Resource{
				Name:      d.GetName(),
				Namespace: d.GetNamespace(),
				Images:    images,
			})
		}

		cntinue = depls.Continue
		if cntinue == "" {
			slog.LogAttrs(
				ctx,
				slog.LevelInfo,
				"received image tags from all deployments",
			)
			break
		}
	}

	return resources, nil
}

func (r Reporter) statefulsets(ctx context.Context) ([]tmpl.Resource, error) {
	var resources []tmpl.Resource
	var cntinue string

	for {
		ss, err := r.kubeClient.
			AppsV1().
			StatefulSets(r.conf.Namespace).
			List(ctx, metav1.ListOptions{
				Continue: cntinue,
				Limit:    30,
			})
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"listing statefulsets",
				slog.String("namespace", r.conf.Namespace),
				slog.String("error", err.Error()),
			)
			return nil, fmt.Errorf("listing statefulsets in ns %q: %w",
				r.conf.Namespace, err)
		}

		for _, s := range ss.Items {
			containers := s.Spec.Template.Spec.Containers
			images := make([]string, len(containers))
			for idx, c := range containers {
				images[idx] = c.Image
			}
			resources = append(resources, tmpl.Resource{
				Name:      s.GetName(),
				Namespace: s.GetNamespace(),
				Images:    images,
			})
		}

		cntinue = ss.Continue
		if cntinue == "" {
			slog.LogAttrs(
				ctx,
				slog.LevelInfo,
				"received image tags from all statefulsets",
			)
			break
		}
	}

	return resources, nil
}
