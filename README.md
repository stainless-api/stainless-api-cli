# Stainless CLI

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

```sh
go install 'github.com/stainless-api/stainless-api-cli/cmd/stl@latest'
```

<!-- x-release-please-end -->

### Running Locally

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
