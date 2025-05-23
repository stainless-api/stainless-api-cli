// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
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
			Name: "wait",
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
	Before:          initAPICommand,
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
	// Log to stderr that we're creating a build (using white text)
	fmt.Fprintf(os.Stderr, "Creating build...\n")
	params := stainlessv0.BuildNewParams{}
	res, err := cc.client.Builds.New(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithRequestBody("application/json", cc.body),
	)
	if err != nil {
		return err
	}

	// Print the build ID to stderr (using white text)
	fmt.Fprintf(os.Stderr, "Build created: %s\n", res.ID)

	if cmd.Bool("wait") {
		fmt.Fprintf(os.Stderr, "Waiting for build to complete...\n")

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
							fmt.Fprintf(os.Stderr, "%s Target %s: %s\n",
								au.BrightGreen("✓"),
								target.name,
								string(target.status))
							anyCompleted = true
						} else if target.status == "failed" {
							// For failures, use red text
							fmt.Fprintf(os.Stderr, "Target %s: %s\n",
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
						fmt.Fprintf(os.Stderr, "%s Build completed successfully\n", au.BrightGreen("✓"))
						break loop
					}
				}

			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if cmd.Bool("pull") {
			fmt.Fprintf(os.Stderr, "Pulling build outputs...\n")
			if err := pullBuildOutputs(context.TODO(), cc.client, *res); err != nil {
				fmt.Fprintf(os.Stderr, "%s Failed to pull outputs: %v\n", au.BrightRed("!"), err)
			} else {
				fmt.Fprintf(os.Stderr, "%s Successfully pulled all outputs\n", au.BrightGreen("✓"))
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

	fmt.Fprintf(os.Stderr, "Found completed targets: %v\n", targets)

	// Pull each target
	for _, target := range targets {
		fmt.Fprintf(os.Stderr, "Pulling output for target: %s\n", target)

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
			fmt.Fprintf(os.Stderr, "%s Error getting output for target %s: %v\n", au.BrightRed("!"), target, err)
			continue
		}

		targetDir := fmt.Sprintf("%s-sdk", target)

		// Handle based on output type
		err = pullOutput(outputRes.Output, outputRes.URL, outputRes.Ref, targetDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s Error pulling output for target %s: %v\n", au.BrightRed("!"), target, err)
			continue
		} else {
			fmt.Fprintf(os.Stderr, "%s Successfully pulled output for target %s\n", au.BrightGreen("✓"), target)
		}
	}

	return nil
}

// pullOutput handles downloading or cloning a build target output
func pullOutput(output, url, ref, targetDir string) error {
	// Create a directory for the output
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", targetDir, err)
	}

	switch output {
	case "git":
		// Clone the repository
		fmt.Fprintf(os.Stderr, "Cloning repository %s (ref: %s) to %s\n", url, ref, targetDir)

		cmd := exec.Command("git", "clone", url, targetDir)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git clone failed: %v", err)
		}

		// Checkout the specific ref
		cmd = exec.Command("git", "-C", targetDir, "checkout", ref)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git checkout failed: %v", err)
		}

		fmt.Fprintf(os.Stderr, "%s Successfully cloned repository to %s\n", au.BrightGreen("✓"), targetDir)

	case "url":
		// Download the tar file
		fmt.Fprintf(os.Stderr, "Downloading from %s to %s\n", url, targetDir)

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
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("tar extraction failed: %v", err)
		}

		fmt.Fprintf(os.Stderr, "%s Successfully downloaded and extracted to %s\n", au.BrightGreen("✓"), targetDir)

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
		option.WithRequestBody("application/json", cc.body),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
