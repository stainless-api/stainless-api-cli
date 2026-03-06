// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mocktest"
	"github.com/stainless-api/stainless-api-cli/internal/requestflag"
)

func TestBuildsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "builds", "create",
			"--api-key", "string",
			"--project", "project",
			"--revision", "string",
			"--allow-empty=true",
			"--branch", "branch",
			"--commit-message", "commit_message",
			"--enable-ai-commit-message=true",
			"--target-commit-messages", "{cli: cli, csharp: csharp, go: go, java: java, kotlin: kotlin, node: node, openapi: openapi, php: php, python: python, ruby: ruby, sql: sql, terraform: terraform, typescript: typescript}",
			"--target", "node",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(buildsCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t, "builds", "create",
			"--api-key", "string",
			"--project", "project",
			"--revision", "string",
			"--allow-empty=true",
			"--branch", "branch",
			"--commit-message", "commit_message",
			"--enable-ai-commit-message=true",
			"--target-commit-messages.cli", "cli",
			"--target-commit-messages.csharp", "csharp",
			"--target-commit-messages.go", "go",
			"--target-commit-messages.java", "java",
			"--target-commit-messages.kotlin", "kotlin",
			"--target-commit-messages.node", "node",
			"--target-commit-messages.openapi", "openapi",
			"--target-commit-messages.php", "php",
			"--target-commit-messages.python", "python",
			"--target-commit-messages.ruby", "ruby",
			"--target-commit-messages.sql", "sql",
			"--target-commit-messages.terraform", "terraform",
			"--target-commit-messages.typescript", "typescript",
			"--target", "node",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"project: project\n" +
			"revision: string\n" +
			"allow_empty: true\n" +
			"branch: branch\n" +
			"commit_message: commit_message\n" +
			"enable_ai_commit_message: true\n" +
			"target_commit_messages:\n" +
			"  cli: cli\n" +
			"  csharp: csharp\n" +
			"  go: go\n" +
			"  java: java\n" +
			"  kotlin: kotlin\n" +
			"  node: node\n" +
			"  openapi: openapi\n" +
			"  php: php\n" +
			"  python: python\n" +
			"  ruby: ruby\n" +
			"  sql: sql\n" +
			"  terraform: terraform\n" +
			"  typescript: typescript\n" +
			"targets:\n" +
			"  - node\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData, "builds", "create",
			"--api-key", "string",
		)
	})
}

func TestBuildsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "builds", "retrieve",
			"--api-key", "string",
			"--build-id", "buildId",
		)
	})
}

func TestBuildsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "builds", "list",
			"--api-key", "string",
			"--max-items", "10",
			"--project", "project",
			"--branch", "branch",
			"--cursor", "cursor",
			"--limit", "1",
			"--revision", "string",
		)
	})
}

func TestBuildsCompare(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "builds", "compare",
			"--api-key", "string",
			"--base", "{branch: branch, revision: string, commit_message: commit_message}",
			"--head", "{branch: branch, revision: string, commit_message: commit_message}",
			"--project", "project",
			"--target", "node",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(buildsCompare)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t, "builds", "compare",
			"--api-key", "string",
			"--base.branch", "branch",
			"--base.revision", "string",
			"--base.commit-message", "commit_message",
			"--head.branch", "branch",
			"--head.revision", "string",
			"--head.commit-message", "commit_message",
			"--project", "project",
			"--target", "node",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"base:\n" +
			"  branch: branch\n" +
			"  revision: string\n" +
			"  commit_message: commit_message\n" +
			"head:\n" +
			"  branch: branch\n" +
			"  revision: string\n" +
			"  commit_message: commit_message\n" +
			"project: project\n" +
			"targets:\n" +
			"  - node\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData, "builds", "compare",
			"--api-key", "string",
		)
	})
}
