# Contributing to Nexlayer CLI

> ⚠️ **Pre-Release Notice**: Nexlayer CLI is currently in early development and the repository is private. Contributions are limited to invited collaborators only. We plan to open the repository for public contributions with our beta v1 release in Q2 2025.

👋 First off, thanks for taking the time to contribute! We're excited to welcome you to the Nexlayer community.

## 🔒 Access & Permissions

During the pre-beta phase:
- Repository access is invite-only
- All contributors must sign an NDA
- Code sharing and redistribution is not permitted
- Direct commits to main branch are disabled
- All changes must go through PR review

## 🚀 Quick Start

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/nexlayer-cli.git
   cd nexlayer-cli
   ```

2. **Install Dependencies**
   ```bash
   # Install Go (1.23.4 or higher required)
   brew install go # macOS
   # or visit https://golang.org/dl/ for other platforms

   # Install development tools
   make setup
   ```

3. **Run Tests**
   ```bash
   make test        # Run unit tests
   make test-short  # Run quick tests
   make coverage    # Run tests with coverage report
   ```

## 💻 Development Workflow

### Setting Up Your Environment

1. **Install Required Tools**
   ```bash
   # Install golangci-lint
   brew install golangci-lint # macOS
   # or
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

   # Install other development dependencies
   make setup
   ```

2. **Configure Git Hooks**
   ```bash
   git config core.hooksPath .githooks
   chmod +x .githooks/*
   ```

### Making Changes

1. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Development Commands**
   ```bash
   make build-dev  # Build for development
   make run        # Run the CLI
   make watch      # Watch for changes and rebuild
   ```

3. **Code Quality**
   ```bash
   make lint       # Run linters
   make fmt        # Format code
   make vet        # Run go vet
   ```

### Before Submitting

1. **Run All Checks**
   ```bash
   make ci         # Runs all CI tasks locally
   ```

2. **Update Documentation**
   - Add inline comments for complex logic
   - Update README.md if adding new features
   - Update API documentation if changing interfaces

## 📝 Coding Standards

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting (automatically run with `make fmt`)
- Keep functions focused and small (max 100 lines recommended)
- Add tests for new functionality

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only changes
- `style`: Changes that don't affect code meaning
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests
- `chore`: Changes to the build process or auxiliary tools

Example:
```
feat(deploy): add support for custom domains

- Add domain validation
- Add DNS verification
- Update documentation

Closes #123
```

### Pull Request Process

1. **Create Pull Request**
   - Use the PR template
   - Link related issues
   - Add labels as appropriate

2. **PR Description**
   - Describe the changes
   - Add testing instructions
   - List breaking changes
   - Add screenshots for UI changes

3. **Code Review**
   - Address review comments
   - Keep the PR focused
   - Squash commits if requested

## 🧪 Testing

### Test Structure

- Place tests in `*_test.go` files
- Use table-driven tests when possible
- Mock external dependencies
- Aim for >80% coverage on new code

### Running Tests

```bash
# Run all tests
make test

# Run specific tests
go test ./pkg/... -run TestSpecificFunction

# Run tests with coverage
make coverage
```

## 📚 Documentation

### Code Documentation

- Add godoc comments for exported functions
- Include examples for complex functionality
- Document error conditions and edge cases

### API Documentation

- Update OpenAPI specs when changing APIs
- Include request/response examples
- Document breaking changes

## 🔧 Tools and Configuration

### Linting

We use golangci-lint with strict settings. Configuration is in `.golangci.yml`.

Key linters enabled:
- `gofmt`: Code formatting
- `govet`: Go best practices
- `staticcheck`: Static analysis
- `gosec`: Security checks
- `errcheck`: Error handling
- `gosimple`: Code simplification

### IDE Setup

VSCode settings:
```json
{
  "go.lintTool": "golangci-lint",
  "go.formatTool": "gofmt",
  "go.useLanguageServer": true
}
```

## 🚨 Reporting Issues

- Use the issue templates
- Include reproduction steps
- Attach logs if relevant
- Tag appropriately

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.
