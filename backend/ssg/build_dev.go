package ssg

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func BuildDevTools() {
	type DevToolsData struct {
		Tools []struct {
			Name string
			Link string
		}
	}

	data := DevToolsData{
		Tools: []struct {
			Name string
			Link string
		}{
			{"Gorani Coding Agent", "https://github.com/johnjallday/gorani-coder"},
			{"Flow Workspace Manager", "https://github.com/johnjallday/flow-workspace"},
		},
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("backend", "ssg", "templates", "dev-tools.html"),
		filepath.Join("backend", "ssg", "templates", "nav.html"),
		filepath.Join("backend", "ssg", "templates", "contact.html"),
	)
	if err != nil {
		log.Fatalf("Error parsing dev tools template: %v", err)
	}

	out, err := os.Create(filepath.Join("public", "dev-tools.html"))
	if err != nil {
		log.Fatalf("Error creating dev-tools.html: %v", err)
	}
	defer out.Close()

	if err := tmpl.Execute(out, data); err != nil {
		log.Fatalf("Error executing dev tools template: %v", err)
	}

	log.Println("public/dev-tools.html built successfully!")
}
