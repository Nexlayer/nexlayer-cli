#!/bin/bash

# Function to fix formatting in a file
fix_file() {
    local file="$1"
    echo "Processing $file..."
    
    # Create a temporary file
    tmp_file="${file}.tmp"
    
    # Fix Printf formatting with newlines
    awk '
        # Inside Printf with newline
        /fmt\.Printf\([^)]*$/ {
            # Replace newline with comma and string concatenation
            gsub(/"\n/, "\", ")
            print
            next
        }
        # Normal line
        { print }
    ' "$file" > "$tmp_file"
    
    # Only replace if the file was changed
    if ! cmp -s "$file" "$tmp_file"; then
        mv "$tmp_file" "$file"
        echo "Fixed formatting in $file"
    else
        rm -f "$tmp_file"
        echo "No changes needed in $file"
    fi
}

# Process all Go files with formatting errors
files=(
    "examples/plugins/hello/main.go"
    "examples/plugins/lint/main.go"
    "pkg/cache/cache.go"
    "pkg/cache/commands.go"
    "pkg/commands/ai_suggest.go"
    "pkg/commands/deploy.go"
    "pkg/commands/domain.go"
    "pkg/commands/info.go"
    "pkg/commands/init.go"
    "pkg/commands/plugin.go"
    "pkg/commands/scale.go"
    "pkg/commands/status.go"
    "pkg/errors/errors.go"
    "pkg/tui/spinner.go"
)

for file in "${files[@]}"; do
    fix_file "$file"
done
