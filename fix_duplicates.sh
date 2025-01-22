#!/bin/bash

# Function to get package name from file path
get_package_name() {
    local file="$1"
    if [[ "$file" == "main.go" ]]; then
        echo "main"
    else
        local dir=$(dirname "$file")
        local pkg=$(basename "$dir")
        echo "$pkg"
    fi
}

# Process each .go file
find . -name "*.go" -not -path "./vendor/*" | while read -r file; do
    echo "Processing $file..."
    pkg=$(get_package_name "$file")
    
    # Create a temporary file
    tmp_file="${file}.tmp"
    
    # Remove duplicate package declarations and format
    awk -v pkg="$pkg" '
        BEGIN { printed_package = 0 }
        /^package/ { 
            if (!printed_package) {
                print "package " pkg
                printed_package = 1
            }
            next
        }
        /^\/\/ Formatted with gofmt/ { next }  # Skip formatting comments
        { print }
    ' "$file" > "$tmp_file"
    
    # Format with gofmt -s
    if gofmt -s -w "$tmp_file" >/dev/null 2>&1; then
        mv "$tmp_file" "$file"
        echo "Successfully formatted $file"
    else
        rm -f "$tmp_file"
        echo "Failed to format $file"
    fi
done
