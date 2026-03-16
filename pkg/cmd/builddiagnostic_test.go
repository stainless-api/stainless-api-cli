// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mocktest"
)

func TestBuildsDiagnosticsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"builds:diagnostics", "list",
			"--max-items", "10",
			"--build-id", "buildId",
			"--cursor", "cursor",
			"--limit", "1",
			"--severity", "fatal",
			"--targets", "targets",
		)
	})
}
