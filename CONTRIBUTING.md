# Contributing to Nexlayer CLI

Thank you for considering contributing to Nexlayer CLI! Your contributions help developers build, deploy, and manage AI-powered applications effortlessly.

## Code Quality Standards

### 1. License Headers
All source files must include the MIT license header:
```go
// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
```

### 2. Package Documentation
Each package must include comprehensive documentation:
- Package purpose and overview
- Key types and interfaces
- Usage examples
- Design decisions and rationale

### 3. Code Style
- All Go files must be formatted with `gofmt`
- Follow Go best practices and idioms
- Use consistent naming conventions
- Group related functions together
- Implement proper interfaces for testability
- Use dependency injection where appropriate

### 4. Error Handling
- Use context support for all API calls
- Implement proper error wrapping
- Include detailed error messages
- Follow proper error logging practices

### 5. Documentation Style
Example of good function documentation:
```go
// GenerateTemplate creates a nexlayer.yaml template for the given project.
// It analyzes the current directory structure to detect components and their types,
// then generates a template with appropriate configuration.
//
// Parameters:
//   - projectName: Name of the project, used as both template name and deployment name
//   - detector: ComponentDetector instance used to identify component types
//
// Returns:
//   - string: The generated YAML template
//   - error: Any error encountered during template generation
//
// Example:
//
//	detector := components.NewComponentDetector()
//	template, err := GenerateTemplate("my-app", detector)
//	if err != nil {
//	    return fmt.Errorf("failed to generate template: %w", err)
//	}
```

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/Nexlayer/nexlayer-cli.git
cd nexlayer-cli
```

2. Install dependencies:
```bash
make setup
```

3. Build the CLI:
```bash
make build
```

4. Run tests:
```bash
make test
make lint
```

## Pull Request Process

1. **Before Starting**
   - For bug fixes: Open an issue with a minimal reproducible example
   - For new features: Start a discussion to align on implementation
   - For refactoring: Explain the benefits to performance/maintainability

2. **PR Requirements**
   - One problem per PR (atomic changes)
   - Include tests or reproduction steps
   - Pass all tests and linting
   - Include documentation updates
   - Follow conventional commit format:
     ```
     feat(cli): add AI-powered YAML validation
     fix(deploy): resolve API timeout issue
     refactor(parser): optimize template loading
     ```

3. **Review Process**
   - Mark work-in-progress PRs as Draft
   - Request review when ready
   - Address feedback promptly
   - Maintainer will merge after approval

## Testing

1. **Unit Tests**
   - Required for all new functionality
   - Must cover edge cases
   - Should be readable and maintainable

2. **Integration Tests**
   - Required for key functionality
   - Should cover real-world use cases
   - Must include CLI interaction tests

## Documentation

1. **README Updates**
   - Keep installation instructions current
   - Update feature documentation
   - Include usage examples

2. **Code Examples**
   - Must be tested and working
   - Should be clear and concise
   - Include comments for complex operations

## Getting Help

- Open an issue for bugs
- Start a discussion for feature requests
- Join our community channels for general help

## License

By contributing to Nexlayer CLI, you agree that your contributions will be licensed under the MIT License.
