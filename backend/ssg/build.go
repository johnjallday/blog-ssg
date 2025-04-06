package ssg

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
)

// RenderMarkdown converts Markdown content to HTML.
func RenderMarkdown(md []byte) (template.HTML, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(md, &buf); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// BentoData holds content for the index page.
type BentoData struct {
	Portfolio  template.HTML
	Blog       template.HTML
	MusicTools template.HTML
	DevTools   template.HTML
}

// BuildHome generates the index page.
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

// BuildPortfolio generates the portfolio page.
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

// PostData holds data for an individual post.
type PostData struct {
	Title   string
	Content template.HTML
}

// BuildPosts converts Markdown posts from content/posts to HTML pages.
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

// BuildBlogIndex creates a static blog index page listing all posts.
func BuildBlogIndex() {
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

	tmpl, err := template.ParseFiles(filepath.Join("backend", "ssg", "templates", "blog_index.html"))
	if err != nil {
		log.Fatalf("Error parsing blog index template: %v", err)
	}

	outFile, err := os.Create(filepath.Join("public", "blog-index.html"))
	if err != nil {
		log.Fatalf("Error creating blog index file: %v", err)
	}
	defer outFile.Close()

	if err := tmpl.Execute(outFile, posts); err != nil {
		log.Fatalf("Error executing blog index template: %v", err)
	}

	log.Println("Built public/blog-index.html successfully!")
}

// BuildBlog generates the blog landing page.
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
			{"Code Snippet Manager", "/assets/tools/snippet-manager.zip"},
			{"Regex Tester", "/assets/tools/regex-tester.zip"},
			{"JSON Formatter", "/assets/tools/json-formatter.zip"},
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

// BuildAll runs all build functions.
func BuildAll() {
	if err := os.MkdirAll("public", os.ModePerm); err != nil {
		log.Fatalf("Error creating public directory: %v", err)
	}

	BuildHome()
	BuildPortfolio()
	BuildBlog()
	BuildPosts()
	BuildBlogIndex()
	BuildMusicTools()
	BuildDevTools()
	AddStyle()
	CopyAssets()
}
