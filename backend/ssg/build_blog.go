package ssg

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
)

func RenderMarkdown(md []byte) (template.HTML, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(md, &buf); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

type PostData struct {
	Title   string
	Content template.HTML
}

func BuildPosts() {
	tmpl, err := template.ParseFiles(filepath.Join("backend", "ssg", "templates", "posts.html"))
	if err != nil {
		log.Fatalf("Error parsing posts template: %v", err)
	}

	outDir := filepath.Join("public", "posts")
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating posts output directory: %v", err)
	}

	postsDir := filepath.Join("content", "posts")
	files, err := os.ReadDir(postsDir)
	if err != nil {
		log.Fatalf("Error reading posts directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			slug := file.Name()[0 : len(file.Name())-3]
			mdPath := filepath.Join(postsDir, file.Name())

			mdContent, err := os.ReadFile(mdPath)
			if err != nil {
				log.Printf("Error reading markdown file %s: %v", file.Name(), err)
				continue
			}

			htmlContent, err := RenderMarkdown(mdContent)
			if err != nil {
				log.Printf("Error converting markdown for %s: %v", file.Name(), err)
				continue
			}

			postData := PostData{
				Title:   slug,
				Content: htmlContent,
			}

			outFilePath := filepath.Join(outDir, slug+".html")
			outFile, err := os.Create(outFilePath)
			if err != nil {
				log.Printf("Error creating output file for %s: %v", slug, err)
				continue
			}

			if err := tmpl.Execute(outFile, postData); err != nil {
				log.Printf("Error executing template for %s: %v", slug, err)
			} else {
				log.Printf("Built post %s", slug)
			}

			outFile.Close()
		}
	}
}

// BlogPost represents a blog post for the index.
type BlogPost struct {
	Slug  string
	Title string
}

func BuildBlog() {
	postsDir := filepath.Join("content", "posts")
	files, err := os.ReadDir(postsDir)
	if err != nil {
		log.Fatalf("Error reading posts directory: %v", err)
	}

	var posts []BlogPost
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			slug := file.Name()[0 : len(file.Name())-3]
			posts = append(posts, BlogPost{
				Slug:  slug,
				Title: slug,
			})
		}
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("backend", "ssg", "templates", "blog.html"),
		filepath.Join("backend", "ssg", "templates", "nav.html"),
		filepath.Join("backend", "ssg", "templates", "contact.html"),
	)
	if err != nil {
		log.Fatalf("Error parsing blog templates: %v", err)
	}

	out, err := os.Create(filepath.Join("public", "blog.html"))
	if err != nil {
		log.Fatalf("Error creating blog.html: %v", err)
	}
	defer out.Close()

	if err := tmpl.Execute(out, posts); err != nil {
		log.Fatalf("Error executing blog template: %v", err)
	}

	log.Println("public/blog.html built successfully!")
}
