#!/bin/bash

# Function to get package name from file path
get_package_name() {
    local file="$1"
    local dir=$(dirname "$file")
    local pkg=$(basename "$dir")
    echo "$pkg"
}

# Function to fix string literals and emojis
fix_file() {
    local file="$1"
    local pkg="$2"
    
    # Create temporary file
    local tmp="${file}.tmp"
    
    # Start with package declaration
    echo "package ${pkg}" > "$tmp"
    echo "" >> "$tmp"
    
    # Process the file line by line
    awk '
        # Skip package declaration
        NR > 1 {
            # Fix unterminated string literals
            gsub(/"""/, "\"")
            gsub(/""/, "\"")
            gsub(/[[:space:]]+$/, "")
            
            # Remove emojis and replace with text equivalents
            gsub(/🔍/, "")  # Remove search emoji
            gsub(/✓/, "")   # Remove checkmark
            gsub(/📋/, "")  # Remove clipboard
            gsub(/📝/, "")  # Remove memo
            gsub(/🎉/, "")  # Remove party
            gsub(/📚/, "")  # Remove books
            gsub(/💡/, "")  # Remove light bulb
            
            print
        }
    ' "$file" >> "$tmp"
    
    # Try to format with gofmt
    if gofmt -s -w "$tmp" >/dev/null 2>&1; then
        mv "$tmp" "$file"
        echo "Successfully formatted $file"
        return 0
    else
        rm -f "$tmp"
        echo "Failed to format $file"
        return 1
    fi
}

# Process problematic files
files=(
    "pkg/tui/model.go"
    "pkg/tui/spinner.go"
    "pkg/tui/editor.go"
    "pkg/api/client.go"
    "pkg/commands/domain.go"
    "pkg/commands/deploy.go"
    "pkg/commands/ai_suggest.go"
    "pkg/commands/list.go"
    "pkg/commands/plugin.go"
    "pkg/commands/status.go"
    "pkg/commands/scale.go"
    "pkg/commands/info.go"
    "pkg/commands/init.go"
    "pkg/errors/errors.go"
)

for file in "${files[@]}"; do
    echo "Processing $file..."
    pkg=$(get_package_name "$file")
    fix_file "$file" "$pkg"
done
