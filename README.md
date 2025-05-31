# Stainless CLI

The official CLI for the Stainless REST API.

It is generated with [Stainless](https://www.stainless.com/).

## Installation

### Installing with Homebrew

```sh
brew tap stainless-api/stl
brew install stl
```

### Installing with Go

<!-- x-release-please-start-version -->

```sh
go install 'github.com/stainless-api/stainless-api-cli'
```

<!-- x-release-please-end -->

## Usage

The CLI follows a resource-based command structure:

```sh
stl [resource] [command] [flags]
```

```sh
stl builds create [--allow-empty] [--project <name>]
```

For details about specific commands, use the `--help` flag.

## Workspace Configuration

The CLI supports workspace configuration to avoid repeatedly specifying the project name. When you run a command, the CLI will:

1. Check if a project name is provided via command-line flag
2. If not, look for a `stainless-workspace.json` file in the current directory or any parent directory
3. Use the project name from the workspace configuration if found

### Initializing a Workspace

You can initialize a workspace configuration with:

```sh
stl workspace init --project your-project-name
```

If you don't provide the `--project` flag, you'll be prompted to enter a project name.

Additionally, when you run a command with a project name flag in an interactive terminal, the CLI will offer to initialize a workspace configuration for you.
