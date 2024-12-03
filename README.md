# Solar Boat CLI 🚀

[![Release](https://github.com/deplotix/solar-boat-cli/actions/workflows/release.yml/badge.svg)](https://github.com/deplotix/solar-boat-cli/actions/workflows/release.yml)

Solar Boat is a command-line interface tool designed for Infrastructure as Code (IaC) and GitOps workflows. It provides intelligent Terraform operations management with automatic dependency detection and stateful/stateless module handling.

## Why "Solar Boat"? ⛵

Inspired by the Ancient Egyptian Solar Boats that carried Pharaohs through their celestial journey, this CLI tool serves as a modern vessel that carries developers through the complexities of operations and infrastructure management. Just as the ancient boats handled the journey through the afterlife so the Pharaoh didn't have to worry about it, Solar Boat CLI handles the operational journey so developers can focus on what they do best - writing code.

## Features ✨

### Current Features
- **Intelligent Terraform Operations**
  - Automatic detection of changed modules
  - Smart handling of stateful and stateless modules
  - Automatic dependency propagation
  - Parallel execution of independent modules
  - Detailed operation reporting

### Coming Soon
- Self-service ephemeral environments on Kubernetes
- Infrastructure management and deployment
- Custom workflow automation

## Installation 📦

### Using Go Install (Recommended)

```bash
# Install the latest version
go install github.com/deplotix/solar-boat-cli@latest

# Install a specific version
go install github.com/deplotix/solar-boat-cli@v0.1.4
```

### Building from Source

```bash
git clone https://github.com/deplotix/solar-boat-cli.git
cd solar-boat-cli
go build -o solarboat
```

## Usage 🛠️

### Basic Commands

```bash
# Plan Terraform changes
solarboat terraform plan

# Plan and save outputs to a specific directory
solarboat terraform plan --output-dir ./terraform-plans

# Apply Terraform changes
solarboat terraform apply
```

### Module Types

Solar Boat CLI recognizes two types of Terraform modules:

- **Stateful Modules**: Modules that manage actual infrastructure state (contain backend configuration)
- **Stateless Modules**: Reusable modules without state (no backend configuration)

When changes are detected in stateless modules, the CLI automatically identifies and processes any stateful modules that depend on them.

### GitHub Actions Integration

Add Solar Boat to your GitHub Actions workflow:

```yaml
name: Infrastructure Management

on:
  pull_request:
    branches: [ main ]
  push:
    branches: [ main ]

jobs:
  infrastructure:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0  # Important for detecting changes

      - name: Plan Infrastructure Changes
        if: github.event_name == 'pull_request'
        uses: deplotix/solar-boat-action@v1
        with:
          command: plan
          github_token: ${{ secrets.GITHUB_TOKEN }}
          
      - name: Apply Infrastructure Changes
        if: github.ref == 'refs/heads/main'
        uses: deplotix/solar-boat-action@v1
        with:
          command: apply
          github_token: ${{ secrets.GITHUB_TOKEN }}
```

This workflow will:
1. Run `terraform plan` on pull requests
2. Save plan artifacts for review
3. Comment on the PR with results
4. Apply changes when merged to main

## Contributing 🤝

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License 📄

This project is licensed under the BSD-3-Clause License - see the [LICENSE](LICENSE) file for details.

## Support 💬

- Issues: [GitHub Issues](https://github.com/deplotix/solar-boat-cli/issues)
- Discussions: [GitHub Discussions](https://github.com/deplotix/solar-boat-cli/discussions)
- Documentation: [Wiki](https://github.com/deplotix/solar-boat-cli/wiki)

## Acknowledgments 🙏

Special thanks to all contributors who help make this project better! Whether you're fixing bugs, improving documentation, or suggesting features, your contributions are greatly appreciated.

~ @devqik (Founder)
