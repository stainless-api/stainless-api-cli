package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v3"
)

type configChangedEvent struct{}

func waitTillConfigChanges(ctx context.Context, cmd *cli.Command, wc WorkspaceConfig) error {
	openapiSpecPath := wc.OpenAPISpec
	if cmd.IsSet("openapi-spec") {
		openapiSpecPath = cmd.String("openapi-spec")
	}
	stainlessConfigPath := wc.StainlessConfig
	if cmd.IsSet("stainless-config") {
		stainlessConfigPath = cmd.String("stainless-config")
	}

	// Get initial file modification times
	openapiSpecInfo, err := os.Stat(openapiSpecPath)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %v", openapiSpecPath, err)
	}
	openapiSpecModTime := openapiSpecInfo.ModTime()

	stainlessConfigInfo, err := os.Stat(stainlessConfigPath)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %v", stainlessConfigPath, err)
	}
	stainlessConfigModTime := stainlessConfigInfo.ModTime()

	// Poll for file changes every 250ms
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check OpenAPI spec file
			if info, err := os.Stat(openapiSpecPath); err == nil {
				if info.ModTime().After(openapiSpecModTime) {
					return nil
				}
			}

			// Check Stainless config file
			if info, err := os.Stat(stainlessConfigPath); err == nil {
				if info.ModTime().After(stainlessConfigModTime) {
					return nil
				}
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
