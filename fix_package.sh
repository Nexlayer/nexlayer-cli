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
    
    # Remove any BOM and -e flags, then add package declaration
    sed -i '' '1s/^.*$/package '"$pkg"'/' "$file"
    
    # Format with gofmt -s
    if gofmt -s -w "$file" >/dev/null 2>&1; then
        echo "Successfully formatted $file"
    else
        echo "Failed to format $file"
    fi
done
