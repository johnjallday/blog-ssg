package ssg

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
)

func BuildPortfolio() {
	type PortfolioData struct {
		Playlist1 template.HTML
		Playlist2 template.HTML
	}
	data := PortfolioData{
		Playlist1: template.HTML(`<div class="playlist">
    <iframe width="560" height="315" src="https://www.youtube.com/embed/videoseries?list=PLXmbGsvEq-1DqQtHyOqedD8dfUGbWLrR5" frameborder="0" allowfullscreen></iframe>
</div>`),
		Playlist2: template.HTML(`<div class="playlist">
    <iframe width="560" height="315" src="https://www.youtube.com/embed/videoseries?list=PLXmbGsvEq-1D37iuSB7Cahw-HuMB2_Wa6" frameborder="0" allowfullscreen></iframe>
</div>`),
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("backend", "ssg", "templates", "portfolio.html"),
		filepath.Join("backend", "ssg", "templates", "nav.html"),
		filepath.Join("backend", "ssg", "templates", "contact.html"),
	)
	if err != nil {
		log.Fatalf("Error parsing portfolio template: %v", err)
	}

	out, err := os.Create(filepath.Join("public", "portfolio.html"))
	if err != nil {
		log.Fatalf("Error creating portfolio.html: %v", err)
	}
	defer out.Close()

	if err := tmpl.Execute(out, data); err != nil {
		log.Fatalf("Error executing portfolio template: %v", err)
	}

	log.Println("public/portfolio.html built successfully!")
}
