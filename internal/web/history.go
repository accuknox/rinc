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

func (s Srv) HistoryPage(c echo.Context) error {
	return render(renderParams{
		Ctx: c,
		Component: layout.Base(
			"History | AccuKnox Reports",
			partial.Navbar(false),
			view.HistoryForm(),
		),
	})
}

type historySearchParams struct {
	Date string `form:"date"`
}

type docWithStamp struct {
	Timestamp time.Time `bson:"timestamp"`
}

func (s Srv) HistorySearch(c echo.Context) error {
	params := new(historySearchParams)
	if err := c.Bind(params); err != nil {
		return render(renderParams{
			Ctx: c,
			Component: view.Error(
				"failed to parse date",
				http.StatusBadRequest,
			),
			Status: http.StatusBadRequest,
		})
	}

	date, err := time.Parse(util.HTMLFormDateLayout, params.Date)
	if err != nil {
		return render(renderParams{
			Ctx: c,
			Component: view.Error(
				"failed to parse date",
				http.StatusBadRequest,
			),
			Status: http.StatusBadRequest,
		})
	}
	eod := date.Add(time.Hour * 24)

	var results []view.SearchResults
	for _, coll := range db.Collections {
		cursor, err := db.
			Database(s.mongo).
			Collection(coll).
			Find(c.Request().Context(), bson.M{
				"timestamp": bson.M{
					"$gte": date,
					"$lt":  eod,
				},
			})
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				continue
			}
			return render(renderParams{
				Ctx: c,
				Component: view.Error(
					err.Error(),
					http.StatusInternalServerError,
				),
				Status: http.StatusInternalServerError,
			})
		}
	Next:
		for cursor.Next(c.Request().Context()) {
			t := new(docWithStamp)
			err := cursor.Decode(t)
			if err != nil {
				cursor.Close(c.Request().Context())
				return render(renderParams{
					Ctx: c,
					Component: view.Error(
						err.Error(),
						http.StatusInternalServerError,
					),
					Status: http.StatusInternalServerError,
				})
			}
			for _, r := range results {
				if r.Timestamp.Equal(t.Timestamp) {
					continue Next
				}
			}
			hr, min, _ := t.Timestamp.UTC().Clock()
			results = append(results, view.SearchResults{
				ID:                     t.Timestamp.Format(util.IsosecLayout),
				Timestamp:              t.Timestamp,
				HumanReadableTimestamp: fmt.Sprintf("%02d:%02d UTC", hr, min),
			})
		}
		cursor.Close(c.Request().Context())
	}

	if len(results) == 0 {
		return render(renderParams{
			Ctx:       c,
			Component: view.HistorySearchResultEmpty(),
		})
	}

	return render(renderParams{
		Ctx:       c,
		Component: view.HistorySearchResult(results),
	})
}
