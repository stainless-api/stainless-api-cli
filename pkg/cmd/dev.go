package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/stainless-api/stainless-api-cli/pkg/components/build"
	"github.com/stainless-api/stainless-api-cli/pkg/components/dev"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-cli/pkg/git"
	"github.com/stainless-api/stainless-api-cli/pkg/workspace"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/stainless-api/stainless-api-go/shared"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var devCommand = cli.Command{
	Name:    "preview",
	Aliases: []string{"dev"},
	Usage:   "Development mode with interactive build monitoring",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "project",
			Aliases: []string{"p"},
			Usage:   "Project name to use for the build",
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
		&cli.StringFlag{
			Name:  "base",
			Value: "HEAD",
			Usage: "Git ref to use as the base revision for comparison",
		},
		&cli.BoolFlag{
			Name:    "watch",
			Aliases: []string{"w"},
			Value:   false,
			Usage:   "Run in 'watch' mode to loop and rebuild when files change.",
		},
	},
	Action: runPreview,
}

func runPreview(ctx context.Context, cmd *cli.Command) error {
	if cmd.Bool("watch") {
		fmt.Print("\033[2J\033[H")
		os.Stdout.Sync()
	}

	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)

	wc := getWorkspace(ctx)

	for {
		if err := runDevBuild(ctx, client, wc, cmd); err != nil {
			if errors.Is(err, build.ErrUserCancelled) {
				return nil
			}
			return err
		}

		if !cmd.Bool("watch") {
			break
		}

		fmt.Print("\nRebuilding...\n\n\033[2J\033[H")
		os.Stdout.Sync()
	}
	return nil
}

// generateEphemeralBranches creates a paired set of ephemeral branch names
// for a compare build: one for the base and one for the head.
func generateEphemeralBranches(branchName string) (baseBranch, headBranch string) {
	now := time.Now()
	randomBytes := make([]byte, 3)
	rand.Read(randomBytes)
	entropy := fmt.Sprintf("%d%02d%02d-%s", now.Year(), now.Month(), now.Day(), base64.RawURLEncoding.EncodeToString(randomBytes))
	baseBranch = fmt.Sprintf("ephemeral-base-%s/%s", entropy, branchName)
	headBranch = fmt.Sprintf("ephemeral-%s/%s", entropy, branchName)
	return
}

// readFileInputMap reads files from disk and returns them as a file input map
// suitable for a build revision.
func readFileInputMap(oasPath, configPath string) (map[string]shared.FileInputUnionParam, error) {
	files := make(map[string]shared.FileInputUnionParam)

	if oasPath != "" {
		content, err := os.ReadFile(oasPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read openapi-spec file: %v", err)
		}
		files["openapi"+path.Ext(oasPath)] = shared.FileInputParamOfFileInputContent(string(content))
	}

	if configPath != "" {
		content, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read stainless-config file: %v", err)
		}
		files["stainless"+path.Ext(configPath)] = shared.FileInputParamOfFileInputContent(string(content))
	}

	return files, nil
}

// gitShowFileInputMap tries to read files at a given git ref and returns them
// as a file input map. Returns nil (not error) if any file can't be read from git.
func gitShowFileInputMap(repoDir, ref, oasPath, configPath string) map[string]shared.FileInputUnionParam {
	files := make(map[string]shared.FileInputUnionParam)

	if oasPath != "" {
		relPath, err := filepath.Rel(repoDir, oasPath)
		if err != nil {
			return nil
		}
		content, err := git.Show(repoDir, ref, relPath)
		if err != nil {
			return nil
		}
		files["openapi"+path.Ext(oasPath)] = shared.FileInputParamOfFileInputContent(string(content))
	}

	if configPath != "" {
		relPath, err := filepath.Rel(repoDir, configPath)
		if err != nil {
			return nil
		}
		content, err := git.Show(repoDir, ref, relPath)
		if err != nil {
			return nil
		}
		files["stainless"+path.Ext(configPath)] = shared.FileInputParamOfFileInputContent(string(content))
	}

	return files
}

// gitRepoRoot returns the top-level directory of the git repo, or "" if not in a repo.
func gitRepoRoot(dir string) string {
	sha, err := git.RevParse(dir, "--show-toplevel")
	if err != nil {
		return ""
	}
	return sha
}

func runDevBuild(ctx context.Context, client stainless.Client, wc workspace.Config, cmd *cli.Command) error {
	projectName := cmd.String("project")
	oasPath := cmd.String("openapi-spec")
	configPath := cmd.String("stainless-config")

	// Determine git state and branch name
	branchName := "main"
	repoDir := gitRepoRoot(".")
	inGitRepo := repoDir != ""

	if inGitRepo {
		if b, err := git.CurrentBranch(repoDir); err == nil {
			branchName = b
		}
	}
	baseBranch, headBranch := generateEphemeralBranches(branchName)

	// Build head revision from current files on disk
	headFiles, err := readFileInputMap(oasPath, configPath)
	if err != nil {
		return err
	}

	// Build base revision: try git show at --base ref, otherwise fall back to "main"
	var baseRevision stainless.BuildCompareParamsBaseRevisionUnion

	baseRef := cmd.String("base")
	if inGitRepo && oasPath != "" {
		files := gitShowFileInputMap(repoDir, baseRef, oasPath, configPath)
		if len(files) > 0 {
			baseRevision.OfFileInputMap = files
		} else {
			baseRevision.OfString = stainless.String("main")
		}
	} else {
		baseRevision.OfString = stainless.String("main")
	}

	compareReq := stainless.BuildCompareParams{
		Project: stainless.String(projectName),
		Base: stainless.BuildCompareParamsBase{
			Branch:   baseBranch,
			Revision: baseRevision,
		},
		Head: stainless.BuildCompareParamsHead{
			Branch: headBranch,
			Revision: stainless.BuildCompareParamsHeadRevisionUnion{
				OfFileInputMap: headFiles,
			},
		},
	}

	downloads := make(map[stainless.Target]string)
	for targetName, targetConfig := range wc.Targets {
		downloads[stainless.Target(targetName)] = targetConfig.OutputPath
	}

	model := dev.NewModel(dev.ModelConfig{
		Client: client,
		Ctx:    ctx,
		Branch: headBranch,
		Start: func() (*stainless.Build, error) {
			options := []option.RequestOption{}
			if cmd.Bool("debug") {
				options = append(options, debugMiddlewareOption)
			}
			resp, err := client.Builds.Compare(
				ctx,
				compareReq,
				options...,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create compare build: %v", err)
			}
			return &resp.Head, nil
		},
		DownloadPaths: downloads,
		Watch:         cmd.Bool("watch"),
	})
	model.Diagnostics.WorkspaceConfig = wc

	p := console.NewProgram(model)
	finalModel, err := p.Run()

	if err != nil {
		return fmt.Errorf("failed to run TUI: %v", err)
	}
	if buildModel, ok := finalModel.(dev.Model); ok {
		return buildModel.Err
	}
	return nil
}

type GenerateSpecParams struct {
	Project string `json:"project"`
	Source  struct {
		Type            string `json:"type"`
		OpenAPISpec     string `json:"openapi_spec"`
		StainlessConfig string `json:"stainless_config"`
	} `json:"source"`
}

func getDiagnostics(ctx context.Context, cmd *cli.Command, client stainless.Client, wc workspace.Config) ([]stainless.BuildDiagnostic, error) {
	var specParams GenerateSpecParams
	if cmd.IsSet("project") {
		specParams.Project = cmd.String("project")
	} else {
		specParams.Project = wc.Project
	}
	specParams.Source.Type = "upload"

	configPath := wc.StainlessConfig
	if cmd.IsSet("stainless-config") {
		configPath = cmd.String("stainless-config")
	} else if configPath == "" {
		return nil, fmt.Errorf("You must provide a stainless configuration file with `--config /path/to/stainless.yml` or run this command from an initialized workspace.")
	}

	stainlessConfig, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Could not read your stainless configuration file:\n%w", err)
	}
	specParams.Source.StainlessConfig = string(stainlessConfig)

	oasPath := wc.OpenAPISpec
	if cmd.IsSet("openapi-spec") {
		oasPath = cmd.String("openapi-spec")
	} else if oasPath == "" {
		return nil, fmt.Errorf("You must provide an OpenAPI specification with `--oas /path/to/openapi.json` or run this command from an initialized workspace.")
	}

	openAPISpec, err := os.ReadFile(oasPath)
	if err != nil {
		return nil, fmt.Errorf("Could not read your stainless configuration file:\n%w", err)
	}
	specParams.Source.OpenAPISpec = string(openAPISpec)

	options := []option.RequestOption{}
	if cmd.Bool("debug") {
		options = append(options, debugMiddlewareOption)
	}
	var result []byte
	err = client.Post(
		ctx,
		"api/generate/spec",
		specParams,
		&result,
		options...,
	)
	if err != nil {
		return nil, err
	}

	transform := "spec.diagnostics.@values.@flatten.#(ignored==false)#"
	jsonObj := gjson.Parse(string(result)).Get(transform)
	var diagnostics []stainless.BuildDiagnostic
	json.Unmarshal([]byte(jsonObj.Raw), &diagnostics)
	return diagnostics, nil
}

func hasBlockingDiagnostic(diagnostics []stainless.BuildDiagnostic) bool {
	for _, d := range diagnostics {
		if !d.Ignored {
			switch d.Level {
			case stainless.BuildDiagnosticLevelFatal:
			case stainless.BuildDiagnosticLevelError:
			case stainless.BuildDiagnosticLevelWarning:
				return true
			case stainless.BuildDiagnosticLevelNote:
				continue
			}
		}
	}
	return false
}
