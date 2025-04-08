package ssg

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
)

// BentoData holds content for the index page.
type BentoData struct {
	Portfolio  template.HTML
	Blog       template.HTML
	MusicTools template.HTML
	DevTools   template.HTML
}

func BuildHome() {
	// Read partial templates from backend/ssg/templates/partials.
	partialsDir := filepath.Join("backend", "ssg", "templates", "partials")

	portfolioContent, err := os.ReadFile(filepath.Join(partialsDir, "portfolio.html"))
	if err != nil {
		log.Fatalf("Error reading portfolio partial: %v", err)
	}
	blogContent, err := os.ReadFile(filepath.Join(partialsDir, "blog.html"))
	if err != nil {
		log.Fatalf("Error reading blog partial: %v", err)
	}
	musicToolsContent, err := os.ReadFile(filepath.Join(partialsDir, "music-tools.html"))
	if err != nil {
		log.Fatalf("Error reading music tools partial: %v", err)
	}
	devToolsContent, err := os.ReadFile(filepath.Join(partialsDir, "dev-tools.html"))
	if err != nil {
		log.Fatalf("Error reading dev tools partial: %v", err)
	}

	data := BentoData{
		Portfolio:  template.HTML(portfolioContent),
		Blog:       template.HTML(blogContent),
		MusicTools: template.HTML(musicToolsContent),
		DevTools:   template.HTML(devToolsContent),
	}

	// Parse master and contact templates.
	tmpl, err := template.ParseFiles(
		filepath.Join("backend", "ssg", "templates", "master.html"),
		filepath.Join("backend", "ssg", "templates", "contact.html"),
	)
	if err != nil {
		log.Fatalf("Error parsing master template: %v", err)
	}

	out, err := os.Create(filepath.Join("public", "index.html"))
	if err != nil {
		log.Fatalf("Error creating index.html: %v", err)
	}
	defer out.Close()

	if err := tmpl.Execute(out, data); err != nil {
		log.Fatalf("Error executing master template: %v", err)
	}

	log.Println("public/index.html built successfully!")
}
