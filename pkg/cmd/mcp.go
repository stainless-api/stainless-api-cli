// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v3"
)

var mcpCommand = cli.Command{
	Name:            "mcp",
	Usage:           "Run Stainless MCP server",
	Description:     "Wrapper around @stainless-api/mcp@latest with environment variables set",
	Action:          handleMCP,
	ArgsUsage:       "[MCP_ARGS...]",
	HideHelpCommand: true,
	SkipFlagParsing: true,
}

func handleMCP(ctx context.Context, cmd *cli.Command) error {
	args := []string{"-y", "@stainless-api/mcp@latest"}

	cc := getAPICommandContext(cmd)

	if cmd.Args().Len() > 0 {
		args = append(args, cmd.Args().Slice()...)
	}

	env := os.Environ()

	// Set STAINLESS_API_KEY if not already in environment
	if apiKey := os.Getenv("STAINLESS_API_KEY"); apiKey == "" {
		authConfig := &AuthConfig{}
		if found, err := authConfig.Find(); err == nil && found && authConfig.AccessToken != "" {
			env = append(env, fmt.Sprintf("STAINLESS_API_KEY=%s", authConfig.AccessToken))
		}
	}

	// Set STAINLESS_PROJECT from workspace config if available
	if cc.workspaceConfig.Project != "" {
		env = append(env, fmt.Sprintf("STAINLESS_PROJECT=%s", cc.workspaceConfig.Project))
	}

	npmCmd := exec.CommandContext(ctx, "npx", args...)
	npmCmd.Env = env
	npmCmd.Stdout = os.Stdout
	npmCmd.Stderr = os.Stderr
	npmCmd.Stdin = os.Stdin

	return npmCmd.Run()
}
