package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// PostHandler serves an individual post from the generated static files.
// It expects a query parameter "slug" that corresponds to a Markdown file's slug.
func PostHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	// Construct the path to the built post file.
	postPath := filepath.Join("public", "posts", slug+".html")
	if _, err := os.Stat(postPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Serve the file directly.
	http.ServeFile(w, r, postPath)
}

// BlogIndexHandler serves the blog index page from the generated static files.
func BlogIndexHandler(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join("public", "blog-index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, indexPath)
}

// RunServer sets up and starts the HTTP server.
func RunServer() {
	// Serve static files from the public folder.
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// Register dynamic routes.
	http.HandleFunc("/post", PostHandler)
	http.HandleFunc("/blog-index", BlogIndexHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
