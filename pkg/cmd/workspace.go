// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stainless-api/stainless-api-cli/pkg/console"

	"github.com/urfave/cli/v3"
)

// Rel returns a relative path similar to filepath.Rel but with custom behavior:
// - If target is empty, returns empty string
// - If relative path doesn't start with "../", it prefixes with "./"
func Rel(basepath, targpath string) string {
	if targpath == "" {
		return ""
	}

	rel, err := filepath.Rel(basepath, targpath)
	if err != nil {
		return targpath
	}

	if !strings.HasPrefix(rel, "../") && !strings.HasPrefix(rel, "./") {
		rel = "./" + rel
	}

	return rel
}

var workspaceInit = cli.Command{
	Name:     "init",
	Usage:    "Initialize Stainless workspace configuration in current directory",
	HideHelp: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "project",
			Usage: "Project name to use for this workspace",
		},
		&cli.StringFlag{
			Name:    "openapi-spec",
			Aliases: []string{"oas"},
			Usage:   "Path to OpenAPI spec file",
		},
		&cli.StringFlag{
			Name:    "stainless-config",
			Aliases: []string{"config"},
			Usage:   "Path to Stainless config file",
		},
		&cli.BoolFlag{
			Name:  "download-config",
			Usage: "Download Stainless config to workspace",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "download-targets",
			Usage: "Download configured targets after build completion",
			Value: true,
		},
	},
	Action:          handleInit,
	HideHelpCommand: true,
}

var workspaceStatus = cli.Command{
	Name:            "status",
	Usage:           "Show workspace configuration status",
	Action:          handleWorkspaceStatus,
	HideHelpCommand: true,
}

func handleWorkspaceStatus(ctx context.Context, cmd *cli.Command) error {
	wc := WorkspaceConfig{}
	if _, err := wc.Find(); err != nil {
		console.Warn("%s", err)
	}

	if wc.ConfigPath == "" {
		group := console.Warn("No workspace configuration found")
		group.Info("Run 'stl workspace init' to initialize a workspace in this directory.")
		return nil
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Get relative path from cwd to config file
	relPath, err := filepath.Rel(cwd, wc.ConfigPath)
	if err != nil {
		relPath = wc.ConfigPath // fallback to absolute path
	}

	group := console.Success("Workspace configuration found")
	group.Property("path", relPath)
	group.Property("project", wc.Project)

	if wc.OpenAPISpec != "" {
		// Check if OpenAPI spec file exists (path is already absolute)
		if _, err := os.Stat(wc.OpenAPISpec); err == nil {
			group.Property("openapi_spec", wc.OpenAPISpec)
		} else {
			group.Property("openapi_spec", wc.OpenAPISpec+" (not found)")
		}
	} else {
		group.Property("openapi_spec", "(not configured)")
	}

	if wc.StainlessConfig != "" {
		// Check if Stainless config file exists (path is already absolute)
		if _, err := os.Stat(wc.StainlessConfig); err == nil {
			group.Property("stainless_config", wc.StainlessConfig)
		} else {
			group.Property("stainless_config", wc.StainlessConfig+" (not found)")
		}
	} else {
		group.Property("stainless_config", "(not configured)")
	}

	return nil
}
