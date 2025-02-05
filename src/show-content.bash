#!/bin/bash

# Check if the current directory is a Git repository
if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "This script must be run inside a Git repository."
    exit 1
fi

# Get the list of files not ignored by .gitignore
# We use -z to correctly handle names with spaces and special characters
git ls-files --cached --others --exclude-standard -z | while IFS= read -r -d '' file; do
    echo "File: $file"
    echo "Content:"
    echo "----------------------------------------"
    cat "$file"
    echo "----------------------------------------"
    echo -e "\n"
done
