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
	"github.com/urfave/cli/v3"
)

// targetInfo holds information about a build target
type targetInfo struct {
	name   string
	status stainless.BuildTargetStatus
}

// getCompletedTargets extracts completed targets from a build response
func getCompletedTargets(buildRes stainless.BuildObject) []targetInfo {
	targets := []targetInfo{}

	// Check each target and add it to the list if it's completed or in postgen
	if buildRes.Targets.JSON.Node.Valid() {
		targets = append(targets, targetInfo{
			name:   "node",
			status: buildRes.Targets.Node.Status,
		})
	}
	if buildRes.Targets.JSON.Typescript.Valid() {
		targets = append(targets, targetInfo{
			name:   "typescript",
			status: buildRes.Targets.Typescript.Status,
		})
	}
	if buildRes.Targets.JSON.Python.Valid() {
		targets = append(targets, targetInfo{
			name:   "python",
			status: buildRes.Targets.Python.Status,
		})
	}
	if buildRes.Targets.JSON.Go.Valid() {
		targets = append(targets, targetInfo{
			name:   "go",
			status: buildRes.Targets.Go.Status,
		})
	}
	if buildRes.Targets.JSON.Cli.Valid() {
		targets = append(targets, targetInfo{
			name:   "cli",
			status: buildRes.Targets.Cli.Status,
		})
	}
	if buildRes.Targets.JSON.Kotlin.Valid() {
		targets = append(targets, targetInfo{
			name:   "kotlin",
			status: buildRes.Targets.Kotlin.Status,
		})
	}
	if buildRes.Targets.JSON.Java.Valid() {
		targets = append(targets, targetInfo{
			name:   "java",
			status: buildRes.Targets.Java.Status,
		})
	}
	if buildRes.Targets.JSON.Ruby.Valid() {
		targets = append(targets, targetInfo{
			name:   "ruby",
			status: buildRes.Targets.Ruby.Status,
		})
	}
	if buildRes.Targets.JSON.Terraform.Valid() {
		targets = append(targets, targetInfo{
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
	cc, err := getAPICommandContextWithWorkspaceDefaults(cmd)
	if err != nil {
		return err
	}
	Progress("Creating build...")
	params := stainless.BuildNewParams{}
	res, err := cc.client.Builds.New(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	// Print the build ID to stderr
	buildGroup := Success("Build created")
	buildGroup.Property("build_id", res.ID)

	if cmd.Bool("wait") {
		waitGroup := Progress("Waiting for build to complete...")

		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		// Track progress for each target
		targetProgress := make(map[string]string)

	loop:
		for {
			select {
			case <-ticker.C:
				res, err = cc.client.Builds.Get(
					context.TODO(),
					res.ID,
				)
				if err != nil {
					return fmt.Errorf("error polling build status: %v", err)
				}

				targets := getCompletedTargets(*res)

				// Update progress for each target
				allCompleted := true
				anyCompleted := false

				// Print status for each target
				for _, target := range targets {
					prevStatus := targetProgress[target.name]

					if string(target.status) != prevStatus {
						// Status changed, update it
						targetProgress[target.name] = string(target.status)

						// Only print completed statuses with a green checkmark
						if isTargetCompleted(target.status) {
							waitGroup.Success("Target %s: %s", target.name, string(target.status))
							anyCompleted = true
						} else if target.status == "failed" {
							// For failures, use red text
							waitGroup.Error("Target %s: %s", target.name, string(target.status))
						}
						// Don't print in-progress status updates
					}

					if !isTargetCompleted(target.status) && target.status != "failed" {
						allCompleted = false
					}
				}

				if (allCompleted || anyCompleted) && len(targets) > 0 {
					if allCompleted {
						waitGroup.Success("Build completed successfully")
						break loop
					}
				}

			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if cmd.Bool("pull") {
			pullGroup := Progress("Pulling build outputs...")
			if err := pullBuildOutputs(context.TODO(), cc.client, *res); err != nil {
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
func pullBuildOutputs(ctx context.Context, client stainless.Client, res stainless.BuildObject) error {
	// Get all targets
	allTargets := getCompletedTargets(res)

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
	for i, target := range targets {
		targetGroup := Progress("[%d/%d] Pulling %s", i+1, len(targets), target)

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
		err = pullOutput(outputRes.Output, outputRes.URL, outputRes.Ref)
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

		if i < len(targets)-1 {
			fmt.Fprintf(os.Stderr, "\n")
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
func pullOutput(output, url, ref string) error {
	switch output {
	case "git":
		// Extract repository name from git URL for directory name
		// Handle formats like:
		// - https://github.com/owner/repo.git
		// - https://github.com/owner/repo
		// - git@github.com:owner/repo.git
		targetDir := filepath.Base(url)

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
		gitGroup := Info("Cloning repository")
		gitGroup.Property("ref", ref)

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
		downloadGroup := Info("Downloading file")
		downloadGroup.Property("url", url)

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
		downloadGroup.Property("filename", filename)

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
