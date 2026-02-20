// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mocktest"
)

func TestProjectsConfigsRetrieve(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"projects:configs", "retrieve",
		"--project", "project",
		"--branch", "branch",
		"--include", "include",
	)
}

func TestProjectsConfigsGuess(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"projects:configs", "guess",
		"--project", "project",
		"--spec", "spec",
		"--branch", "branch",
	)
}
