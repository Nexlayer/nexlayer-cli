# Detector Migration Guide

## Overview

We're transitioning from individual technology detectors to a unified stack detection approach. This change brings several benefits:

- **Improved Performance**: Single-pass scanning reduces I/O operations
- **Better Accuracy**: Holistic stack analysis with confidence scoring
- **Simplified Integration**: One detector for all technology stacks
- **Extensible Design**: Easy to add new stack patterns
- **Reduced Maintenance**: Centralized pattern definitions

## Timeline

- **Current version (2.x)**: 
  - Deprecated detectors are marked but still functional
  - Warning logs will appear in development
  - New unified detector available for testing

- **Next minor release (2.x+1)**:
  - Deprecated detectors will emit warning logs in production
  - Documentation will focus on unified detector
  - Performance improvements for unified detector

- **Next major release (3.0)**:
  - Deprecated detectors will be removed
  - Only unified stack detection available
  - Breaking changes in detector interfaces

## Deprecated Detectors

### Combined Stack Detectors (High Priority Migration)
- `NextjsSupabaseLangchainDetector`
- `NextjsSupabaseOpenaiDetector`

### Framework Detectors
- `NextjsDetector`
- `BunDetector`
- `SvelteDetector`
- `VueDetector`

### Database Detectors
- `PostgresqlDetector`
- `PgvectorDetector`
- `SupabaseDetector`

### AI/ML Detectors
- `LangchainDetector`
- `OpenaiDetector`
- `GeminiDetector`
- `AIDetector`

### UI/Payment Detectors
- `TailwindDetector`
- `StripeDetector`

## Migration Steps

### Using the DetectorRegistry (Recommended)

If you're using the standard `DetectorRegistry`, no immediate changes are needed:

```go
// Before and after migration - your code stays the same
registry := detection.NewDetectorRegistry()
info, err := registry.Detect(dir)
```

The registry now includes the unified `StackDetector` which handles all technology detection.

### Direct Detector Usage

Replace individual detector instantiations with the unified `StackDetector`:

```go
// BEFORE - Multiple detectors
nextjsDetector := &detectors.NextjsDetector{}
supabaseDetector := &detectors.SupabaseDetector{}
langchainDetector := &detectors.LangchainDetector{}

// AFTER - Single unified detector
detector := detection.NewStackDetector()
info, err := detector.Detect(dir)
```

### Working with Detection Results

The unified detector provides richer metadata about detected stacks:

```go
// Get detection results
info, err := detector.Detect(dir)
if err != nil {
    log.Fatal(err)
}

// Access basic info
projectType := info.Type
confidence := info.Metadata["confidence"].(float64)

// Access stack components
if stack, ok := info.Metadata["stack_components"].(map[string]interface{}); ok {
    frontend := stack["frontend"].([]string)
    backend := stack["backend"].([]string)
    database := stack["database"].([]string)
    ai := stack["ai"].([]string)
}

// Check specific technologies
hasNextjs := info.Dependencies["nextjs"] != ""
hasSupabase := info.Dependencies["supabase"] != ""
hasLangchain := info.Dependencies["langchain"] != ""
```

### Custom Stack Patterns

You can extend the unified detector with your own stack patterns:

```go
customPatterns := map[string]detection.StackDefinition{
    "my-custom-stack": {
        Name: "Custom Stack",
        Components: detection.Components{
            Frontend: []string{"custom-ui"},
            Backend:  []string{"custom-api"},
        },
        MainPatterns: []detection.DetectionPattern{
            {
                Type:       detection.PatternDependency,
                Pattern:    "custom-package",
                Path:      "package.json",
                Confidence: 0.6,
            },
        },
    },
}

detector := detection.NewStackDetector()
detector.AddPatterns(customPatterns)
```

## Testing During Migration

1. **Parallel Testing**
   ```go
   // Run both detectors and compare results
   oldInfo, _ := oldDetector.Detect(dir)
   newInfo, _ := detection.NewStackDetector().Detect(dir)
   
   if oldInfo.Type != newInfo.Type {
       log.Printf("Detection mismatch: old=%s, new=%s", oldInfo.Type, newInfo.Type)
   }
   ```

2. **Confidence Validation**
   ```go
   // Check confidence levels
   if confidence, ok := newInfo.Metadata["confidence"].(float64); ok {
       if confidence < 0.8 {
           log.Printf("Low confidence detection: %f", confidence)
       }
   }
   ```

3. **Component Verification**
   ```go
   // Verify all required components are detected
   requiredTech := []string{"nextjs", "supabase", "langchain"}
   for _, tech := range requiredTech {
       if _, exists := newInfo.Dependencies[tech]; !exists {
           log.Printf("Missing required technology: %s", tech)
       }
   }
   ```

## Performance Considerations

The unified detector is optimized for performance:
- Parallel file scanning with goroutines
- File content caching
- Early exit on high-confidence matches
- Optimized pattern matching

## Best Practices

1. **Always check error returns**
   ```go
   info, err := detector.Detect(dir)
   if err != nil {
       // Handle error appropriately
       return err
   }
   ```

2. **Validate confidence scores**
   ```go
   if confidence, ok := info.Metadata["confidence"].(float64); ok {
       if confidence < 0.5 {
           // Consider manual verification
       }
   }
   ```

3. **Handle multiple possible stacks**
   ```go
   if stacks, ok := info.Metadata["detected_stacks"].([]string); ok {
       if len(stacks) > 1 {
           // Multiple stacks detected, may need manual selection
       }
   }
   ```

## Questions and Support

- **GitHub Issues**: Open an issue for bugs or feature requests
- **Documentation**: Full API reference available at [docs.nexlayer.dev](https://docs.nexlayer.dev)
- **Community**: Join our Discord for real-time support
- **Examples**: Sample code available in [examples/](https://github.com/Nexlayer/nexlayer-cli/tree/main/examples)

## Future Roadmap

- **3.0 Release**: Complete removal of deprecated detectors
- **3.1 Release**: Enhanced pattern matching with ML support
- **3.2 Release**: Real-time pattern updates from central registry 