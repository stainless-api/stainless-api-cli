# Stainless CLI

> [!CAUTION]
>
> The CLI is unstable and the API may change. Please feel free to use it locally, but don't build scripts against it.

The official CLI for the [Stainless REST API](https://www.stainless.com/docs/getting-started/quickstart-cli).

It is generated with [Stainless](https://www.stainless.com/).

<!-- x-release-please-start-version -->

## Installation

### Installing with Homebrew

```sh
brew tap stainless-api/tap
brew install stl
```

### Installing with Go

To test or install the CLI locally, you need [Go](https://go.dev/doc/install) version 1.22 or later installed.

```sh
go install 'github.com/stainless-api/stainless-api-cli/cmd/stl@latest'
```

Once you have run `go install`, the binary is placed in your Go bin directory:

- **Default location**: `$HOME/go/bin` (or `$GOPATH/bin` if GOPATH is set)
- **Check your path**: Run `go env GOPATH` to see the base directory

If commands aren't found after installation, add the Go bin directory to your PATH:

```sh
# Add to your shell profile (.zshrc, .bashrc, etc.)
export PATH="$PATH:$(go env GOPATH)/bin"
```

<!-- x-release-please-end -->

### Running Locally

After cloning the git repository for this project, you can use the
`scripts/run` script to run the tool locally:

```sh
./scripts/run args...
```

## Usage

The CLI follows a resource-based command structure:

```sh
stl [resource] [command] [flags]
```

```sh
stl builds create \
  --project stainless \
  --revision main \
  --allow-empty \
  --branch branch \
  --commit-message commit_message \
  --enable-ai-commit-message \
  --target-commit-messages '{cli: cli, csharp: csharp, go: go, java: java, kotlin: kotlin, node: node, openapi: openapi, php: php, python: python, ruby: ruby, sql: sql, terraform: terraform, typescript: typescript}' \
  --target node
```

For details about specific commands, use the `--help` flag.

## Global Flags

- `--help` - Show command line usage
- `--debug` - Enable debug logging (includes HTTP request/response details)
- `--version`, `-v` - Show the CLI version

- `--base-url` - Use a custom API backend URL
- `--format` - Change the output format (`auto`, `explore`, `json`, `jsonl`, `pretty`, `raw`, `yaml`)
- `--format-error` - Change the output format for errors (`auto`, `explore`, `json`, `jsonl`, `pretty`, `raw`, `yaml`)
- `--transform` - Transform the data output using [GJSON syntax](https://github.com/tidwall/gjson/blob/master/SYNTAX.md)
- `--transform-error` - Transform the error output using [GJSON syntax](https://github.com/tidwall/gjson/blob/master/SYNTAX.md)
## Workspace Configuration

The CLI supports workspace configuration to avoid repeatedly specifying the project name. When you run a command, the CLI will:

1. Check if a project name is provided via command-line flag
2. If not, look for a `.stainless/workspace.json` file in the current directory or any parent directory
3. Use the project name from the workspace configuration if found

### Initializing a Workspace

You can initialize a workspace configuration with:

```sh
stl workspace init --project your-project-name
```

If you don't provide the `--project` flag, you'll be prompted to enter a project name.

Additionally, when you run a command with a project name flag in an interactive terminal, the CLI will offer to initialize a workspace configuration for you.
