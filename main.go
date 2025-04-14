package main

import (
	"blog/backend/ssg"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Build the entire static site.
	fmt.Println("Building the static site...")
	ssg.BuildAll()
	fmt.Println("Build completed.")

	// Commit and push the changes to Git.
	if err := GitCommitAndPush(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// GitCommitAndPush automates git add, commit, and push
func GitCommitAndPush() error {
	reader := bufio.NewReader(os.Stdin)

	// Ask for a commit message
	fmt.Print("Enter commit message: ")
	commitMessage, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read commit message: %w", err)
	}
	commitMessage = strings.TrimSpace(commitMessage)

	// Run git add .
	if err := runCommand("git", "add", "."); err != nil {
		return fmt.Errorf("failed to run 'git add .': %w", err)
	}

	// Run git commit -m "<message>"
	if err := runCommand("git", "commit", "-m", commitMessage); err != nil {
		return fmt.Errorf("failed to run 'git commit': %w", err)
	}

	// Run git push
	if err := runCommand("git", "push"); err != nil {
		return fmt.Errorf("failed to run 'git push': %w", err)
	}

	fmt.Println("Changes have been pushed successfully!")
	return nil
}

// Helper function to run a command
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
