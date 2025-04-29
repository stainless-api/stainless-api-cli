# Stainless API CLI

A command-line interface for interacting with the Stainless API to manage projects, builds, and SDK generation.

## Installation

```bash
go install github.com/stainless-api/stainless-api-cli@latest
```

## Authentication

Before using the CLI, you need to authenticate with the Stainless API.

```bash
# Log in using OAuth device flow
stainless-api-cli auth login

# Check authentication status
stainless-api-cli auth status

# Log out and remove saved credentials
stainless-api-cli auth logout
```

You can also authenticate by setting the `STAINLESS_API_KEY` environment variable, which takes precedence over saved credentials.

## Projects

Manage your Stainless projects.

```bash
# Retrieve a project
stainless-api-cli projects retrieve --project-name <project-name>

# Update a project
stainless-api-cli projects update --project-name <project-name> --display-name "New Project Name"
```

### Project Branches

```bash
# Create a new branch
stainless-api-cli projects:branches create --project <project-name> --branch <branch-name> --branch-from main

# Retrieve a branch
stainless-api-cli projects:branches retrieve --project <project-name> --branch <branch-name>
```

### Project Configs

```bash
# Retrieve project configuration
stainless-api-cli projects:configs retrieve --project <project-name> --branch <branch-name>

# Guess project configuration based on OpenAPI spec
stainless-api-cli projects:configs guess --project <project-name> --spec <path-to-spec>
```

## Builds

Create and manage builds for your projects.

```bash
# Create a new build
stainless-api-cli builds create --project <project-name> --revision <revision> --openapi-spec <path-to-spec> --stainless-config <path-to-config>

# Create a build and wait for completion
stainless-api-cli builds create --project <project-name> --revision <revision> --openapi-spec <path-to-spec> --wait

# Create a build, wait for completion, and pull outputs
stainless-api-cli builds create --project <project-name> --revision <revision> --openapi-spec <path-to-spec> --wait --pull

# Retrieve a build
stainless-api-cli builds retrieve --build-id <build-id>

# List builds for a project
stainless-api-cli builds list --project <project-name> --branch <branch-name>
```

## Build Target Outputs

Retrieve and pull build target outputs.

```bash
# Retrieve build target output
stainless-api-cli build_target_outputs retrieve --build-id <build-id> --target <target> --type <type> --output <output>

# Pull build target output
stainless-api-cli build_target_outputs pull --build-id <build-id> --target <target> --type <type> --output <output>
```

## Environment Variables

- `STAINLESS_API_KEY`: API key for authentication (takes precedence over saved credentials)
- `NO_COLOR`: Disable colored output
- `FORCE_COLOR`: Force colored output (`1` to enable, `0` to disable)

## Examples

```bash
# Generate a TypeScript SDK for your API
stainless-api-cli builds create --project my-project --branch main --openapi-spec ./openapi.yml --wait --pull --targets typescript

# Generate multiple SDKs at once
stainless-api-cli builds create --project my-project --branch main --openapi-spec ./openapi.yml --wait --pull --targets typescript,python,go
```

## Shell Completion

The CLI supports shell completion. To enable it:

```bash
# For bash
stainless-api-cli completion bash > /etc/bash_completion.d/stainless-api-cli

# For zsh
stainless-api-cli completion zsh > "${fpath[1]}/_stainless-api-cli"

# For fish
stainless-api-cli completion fish > ~/.config/fish/completions/stainless-api-cli.fish
```

## Available SDK Targets

The Stainless API CLI can generate SDKs for multiple languages:

- `typescript` - TypeScript SDK
- `node` - Node.js SDK
- `python` - Python SDK
- `go` - Go SDK
- `ruby` - Ruby SDK
- `java` - Java SDK
- `kotlin` - Kotlin SDK
- `cli` - Command-line interface
- `terraform` - Terraform provider
