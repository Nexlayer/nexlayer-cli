#!/bin/bash

# Function to fix formatting in a file
fix_file() {
    local file="$1"
    echo "Processing $file..."
    
    # Create a temporary file
    tmp_file="${file}.tmp"
    
    # Fix Printf and Fprintf formatting with newlines
    sed -E '
        # Fix Printf/Fprintf with newlines
        /fmt\.(Printf|Fprintf)\([^)]*$/ {
            # Join next line
            N
            # Replace newline with comma
            s/("[^"]*)"[\n\r]*"([^"]*)"/"\1\2"/g
            # Add missing comma if needed
            s/("[^"]*")\s*,?\s*([^,)]+)$/\1, \2/g
        }
        # Fix other common patterns
        s/fmt\.Println\("([^"]*)"[\n\r]*"([^"]*)"\)/fmt.Println("\1\2")/g
        s/fmt\.Printf\("([^"]*)"[\n\r]*"([^"]*)"\)/fmt.Printf("\1\2")/g
    ' "$file" > "$tmp_file"
    
    # Only replace if the file was changed
    if ! cmp -s "$file" "$tmp_file"; then
        mv "$tmp_file" "$file"
        echo "Fixed formatting in $file"
    else
        rm -f "$tmp_file"
        echo "No changes needed in $file"
    fi
    
    # Run gofmt on the file
    gofmt -w "$file"
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
