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

// parseTargetPaths processes target flags to extract target:path syntax with workspace config
// Returns a map of target names to their custom paths
func parseTargetPaths(workspaceConfig WorkspaceConfig) map[string]string {
	targetPaths := make(map[string]string)

	// Check workspace configuration for target paths if loaded
	for targetName, targetConfig := range workspaceConfig.Targets {
		if targetConfig.OutputPath != "" {
			targetPaths[targetName] = targetConfig.OutputPath
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
func waitForBuildCompletion(ctx context.Context, client stainless.Client, build *stainless.BuildObject, waitGroup *Group) (*stainless.BuildObject, error) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	targetProgress := make(map[string]string)

	for {
		targets := getBuildTargetInfo(*build)
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
				return build, nil
			}
		}

		select {
		case <-ticker.C:
			var err error
			build, err = client.Builds.Get(ctx, build.ID)
			if err != nil {
				waitGroup.Error("Error polling build status: %v", err)
				return nil, fmt.Errorf("build polling failed: %v", err)
			}

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

var buildsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a build, on top of a project branch, against a given input revision.",
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
	Usage: "Retrieve a build by its ID.",
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
	Usage: "List user-triggered builds for a given project.",
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
	Usage: "Create two builds whose outputs can be directly compared with each other.",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name: "base.branch",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "base.branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "base.revision",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "base.revision",
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
			Name: "head.branch",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "head.branch",
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
	cc := getAPICommandContext(cmd)

	// Handle file flags by reading files and mutating JSON body
	if err := applyFileFlag(cmd, "openapi-spec", "revision.openapi\\.yml.content"); err != nil {
		return err
	}
	if err := applyFileFlag(cmd, "stainless-config", "revision.openapi\\.stainless\\.yml.content"); err != nil {
		return err
	}

	// Parse target paths using cached workspace config
	targetPaths := parseTargetPaths(cc.workspaceConfig)
	buildGroup := Info("Creating build...")
	params := stainless.BuildNewParams{}
	var res []byte
	_, err := cc.client.Builds.New(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	buildGroup.Property("build_id", res.ID)

	if cmd.Bool("wait") {
		waitGroup := Info("Waiting for latest build to complete...")

		res, err = waitForBuildCompletion(context.TODO(), cc.client, res, &waitGroup)
		if err != nil {
			return err
		}

		// Pull if explicitly set via --pull flag, or if workspace has configured targets and --pull wasn't explicitly set to false
		shouldPull := cmd.Bool("pull") || (cc.HasWorkspaceTargets() && !cmd.IsSet("pull"))

		if shouldPull {
			pullGroup := Info("Downloading build outputs...")
			if err := pullBuildOutputs(context.TODO(), cc.client, *res, targetPaths, &pullGroup); err != nil {
				pullGroup.Error("Failed to download outputs: %v", err)
			} else {
				pullGroup.Success("Successfully downloaded all outputs")
			}
		}
	}

	format := cmd.Root().String("format")
	return ShowJSON("builds create", string(res), format)
}

func handleBuildsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	var res []byte
	_, err := cc.client.Builds.Get(
		context.TODO(),
		cmd.Value("build-id").(string),
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("builds retrieve", string(res), format)
}

// pullBuildOutputs pulls the outputs for a completed build
func pullBuildOutputs(ctx context.Context, client stainless.Client, res stainless.BuildObject, targetPaths map[string]string, pullGroup *Group) error {
	// Get all targets
	allTargets := getBuildTargetInfo(res)

	// Filter to only completed targets without fatal conclusions
	var targets []string
	for _, target := range allTargets {
		if isTargetCompleted(target.status) && !hasFailedCommitStep(res, stainless.Target(target.name)) {
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
			targetGroup.Success("Successfully downloaded")
		} else {
			filename := extractFilenameFromURL(outputRes.URL)
			targetGroup.Success("Successfully downloaded %s", filename)
		}
	}

	return nil
}

// hasFailedCommitStep checks if a target has a fatal commit conclusion
func hasFailedCommitStep(build stainless.BuildObject, target stainless.Target) bool {
	buildObj := NewBuildObject(&build)
	buildTarget := buildObj.BuildTarget(target)
	if buildTarget == nil {
		return false
	}

	status, _, conclusion := buildTarget.StepInfo("commit")
	if status == "completed" && conclusion == "fatal" {
		return true
	}

	return false
}

// stripHTTPAuth removes HTTP authentication credentials from a URL for display purposes
func stripHTTPAuth(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	// Remove user info (username:password)
	parsedURL.User = nil
	return parsedURL.String()
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
		targetDir = strings.TrimSuffix(targetDir, ".git")

		// Handle empty or invalid names
		if targetDir == "" || targetDir == "." || targetDir == "/" {
			targetDir = "repository"
		}

		// Check if directory exists
		if _, err := os.Stat(targetDir); err == nil {
			// Check if it's a git directory
			if _, err := os.Stat(filepath.Join(targetDir, ".git")); err != nil {
				// Not a git directory, return error
				return fmt.Errorf("directory %s already exists and is not a git repository", targetDir)
			}
		} else {
			// Directory doesn't exist, create it and initialize git repo
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", targetDir, err)
			}

			// Initialize git repository
			cmd := exec.Command("git", "-C", targetDir, "init")
			var stderr bytes.Buffer
			cmd.Stdout = nil
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("git init failed: %v\nGit error: %s", err, stderr.String())
			}
		}

		{
			// Check if origin remote exists, add it if not present
			cmd := exec.Command("git", "-C", targetDir, "remote", "get-url", "origin")
			var stderr bytes.Buffer
			cmd.Stdout = nil
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				// Origin doesn't exist, add it with stripped auth
				targetGroup.Property("adding remote origin", stripHTTPAuth(url))
				addCmd := exec.Command("git", "-C", targetDir, "remote", "add", "origin", stripHTTPAuth(url))
				var addStderr bytes.Buffer
				addCmd.Stdout = nil
				addCmd.Stderr = &addStderr
				if err := addCmd.Run(); err != nil {
					return fmt.Errorf("git remote add failed: %v\nGit error: %s", err, addStderr.String())
				}
			}

			targetGroup.Property("fetching from", stripHTTPAuth(url))
			cmd = exec.Command("git", "-C", targetDir, "fetch", url, ref)
			cmd.Stdout = nil
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("git fetch failed: %v\nGit error: %s", err, stderr.String())
			}
		}

		// Checkout the specific ref
		{
			targetGroup.Property("checking out ref", ref)
			cmd := exec.Command("git", "-C", targetDir, "checkout", ref)
			var stderr bytes.Buffer
			cmd.Stdout = nil // Suppress git checkout output
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("git checkout failed: %v\nGit error: %s", err, stderr.String())
			}
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
	var res []byte
	_, err := cc.client.Builds.List(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("builds list", string(res), format)
}

func handleBuildsCompare(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.BuildCompareParams{}
	var res []byte
	_, err := cc.client.Builds.Compare(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("builds compare", string(res), format)
}
