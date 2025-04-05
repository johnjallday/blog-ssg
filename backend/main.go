package main

import (
	"blog/backend/server"
	"blog/backend/ssg"
)

func main() {
	// Build the entire static site.
	ssg.BuildAll()

	// Start serving the generated site.
	server.RunServer()
}
