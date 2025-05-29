// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

// targetInfo holds information about a build target
type targetInfo struct {
	name   string
	status stainlessv0.BuildTargetStatus
}

// getCompletedTargets extracts completed targets from a build response
func getCompletedTargets(buildRes stainlessv0.BuildObject) []targetInfo {
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
func isTargetCompleted(status stainlessv0.BuildTargetStatus) bool {
	return status == "completed" || status == "postgen"
}

var buildsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a new build",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "project",
			Action: getAPIFlagAction[string]("body", "project"),
		},
		&cli.StringFlag{
			Name:   "revision",
			Action: getAPIFlagAction[string]("body", "revision"),
		},
		&cli.StringFlag{
			Name:    "openapi-spec",
			Aliases: []string{"oas"},
			Action:  getAPIFlagFileAction("body", "revision.openapi\\.yml.content"),
		},
		&cli.StringFlag{
			Name:    "stainless-config",
			Aliases: []string{"config"},
			Action:  getAPIFlagFileAction("body", "revision.openapi\\.stainless\\.yml.content"),
		},
		&cli.BoolFlag{
			Name:   "allow-empty",
			Action: getAPIFlagAction[bool]("body", "allow_empty"),
		},
		&cli.BoolFlag{
			Name:  "wait",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "pull",
			Usage: "Pull the build outputs after completion (only works with --wait)",
		},
		&cli.StringFlag{
			Name:   "branch",
			Action: getAPIFlagAction[string]("body", "branch"),
		},
		&cli.StringFlag{
			Name:   "commit-message",
			Action: getAPIFlagAction[string]("body", "commit_message"),
		},
		&cli.StringFlag{
			Name:   "targets",
			Action: getAPIFlagAction[string]("body", "targets.#"),
		},
		&cli.StringFlag{
			Name:   "+target",
			Action: getAPIFlagAction[string]("body", "targets.-1"),
		},
	},
	Before:          initAPICommandWithWorkspaceDefaults,
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
	Before:          initAPICommand,
	Action:          handleBuildsRetrieve,
	HideHelpCommand: true,
}

var buildsList = cli.Command{
	Name:  "list",
	Usage: "List builds for a project",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "project",
			Action: getAPIFlagAction[string]("query", "project"),
		},
		&cli.StringFlag{
			Name:   "branch",
			Action: getAPIFlagAction[string]("query", "branch"),
		},
		&cli.StringFlag{
			Name:   "cursor",
			Action: getAPIFlagAction[string]("query", "cursor"),
		},
		&cli.FloatFlag{
			Name:   "limit",
			Action: getAPIFlagAction[float64]("query", "limit"),
		},
		&cli.StringFlag{
			Name:   "revision",
			Action: getAPIFlagAction[string]("query", "revision"),
		},
	},
	Before:          initAPICommand,
	Action:          handleBuildsList,
	HideHelpCommand: true,
}

var buildsCompare = cli.Command{
	Name:  "compare",
	Usage: "Creates two builds whose outputs can be compared directly",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "base.revision",
			Action: getAPIFlagAction[string]("body", "base.revision"),
		},
		&cli.StringFlag{
			Name:   "base.branch",
			Action: getAPIFlagAction[string]("body", "base.branch"),
		},
		&cli.StringFlag{
			Name:   "base.commit_message",
			Action: getAPIFlagAction[string]("body", "base.commit_message"),
		},
		&cli.StringFlag{
			Name:   "head.revision",
			Action: getAPIFlagAction[string]("body", "head.revision"),
		},
		&cli.StringFlag{
			Name:   "head.branch",
			Action: getAPIFlagAction[string]("body", "head.branch"),
		},
		&cli.StringFlag{
			Name:   "head.commit_message",
			Action: getAPIFlagAction[string]("body", "head.commit_message"),
		},
		&cli.StringFlag{
			Name:   "project",
			Action: getAPIFlagAction[string]("body", "project"),
		},
		&cli.StringFlag{
			Name:   "targets",
			Action: getAPIFlagAction[string]("body", "targets.#"),
		},
		&cli.StringFlag{
			Name:   "+target",
			Action: getAPIFlagAction[string]("body", "targets.-1"),
		},
	},
	Before:          initAPICommand,
	Action:          handleBuildsCompare,
	HideHelpCommand: true,
}

func handleBuildsCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	// Log to stderr that we're creating a build
	fmt.Fprintf(os.Stderr, "%s Creating build...\n", au.BrightCyan("✱"))
	params := stainlessv0.BuildNewParams{}
	res, err := cc.client.Builds.New(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	// Print the build ID to stderr
	fmt.Fprintf(os.Stderr, "  %s Build created: %s\n", au.BrightGreen("•"), au.Bold(res.ID))

	if cmd.Bool("wait") {
		fmt.Fprintf(os.Stderr, "%s Waiting for build to complete...\n", au.BrightCyan("✱"))

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
							fmt.Fprintf(os.Stderr, "  %s Target %s: %s\n",
								au.BrightGreen("•"),
								target.name,
								string(target.status))
							anyCompleted = true
						} else if target.status == "failed" {
							// For failures, use red text
							fmt.Fprintf(os.Stderr, "  %s Target %s: %s\n",
								au.BrightRed("•"),
								target.name,
								au.BrightRed(string(target.status)))
						}
						// Don't print in-progress status updates
					}

					if !isTargetCompleted(target.status) && target.status != "failed" {
						allCompleted = false
					}
				}

				if (allCompleted || anyCompleted) && len(targets) > 0 {
					if allCompleted {
						fmt.Fprintf(os.Stderr, "  %s Build completed successfully\n", au.BrightGreen("✱"))
						break loop
					}
				}

			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if cmd.Bool("pull") {
			fmt.Fprintf(os.Stderr, "%s Pulling build outputs...\n", au.BrightCyan("✱"))
			if err := pullBuildOutputs(context.TODO(), cc.client, *res); err != nil {
				fmt.Fprintf(os.Stderr, "%s Failed to pull outputs: %v\n", au.BrightRed("✱"), err)
			} else {
				fmt.Fprintf(os.Stderr, "%s Successfully pulled all outputs\n", au.BrightGreen("✱"))
			}
		}
	}

	// Print the actual JSON response to stdout for piping
	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBuildsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	res, err := cc.client.Builds.Get(
		context.TODO(),
		cmd.Value("build-id").(string),
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

// pullBuildOutputs pulls the outputs for a completed build
func pullBuildOutputs(ctx context.Context, client stainlessv0.Client, res stainlessv0.BuildObject) error {
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
		targetDir := fmt.Sprintf("%s-%s", res.Project, target)

		fmt.Fprintf(os.Stderr, "%s [%d/%d] Pulling %s → %s\n",
			au.BrightCyan("✱"), i+1, len(targets), au.Bold(target), au.Cyan(targetDir))

		// Get the output details
		outputRes, err := client.BuildTargetOutputs.Get(
			ctx,
			stainlessv0.BuildTargetOutputGetParams{
				BuildID: res.ID,
				Target:  stainlessv0.BuildTargetOutputGetParamsTarget(target),
				Type:    "source",
				Output:  "git",
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s Failed to get output details for %s: %v\n",
				au.BrightRed("✱"), target, err)
			continue
		}

		// Handle based on output type
		err = pullOutput(outputRes.Output, outputRes.URL, outputRes.Ref, targetDir, target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s Failed to pull %s: %v\n",
				au.BrightRed("✱"), target, err)
			continue
		}

		fmt.Fprintf(os.Stderr, "  %s Successfully pulled to %s\n",
			au.BrightBlack("•"), au.Cyan(targetDir))

		if i < len(targets)-1 {
			fmt.Fprintf(os.Stderr, "\n")
		}
	}

	return nil
}

// pullOutput handles downloading or cloning a build target output
func pullOutput(output, url, ref, targetDir, target string) error {
	// Remove existing directory if it exists
	if _, err := os.Stat(targetDir); err == nil {
		fmt.Fprintf(os.Stderr, "  %s Removing existing directory %s\n", au.BrightBlack("•"), targetDir)
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("failed to remove existing directory %s: %v", targetDir, err)
		}
	}

	// Create a fresh directory for the output
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", targetDir, err)
	}

	switch output {
	case "git":
		// Clone the repository
		fmt.Fprintf(os.Stderr, "  %s Cloning repository\n", au.BrightBlack("•"))
		fmt.Fprintf(os.Stderr, "  %s Checking out ref %s\n", au.BrightBlack("•"), au.Bold(ref))

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
		// Download the tar file
		fmt.Fprintf(os.Stderr, "  %s Downloading archive %s\n", au.BrightBlack("•"), au.Underline(url))
		fmt.Fprintf(os.Stderr, "  %s Extracting to %s\n", au.BrightBlack("•"), targetDir)

		// Create a temporary file for the tar download
		tmpFile, err := os.CreateTemp("", "stainless-*.tar.gz")
		if err != nil {
			return fmt.Errorf("failed to create temporary file: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		// Download the file
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to download file: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("download failed with status: %s", resp.Status)
		}

		// Copy the response body to the temporary file
		_, err = io.Copy(tmpFile, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to save downloaded file: %v", err)
		}
		tmpFile.Close()

		// Extract the tar file
		cmd := exec.Command("tar", "-xzf", tmpFile.Name(), "-C", targetDir)
		cmd.Stdout = nil // Suppress tar output
		cmd.Stderr = nil
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("tar extraction failed: %v", err)
		}

	default:
		return fmt.Errorf("unsupported output type: %s. Supported types are 'git' and 'url'", output)
	}

	return nil
}

func handleBuildsList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := stainlessv0.BuildListParams{}
	res, err := cc.client.Builds.List(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBuildsCompare(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := stainlessv0.BuildCompareParams{}
	res, err := cc.client.Builds.Compare(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

// initAPICommandWithWorkspaceDefaults applies workspace defaults before initializing API command
func initAPICommandWithWorkspaceDefaults(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	cc, err := initAPICommand(ctx, cmd)
	if err != nil {
		return nil, err
	}
	config, _, err := FindWorkspaceConfig()
	if err == nil && config != nil {
		if !cmd.IsSet("openapi-spec") && !cmd.IsSet("oas") && config.OpenAPISpec != "" {
			fileAction := getAPIFlagFileAction("body", "revision.openapi\\.yml.content")
			if err := fileAction(cc, cmd, config.OpenAPISpec); err != nil {
				return nil, fmt.Errorf("failed to load OpenAPI spec from workspace config: %v", err)
			}
		}

		if !cmd.IsSet("stainless-config") && !cmd.IsSet("config") && config.StainlessConfig != "" {
			fileAction := getAPIFlagFileAction("body", "revision.openapi\\.stainless\\.yml.content")
			if err := fileAction(cc, cmd, config.StainlessConfig); err != nil {
				return nil, fmt.Errorf("failed to load Stainless config from workspace config: %v", err)
			}
		}
	}
	return cc, err
}
