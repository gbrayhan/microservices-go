#!/usr/bin/env bash
exec > project-content.txt 2>&1

BLACKLIST_DIRS=(
  ".github"
  "docs"
  "postman"
  "path/another_directory"
)

BLACKLIST_FILES=(
  ".dockerignore"
  "CODE_OF_CONDUCT.md"
  "README.md"
  "SECURITY.md"
  "CONTRIBUTING.md"
  "LICENSE"
  "show-content.bash"
  "go.mod"
  "go.sum"
  ".gitignore"
)

BLACKLIST_EXTENSIONS=(
  "md"
  "txt"
  "pdf"
  "log"
)

is_in_blacklisted_dir() {
  local f="$1"
  for d in "${BLACKLIST_DIRS[@]}"; do
    [[ "$f" == "$d/"* ]] && return 0
  done
  return 1
}

is_blacklisted_file() {
  local f="$1"
  for b in "${BLACKLIST_FILES[@]}"; do
    [[ "$f" == "$b" ]] && return 0
  done
  return 1
}

is_blacklisted_extension() {
  local f="$1"
  [[ "$f" != *.* ]] && return 1
  local ext="${f##*.}"
  ext="$(printf '%s' "$ext" | tr '[:upper:]' '[:lower:]')"
  for b in "${BLACKLIST_EXTENSIONS[@]}"; do
    [[ "$ext" == "$b" ]] && return 0
  done
  return 1
}

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "This script must be run inside a Git repository."
  exit 1
fi

git ls-files --cached --others --exclude-standard -z |
while IFS= read -r -d '' file; do
  if is_in_blacklisted_dir "$file" || is_blacklisted_file "$file" || is_blacklisted_extension "$file"; then
    continue
  fi

  echo "File: $file"
  echo "Content:"
  echo "---"
  cat "$file"
  echo "---"
  echo
done
