# solar-boat-cli

Solar Boat is a command-line interface tool designed for Infrastructure as Code (IaC) and GitOps workflows. It provides a wide range of Developer Experience (DX) capabilities for managing infrastructure and deployments.

## Why "Solar Boat"?

Inspired by the Ancient Egyptian Solar Boats that carried Pharaohs through their celestial journey, this CLI tool serves as a modern vessel that carries developers through the complexities of operations and infrastructure management. Just as the ancient boats handled the journey through the afterlife so the Pharaoh didn't have to worry about it, Solar Boat CLI handles the operational journey so developers can focus on what they do best - writing code.

## Features

- Terraform GitOps and operations management
  - Detect and plan changes in Terraform modules
  - Apply changes to affected modules
  - Automatic dependency detection
- Self-service ephemeral environments on Kubernetes (coming soon)
- Infrastructure management and deployment (coming soon)

## Usage

### Basic Commands

```bash
# Plan Terraform changes
solarboat terraform plan

# Apply Terraform changes
solarboat terraform apply

# Create ephemeral environment (coming soon)
solarboat env create

# Delete ephemeral environment (coming soon)
solarboat env delete
```

### Configuration

To install solarboat as a github action in your pipeline, you can use the following example:

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
      
      - name: Setup Solar Boat CLI
        uses: deplotix/solar-boat-action@v1
        with:
          version: 'latest'  # or specify a version like 'v1.0.0'
          
      - name: Plan Infrastructure Changes
        run: solarboat terraform plan
        
      - name: Apply Infrastructure Changes
        if: github.ref == 'refs/heads/main'
        run: solarboat terraform apply
```

This example shows how to:
1. Set up Solar Boat CLI in a GitHub Actions workflow
2. Run terraform plan on pull requests
3. Apply changes automatically when merging to main branch

You can customize the workflow based on your specific needs and security requirements.

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the BSD-3-Clause License - see the [LICENSE](LICENSE) file for details.

## Support

- Issues: [GitHub Issues](https://github.com/deplotix/solar-boat-cli/issues)
- Discussions: [GitHub Discussions](https://github.com/deplotix/solar-boat-cli/discussions)

## Acknowledgments

There aren't any contributors yet (awkward silence 😅), but special thanks in advance to all future contributors who share my passion for making developers' lives easier! Can't wait to build something amazing together!

~ @devqik (Founder, and currently lonely contributor)
