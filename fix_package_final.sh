#!/bin/bash

# Function to get package name from file path
get_package_name() {
    local file="$1"
    if [[ "$file" == "./main.go" ]]; then
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
    
    # Remove everything up to and including the first package declaration
    sed -i.bak '1,/^package.*$/d' "$file"
    
    # Add the correct package declaration at the start
    echo -e "package $pkg\n$(cat $file)" > "$file"
    
    # Remove backup file
    rm -f "${file}.bak"
    
    # Format with gofmt -s
    if gofmt -s -w "$file" >/dev/null 2>&1; then
        echo "Successfully formatted $file"
    else
        echo "Failed to format $file"
    fi
done
