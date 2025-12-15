# Contributing to SQLC-Wizard

Thank you for your interest in contributing to SQLC-Wizard! This document provides guidelines and information for contributors.

## 🚀 Quick Start

1. Fork the repository
2. Clone your fork locally
3. Create a feature branch
4. Make your changes
5. Test your changes thoroughly
6. Submit a pull request

## 📋 Prerequisites

- Go 1.25 or later
- Git
- Make or Just (for build automation)
- Docker (optional, for testing)

## 🛠️ Development Setup

```bash
# Clone the repository
git clone https://github.com/LarsArtmann/SQLC-Wizzard.git
cd SQLC-Wizzard

# Install dependencies
just deps

# Run all tests
just test

# Build the binary
just build

# Run development workflow
just dev
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests with coverage
just test

# Run tests for a specific package
go test ./internal/wizard -v

# Run tests with race detection
go test -race ./...

# Run integration tests
go test -tags=integration ./...
```

### Test Coverage

We aim for >80% test coverage. Check your coverage:

```bash
# Run tests with coverage
go test -coverprofile=coverage.txt ./...

# View coverage report
go tool cover -html=coverage.txt
```

### Writing Tests

- Use BDD style with Ginkgo/Gomega
- Write both unit and integration tests
- Test edge cases and error conditions
- Keep tests focused and maintainable

```go
var _ = Describe("Wizard", func() {
    It("should create wizard with default values", func() {
        wiz := wizard.NewWizard()
        Expect(wiz).NotTo(BeNil())
        Expect(wiz.GetResult().GenerateQueries).To(BeTrue())
    })
})
```

## 📝 Code Style

### Go Conventions

- Follow Go standard formatting (`go fmt`)
- Use `golangci-lint` for linting
- Keep lines under 120 characters
- Use meaningful variable names
- Add comprehensive comments for public APIs

### Architecture

- Follow Domain-Driven Design principles
- Maintain clean architecture layers
- Use dependency injection
- Write clear, focused functions
- Prefer composition over inheritance

### Example

```go
// NewWizard creates a new wizard instance with defaults
func NewWizard() *Wizard {
    theme := huh.ThemeBase()
    ui := NewUIHelper()

    return &Wizard{
        theme: theme,
        ui:    ui,
        result: &WizardResult{
            GenerateQueries: true,
            GenerateSchema:  true,
        },
    }
}
```

## 🔧 Project Structure

```
SQLC-Wizzard/
├── cmd/                    # CLI entrypoints
├── internal/               # Private application code
│   ├── commands/          # CLI commands
│   ├── wizard/            # Interactive wizard
│   ├── templates/         # Template system
│   ├── generators/        # File generation
│   └── domain/           # Business logic
├── pkg/                   # Public packages
├── generated/             # TypeSpec-generated types
└── templates/             # SQL template files
```

## 🎯 Contribution Types

### 🐛 Bug Reports

- Use the issue template for bug reports
- Include reproduction steps
- Provide environment details
- Add error logs and screenshots

### ✨ Feature Requests

- Open an issue for discussion first
- Describe the use case clearly
- Consider backward compatibility
- Provide examples if possible

### 📚 Documentation

- Improve README.md
- Add code comments
- Write examples and tutorials
- Update architecture documentation

### 🧹 Code Quality

- Fix linting issues
- Improve test coverage
- Refactor complex code
- Optimize performance

## 🔄 Pull Request Process

### Before Submitting

1. **Test thoroughly**

   ```bash
   just test
   just lint
   just build
   ```

2. **Update documentation**
   - README.md changes
   - Code comments
   - Architecture docs

3. **Commit convention**

   ```
   type(scope): description

   [optional body]

   [optional footer]
   ```

   Examples:

   ```
   feat(wizard): add new project type
   fix(init): handle missing sqlc binary
   docs(readme): update installation instructions
   ```

### Submitting PR

1. Create a feature branch from `main`
2. Make your changes with small, focused commits
3. Ensure all tests pass
4. Update documentation
5. Submit PR with clear description
6. Link related issues
7. Request review from maintainers

### PR Template

```markdown
## Description
Brief description of changes made.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] No breaking changes (or documented)
```

## 🏷️ Release Process

1. Update version in `go.mod`
2. Update CHANGELOG.md
3. Create git tag
4. GitHub Actions will automatically:
   - Run tests
   - Build binaries
   - Create GitHub release
   - Update Homebrew formula
   - Build Docker image

## 🤝 Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Maintain professional communication

### Getting Help

- Open GitHub Discussions for questions
- Join our Discord community
- Check existing issues before creating new ones
- Tag maintainers for urgent issues

## 🎉 Recognition

Contributors are recognized in:

- README.md contributors section
- Release notes
- GitHub contributors graph
- Special mentions in blog posts

## 📞 Contact

- **Maintainer**: Lars Artmann
- **Email**: lars@artmann.email
- **GitHub**: @LarsArtmann
- **Discord**: [Link to Discord server]

---

Thank you for contributing to SQLC-Wizard! Your contributions make this tool better for everyone. 🙏
