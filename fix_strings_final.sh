#!/bin/bash

# Process each .go file
find . -name "*.go" -not -path "./vendor/*" | while read -r file; do
    echo "Processing $file..."
    
    # Fix string literals with newlines and BOM
    sed -i.bak -E '
        # Fix string literals with newlines
        s/"\([^"]*\)\n"/"\1\\n"/g
        
        # Fix package declarations (remove BOM and duplicates)
        1 {
            s/^[^p]*package/package/
            s/^package[[:space:]]+([[:alnum:]]+).*$/package \1/
        }
        
        # Remove duplicate package declarations and formatting comments
        /^\/\/ Formatted with gofmt/d
        2,$ {
            /^package[[:space:]]/d
        }
    ' "$file"
    
    # Remove backup file
    rm -f "${file}.bak"
    
    # Format with gofmt -s
    if gofmt -s -w "$file" >/dev/null 2>&1; then
        echo "Successfully formatted $file"
    else
        echo "Failed to format $file"
    fi
done
