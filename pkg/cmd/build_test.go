// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mocktest"
)

func TestBuildsCreate(t *testing.T) {
	t.Skip("Prism tests are disabled")
	mocktest.TestRunMockTestWithFlags(
		t,
		"builds", "create",
		"--project", "project",
		"--revision", "string",
		"--allow-empty",
		"--branch", "branch",
		"--commit-message", "commit_message",
		"--target-commit-messages", "{cli: cli, csharp: csharp, go: go, java: java, kotlin: kotlin, node: node, openapi: openapi, php: php, python: python, ruby: ruby, sql: sql, terraform: terraform, typescript: typescript}",
		"--target", "node",
	)
}

func TestBuildsRetrieve(t *testing.T) {
	t.Skip("Prism tests are disabled")
	mocktest.TestRunMockTestWithFlags(
		t,
		"builds", "retrieve",
		"--build-id", "buildId",
	)
}

func TestBuildsList(t *testing.T) {
	t.Skip("Prism tests are disabled")
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
	t.Skip("Prism tests are disabled")
	mocktest.TestRunMockTestWithFlags(
		t,
		"builds", "compare",
		"--base", "{branch: branch, revision: string, commit_message: commit_message}",
		"--head", "{branch: branch, revision: string, commit_message: commit_message}",
		"--project", "project",
		"--target", "node",
	)
}
