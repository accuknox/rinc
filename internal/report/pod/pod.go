package pod

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/db"
	"github.com/accuknox/rinc/internal/report"
	types "github.com/accuknox/rinc/types/pod"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// Reporter is the pod status reporter.
type Reporter struct {
	kubeClient *kubernetes.Clientset
	conf       conf.PodStatus
	mongo      *mongo.Client
}

// NewReporter creates a new pod status reporter.
func NewReporter(c conf.PodStatus, k *kubernetes.Clientset, mongo *mongo.Client) Reporter {
	return Reporter{
		conf:       c,
		kubeClient: k,
		mongo:      mongo,
	}
}

// Report satisfies the report.Reporter interface by fetching the status of
// pods from the Kubernetes API server, and writes the report to the database.
func (r Reporter) Report(ctx context.Context, now time.Time) error {
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

	ss, err := r.statefulsets(ctx)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"fetching statefulset resources",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("fetching statefulsets: %w", err)
	}

	metrics := types.Metrics{
		Timestamp:    now,
		Deployments:  depls,
		Statefulsets: ss,
	}

	result, err := db.Database(r.mongo).
		Collection(db.CollectionPodStatus).
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
		"podStatus: inserted document into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	alerts := report.SoftEvaluateAlerts(ctx, r.conf.Alerts, metrics)
	result, err = db.
		Database(r.mongo).
		Collection(db.CollectionAlerts).
		InsertOne(ctx, bson.M{
			"timestamp": now,
			"from":      db.CollectionPodStatus,
			"alerts":    alerts,
		})
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"podStatus: inserting alerts into mongodb",
			slog.Time("timestamp", now),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("inserting alerts into mongodb: %w", err)
	}
	slog.LogAttrs(
		ctx,
		slog.LevelDebug,
		"podStatus: inserted alerts into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	return nil
}

func (r Reporter) deployments(ctx context.Context) ([]types.Resource, error) {
	var deployments []types.Resource
	var cntinue string

	for {
		depls, err := r.kubeClient.
			AppsV1().
			Deployments(r.conf.Namespace).
			List(ctx, metav1.ListOptions{
				Continue: cntinue,
				Limit:    15,
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
			podList, err := r.pods(ctx, d.Namespace, d.Spec.Selector.MatchLabels)
			if err != nil {
				slog.LogAttrs(
					ctx,
					slog.LevelError,
					"listing pods",
					slog.String("kind", d.Kind),
					slog.String("for", d.Name),
					slog.String("namespace", d.Namespace),
					slog.String("error", err.Error()),
				)
				continue
			}
			pods := make([]types.Pod, len(podList.Items))
			for idx, pod := range podList.Items {
				var containers []types.Container
				for _, c := range pod.Status.InitContainerStatuses {
					var lastTermState string
					if c.LastTerminationState.Terminated != nil {
						lastTermState = c.LastTerminationState.Terminated.Reason
					}
					containers = append(containers, types.Container{
						Name:                 c.Name,
						IsInit:               true,
						Ready:                c.Ready,
						State:                containerState(c.State),
						RestartCount:         c.RestartCount,
						LastTerminationState: lastTermState,
					})
				}
				for _, c := range pod.Status.ContainerStatuses {
					var lastTermState string
					if c.LastTerminationState.Terminated != nil {
						lastTermState = c.LastTerminationState.Terminated.Reason
					}
					containers = append(containers, types.Container{
						Name:                 c.Name,
						Ready:                c.Ready,
						State:                containerState(c.State),
						RestartCount:         c.RestartCount,
						LastTerminationState: lastTermState,
					})
				}
				pods[idx] = types.Pod{
					Name:       pod.Name,
					Status:     podStatus(pod.Status),
					QOSClass:   string(pod.Status.QOSClass),
					StartTime:  pod.Status.StartTime.Time,
					Containers: containers,
				}
			}
			deployments = append(deployments, types.Resource{
				Name:      d.Name,
				Namespace: d.Namespace,
				Pods:      pods,
			})
		}

		cntinue = depls.Continue
		if cntinue == "" {
			break
		}
	}

	return deployments, nil
}

func (r Reporter) statefulsets(ctx context.Context) ([]types.Resource, error) {
	var statefulsets []types.Resource
	var cntinue string

	for {
		ss, err := r.kubeClient.
			AppsV1().
			StatefulSets(r.conf.Namespace).
			List(ctx, metav1.ListOptions{
				Continue: cntinue,
				Limit:    15,
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

		for _, d := range ss.Items {
			podList, err := r.pods(ctx, d.Namespace, d.Spec.Selector.MatchLabels)
			if err != nil {
				slog.LogAttrs(
					ctx,
					slog.LevelError,
					"listing pods",
					slog.String("kind", d.Kind),
					slog.String("for", d.Name),
					slog.String("namespace", d.Namespace),
					slog.String("error", err.Error()),
				)
				continue
			}
			pods := make([]types.Pod, len(podList.Items))
			for idx, pod := range podList.Items {
				var containers []types.Container
				for _, c := range pod.Status.InitContainerStatuses {
					var lastTermState string
					if c.LastTerminationState.Terminated != nil {
						lastTermState = c.LastTerminationState.Terminated.Reason
					}
					containers = append(containers, types.Container{
						Name:                 c.Name,
						IsInit:               true,
						Ready:                c.Ready,
						State:                containerState(c.State),
						RestartCount:         c.RestartCount,
						LastTerminationState: lastTermState,
					})
				}
				for _, c := range pod.Status.ContainerStatuses {
					var lastTermState string
					if c.LastTerminationState.Terminated != nil {
						lastTermState = c.LastTerminationState.Terminated.Reason
					}
					containers = append(containers, types.Container{
						Name:                 c.Name,
						Ready:                c.Ready,
						State:                containerState(c.State),
						RestartCount:         c.RestartCount,
						LastTerminationState: lastTermState,
					})
				}
				pods[idx] = types.Pod{
					Name:       pod.Name,
					Status:     podStatus(pod.Status),
					QOSClass:   string(pod.Status.QOSClass),
					StartTime:  pod.Status.StartTime.Time,
					Containers: containers,
				}
			}
			statefulsets = append(statefulsets, types.Resource{
				Name:      d.Name,
				Namespace: d.Namespace,
				Pods:      pods,
			})
		}

		cntinue = ss.Continue
		if cntinue == "" {
			break
		}
	}

	return statefulsets, nil
}

func (r Reporter) pods(ctx context.Context, ns string, labels labels.Set) (*corev1.PodList, error) {
	selector := metav1.FormatLabelSelector(metav1.SetAsLabelSelector(labels))
	return r.kubeClient.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{
		LabelSelector: selector,
	})
}

func containerState(s corev1.ContainerState) string {
	if s.Running != nil {
		return "RUNNING"
	}
	out := "%s: Reason=%s"
	if s.Waiting != nil {
		return fmt.Sprintf(out, "WAITING", s.Waiting.Reason)
	}
	if s.Terminated != nil {
		return fmt.Sprintf(out, "TERMINATED", s.Terminated.Reason)
	}
	return ""
}

func podStatus(s corev1.PodStatus) string {
	// check if the pod was evicted
	if s.Phase == corev1.PodFailed && s.Reason == "Evicted" {
		return "Evicted"
	}

	hasWaiting := false
	hasRunning := false
	hasTerminatedWithError := false
	allTerminatedCompleted := true

	for _, cs := range s.ContainerStatuses {
		if cs.State.Waiting != nil {
			hasWaiting = true
			if cs.State.Waiting.Reason != "" {
				return cs.State.Waiting.Reason // e.g., "CrashLoopBackOff"
			}
			continue
		}

		if cs.State.Running != nil {
			hasRunning = true
			allTerminatedCompleted = false
			continue
		}

		if cs.State.Terminated != nil {
			switch cs.State.Terminated.Reason {
			case "Error":
				hasTerminatedWithError = true
				allTerminatedCompleted = false
			case "Completed":
				// no action needed, just mark as completed
			default:
				allTerminatedCompleted = false
			}
		}
	}

	if hasWaiting {
		return "Pending"
	}
	if hasTerminatedWithError {
		return "Error"
	}
	if allTerminatedCompleted {
		return "Completed"
	}
	if hasRunning {
		return "Running"
	}
	return "Unknown" // fallback status if none of the above applies
}
