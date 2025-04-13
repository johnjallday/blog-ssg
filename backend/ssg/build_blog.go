package ssg

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func RenderMarkdown(md []byte) (template.HTML, error) {
	var buf bytes.Buffer

	mdParser := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // GitHub-Flavored Markdown (includes checkboxes, tables, etc.)
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // adds id="..." to headings
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // turns \n into <br> in paragraphs
			html.WithXHTML(),
		),
	)

	if err := mdParser.Convert(md, &buf); err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil
}

// PostData represents a blog post with hashtags and a posted date.
type PostData struct {
	Title    string
	Content  template.HTML
	Hashtags []string
	Date     string // New field for the posted date
}

// ExtractHashtags parses hashtags from the Markdown content.
func ExtractHashtags(content string) []string {
	re := regexp.MustCompile(`#\w+`)
	matches := re.FindAllString(content, -1)
	hashtags := make([]string, len(matches))
	for i, tag := range matches {
		hashtags[i] = strings.TrimPrefix(tag, "#")
	}
	return hashtags
}

// ExtractMetadata parses the metadata (e.g., Title, Date) from the Markdown content.
func ExtractMetadata(mdContent string) (string, string) {
	var title, date string
	lines := strings.Split(mdContent, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Title:") {
			title = strings.TrimSpace(strings.TrimPrefix(line, "Title:"))
		} else if strings.HasPrefix(line, "Date:") {
			date = strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
		}
		if title != "" && date != "" {
			break
		}
	}
	return title, date
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

			mdContentBytes, err := os.ReadFile(mdPath)
			if err != nil {
				log.Printf("Error reading markdown file %s: %v", file.Name(), err)
				continue
			}
			mdContent := string(mdContentBytes)

			title, postDate := ExtractMetadata(mdContent)
			if title == "" {
				title = slug // Fallback to slug if Title is missing
			}

			// Remove the first two lines (Title and Date)
			mdContentLines := strings.Split(mdContent, "\n")
			if len(mdContentLines) > 2 {
				mdContent = strings.Join(mdContentLines[2:], "\n")
			}

			// Extract hashtags and remove them from the content
			hashtags := ExtractHashtags(mdContent)
			mdContentLines = strings.Split(mdContent, "\n")
			for i := len(mdContentLines) - 1; i >= 0; i-- {
				if strings.HasPrefix(mdContentLines[i], "#") {
					mdContentLines = mdContentLines[:i]
				} else {
					break
				}
			}
			mdContent = strings.Join(mdContentLines, "\n")

			htmlContent, err := RenderMarkdown([]byte(mdContent))
			if err != nil {
				log.Printf("Error converting markdown for %s: %v", file.Name(), err)
				continue
			}

			postData := PostData{
				Title:    title,
				Content:  htmlContent,
				Hashtags: hashtags, // Pass extracted hashtags
				Date:     postDate, // Use extracted post date
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
	Date  string // Add Date field to BlogPost
}

type BlogPageData struct {
	Posts []BlogPost
	Post  *PostData
}

func BuildBlog() {
	postsDir := filepath.Join("content", "posts")
	files, err := os.ReadDir(postsDir)
	if err != nil {
		log.Fatalf("Error reading posts directory: %v", err)
	}

	// Collect files with mod time
	entries := []struct {
		Info os.DirEntry
		Mod  int64
	}{}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			info, err := os.Stat(filepath.Join(postsDir, file.Name()))
			if err != nil {
				log.Printf("Error stating file %s: %v", file.Name(), err)
				continue
			}
			entries = append(entries, struct {
				Info os.DirEntry
				Mod  int64
			}{
				Info: file,
				Mod:  info.ModTime().Unix(),
			})
		}
	}

	// Sort newest first
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Mod > entries[j].Mod
	})

	var posts []BlogPost
	var latestPost *PostData

	for i, entry := range entries {
		slug := entry.Info.Name()[:len(entry.Info.Name())-3]

		mdPath := filepath.Join(postsDir, entry.Info.Name())
		mdContentBytes, err := os.ReadFile(mdPath)
		if err != nil {
			log.Printf("Error reading post content for %s: %v", slug, err)
			continue
		}
		mdContent := string(mdContentBytes)

		title, postDate := ExtractMetadata(mdContent)
		if title == "" {
			title = slug // Fallback to slug if Title is missing
		}

		// Remove the first two lines (Title and Date)
		mdContentLines := strings.Split(mdContent, "\n")
		if len(mdContentLines) > 2 {
			mdContent = strings.Join(mdContentLines[2:], "\n")
		}

		// Extract hashtags and remove them from the content
		hashtags := ExtractHashtags(mdContent)
		mdContentLines = strings.Split(mdContent, "\n")
		for i := len(mdContentLines) - 1; i >= 0; i-- {
			if strings.HasPrefix(mdContentLines[i], "#") {
				mdContentLines = mdContentLines[:i]
			} else {
				break
			}
		}
		mdContent = strings.Join(mdContentLines, "\n")

		posts = append(posts, BlogPost{
			Slug:  slug,
			Title: title,
			Date:  postDate, // Pass extracted postDate
		})

		// Render the latest post content for the main page
		if i == 0 {
			htmlContent, err := RenderMarkdown([]byte(mdContent))
			if err != nil {
				log.Printf("Error rendering latest post: %v", err)
				break
			}

			latestPost = &PostData{
				Title:    title,
				Content:  htmlContent,
				Hashtags: hashtags, // Pass extracted hashtags
				Date:     postDate, // Use extracted post date
			}
		}
	}

	// Load templates
	tmpl, err := template.ParseFiles(
		filepath.Join("backend", "ssg", "templates", "blog.html"),
		filepath.Join("backend", "ssg", "templates", "nav.html"),
		filepath.Join("backend", "ssg", "templates", "contact.html"),
	)
	if err != nil {
		log.Fatalf("Error parsing blog templates: %v", err)
	}

	// Create output file
	out, err := os.Create(filepath.Join("public", "blog.html"))
	if err != nil {
		log.Fatalf("Error creating blog.html: %v", err)
	}
	defer out.Close()

	pageData := BlogPageData{
		Posts: posts,
		Post:  latestPost,
	}

	if err := tmpl.Execute(out, pageData); err != nil {
		log.Fatalf("Error executing blog template: %v", err)
	}

	log.Println("public/blog.html built successfully!")
}
