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
    
    # Create a new file with UTF-8 encoding
    {
        # Write package declaration
        echo "package $pkg"
        echo ""
        
        # Extract imports and rest of the file
        awk '
            BEGIN { in_imports = 0; printed_imports = 0 }
            /^import/ { if (!printed_imports) { print "import ("; in_imports = 1; printed_imports = 1; next } }
            /^)/ && in_imports { print ")"; print ""; in_imports = 0; next }
            in_imports { if ($0 ~ /^[[:space:]]*"/) print "\t" $0; next }
            !/^package/ && !/^$/ { print }
        ' "$file" | sed 's/[[:space:]]*$//'
    } > "${file}.new"
    
    # Format with gofmt -s
    if gofmt -s -w "${file}.new" >/dev/null 2>&1; then
        # Only replace if gofmt succeeded
        mv "${file}.new" "$file"
        echo "Successfully formatted $file"
    else
        rm -f "${file}.new"
        echo "Failed to format $file"
    fi
done
