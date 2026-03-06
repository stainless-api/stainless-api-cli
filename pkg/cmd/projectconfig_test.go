// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mocktest"
)

func TestProjectsConfigsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "projects:configs", "retrieve",
			"--api-key", "string",
			"--project", "project",
			"--branch", "branch",
			"--include", "include",
		)
	})
}

func TestProjectsConfigsGuess(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t, "projects:configs", "guess",
			"--api-key", "string",
			"--project", "project",
			"--spec", "spec",
			"--branch", "branch",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"spec: spec\n" +
			"branch: branch\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData, "projects:configs", "guess",
			"--api-key", "string",
			"--project", "project",
		)
	})
}
