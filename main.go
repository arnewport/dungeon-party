package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	//"github.com/labstack/echo/v4/middleware"
)

// Echo renderer wrapper around html/template
type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Global templates (kept for parity with your original code; now lives in renderer)
var tmpl *template.Template

func main() {
	// Parse all templates in the templates folder
	var err error
	tmpl, err = template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	e := echo.New()

	// Attach renderer so c.Render works
	e.Renderer = &TemplateRenderer{templates: tmpl}

	// Static files (e.g. /static/js/main.js)
	e.Static("/static", "static")

	// Routes
	e.GET("/", indexHandler)

	log.Println("Server running on http://localhost:8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Handlers
func indexHandler(c echo.Context) error {
	data := map[string]any{
		"Title": "Welcome to Dungeon Party",
	}
	return c.Render(http.StatusOK, "index.html", data)
}
