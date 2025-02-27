# Detector Migration Guide

## Overview

We're transitioning from individual technology detectors to a unified stack detection approach. This document provides guidance on how to migrate from deprecated detectors to the new unified `StackDetector`.

## Timeline

- **Current version**: Deprecated detectors are marked but still functional
- **Next minor release**: Deprecated detectors will emit warning logs when used
- **Next major release**: Deprecated detectors will be removed

## Deprecated Detectors

The following detectors are deprecated and will be removed in a future version:

### Combined Stack Detectors (High Priority)
- `NextjsSupabaseLangchainDetector`
- `NextjsSupabaseOpenaiDetector`

### Specialized Component Detectors
- `PgvectorDetector`
- `LangchainDetector`
- `OpenaiDetector`
- `StripeDetector`
- `GeminiDetector`
- `TailwindDetector`

## Migration Steps

### If you're using the DetectorRegistry

If you're using the standard `DetectorRegistry` through `detection.NewDetectorRegistry()`, no changes are needed! The registry now includes the unified `StackDetector` which will automatically handle technology stack detection.

### If you're directly instantiating detectors

Replace:
```go
detector := &detectors.NextjsSupabaseLangchainDetector{}
info, err := detector.Detect(dir)
```

With:
```go
detector := detection.NewStackDetector()
info, err := detector.Detect(dir)
```

### Working with detection results

The detection results structure remains the same (`*types.ProjectInfo`), but the unified detector provides more comprehensive metadata about technology stacks. Check the `Dependencies` map in the result for detailed component information:

```go
// Access stack components
frontend := info.Dependencies["frontend"]
backend := info.Dependencies["backend"]
database := info.Dependencies["database"]
aiComponents := info.Dependencies["ai"]
```

## Testing

We recommend testing your application with both the old and new detection approach during the transition period to ensure compatibility.

## Questions and Support

If you encounter any issues during migration, please open an issue in our GitHub repository or contact our support team. 