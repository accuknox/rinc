package web

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/accuknox/rinc/internal/db"
	"github.com/accuknox/rinc/internal/util"
	types "github.com/accuknox/rinc/types/connectivity"
	"github.com/accuknox/rinc/view"
	tmpl "github.com/accuknox/rinc/view/connectivity"
	"github.com/accuknox/rinc/view/layout"
	"github.com/accuknox/rinc/view/partial"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s Srv) Connectivity(c echo.Context) error {
	id := c.Param("id")
	title := fmt.Sprintf("%s - Connectivity | AccuKnox Reports", id)
	timestamp, err := time.Parse(util.IsosecLayout, id)
	if err != nil {
		return render(renderParams{
			Ctx: c,
			Component: layout.Base(
				title,
				partial.Navbar(false),
				view.Error(
					"failed to parse timestamp",
					http.StatusBadRequest,
				),
			),
			Status: http.StatusBadRequest,
		})
	}

	result := db.
		Database(s.mongo).
		Collection(db.CollectionConnectivity).
		FindOne(c.Request().Context(), bson.M{
			"timestamp": timestamp,
		})
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return render(renderParams{
				Ctx: c,
				Component: layout.Base(
					title,
					partial.Navbar(false),
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
				partial.Navbar(false),
				view.Error(
					err.Error(),
					http.StatusInternalServerError,
				),
			),
			Status: http.StatusInternalServerError,
		})
	}

	metrics := new(types.Metrics)
	if err := result.Decode(metrics); err != nil {
		return render(renderParams{
			Ctx: c,
			Component: layout.Base(
				title,
				partial.Navbar(false),
				view.Error(
					err.Error(),
					http.StatusInternalServerError,
				),
			),
			Status: http.StatusInternalServerError,
		})
	}

	result = db.
		Database(s.mongo).
		Collection(db.CollectionAlerts).
		FindOne(c.Request().Context(), bson.M{
			"timestamp": timestamp,
			"from":      db.CollectionConnectivity,
		})
	err = result.Err()
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return render(renderParams{
			Ctx: c,
			Component: layout.Base(
				title,
				partial.Navbar(false),
				view.Error(
					err.Error(),
					http.StatusInternalServerError,
				),
			),
			Status: http.StatusInternalServerError,
		})
	}

	alerts := new(db.AlertDocument)

	if err == nil {
		err := result.Decode(&alerts)
		if err != nil {
			return render(renderParams{
				Ctx: c,
				Component: layout.Base(
					title,
					partial.Navbar(false),
					view.Error(
						err.Error(),
						http.StatusInternalServerError,
					),
				),
				Status: http.StatusInternalServerError,
			})
		}
	}

	return render(renderParams{
		Ctx: c,
		Component: layout.Base(
			title,
			partial.Navbar(false),
			tmpl.Report(*metrics, alerts.Alerts),
		),
	})
}
