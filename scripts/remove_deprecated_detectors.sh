#!/bin/bash
# Script to remove deprecated detectors in a future major version release
# This script should be run during the preparation of the next major version

echo "Removing deprecated detectors..."

# List of deprecated detector files to remove
DEPRECATED_DETECTORS=(
  "pkg/detection/detectors/nextjs_supabase_langchain_detector.go"
  "pkg/detection/detectors/nextjs_supabase_openai_detector.go"
  "pkg/detection/detectors/pgvector_detector.go"
  "pkg/detection/detectors/langchain_detector.go"
  "pkg/detection/detectors/openai_detector.go"
  "pkg/detection/detectors/stripe_detector.go"
  "pkg/detection/detectors/gemini_detector.go"
  "pkg/detection/detectors/tailwind_detector.go"
)

# Remove each deprecated detector
for detector in "${DEPRECATED_DETECTORS[@]}"; do
  if [ -f "$detector" ]; then
    echo "Removing $detector"
    rm "$detector"
  else
    echo "Warning: $detector not found"
  fi
done

echo "Deprecated detectors removed. Please build and test the project to ensure everything works correctly."
echo "Don't forget to update the MIGRATION_GUIDE.md to reflect that these detectors have been removed." 