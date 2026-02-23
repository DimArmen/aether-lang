# Contributing to Aether

Thank you for your interest in contributing to Aether! This document provides guidelines and instructions for contributing.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/aether.git
   cd aether
   ```
3. **Install dependencies**:
   ```bash
   make deps
   ```
4. **Build the project**:
   ```bash
   make build
   ```

## Development Workflow

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and test them:
   ```bash
   make test
   make fmt
   make vet
   ```

3. Commit your changes with a descriptive commit message:
   ```bash
   git commit -m "Add feature: your feature description"
   ```

4. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

5. Open a Pull Request on GitHub

## Code Style

- Follow Go conventions and idiomatic Go code
- Run `make fmt` before committing
- Run `make vet` to catch common issues
- Add comments for exported functions and types
- Write clear, self-documenting code

## Testing

- Write tests for new functionality
- Ensure all tests pass: `make test`
- Aim for >80% code coverage
- Include both unit tests and integration tests where appropriate

## Pull Request Guidelines

### PR Title Format
- `feat: Add new feature`
- `fix: Fix bug description`
- `docs: Update documentation`
- `refactor: Refactor component`
- `test: Add tests`
- `chore: Update dependencies`

### PR Description
Include:
- **What**: Brief description of changes
- **Why**: Reason for the changes
- **How**: Implementation approach
- **Testing**: How you tested the changes
- **Screenshots**: If applicable

### Review Process
- All PRs require at least one approval
- CI tests must pass
- Code review feedback must be addressed
- Keep PRs focused and reasonably sized

## Areas for Contribution

### High Priority
- **Core Language**: Lexer, parser, type checker improvements
- **Provider Implementation**: AWS, Azure, GCP resource types
- **CLI Commands**: Implementing command functionality
- **Testing**: Unit tests, integration tests
- **Documentation**: Examples, tutorials, API docs

### Medium Priority
- **State Management**: Backend implementations
- **AI Agents**: Analyzer and autonomous agent logic
- **Module System**: Module registry, packaging
- **VS Code Extension**: Language support, syntax highlighting

### Future/Research
- **JIT Compilation**: Performance optimizations
- **Advanced AI**: Multi-agent coordination
- **Visual Tools**: Infrastructure editor
- **Migration Tools**: Terraform/CloudFormation conversion

## Bug Reports

Good bug reports include:
- **Clear title** describing the issue
- **Steps to reproduce** the bug
- **Expected behavior** vs actual behavior
- **Environment** (OS, Go version, Aether version)
- **Code samples** or minimal reproduction
- **Logs/error messages** if applicable

Use the bug report template on GitHub Issues.

## Feature Requests

Feature requests should include:
- **Use case**: Why do you need this feature?
- **Proposed solution**: How should it work?
- **Alternatives**: Other approaches considered
- **Examples**: Code samples showing usage

Use the feature request template on GitHub Issues.

## Questions and Discussion

- **Discord**: Join our [Discord server](https://discord.gg/aether-lang) for real-time discussion
- **GitHub Discussions**: For longer-form questions and design discussions
- **Stack Overflow**: Tag questions with `aether-lang`

## Code of Conduct

### Our Standards

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on what is best for the community
- Show empathy towards other community members

### Unacceptable Behavior

- Harassment, discrimination, or intimidation
- Trolling, insulting, or derogatory comments
- Publishing others' private information
- Other conduct inappropriate in a professional setting

### Enforcement

Violations of the code of conduct may result in:
1. Warning from maintainers
2. Temporary ban from community spaces
3. Permanent ban in severe cases

Report issues to: conduct@aether-lang.org

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Eligible for contributor swag (coming soon!)

Major contributors may be invited to join the core team.

## Development Setup Tips

### Recommended Tools
- **IDE**: VS Code with Go extension
- **Debugger**: Delve
- **Linter**: golangci-lint
- **Git hooks**: pre-commit for automated checks

### Project Structure
```
aether/
├── cmd/           # Command-line interface
├── pkg/           # Core packages
│   ├── lexer/     # Tokenization
│   ├── parser/    # Parsing
│   ├── ast/       # Abstract syntax tree
│   ├── types/     # Type system
│   └── ...
├── providers/     # Cloud providers
├── examples/      # Example code
└── docs/          # Documentation
```

### Running During Development
```bash
# Quick build and run
make run

# Run with an example
make run-example

# Watch mode (requires entr or similar)
ls **/*.go | entr -r make run

# Debug mode
dlv debug ./cmd/aether
```

## License

By contributing to Aether, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Aether! 🚀
