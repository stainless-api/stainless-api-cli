// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mocktest"
)

func TestOrgsRetrieve(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"orgs", "retrieve",
		"--org", "org",
	)
}

func TestOrgsList(t *testing.T) {
	mocktest.TestRunMockTestWithFlags(
		t,
		"orgs", "list",
	)
}
