#!/bin/bash

# Check if gh CLI is installed
if ! command -v gh &>/dev/null; then
  echo "GitHub CLI (gh) not found. Please install it first."
  exit 1
fi

# Check if user is logged in
gh auth status || gh auth login

# Initialize git repository if not already initialized
if [ ! -d .git ]; then
  git init
  echo "Git repository initialized"
fi

# Add files to git
git add .

# Make initial commit
git commit -m "Initial commit"

# Create GitHub repository and push
read -p "Enter repository name (default: blog-ssg): " repo_name
repo_name=${repo_name:-blog-ssg}

description="A static site generator for Allday Music Blog built with Go"

if [ -n "$description" ]; then
  gh repo create "$repo_name" --public --description "$description" --source=. --push
else
  gh repo create "$repo_name" --public --source=. --push
fi

echo "Repository successfully created and code pushed to GitHub!"
