#!/bin/bash

# Define the blacklist of directories (relative to the repository root)
BLACKLIST_DIRS=(".github" "docs" "postman" "path/another_directory")

# Define the blacklist of files (relative to the repository root)
BLACKLIST_FILES=(".dockerignore" "CODE_OF_CONDUCT.md" "README.md" "SECURITY.md" "CONTRIBUTING.md" "LICENSE" "show-content.bash" "go.mod" "go.sum" ".gitignore")

# Function to check if a file is located in one of the blacklisted directories
is_in_blacklisted_dir() {
    local file="$1"
    for dir in "${BLACKLIST_DIRS[@]}"; do
        # If the file path starts with the directory followed by a slash, it is considered inside the directory
        if [[ "$file" == "$dir/"* ]]; then
            return 0
        fi
    done
    return 1
}

# Function to check if a file is in the blacklist of files
is_blacklisted_file() {
    local file="$1"
    for bl_file in "${BLACKLIST_FILES[@]}"; do
        if [[ "$file" == "$bl_file" ]]; then
            return 0
        fi
    done
    return 1
}

# Check that the current directory is a Git repository
if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "This script must be run inside a Git repository."
    exit 1
fi

# Get the list of files not ignored by .gitignore
# We use -z to correctly handle names with spaces and special characters
git ls-files --cached --others --exclude-standard -z | while IFS= read -r -d '' file; do
    # Check if the file is in either blacklist
    if is_in_blacklisted_dir "$file" || is_blacklisted_file "$file"; then
        continue  # Skip the file if it is blacklisted
    fi

    echo "File: $file"
    echo "Content:"
    echo "---"
    cat "$file"
    echo "---"
    echo -e "\n"
done
