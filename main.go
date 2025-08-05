package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

// Global templates
var tmpl *template.Template

func main() {
	// Parse all templates in the templates folder
	var err error
	tmpl, err = template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	// Set up router
	r := chi.NewRouter()

	// Static file server (e.g. /static/js/main.js)
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Routes
	r.Get("/", indexHandler)

	// Start server
	log.Println("Server running on http://localhost:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Handlers
func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "Welcome to Dungeon Party",
	}
	err := tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
