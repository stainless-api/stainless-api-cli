// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mocktest"
	"github.com/stainless-api/stainless-api-cli/internal/requestflag"
)

func TestBuildsCreate(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"builds", "create",
		"--project", "project",
		"--revision", "string",
		"--allow-empty=true",
		"--branch", "branch",
		"--commit-message", "commit_message",
		"--enable-ai-commit-message=true",
		"--target-commit-messages", "{cli: cli, csharp: csharp, go: go, java: java, kotlin: kotlin, node: node, openapi: openapi, php: php, python: python, ruby: ruby, sql: sql, terraform: terraform, typescript: typescript}",
		"--target", "node",
	)

	// Check that inner flags have been set up correctly
	requestflag.CheckInnerFlags(buildsCreate)

	// Alternative argument passing style using inner flags
	mocktest.TestRunMockTestWithFlags(
		t,
		"builds", "create",
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
}

func TestBuildsRetrieve(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"builds", "retrieve",
		"--build-id", "buildId",
	)
}

func TestBuildsList(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"builds", "list",
		"--project", "project",
		"--branch", "branch",
		"--cursor", "cursor",
		"--limit", "1",
		"--revision", "string",
	)
}

func TestBuildsCompare(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"builds", "compare",
		"--base", "{branch: branch, revision: string, commit_message: commit_message}",
		"--head", "{branch: branch, revision: string, commit_message: commit_message}",
		"--project", "project",
		"--target", "node",
	)

	// Check that inner flags have been set up correctly
	requestflag.CheckInnerFlags(buildsCompare)

	// Alternative argument passing style using inner flags
	mocktest.TestRunMockTestWithFlags(
		t,
		"builds", "compare",
		"--base.branch", "branch",
		"--base.revision", "string",
		"--base.commit-message", "commit_message",
		"--head.branch", "branch",
		"--head.revision", "string",
		"--head.commit-message", "commit_message",
		"--project", "project",
		"--target", "node",
	)
}
