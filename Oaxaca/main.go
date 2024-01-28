package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"TeamProject30/backend/database" // Update with your actual project path
)

//go:embed menu.html
var content embed.FS

var templates = template.Must(template.ParseFS(content, "menu.html"))

func main() {
	http.HandleFunc("/menu", menuHandler)
	http.ListenAndServe(":8080", nil)
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch menu items from database
	menuItems := database.QueryMenu()

	// Render menu page template with menu items
	renderTemplate(w, menuItems)
}

func renderTemplate(w http.ResponseWriter, data interface{}) {
	err := templates.ExecuteTemplate(w, "menu.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
