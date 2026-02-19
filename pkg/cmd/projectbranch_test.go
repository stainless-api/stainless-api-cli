// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mocktest"
)

func TestProjectsBranchesCreate(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	mocktest.TestRunMockTestWithFlags(
		t,
		"projects:branches", "create",
		"--project", "project",
		"--branch", "branch",
		"--branch-from", "branch_from",
		"--force=true",
	)
}

func TestProjectsBranchesRetrieve(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	mocktest.TestRunMockTestWithFlags(
		t,
		"projects:branches", "retrieve",
		"--project", "project",
		"--branch", "branch",
	)
}

func TestProjectsBranchesList(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	mocktest.TestRunMockTestWithFlags(
		t,
		"projects:branches", "list",
		"--project", "project",
		"--cursor", "cursor",
		"--limit", "1",
	)
}

func TestProjectsBranchesDelete(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	mocktest.TestRunMockTestWithFlags(
		t,
		"projects:branches", "delete",
		"--project", "project",
		"--branch", "branch",
	)
}

func TestProjectsBranchesRebase(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	mocktest.TestRunMockTestWithFlags(
		t,
		"projects:branches", "rebase",
		"--project", "project",
		"--branch", "branch",
		"--base", "base",
	)
}

func TestProjectsBranchesReset(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	mocktest.TestRunMockTestWithFlags(
		t,
		"projects:branches", "reset",
		"--project", "project",
		"--branch", "branch",
		"--target-config-sha", "target_config_sha",
	)
}
