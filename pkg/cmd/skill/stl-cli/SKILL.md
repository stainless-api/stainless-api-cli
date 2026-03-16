---
name: stl-cli
description: Use when working with the Stainless CLI (stl) to manage SDK generation projects, builds, branches, configs, and authentication. Covers all stl commands including auth, workspace, init, dev, lint, and API resources (projects, builds, orgs, user).
---

# Stainless CLI (`stl`)

## Overview

`stl` is the CLI for the Stainless API — a platform that generates and manages SDKs from OpenAPI specs. Version `0.1.0-alpha.79`.

## Global Options

These apply to all API resource commands:

| Flag | Description |
|------|-------------|
| `--debug` | Enable debug logging |
| `--base-url string` | Override base URL for API requests |
| `--format string` | Response format: `auto`, `explore`, `json`, `jsonl`, `pretty`, `raw`, `yaml` (default: `auto`) |
| `--format-error string` | Error format (same options as `--format`) |
| `--transform string` | GJSON transformation for data output |
| `--transform-error string` | GJSON transformation for errors |
| `--environment string` | Set the environment for API requests |

## Auth

```
stl auth login              # OAuth browser-based login
stl auth logout             # Remove saved credentials
stl auth status             # Check current auth status
```

`login` accepts `--client-id string` to override the OAuth client ID.

## Workspace

```
stl workspace init          # Initialize workspace config in current directory
stl workspace status        # Show workspace configuration status
```

Note: `workspace init` and `workspace status` do not accept `--help`/`-h` flags.

## Init

Set up a new Stainless project:

```
stl init [options]
```

| Flag | Description |
|------|-------------|
| `--org string` | Organization name |
| `--project string` | Project slug |
| `--display-name string` | Project display name |
| `--targets string` | Comma-separated target languages |
| `--openapi-spec / --oas string` | Path to OpenAPI spec file |
| `--workspace-init` | Also initialize workspace configuration |
| `--download-config` | Download Stainless config to workspace |
| `--download-targets` | Download and configure SDK targets |

## Dev / Preview

Interactive development mode with build monitoring:

```
stl dev [options]
stl preview [options]       # alias
```

| Flag | Description |
|------|-------------|
| `-p / --project string` | Project name |
| `--oas / --openapi-spec string` | Path to OpenAPI spec |
| `--config / --stainless-config string` | Path to Stainless config |
| `-b / --branch string` | Branch to use |
| `-t / --target string` | Target language(s) — repeatable |
| `-w / --watch` | Watch mode: rebuild when files change |

## Lint

```
stl lint [options]
```

| Flag | Description |
|------|-------------|
| `-p / --project string` | Project name |
| `--oas / --openapi-spec string` | Path to OpenAPI spec |
| `--config / --stainless-config string` | Path to Stainless config |
| `-w / --watch` | Watch for file changes and re-run |

## MCP Server

Run Stainless as an MCP server:

```
stl mcp [options]
```

| Flag | Description |
|------|-------------|
| `--transport string` | `stdio` (default, local) or `http` (remote) |
| `--port number` | Port for http transport (default: 3000) |
| `--socket string` | Unix socket for http transport |
| `--tools array` | Explicitly enable tools: `code`, `docs` |
| `--no-tools array` | Explicitly disable tools: `code`, `docs` |
| `--code-execution-mode string` | `stainless-sandbox` (default) or `local` |
| `--code-allow-http-gets` | Allow GET-mapped code tool methods |
| `--code-allowed-methods array` | Regex patterns of allowed methods |
| `--code-blocked-methods array` | Regex patterns of blocked methods |
| `--stainless-api-key string` | API key for Stainless-hosted tool endpoints |

---

## API Resources

### `stl projects`

| Command | Description |
|---------|-------------|
| `create` | Create a new project |
| `retrieve` | Retrieve a project by name |
| `update` | Update a project's properties |
| `list` | List projects in an org (oldest→newest) |
| `generate-commit-message` | AI commit message by comparing two git refs |

**create** flags: `--display-name`, `--org`, `--slug`, `--target` (repeatable), `--revision string=any` (repeatable, file contents)

**retrieve** flags: `--project`

**update** flags: `--project`, `--display-name`

**list** flags: `--org`, `--cursor`, `--limit` (default 10, max 100)

**generate-commit-message** flags: `--project`, `--target`, `--base-ref`, `--head-ref`

---

### `stl projects:branches`

| Command | Description |
|---------|-------------|
| `create` | Create a new branch |
| `retrieve` | Retrieve a branch by name |
| `list` | List branches for a project |
| `delete` | Delete a branch |
| `rebase` | Rebase a branch onto another |
| `reset` | Reset a branch |

**create** flags: `--project`, `--branch`, `--branch-from`, `--force` (don't error if branch exists)

**retrieve / delete** flags: `--project`, `--branch`

**list** flags: `--project`, `--cursor`, `--limit` (default 10, max 100)

**rebase** flags: `--project`, `--branch`, `--base` (default: `main`)

**reset** flags: `--project`, `--branch`, `--target-config-sha` (required when resetting main)

---

### `stl projects:configs`

| Command | Description |
|---------|-------------|
| `retrieve` | Get config files for a project |
| `guess` | AI-generated config suggestions from an OpenAPI spec |

**retrieve** flags: `--project`, `--branch` (default: main), `--include`

**guess** flags: `--project`, `--spec` (OpenAPI spec), `--branch` (default: main)

---

### `stl builds`

| Command | Description |
|---------|-------------|
| `create` | Create a build against an input revision |
| `retrieve` | Retrieve a build by ID |
| `list` | List user-triggered builds for a project |
| `compare` | Create two comparable builds (base vs head) |

**create** key flags:

| Flag | Description |
|------|-------------|
| `--oas / --openapi-spec` | Path to OpenAPI spec |
| `--config / --stainless-config` | Path to Stainless config |
| `--project` | Project name |
| `--branch` | Project branch |
| `--target` (repeatable) | SDK targets to build (default: all) |
| `--revision` | Branch name, commit SHA, merge command (`base..head`), or file contents |
| `--wait string` | `all` (default), `commit`, or `none` |
| `--pull` | Pull build outputs after completion (requires `--wait`) |
| `--commit-message` | Commit message for new commit |
| `--enable-ai-commit-message` | Generate AI commit messages |
| `--target-commit-messages` (repeatable) | Per-SDK commit messages |
| `--allow-empty` | Allow empty commits |

**retrieve** flags: `--build-id`

**list** flags: `--project`, `--branch`, `--cursor`, `--limit` (default 10, max 100), `--revision`

**compare** flags: `--project`, `--target` (repeatable), `--base string=any` (repeatable), `--head string=any` (repeatable)

---

### `stl builds:diagnostics`

```
stl builds:diagnostics list
```

| Flag | Description |
|------|-------------|
| `--build-id` | Build ID |
| `--cursor` | Pagination cursor |
| `--limit` | Default 100, max 100 |
| `--severity` | Min severity: `fatal > error > warning > note` |
| `--targets` | Comma-delimited language targets to filter |

---

### `stl builds:target-outputs`

```
stl builds:target-outputs retrieve [options]
```

| Flag | Description |
|------|-------------|
| `--build-id` | Build ID |
| `--project` | Project name (required if no build-id) |
| `--branch` | Branch name (default: main) |
| `--target` (repeatable) | SDK language target name(s) |
| `--type` | Output type |
| `--output` | Format: `url` (download URL) or `git` (temp access token) |
| `--pull` | Pull the outputs |

---

### `stl orgs`

```
stl orgs retrieve --org <name>
stl orgs list
```

---

### `stl user`

```
stl user retrieve            # Get current authenticated user info
```

---

## Common Workflows

**Initial project setup:**
```bash
stl auth login
stl init --org myorg --project myapi --targets typescript,python --oas ./openapi.yaml
```

**Trigger a build and wait:**
```bash
stl builds create --project myapi --oas ./openapi.yaml --wait all --pull
```

**Local dev loop:**
```bash
stl dev --project myapi --oas ./openapi.yaml --watch
```

**Lint config:**
```bash
stl lint --oas ./openapi.yaml --config ./stainless.yaml --watch
```

**Check a build's diagnostics:**
```bash
stl builds:diagnostics list --build-id <id> --severity warning
```

**Download SDK output:**
```bash
stl builds:target-outputs retrieve --project myapi --target typescript --output git
```

**Branch management:**
```bash
stl projects:branches create --project myapi --branch feature/new-endpoints --branch-from main
stl projects:branches rebase --project myapi --branch feature/new-endpoints --base main
stl projects:branches delete --project myapi --branch feature/new-endpoints
```
