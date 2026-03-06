// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mocktest"
)

func TestProjectsBranchesCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "projects:branches", "create",
			"--api-key", "string",
			"--project", "project",
			"--branch", "branch",
			"--branch-from", "branch_from",
			"--force=true",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"branch: branch\n" +
			"branch_from: branch_from\n" +
			"force: true\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData, "projects:branches", "create",
			"--api-key", "string",
			"--project", "project",
		)
	})
}

func TestProjectsBranchesRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "projects:branches", "retrieve",
			"--api-key", "string",
			"--project", "project",
			"--branch", "branch",
		)
	})
}

func TestProjectsBranchesList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "projects:branches", "list",
			"--api-key", "string",
			"--max-items", "10",
			"--project", "project",
			"--cursor", "cursor",
			"--limit", "1",
		)
	})
}

func TestProjectsBranchesDelete(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "projects:branches", "delete",
			"--api-key", "string",
			"--project", "project",
			"--branch", "branch",
		)
	})
}

func TestProjectsBranchesRebase(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "projects:branches", "rebase",
			"--api-key", "string",
			"--project", "project",
			"--branch", "branch",
			"--base", "base",
		)
	})
}

func TestProjectsBranchesReset(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "projects:branches", "reset",
			"--api-key", "string",
			"--project", "project",
			"--branch", "branch",
			"--target-config-sha", "target_config_sha",
		)
	})
}
