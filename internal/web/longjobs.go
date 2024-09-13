package web

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/accuknox/rinc/view"
	"github.com/accuknox/rinc/view/layout"

	"github.com/labstack/echo/v4"
)

func (s Srv) LongRunningJobs(c echo.Context) error {
	id := c.Param("id")
	loc := filepath.Join(s.conf.Output, id, "longrunningjobs.html")
	t, err := template.ParseFiles(loc)
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
