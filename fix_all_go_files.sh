#!/bin/bash

# Function to get package name from file path
get_package_name() {
    local file="$1"
    local dir=$(dirname "$file")
    local pkg=$(basename "$dir")
    echo "$pkg"
}

# Process each .go file
find . -name "*.go" -not -path "./vendor/*" | while read -r file; do
    echo "Processing $file..."
    pkg=$(get_package_name "$file")
    
    # Create temporary file with proper package declaration
    awk -v pkg="$pkg" '
        NR==1 { print "package " pkg }  # First line is package declaration
        NR>1  { print }                 # Print rest of the file as is
    ' "$file" > "${file}.tmp"
    
    # Format with gofmt -s
    if gofmt -s -w "${file}.tmp"; then
        mv "${file}.tmp" "$file"
        echo "Successfully formatted $file"
    else
        rm "${file}.tmp"
        echo "Failed to format $file"
    fi
done
