package web

import (
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
	timestamp, err := time.Parse(util.IsosecLayout, id)
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
				"timestamp": timestamp,
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
		switch coll {
		case db.CollectionRabbitmq:
			statuses = append(statuses, view.OverviewStatus{
				Name: "RabbitMQ",
				Slug: "rabbitmq",
				ID:   id,
			})
		case db.CollectionCeph:
			statuses = append(statuses, view.OverviewStatus{
				Name: "CEPH",
				Slug: "ceph",
				ID:   id,
			})
		case db.CollectionDass:
			statuses = append(statuses, view.OverviewStatus{
				Name: "Deployment & Statefulset Status",
				Slug: "deployment-and-statefulset-status",
				ID:   id,
			})
		case db.CollectionLongJobs:
			statuses = append(statuses, view.OverviewStatus{
				Name: "Long Running Jobs",
				Slug: "longjobs",
				ID:   id,
			})
		case db.CollectionImageTag:
			statuses = append(statuses, view.OverviewStatus{
				Name: "Image Tags",
				Slug: "imagetags",
				ID:   id,
			})
		}
	}

	if len(statuses) == 0 {
		return render(renderParams{
			Ctx: c,
			Component: layout.Base(
				title,
				partial.Navbar(true, true),
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
			partial.Navbar(true, true),
			view.Overview(statuses),
			partial.Footer(timestamp),
		),
	})
}
