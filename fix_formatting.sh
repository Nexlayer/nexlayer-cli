#!/bin/bash

# List of files to format
files=(
    "pkg/plugin/manager.go"
    "cmd/root.go"
    "pkg/api/client.go"
    "pkg/tui/wizard.go"
    "pkg/commands/scale.go"
    "pkg/api/types.go"
    "pkg/tui/model.go"
    "pkg/commands/ai_suggest.go"
    "pkg/commands/domain.go"
    "pkg/commands/info.go"
    "pkg/cache/commands.go"
    "pkg/commands/deploy.go"
    "pkg/commands/init.go"
    "pkg/plugin/loader.go"
    "pkg/errors/errors.go"
    "pkg/ai/ai.go"
    "pkg/ai/factory.go"
    "pkg/ai/openai.go"
    "pkg/commands/list.go"
    "pkg/api/client_test.go"
    "pkg/ai/claude.go"
    "pkg/commands/status.go"
)

for file in "${files[@]}"; do
    echo "Processing $file..."
    
    # Create a temporary file
    temp_file="${file}.tmp"
    
    # Copy original file to temp file, stripping any BOM or special characters
    tr -cd '[:print:]\n' < "$file" > "$temp_file"
    
    # Format with gofmt -s
    gofmt -s "$temp_file" > "${file}.formatted"
    
    # Check if formatting was successful
    if [ $? -eq 0 ]; then
        # Replace original with formatted version
        mv "${file}.formatted" "$file"
        rm "$temp_file"
        echo "Successfully formatted $file"
    else
        echo "Failed to format $file"
        rm "$temp_file" "${file}.formatted"
    fi
done
