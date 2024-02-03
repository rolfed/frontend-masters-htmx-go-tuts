package main

import (
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Templates struct {
    tmpl *template.Template
}

func newTemplate() *Templates {
    return &Templates{
        tmpl: template.Must(template.ParseGlob("views/*.html")),
    }
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.tmpl.ExecuteTemplate(w, name, data)
}

type Count struct {
    Count int
}

func main() {
    e := echo.New()
    count := Count{Count: 0}
    e.Renderer = newTemplate()

    // Debug Mode
    e.Debug = true
    e.Use(middleware.Logger())

    // Handler
    e.GET("/", func(c echo.Context) error {
        count.Count++
        return c.Render(http.StatusOK, "index.html", count)
    })

    // Start server
    e.Logger.Fatal(e.Start(":42069"))
}
