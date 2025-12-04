# Stainless CLI

The official CLI for the [Stainless REST API](https://www.stainless.com/docs/getting-started/quickstart-cli).

It is generated with [Stainless](https://www.stainless.com/).

## Installation

### Installing with Homebrew

```sh
brew tap stainless-api/tap
brew install stl
```

### Installing with Go

<!-- x-release-please-start-version -->

```sh
go install 'github.com/stainless-api/stainless-api-cli/cmd/stl@latest'
```

### Running Locally

<!-- x-release-please-start-version -->

```sh
go run cmd/stl/main.go
```

<!-- x-release-please-end -->

## Usage

The CLI follows a resource-based command structure:

```sh
stl [resource] [command] [flags]
```

```sh
stl builds create \
  --revision main \
<<JSON
{
  "project": "stainless"
}
JSON
```

For details about specific commands, use the `--help` flag.

## Global Flags

- `--debug` - Enable debug logging (includes HTTP request/response details)
- `--version`, `-v` - Show the CLI version
