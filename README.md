# Aether

[![GitHub](https://img.shields.io/badge/GitHub-DimArmen%2Faether--lang-blue?logo=github)](https://github.com/DimArmen/aether-lang)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev)

A next-generation Infrastructure as Code language with AI-powered intelligence.

> **Status:** Early development (v0.1.0-alpha) - Core lexer, parser, and CLI framework complete. Actively implementing runtime and provider integrations.

## Quick Start

```bash
# Initialize a new project
aether init my-infrastructure

# Create your first resource
cat > main.ae <<EOF
resource compute.instance "web" {
  machine_type = "medium"
  os_image = "ubuntu-22.04"
  region = "us-east"
  
  tags = {
    Name = "web-server"
    Environment = "production"
  }
}

output "server_ip" {
  value = instance.web.public_ip
}
EOF

# Preview changes
aether plan

# Apply changes
aether apply
```

## Installation

### From Source

```bash
git clone https://github.com/aether-lang/aether.git
cd aether
go build -o aether ./cmd/aether
sudo mv aether /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/aether-lang/aether/cmd/aether@latest
```

## Documentation

- [Vision](VISION.md) - Project goals and philosophy
- [Design](DESIGN.md) - Technical design and architecture
- [Syntax](docs/syntax.md) - Language syntax reference
- [Providers](docs/providers.md) - Multi-cloud provider system
- [AWS Guide](docs/aws-guide.md) - Comprehensive AWS provider documentation
- [AI Agents](docs/ai-agents.md) - AI agent framework
- [Roadmap](ROADMAP.md) - Development roadmap
- [Examples](EXAMPLES.md) - Code examples
- [AWS Examples](examples/README-AWS.md) - AWS-specific examples and tutorials

## Features

- **Hybrid Syntax**: Declarative resources with embedded scripting for complex logic
- **Multi-Cloud**: Deploy to AWS, Azure, GCP from the same code
- **AI Agents**: Three-tier AI system (Assistant, Analyzers, Autonomous)
- **Type Safe**: Strong type system catches errors before deployment
- **Built-in Testing**: Unit tests, integration tests, property-based tests
- **Intelligent State**: Automatic locking, encryption, drift detection

## Project Structure

```
aether/
├── cmd/
│   └── aether/          # CLI entry point
├── pkg/
│   ├── lexer/           # Tokenizer
│   ├── parser/          # Parser and AST
│   ├── types/           # Type system
│   ├── interpreter/     # Interpreter
│   ├── runtime/         # Resource management
│   ├── provider/        # Provider interface
│   ├── state/           # State management
│   ├── agent/           # AI agents
│   └── cli/             # CLI commands
├── providers/
│   ├── aws/             # AWS provider
│   ├── azure/           # Azure provider
│   └── gcp/             # GCP provider
├── stdlib/              # Standard library modules
├── examples/            # Example projects
└── docs/                # Documentation
```

## Development Status

**Version**: 0.1.0-alpha  
**Status**: Active Development

See [ROADMAP.md](ROADMAP.md) for details.

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

- Report bugs and request features through [GitHub Issues](https://github.com/aether-lang/aether/issues)
- Submit pull requests following our coding standards
- Join our [Discord](https://discord.gg/aether-lang) for discussions

## License

MIT License - see [LICENSE](LICENSE) for details

## Community

- **Discord**: https://discord.gg/aether-lang
- **Twitter**: @aetherlang
- **Website**: https://aether-lang.org

---

*"Infrastructure should be intelligent, safe, and invisible. Aether makes it so."*
