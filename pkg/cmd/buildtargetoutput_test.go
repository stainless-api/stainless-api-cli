// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mocktest"
)

func TestBuildsTargetOutputsRetrieve(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"builds:target-outputs", "retrieve",
		"--build-id", "build_id",
		"--target", "node",
		"--type", "source",
		"--output", "url",
	)
}
