package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// PostHandler serves individual blog posts like /blog/posts/test
func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Trim the route prefix to get the slug
	slug := filepath.Base(r.URL.Path)

	if slug == "" || slug == "posts" {
		http.NotFound(w, r)
		return
	}

	postPath := filepath.Join("public", "posts", slug+".html")
	if _, err := os.Stat(postPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, postPath)
}

// BlogHandler serves the main blog page.
func BlogHandler(w http.ResponseWriter, r *http.Request) {
	blogPath := filepath.Join("public", "blog.html")
	if _, err := os.Stat(blogPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, blogPath)
}

// BlogIndexHandler serves the blog index page.
func BlogIndexHandler(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join("public", "blog-index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, indexPath)
}

// PortfolioHandler serves the portfolio page.
func PortfolioHandler(w http.ResponseWriter, r *http.Request) {
	portfolioPath := filepath.Join("public", "portfolio.html")
	if _, err := os.Stat(portfolioPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, portfolioPath)
}

// MusicToolsHandler serves the music tools page.
func MusicToolsHandler(w http.ResponseWriter, r *http.Request) {
	musicPath := filepath.Join("public", "music-tools.html")
	if _, err := os.Stat(musicPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, musicPath)
}

// DevToolsHandler serves the dev tools page.
func DevToolsHandler(w http.ResponseWriter, r *http.Request) {
	devPath := filepath.Join("public", "dev-tools.html")
	if _, err := os.Stat(devPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, devPath)
}

// RunServer sets up and starts the HTTP server.
func RunServer() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// Updated route
	http.HandleFunc("/blog/posts/", PostHandler)

	// Other routes
	http.HandleFunc("/blog", BlogHandler)
	http.HandleFunc("/blog-index", BlogIndexHandler)
	http.HandleFunc("/portfolio", PortfolioHandler)
	http.HandleFunc("/music-tools", MusicToolsHandler)
	http.HandleFunc("/dev-tools", DevToolsHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
