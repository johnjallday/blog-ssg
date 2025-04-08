package ssg

import (
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
)

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

// copyDir recursively copies a directory from src to dst.
func copyDir(src string, dst string) error {
	// Create destination directory.
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyAssets moves the static assets from content/assets to public/assets.
func CopyAssets() {
	src := filepath.Join("content", "assets")
	dst := filepath.Join("public", "assets")
	if err := copyDir(src, dst); err != nil {
		log.Fatalf("Error copying assets: %v", err)
	}
	log.Println("Assets copied successfully!")
}

// AddStyle copies the styles.css file to the public directory.
func AddStyle() {
	src := filepath.Join("backend", "ssg", "templates", "styles.css")
	dst := filepath.Join("public", "styles.css")

	if err := copyFile(src, dst); err != nil {
		log.Fatalf("Error copying styles.css: %v", err)
	}
	log.Println("styles.css copied successfully!")
}

// BuildMusicTools generates the music tools page.
func BuildMusicTools() {
	// Optional: define some dynamic data if needed in the future
	type MusicToolsData struct {
		Tools []struct {
			Name string
			Link string
		}
	}

	// Static for now — feel free to make it dynamic later
	data := MusicToolsData{
		Tools: []struct {
			Name string
			Link string
		}{
			{"BeatMaker", "/assets/tools/beatmaker.zip"},
			{"Chord Generator", "/assets/tools/chordgen.zip"},
			{"Synth Pack", "/assets/tools/synthpack.zip"},
		},
	}

	// Parse the required templates
	tmpl, err := template.ParseFiles(
		filepath.Join("backend", "ssg", "templates", "music-tools.html"),
		filepath.Join("backend", "ssg", "templates", "nav.html"),
		filepath.Join("backend", "ssg", "templates", "contact.html"),
	)
	if err != nil {
		log.Fatalf("Error parsing music tools template: %v", err)
	}

	// Create the output file
	out, err := os.Create(filepath.Join("public", "music-tools.html"))
	if err != nil {
		log.Fatalf("Error creating music-tools.html: %v", err)
	}
	defer out.Close()

	// Render the page with the provided data
	if err := tmpl.Execute(out, data); err != nil {
		log.Fatalf("Error executing music tools template: %v", err)
	}

	log.Println("public/music-tools.html built successfully!")
}

// BuildDevTools generates the dev tools page.
// BuildAll runs all build functions.
func BuildAll() {
	if err := os.MkdirAll("public", os.ModePerm); err != nil {
		log.Fatalf("Error creating public directory: %v", err)
	}

	BuildHome()
	BuildPortfolio()
	BuildBlog()
	BuildPosts()
	BuildMusicTools()
	BuildDevTools()
	AddStyle()
	CopyAssets()
}
