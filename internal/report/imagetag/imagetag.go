package imagetag

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/db"
	"github.com/accuknox/rinc/internal/report"
	types "github.com/accuknox/rinc/types/imagetag"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Reporter is the image tag reporter.
type Reporter struct {
	kubeClient *kubernetes.Clientset
	conf       conf.ImageTag
	mongo      *mongo.Client
}

// NewReporter creates a new image tag reporter.
func NewReporter(c conf.ImageTag, k *kubernetes.Clientset, mongo *mongo.Client) Reporter {
	return Reporter{
		conf:       c,
		kubeClient: k,
		mongo:      mongo,
	}
}

// Report satisfies the report.Reporter interface by fetching the image tags of
// deployments and statefulsets from the Kubernetes API server, and writes the
// report to the database.
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

	metrics := types.Metrics{
		Timestamp:    now,
		Deployments:  depls,
		Statefulsets: statefulsets,
	}

	result, err := db.Database(r.mongo).
		Collection(db.CollectionImageTag).
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
		"imagetag: inserted document into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	alerts := report.SoftEvaluateAlerts(ctx, r.conf.Alerts, metrics)
	result, err = db.
		Database(r.mongo).
		Collection(db.CollectionAlerts).
		InsertOne(ctx, bson.M{
			"timestamp": now,
			"from":      db.CollectionImageTag,
			"alerts":    alerts,
		})
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"imagetag: inserting alerts into mongodb",
			slog.Time("timestamp", now),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("inserting alerts into mongodb: %w", err)
	}
	slog.LogAttrs(
		ctx,
		slog.LevelDebug,
		"imagetag: inserted alerts into mongodb",
		slog.Any("insertedId", result.InsertedID),
	)

	return nil
}

func (r Reporter) deployments(ctx context.Context) ([]types.Resource, error) {
	var resources []types.Resource
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
			containers := append(
				d.Spec.Template.Spec.InitContainers,
				d.Spec.Template.Spec.Containers...,
			)
			images := make([]types.Image, len(containers))
			for idx, c := range containers {
				images[idx] = types.Image{
					Name:              c.Image,
					FromInitContainer: idx < len(d.Spec.Template.Spec.InitContainers),
				}
			}
			resources = append(resources, types.Resource{
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

func (r Reporter) statefulsets(ctx context.Context) ([]types.Resource, error) {
	var resources []types.Resource
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
			containers := append(
				s.Spec.Template.Spec.InitContainers,
				s.Spec.Template.Spec.Containers...,
			)
			images := make([]types.Image, len(containers))
			for idx, c := range containers {
				images[idx] = types.Image{
					Name:              c.Image,
					FromInitContainer: idx < len(s.Spec.Template.Spec.InitContainers),
				}
			}
			resources = append(resources, types.Resource{
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
