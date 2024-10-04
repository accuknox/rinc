package web

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/accuknox/rinc/internal/util"
	"github.com/accuknox/rinc/view"
	"github.com/accuknox/rinc/view/layout"

	"github.com/labstack/echo/v4"
)

func (s Srv) Index(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		entries, err := os.ReadDir(s.conf.Output)
		if err != nil {
			return render(renderParams{
				Ctx: c,
				Component: layout.Base(
					"AccuKnox Reports",
					view.Error(
						"failed to read latest report",
						http.StatusInternalServerError,
					),
				),
				Status: http.StatusInternalServerError,
			})
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			if !util.IsIsosec(entry.Name()) {
				continue
			}
			if entry.Name() > id {
				id = entry.Name()
			}
		}
	}

	if id == "" {
		return render(renderParams{
			Ctx:       c,
			Component: layout.Base("AccuKnox Reports", view.EmptyIndex()),
			Status:    http.StatusOK,
		})
	}

	t, err := template.ParseFiles(filepath.Join(s.conf.Output, id, "index.html"))
	if err != nil {
		return render(renderParams{
			Ctx: c,
			Component: layout.Base(
				"AccuKnox Reports",
				view.Error(
					"failed to parse template",
					http.StatusInternalServerError,
				),
			),
			Status: http.StatusInternalServerError,
		})
	}

	err = t.Execute(c.Response(), nil)
	if err != nil {
		return render(renderParams{
			Ctx: c,
			Component: layout.Base(
				"AccuKnox Reports",
				view.Error(
					"failed to execute template",
					http.StatusInternalServerError,
				),
			),
			Status: http.StatusInternalServerError,
		})
	}

	return nil
}
