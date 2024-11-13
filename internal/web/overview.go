package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/accuknox/rinc/internal/db"
	"github.com/accuknox/rinc/internal/util"
	"github.com/accuknox/rinc/view"
	"github.com/accuknox/rinc/view/layout"
	"github.com/accuknox/rinc/view/partial"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s Srv) Overview(c echo.Context) error {
	id := c.Param("id")
	title := fmt.Sprintf("%s - Overview | Accuknox Reports", id)
	at, err := time.Parse(util.IsosecLayout, id)
	if err != nil {
		return render(renderParams{
			Ctx: c,
			Component: layout.Base(
				title,
				view.Error(
					"failed to parse timestamp",
					http.StatusBadRequest,
				),
			),
			Status: http.StatusBadRequest,
		})
	}

	var statuses []view.OverviewStatus

	for _, coll := range db.Collections {
		result := db.
			Database(s.mongo).
			Collection(coll).
			FindOne(c.Request().Context(), bson.M{
				"timestamp": at,
			})
		if err := result.Err(); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				continue
			}
			return render(renderParams{
				Ctx: c,
				Component: layout.Base(
					"AccuKnox Reports",
					view.Error(
						err.Error(),
						http.StatusInternalServerError,
					),
				),
				Status: http.StatusInternalServerError,
			})
		}
		count, err := s.fetchAlertsCount(c.Request().Context(), coll, at)
		if err != nil {
			return render(renderParams{
				Ctx: c,
				Component: layout.Base(
					"AccuKnox Reports",
					view.Error(
						err.Error(),
						http.StatusInternalServerError,
					),
				),
				Status: http.StatusInternalServerError,
			})
		}
		switch coll {
		case db.CollectionRabbitmq:
			statuses = append(statuses, view.OverviewStatus{
				Name:        "RabbitMQ",
				Slug:        "rabbitmq",
				ID:          id,
				AlertsCount: count,
			})
		case db.CollectionCeph:
			statuses = append(statuses, view.OverviewStatus{
				Name:        "CEPH",
				Slug:        "ceph",
				ID:          id,
				AlertsCount: count,
			})
		case db.CollectionDass:
			statuses = append(statuses, view.OverviewStatus{
				Name:        "Deployment & Statefulset Status",
				Slug:        "deployment-and-statefulset-status",
				ID:          id,
				AlertsCount: count,
			})
		case db.CollectionLongJobs:
			statuses = append(statuses, view.OverviewStatus{
				Name:        "Long Running Jobs",
				Slug:        "longjobs",
				ID:          id,
				AlertsCount: count,
			})
		case db.CollectionImageTag:
			statuses = append(statuses, view.OverviewStatus{
				Name:        "Image Tags",
				Slug:        "imagetags",
				ID:          id,
				AlertsCount: count,
			})
		case db.CollectionResourceUtilization:
			statuses = append(statuses, view.OverviewStatus{
				Name:        "Resource Utilization",
				Slug:        "resource-utilization",
				ID:          id,
				AlertsCount: count,
			})
		}
	}

	if len(statuses) == 0 {
		return render(renderParams{
			Ctx: c,
			Component: layout.Base(
				title,
				partial.Navbar(true),
				view.Error(
					"Kindly make sure that the URL is correct",
					http.StatusNotFound,
				),
			),
			Status: http.StatusNotFound,
		})
	}

	return render(renderParams{
		Ctx: c,
		Component: layout.Base(
			title,
			partial.Navbar(true),
			view.Overview(statuses),
			partial.Footer(at),
		),
	})
}

func (s Srv) fetchAlertsCount(ctx context.Context, from string, at time.Time) (view.AlertsCount, error) {
	cursor, err := db.Database(s.mongo).Collection("alerts").Find(ctx, bson.M{
		"from":      from,
		"timestamp": at,
	})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("finding alerts from %q at %v: %w", from, at, err)
	}
	defer cursor.Close(ctx)
	count := make(view.AlertsCount, 3)
	for cursor.Next(ctx) {
		alerts := new(db.AlertDocument)
		err := cursor.Decode(alerts)
		if err != nil {
			return nil, fmt.Errorf("decoding document at cursor: %w", err)
		}
		for _, alert := range alerts.Alerts {
			count[alert.Severity]++
		}
	}
	return count, nil
}
