#!/bin/bash

# Function to fix string literals in a file
fix_file() {
    local file="$1"
    echo "Processing $file..."
    
    # Create a temporary file
    tmp_file="${file}.tmp"
    
    # Fix string literals with newlines
    awk '
        # Inside string literal
        /^[^"]*"[^"]*$/ {
            # Replace raw newlines with \n
            gsub(/\n/, "\\n")
            # Add closing quote
            print $0 "\""
            next
        }
        # Normal line
        { print }
    ' "$file" > "$tmp_file"
    
    # Only replace if the file was changed
    if ! cmp -s "$file" "$tmp_file"; then
        mv "$tmp_file" "$file"
        echo "Fixed string literals in $file"
    else
        rm -f "$tmp_file"
        echo "No changes needed in $file"
    fi
}

# Process all Go files with string literal errors
files=(
    "examples/plugins/hello/main.go"
    "examples/plugins/lint/main.go"
    "pkg/api/client.go"
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
    "pkg/tui/editor.go"
    "pkg/tui/model.go"
    "pkg/tui/spinner.go"
)

for file in "${files[@]}"; do
    fix_file "$file"
done
