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
    
    # Create a new file with UTF-8 encoding
    echo -n "package $pkg" > "${file}.new"
    echo "" >> "${file}.new"
    echo "" >> "${file}.new"
    
    # Copy the rest of the file, starting from the imports
    tail -n +2 "$file" | grep -v "^package" >> "${file}.new"
    
    # Format with gofmt -s
    if gofmt -s -w "${file}.new"; then
        mv "${file}.new" "$file"
        echo "Successfully formatted $file"
    else
        rm "${file}.new"
        echo "Failed to format $file"
    fi
done
