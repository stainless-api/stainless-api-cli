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
stl builds create \
  --revision string
```

For details about specific commands, use the `--help` flag.

## Global Flags

- `--debug` - Enable debug logging (includes HTTP request/response details)
- `--version`, `-v` - Show the CLI version
