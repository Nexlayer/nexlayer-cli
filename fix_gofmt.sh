#!/bin/bash

# List of files to format
files=(
    "pkg/ai/factory.go"
    "pkg/commands/ai_suggest.go"
    "pkg/commands/init.go"
    "pkg/plugin/manager.go"
    "pkg/tui/model.go"
    "pkg/cache/commands.go"
    "pkg/errors/errors.go"
    "pkg/commands/domain.go"
    "pkg/commands/list.go"
    "pkg/commands/scale.go"
    "pkg/commands/status.go"
    "pkg/tui/wizard.go"
    "cmd/root.go"
    "pkg/api/client.go"
    "pkg/api/client_test.go"
    "pkg/ai/openai.go"
    "pkg/commands/deploy.go"
    "pkg/ai/ai.go"
    "pkg/ai/claude.go"
    "pkg/api/types.go"
    "pkg/commands/info.go"
    "pkg/plugin/loader.go"
)

for file in "${files[@]}"; do
    echo "Processing $file..."
    
    # Create a backup
    cp "$file" "$file.bak"
    
    # Remove any potential BOM or special characters and ensure Unix line endings
    cat "$file" | tr -d '\r' | tr -cd '[:print:]\n' > "$file.clean"
    
    # Format with gofmt -s
    gofmt -s -w "$file.clean"
    
    # Check if formatting succeeded
    if [ $? -eq 0 ]; then
        # Add package declaration at the top if missing
        package_name=$(dirname "$file" | sed 's/.*\///' | tr '/' '_')
        if ! grep -q "^package" "$file.clean"; then
            echo -e "package ${package_name}\n$(cat "$file.clean")" > "$file.clean"
        fi
        
        # Move the formatted file back
        mv "$file.clean" "$file"
        echo "Successfully formatted $file"
    else
        # Restore from backup if formatting failed
        mv "$file.bak" "$file"
        echo "Failed to format $file"
    fi
    
    # Clean up
    rm -f "$file.bak" "$file.clean"
done
