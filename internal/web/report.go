package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/accuknox/rinc/internal/util"
	"github.com/accuknox/rinc/view"
	"github.com/accuknox/rinc/view/layout"

	"github.com/labstack/echo/v4"
)

func (s Srv) Report(c echo.Context) error {
	id := c.Param("id")
	tmpl := c.Param("template")
	title := fmt.Sprintf("%s - %s | AccuKnox Reports", id, strings.ToTitle(tmpl))

	loc := filepath.Join(s.conf.Output, id, tmpl+".html")
	exists, err := util.FileExists(loc)
	if err != nil {
		return render(renderParams{
			Ctx: c,
			Component: layout.Base(
				title,
				view.Error(
					err.Error(),
					http.StatusInternalServerError,
				),
			),
			Status: http.StatusInternalServerError,
		})
	}

	if !exists {
		return render(renderParams{
			Ctx: c,
			Component: layout.Base(
				title,
				view.Error(
					"requested report is not available",
					http.StatusNotFound,
				),
			),
			Status: http.StatusNotFound,
		})
	}

	t, err := template.ParseFiles(loc)
	if err != nil {
		return render(renderParams{
			Ctx: c,
			Component: layout.Base(
				title,
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
				title,
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
