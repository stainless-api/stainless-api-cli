// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mocktest"
)

func TestProjectsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"projects", "create",
			"--display-name", "display_name",
			"--org", "org",
			"--revision", "{foo: {content: content}}",
			"--slug", "slug",
			"--target", "node",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"display_name: display_name\n" +
			"org: org\n" +
			"revision:\n" +
			"  foo:\n" +
			"    content: content\n" +
			"slug: slug\n" +
			"targets:\n" +
			"  - node\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"projects", "create",
		)
	})
}

func TestProjectsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"projects", "retrieve",
			"--project", "project",
		)
	})
}

func TestProjectsUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"projects", "update",
			"--project", "project",
			"--display-name", "display_name",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("display_name: display_name")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"projects", "update",
			"--project", "project",
		)
	})
}

func TestProjectsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"projects", "list",
			"--max-items", "10",
			"--cursor", "cursor",
			"--limit", "1",
			"--org", "org",
		)
	})
}

func TestProjectsGenerateCommitMessage(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"projects", "generate-commit-message",
			"--project", "project",
			"--target", "python",
			"--base-ref", "base_ref",
			"--head-ref", "head_ref",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"base_ref: base_ref\n" +
			"head_ref: head_ref\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"projects", "generate-commit-message",
			"--project", "project",
			"--target", "python",
		)
	})
}
