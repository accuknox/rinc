package web

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/accuknox/rinc/internal/util"
	"github.com/accuknox/rinc/view"
	"github.com/accuknox/rinc/view/layout"
	"github.com/accuknox/rinc/view/partial"

	"github.com/labstack/echo/v4"
)

func (s Srv) HistoryPage(c echo.Context) error {
	return render(renderParams{
		Ctx: c,
		Component: layout.Base(
			"History | AccuKnox Reports",
			partial.Navbar(false, false),
			view.HistoryForm(),
		),
		Status: http.StatusOK,
	})
}

type historySearchParams struct {
	Date string `form:"date"`
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

	params.Date = strings.ReplaceAll(params.Date, "-", "")
	if params.Date == "" {
		return render(renderParams{
			Ctx: c,
			Component: view.Error(
				"date cannot be blank",
				http.StatusBadRequest,
			),
			Status: http.StatusBadRequest,
		})
	}

	entries, err := os.ReadDir(s.conf.Output)
	if err != nil {
		return render(renderParams{
			Ctx: c,
			Component: view.Error(
				fmt.Sprintf("failed to read output directory %q", s.conf.Output),
				http.StatusInternalServerError,
			),
			Status: http.StatusInternalServerError,
		})
	}

	var targets []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if !util.IsIsosec(entry.Name()) {
			continue
		}
		if strings.HasPrefix(entry.Name(), params.Date) {
			targets = append(targets, entry.Name())
		}
	}

	if len(targets) == 0 {
		return render(renderParams{
			Ctx:       c,
			Component: view.HistorySearchResultEmpty(),
			Status:    http.StatusOK,
		})
	}

	results := make([]view.SearchResults, len(targets))
	for idx, target := range targets {
		t, err := time.Parse(util.IsosecLayout, target)
		if err != nil {
			continue
		}
		hr, min, _ := t.UTC().Clock()
		results[idx] = view.SearchResults{
			ID:                     target,
			Timestamp:              t,
			HumanReadableTimestamp: fmt.Sprintf("%02d:%02d UTC", hr, min),
		}
	}

	return render(renderParams{
		Ctx:       c,
		Component: view.HistorySearchResult(results),
		Status:    http.StatusOK,
	})
}
