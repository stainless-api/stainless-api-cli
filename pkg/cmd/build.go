// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/urfave/cli/v3"
)

// BuildTargetInfo holds information about a build target
type BuildTargetInfo struct {
	name   string
	status stainless.BuildTargetStatus
}

// parseTargetPaths processes target flags to extract target:path syntax
// Returns a map of target names to their custom paths
func parseTargetPaths() map[string]string {
	targetPaths := make(map[string]string)

	// First, check workspace configuration for target paths
	var config WorkspaceConfig
	found, err := config.Find()
	if err == nil && found && config.Targets != nil {
		for targetName, targetConfig := range config.Targets {
			if targetConfig.OutputPath != "" {
				targetPaths[targetName] = targetConfig.OutputPath
			}
		}
	}

	// Get the current JSON body with all mutations applied
	body, _, _, err := jsonflag.ApplyMutations([]byte("{}"), []byte("{}"), []byte("{}"))
	if err != nil {
		return targetPaths // If we can't parse, return map with workspace paths
	}

	// Check if there are any targets in the body
	targetsResult := gjson.GetBytes(body, "targets")
	if !targetsResult.Exists() {
		return targetPaths
	}

	// Process the targets array
	var cleanTargets []string
	for _, targetResult := range targetsResult.Array() {
		target := targetResult.String()
		cleanTarget, path := processSingleTarget(target)
		cleanTargets = append(cleanTargets, cleanTarget)
		if path != "" {
			// Command line target:path overrides workspace configuration
			targetPaths[cleanTarget] = path
		}
	}

	// Update the JSON body with cleaned targets
	if len(cleanTargets) > 0 {
		body, err = sjson.SetBytes(body, "targets", cleanTargets)
		if err != nil {
			return targetPaths
		}

		// Clear mutations and re-apply the cleaned JSON
		jsonflag.ClearMutations()

		// Re-register the cleaned body
		bodyObj := gjson.ParseBytes(body)
		bodyObj.ForEach(func(key, value gjson.Result) bool {
			jsonflag.Mutate(jsonflag.Body, key.String(), value.Value())
			return true
		})
	}

	return targetPaths
}

// processSingleTarget extracts path from target:path and returns clean target name and path
func processSingleTarget(target string) (string, string) {
	target = strings.TrimSpace(target)
	if !strings.Contains(target, ":") {
		return target, ""
	}

	parts := strings.SplitN(target, ":", 2)
	if len(parts) != 2 {
		return target, ""
	}

	targetName := strings.TrimSpace(parts[0])
	targetPath := strings.TrimSpace(parts[1])
	return targetName, targetPath
}

// getBuildTargetInfo extracts completed targets from a build response
func getBuildTargetInfo(buildRes stainless.BuildObject) []BuildTargetInfo {
	targets := []BuildTargetInfo{}

	// Check each target and add it to the list if it's completed or in postgen
	if buildRes.Targets.JSON.Node.Valid() {
		targets = append(targets, BuildTargetInfo{
			name:   "node",
			status: buildRes.Targets.Node.Status,
		})
	}
	if buildRes.Targets.JSON.Typescript.Valid() {
		targets = append(targets, BuildTargetInfo{
			name:   "typescript",
			status: buildRes.Targets.Typescript.Status,
		})
	}
	if buildRes.Targets.JSON.Python.Valid() {
		targets = append(targets, BuildTargetInfo{
			name:   "python",
			status: buildRes.Targets.Python.Status,
		})
	}
	if buildRes.Targets.JSON.Go.Valid() {
		targets = append(targets, BuildTargetInfo{
			name:   "go",
			status: buildRes.Targets.Go.Status,
		})
	}
	if buildRes.Targets.JSON.Cli.Valid() {
		targets = append(targets, BuildTargetInfo{
			name:   "cli",
			status: buildRes.Targets.Cli.Status,
		})
	}
	if buildRes.Targets.JSON.Kotlin.Valid() {
		targets = append(targets, BuildTargetInfo{
			name:   "kotlin",
			status: buildRes.Targets.Kotlin.Status,
		})
	}
	if buildRes.Targets.JSON.Java.Valid() {
		targets = append(targets, BuildTargetInfo{
			name:   "java",
			status: buildRes.Targets.Java.Status,
		})
	}
	if buildRes.Targets.JSON.Ruby.Valid() {
		targets = append(targets, BuildTargetInfo{
			name:   "ruby",
			status: buildRes.Targets.Ruby.Status,
		})
	}
	if buildRes.Targets.JSON.Terraform.Valid() {
		targets = append(targets, BuildTargetInfo{
			name:   "terraform",
			status: buildRes.Targets.Terraform.Status,
		})
	}

	return targets
}

// isTargetCompleted checks if a target is in a completed state
func isTargetCompleted(status stainless.BuildTargetStatus) bool {
	return status == "completed" || status == "postgen"
}

// waitForBuildCompletion polls a build until completion and shows progress updates
func waitForBuildCompletion(ctx context.Context, client stainless.Client, buildID string, waitGroup *Group) (*stainless.BuildObject, error) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	targetProgress := make(map[string]string)

	for {
		select {
		case <-ticker.C:
			buildRes, err := client.Builds.Get(ctx, buildID)
			if err != nil {
				waitGroup.Error("Error polling build status: %v", err)
				return nil, fmt.Errorf("build polling failed: %v", err)
			}

			targets := getBuildTargetInfo(*buildRes)
			allCompleted := true

			for _, target := range targets {
				prevStatus := targetProgress[target.name]

				if string(target.status) != prevStatus {
					targetProgress[target.name] = string(target.status)

					if isTargetCompleted(target.status) {
						waitGroup.Success("%s: %s", target.name, "completed")
					} else if target.status == "failed" {
						waitGroup.Error("%s: %s", target.name, string(target.status))
					}
				}

				if !isTargetCompleted(target.status) && target.status != "failed" {
					allCompleted = false
				}
			}

			if allCompleted && len(targets) > 0 {
				if allCompleted {
					waitGroup.Success("Build completed successfully")
					return buildRes, nil
				}
			}

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

var buildsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a new build",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name: "project",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "project",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "revision",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "revision",
			},
		},
		&jsonflag.JSONFileFlag{
			Name:    "openapi-spec",
			Aliases: []string{"oas"},
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "revision.openapi\\.yml.content",
			},
		},
		&jsonflag.JSONFileFlag{
			Name:    "stainless-config",
			Aliases: []string{"config"},
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "revision.openapi\\.stainless\\.yml.content",
			},
		},
		&cli.BoolFlag{
			Name:  "wait",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "pull",
			Usage: "Pull the build outputs after completion (only works with --wait)",
		},
		&jsonflag.JSONBoolFlag{
			Name: "allow-empty",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "allow_empty",
				SetValue: true,
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "branch",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "branch",
			},
			Required: true,
		},
		&jsonflag.JSONStringFlag{
			Name: "commit-message",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "commit_message",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "targets",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "target",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.-1",
			},
			Hidden: true,
		},
		&jsonflag.JSONStringFlag{
			Name: "+target",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.-1",
			},
		},
	},
	Action:          handleBuildsCreate,
	HideHelpCommand: true,
}

var buildsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a build by ID",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "build-id",
		},
	},
	Action:          handleBuildsRetrieve,
	HideHelpCommand: true,
}

var buildsList = cli.Command{
	Name:  "list",
	Usage: "List builds for a project",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name: "project",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "project",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "branch",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "cursor",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "cursor",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name: "limit",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "limit",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "revision",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "revision",
			},
		},
	},
	Action:          handleBuildsList,
	HideHelpCommand: true,
}

var buildsCompare = cli.Command{
	Name:  "compare",
	Usage: "Creates two builds whose outputs can be compared directly",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name: "base.revision",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "base.revision",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "base.branch",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "base.branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "base.commit_message",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "base.commit_message",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "head.revision",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "head.revision",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "head.branch",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "head.branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "head.commit_message",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "head.commit_message",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "project",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "project",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "targets",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "+target",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.-1",
			},
		},
	},
	Action:          handleBuildsCompare,
	HideHelpCommand: true,
}

func handleBuildsCreate(ctx context.Context, cmd *cli.Command) error {
	targetPaths := parseTargetPaths()

	cc, err := getAPICommandContextWithWorkspaceDefaults(cmd)
	if err != nil {
		return err
	}
	buildGroup := Info("Creating build...")
	params := stainless.BuildNewParams{}
	res, err := cc.client.Builds.New(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	buildGroup.Property("build_id", res.ID)

	if cmd.Bool("wait") {
		waitGroup := Info("Waiting for build to complete...")

		res, err = waitForBuildCompletion(context.TODO(), cc.client, res.ID, &waitGroup)
		if err != nil {
			return err
		}

		if cmd.Bool("pull") {
			pullGroup := Info("Pulling build outputs...")
			if err := pullBuildOutputs(context.TODO(), cc.client, *res, targetPaths, &pullGroup); err != nil {
				pullGroup.Error("Failed to pull outputs: %v", err)
			} else {
				pullGroup.Success("Successfully pulled all outputs")
			}
		}
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBuildsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	res, err := cc.client.Builds.Get(
		context.TODO(),
		cmd.Value("build-id").(string),
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

// pullBuildOutputs pulls the outputs for a completed build
func pullBuildOutputs(ctx context.Context, client stainless.Client, res stainless.BuildObject, targetPaths map[string]string, pullGroup *Group) error {
	// Get all targets
	allTargets := getBuildTargetInfo(res)

	// Filter to only completed targets
	var targets []string
	for _, target := range allTargets {
		if isTargetCompleted(target.status) {
			targets = append(targets, target.name)
		}
	}

	if len(targets) == 0 {
		return fmt.Errorf("no completed targets found in build %s", res.ID)
	}

	// Pull each target
	for _, target := range targets {
		// Use custom path if specified, otherwise use default
		var targetDir string
		if customPath, exists := targetPaths[target]; exists {
			targetDir = customPath
		}

		targetGroup := pullGroup.Progress("downloading %s → %s", target, targetDir)

		// Get the output details
		outputRes, err := client.Builds.TargetOutputs.Get(
			ctx,
			stainless.BuildTargetOutputGetParams{
				BuildID: res.ID,
				Target:  stainless.BuildTargetOutputGetParamsTarget(target),
				Type:    "source",
				Output:  "git",
			},
		)
		if err != nil {
			targetGroup.Error("Failed to get output details for %s: %v", target, err)
			continue
		}

		// Handle based on output type
		err = pullOutput(outputRes.Output, outputRes.URL, outputRes.Ref, targetDir, &targetGroup)
		if err != nil {
			targetGroup.Error("Failed to pull %s: %v", target, err)
			continue
		}

		// Get the appropriate success message based on output type
		if outputRes.Output == "git" {
			// Extract repository name from git URL for success message
			repoName := filepath.Base(outputRes.URL)
			if strings.HasSuffix(repoName, ".git") {
				repoName = strings.TrimSuffix(repoName, ".git")
			}
			if repoName == "" || repoName == "." || repoName == "/" {
				repoName = "repository"
			}
			targetGroup.Success("Successfully pulled to %s", repoName)
		} else {
			filename := extractFilenameFromURL(outputRes.URL)
			targetGroup.Success("Successfully downloaded %s", filename)
		}
	}

	return nil
}

// extractFilenameFromURL extracts the filename from just the URL path (without query parameters)
func extractFilenameFromURL(urlStr string) string {
	// Parse URL to remove query parameters
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		// If URL parsing fails, use the original approach
		filename := filepath.Base(urlStr)
		if filename == "." || filename == "/" || filename == "" {
			return "download"
		}
		return filename
	}

	// Extract filename from URL path (without query parameters)
	filename := filepath.Base(parsedURL.Path)
	if filename == "." || filename == "/" || filename == "" {
		return "download"
	}

	return filename
}

// extractFilename extracts the filename from a URL and HTTP response headers
func extractFilename(urlStr string, resp *http.Response) string {
	// First, try to get filename from Content-Disposition header
	if contentDisp := resp.Header.Get("Content-Disposition"); contentDisp != "" {
		// Parse Content-Disposition header for filename
		// Format: attachment; filename="example.txt" or attachment; filename=example.txt
		if strings.Contains(contentDisp, "filename=") {
			parts := strings.Split(contentDisp, "filename=")
			if len(parts) > 1 {
				filename := strings.TrimSpace(parts[1])
				// Remove quotes if present
				filename = strings.Trim(filename, `"`)
				// Remove any additional parameters after semicolon
				if idx := strings.Index(filename, ";"); idx != -1 {
					filename = filename[:idx]
				}
				filename = strings.TrimSpace(filename)
				if filename != "" {
					return filename
				}
			}
		}
	}

	// Fallback to URL path parsing
	return extractFilenameFromURL(urlStr)
}

// pullOutput handles downloading or cloning a build target output
func pullOutput(output, url, ref, targetDir string, targetGroup *Group) error {
	switch output {
	case "git":
		// Extract repository name from git URL for directory name
		// Handle formats like:
		// - https://github.com/owner/repo.git
		// - https://github.com/owner/repo
		// - git@github.com:owner/repo.git
		if targetDir == "" {
			targetDir = filepath.Base(url)
		}

		// Remove .git suffix if present
		if strings.HasSuffix(targetDir, ".git") {
			targetDir = strings.TrimSuffix(targetDir, ".git")
		}

		// Handle empty or invalid names
		if targetDir == "" || targetDir == "." || targetDir == "/" {
			targetDir = "repository"
		}

		// Remove existing directory if it exists
		if _, err := os.Stat(targetDir); err == nil {
			Info("Removing existing directory %s", targetDir)
			if err := os.RemoveAll(targetDir); err != nil {
				return fmt.Errorf("failed to remove existing directory %s: %v", targetDir, err)
			}
		}

		// Create a fresh directory for the output
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", targetDir, err)
		}
		// Clone the repository
		targetGroup.Property("cloning ref", ref)

		cmd := exec.Command("git", "clone", url, targetDir)
		var stderr bytes.Buffer
		cmd.Stdout = nil // Suppress git clone output
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git clone failed: %v\nGit error: %s", err, stderr.String())
		}

		// Checkout the specific ref
		cmd = exec.Command("git", "-C", targetDir, "checkout", ref)
		stderr.Reset()
		cmd.Stdout = nil // Suppress git checkout output
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git checkout failed: %v\nGit error: %s", err, stderr.String())
		}

	case "url":
		// Download the file directly to current directory
		targetGroup.Property("downloading url", url)

		// Download the file
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to download file: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("download failed with status: %s", resp.Status)
		}

		// Extract filename from URL and Content-Disposition header
		filename := extractFilename(url, resp)
		targetGroup.Property("downloaded", filename)

		// Create the output file in current directory
		outFile, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		defer outFile.Close()

		// Copy the response body to the output file
		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to save downloaded file: %v", err)
		}

	default:
		return fmt.Errorf("unsupported output type: %s. Supported types are 'git' and 'url'", output)
	}

	return nil
}

func handleBuildsList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.BuildListParams{}
	res, err := cc.client.Builds.List(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBuildsCompare(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.BuildCompareParams{}
	res, err := cc.client.Builds.Compare(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

// getAPICommandWithWorkspaceDefaults applies workspace defaults before initializing API command
func getAPICommandContextWithWorkspaceDefaults(cmd *cli.Command) (*apiCommandContext, error) {
	cc := getAPICommandContext(cmd)
	var config WorkspaceConfig
	found, err := config.Find()
	if err == nil && found {
		// Get the directory containing the workspace config file
		configDir := filepath.Dir(config.ConfigPath)

		if !cmd.IsSet("openapi-spec") && !cmd.IsSet("oas") && config.OpenAPISpec != "" {
			// Resolve OpenAPI spec path relative to workspace config directory
			openAPIPath := filepath.Join(configDir, config.OpenAPISpec)
			content, err := os.ReadFile(openAPIPath)
			if err != nil {
				return nil, fmt.Errorf("failed to load OpenAPI spec from workspace config: %v", err)
			}
			jsonflag.Mutate(jsonflag.Body, "revision.openapi\\.yml.content", string(content))
		}

		if !cmd.IsSet("stainless-config") && !cmd.IsSet("config") && config.StainlessConfig != "" {
			// Resolve Stainless config path relative to workspace config directory
			stainlessConfigPath := filepath.Join(configDir, config.StainlessConfig)
			content, err := os.ReadFile(stainlessConfigPath)
			if err != nil {
				return nil, fmt.Errorf("failed to load Stainless config from workspace config: %v", err)
			}
			jsonflag.Mutate(jsonflag.Body, "revision.openapi\\.stainless\\.yml.content", string(content))
		}
	}
	return cc, err
}
