package dass

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/util"
	tmpl "github.com/accuknox/rinc/view/dass"
	"github.com/accuknox/rinc/view/layout"
	"github.com/accuknox/rinc/view/partial"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Reporter is the deployment and statefulset status (DaSS) reporter.
type Reporter struct {
	kubeClient *kubernetes.Clientset
	conf       conf.DaSS
}

// NewReporter creates a new deployment and statefulset status (DaSS) reporter.
func NewReporter(c conf.DaSS, kubeClient *kubernetes.Clientset) Reporter {
	return Reporter{
		conf:       c,
		kubeClient: kubeClient,
	}
}

// Report satisfies the report.Reporter interface by fetching the status of
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

	ss, err := r.statefulset(ctx)
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
		fmt.Sprintf("Deployment & Statefulset Status - %s | AccuKnox Reports", stamp),
		partial.Navbar(false, false),
		tmpl.Report(tmpl.Data{
			Timestamp:    now,
			Deployments:  depls,
			Statefulsets: ss,
		}),
	)
	err = c.Render(ctx, to)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"rendering DaSS report template",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("rendering DaSS report template: %w", err)
	}

	return nil
}

func (r Reporter) deployments(ctx context.Context) ([]tmpl.Resource, error) {
	var deployments []tmpl.Resource
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
			events, err := r.events(ctx, d.Name, d.Kind)
			if err != nil {
				slog.LogAttrs(
					ctx,
					slog.LevelError,
					"fetching events",
					slog.String("kind", d.Kind),
					slog.String("for", d.Name),
					slog.String("namespace", d.Namespace),
					slog.String("error", err.Error()),
				)
				return nil, fmt.Errorf("fetching events for %q: %w", d.Name, err)
			}
			var desiredReplicas int32
			if d.Spec.Replicas != nil {
				desiredReplicas = *d.Spec.Replicas
			}
			deployments = append(deployments, tmpl.Resource{
				Name:              d.Name,
				Namespace:         d.Namespace,
				Age:               time.Since(d.CreationTimestamp.Time),
				DesiredReplicas:   desiredReplicas,
				ReadyReplicas:     d.Status.ReadyReplicas,
				AvailableReplicas: d.Status.AvailableReplicas,
				UpdatedReplicas:   d.Status.UpdatedReplicas,
				Events:            events,
				IsReplicaFailure:  deploymentHasReplicaFailure(d.Status.Conditions),
				IsAvailable:       isDeploymentAvailable(d.Status.Conditions),
			})
			slog.LogAttrs(
				ctx,
				slog.LevelDebug,
				"collected deployment",
				slog.String("name", d.Name),
				slog.String("namespace", d.Namespace),
				slog.String("kind", d.Kind),
				slog.Int("desiredReplicas", int(desiredReplicas)),
				slog.Int("readyReplicas", int(d.Status.ReadyReplicas)),
				slog.Int("availableReplicas", int(d.Status.AvailableReplicas)),
				slog.Int("updatedReplicas", int(d.Status.UpdatedReplicas)),
			)
		}

		cntinue = depls.Continue
		if cntinue == "" {
			break
		}
	}

	return deployments, nil
}

func (r Reporter) statefulset(ctx context.Context) ([]tmpl.Resource, error) {
	var statefulsets []tmpl.Resource
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
			events, err := r.events(ctx, s.Name, s.Kind)
			if err != nil {
				slog.LogAttrs(
					ctx,
					slog.LevelError,
					"fetching events",
					slog.String("kind", s.Kind),
					slog.String("for", s.Name),
					slog.String("namespace", s.Namespace),
					slog.String("error", err.Error()),
				)
				return nil, fmt.Errorf("fetching events for %q: %w", s.Name, err)
			}
			var desiredReplicas int32
			if s.Spec.Replicas != nil {
				desiredReplicas = *s.Spec.Replicas
			}
			statefulsets = append(statefulsets, tmpl.Resource{
				Name:              s.Name,
				Namespace:         s.Namespace,
				Age:               time.Since(s.CreationTimestamp.Time),
				DesiredReplicas:   desiredReplicas,
				ReadyReplicas:     s.Status.ReadyReplicas,
				AvailableReplicas: s.Status.AvailableReplicas,
				UpdatedReplicas:   s.Status.UpdatedReplicas,
				Events:            events,
			})
			slog.LogAttrs(
				ctx,
				slog.LevelDebug,
				"collected statefulset",
				slog.String("name", s.Name),
				slog.String("namespace", s.Namespace),
				slog.String("kind", s.Kind),
				slog.Int("desiredReplicas", int(desiredReplicas)),
				slog.Int("readyReplicas", int(s.Status.ReadyReplicas)),
				slog.Int("availableReplicas", int(s.Status.AvailableReplicas)),
				slog.Int("updatedReplicas", int(s.Status.UpdatedReplicas)),
			)
		}

		cntinue = ss.Continue
		if cntinue == "" {
			break
		}
	}

	return statefulsets, nil
}

func (r Reporter) events(ctx context.Context, name, kind string) ([]tmpl.Event, error) {
	var events []tmpl.Event
	evList, err := r.kubeClient.
		CoreV1().
		Events(r.conf.Namespace).
		List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("involvedObject.name=%s", name),
			TypeMeta:      metav1.TypeMeta{Kind: kind},
		})
	if err != nil {
		return nil, err
	}
	for _, ev := range evList.Items {
		events = append(events, tmpl.Event{
			Type:    ev.Type,
			Reason:  ev.Reason,
			Message: ev.Message,
		})
	}
	return events, nil
}
