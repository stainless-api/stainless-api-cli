# Stainless V0 CLI

The official CLI for the Stainless V0 REST API.

It is generated with [Stainless](https://www.stainless.com/).

## Installation

### Installing with Go

<!-- x-release-please-start-version -->

```sh
go install 'github.com/stainless-api/stainless-api-cli'
```

<!-- x-release-please-end -->

## Usage

The CLI follows a resource-based command structure:

```sh
stainless-api-cli [resource] [command] [flags]
```

```sh
stainless-api-cli builds create \
  --revision string
```

For details about specific commands, use the `--help` flag.
