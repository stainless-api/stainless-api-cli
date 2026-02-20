// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mocktest"
)

func TestProjectsCreate(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"projects", "create",
		"--display-name", "display_name",
		"--org", "org",
		"--revision", "{foo: {content: content}}",
		"--slug", "slug",
		"--target", "node",
	)
}

func TestProjectsRetrieve(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"projects", "retrieve",
		"--project", "project",
	)
}

func TestProjectsUpdate(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"projects", "update",
		"--project", "project",
		"--display-name", "display_name",
	)
}

func TestProjectsList(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"projects", "list",
		"--cursor", "cursor",
		"--limit", "1",
		"--org", "org",
	)
}

func TestProjectsGenerateCommitMessage(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"projects", "generate-commit-message",
		"--project", "project",
		"--target", "python",
		"--base-ref", "base_ref",
		"--head-ref", "head_ref",
	)
}
