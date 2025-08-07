package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var ErrUserCancelled = errors.New("user cancelled")

// BuildModel represents the bubbletea model for build monitoring
type BuildModel struct {
	start       func() (*stainless.BuildObject, error)
	started     time.Time
	ended       *time.Time
	build       *stainless.BuildObject
	branch      string
	diagnostics []stainless.BuildDiagnosticListResponse
	downloads   map[stainless.Target]struct {
		status string
		path   string
	}
	view string

	cc          *apiCommandContext
	ctx         context.Context
	err         error
	isCompleted bool
}

type tickMsg time.Time
type fetchBuildMsg *stainless.BuildObject
type fetchDiagnosticsMsg []stainless.BuildDiagnosticListResponse
type errorMsg error
type downloadMsg stainless.Target
type triggerNewBuildMsg struct{}

func NewBuildModel(cc *apiCommandContext, ctx context.Context, branch string, fn func() (*stainless.BuildObject, error)) BuildModel {
	return BuildModel{
		start:   fn,
		started: time.Now(),
		cc:      cc,
		ctx:     ctx,
		branch:  branch,
	}
}

func (m BuildModel) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
		func() tea.Msg {
			build, err := m.start()
			if err != nil {
				return errorMsg(err)
			}
			return fetchBuildMsg(build)
		},
	)
}

func (m BuildModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.err = ErrUserCancelled
			cmds = append(cmds, tea.Quit)
		case "enter":
			cmds = append(cmds, tea.Quit)
		}

	case downloadMsg:
		download := m.downloads[stainless.Target(msg)]
		download.status = "completed"
		m.downloads[stainless.Target(msg)] = download

	case tickMsg:
		if m.build != nil {
			cmds = append(cmds, m.fetchBuildStatus())
		}
		m.getBuildDuration()
		cmds = append(cmds, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}))

	case fetchBuildMsg:
		if m.build == nil {
			m.build = msg
			m.downloads = make(map[stainless.Target]struct {
				status string
				path   string
			})
			for targetName, targetConfig := range m.cc.workspaceConfig.Targets {
				m.downloads[stainless.Target(targetName)] = struct {
					status string
					path   string
				}{
					status: "not started",
					path:   targetConfig.OutputPath,
				}
			}
			cmds = append(cmds, m.updateView("header"))
		}

		m.build = msg
		buildObj := NewBuildObject(m.build)
		if !m.isCompleted {
			// Check if all commit steps are completed
			allCommitsCompleted := true
			for _, target := range buildObj.Languages() {
				buildTarget := buildObj.BuildTarget(target)
				if buildTarget != nil && !buildTarget.IsCommitCompleted() {
					allCommitsCompleted = false
					break
				}
			}
			if allCommitsCompleted {
				m.isCompleted = true
				cmds = append(cmds, m.fetchDiagnostics())
			}
		}
		languages := buildObj.Languages()
		for _, target := range languages {
			buildTarget := buildObj.BuildTarget(target)
			if buildTarget == nil {
				continue
			}
			status, _, conclusion := buildTarget.StepInfo("commit")
			if status == "completed" && conclusion != "fatal" {
				if download, ok := m.downloads[target]; ok && download.status == "not started" {
					download.status = "started"
					cmds = append(cmds, m.downloadTarget(target))
					m.downloads[target] = download
				}
			}
		}

	case fetchDiagnosticsMsg:
		if m.diagnostics == nil {
			m.diagnostics = msg
			cmds = append(cmds, m.updateView("diagnostics"))
		}

	case errorMsg:
		m.err = msg
		cmds = append(cmds, tea.Quit)
	}
	return m, tea.Sequence(cmds...)
}

func (m BuildModel) downloadTarget(target stainless.Target) tea.Cmd {
	return func() tea.Msg {
		if m.build == nil {
			return errorMsg(fmt.Errorf("no current build to download target from"))
		}
		params := stainless.BuildTargetOutputGetParams{
			BuildID: m.build.ID,
			Target:  stainless.BuildTargetOutputGetParamsTarget(target),
			Type:    "source",
			Output:  "git",
		}
		outputRes, err := m.cc.client.Builds.TargetOutputs.Get(
			context.TODO(),
			params,
		)
		if err != nil {
			return errorMsg(err)
		}
		err = pullOutput(outputRes.Output, outputRes.URL, outputRes.Ref, m.downloads[target].path, &Group{silent: true})
		if err != nil {
			return errorMsg(err)
		}
		return downloadMsg(target)
	}
}

func (m BuildModel) fetchBuildStatus() tea.Cmd {
	return func() tea.Msg {
		if m.build == nil {
			return errorMsg(fmt.Errorf("no current build to fetch status for"))
		}
		build, err := m.cc.client.Builds.Get(m.ctx, m.build.ID)
		if err != nil {
			return errorMsg(fmt.Errorf("failed to get build status: %v", err))
		}
		return fetchBuildMsg(build)
	}
}

func (m BuildModel) fetchDiagnostics() tea.Cmd {
	return func() tea.Msg {
		if m.build == nil {
			return errorMsg(fmt.Errorf("no current build to fetch diagnostics for"))
		}
		diags := []stainless.BuildDiagnosticListResponse{}
		diagnostics := m.cc.client.Builds.Diagnostics.ListAutoPaging(m.ctx, m.build.ID, stainless.BuildDiagnosticListParams{
			Limit: stainless.Float(100),
		})
		if diagnostics.Next() {
			diags = append(diags, diagnostics.Current())
		}
		return fetchDiagnosticsMsg(diags)
	}
}

func (m *BuildModel) getBuildDuration() time.Duration {
	if m.build == nil {
		return time.Since(m.started)
	}

	buildObj := NewBuildObject(m.build)
	if buildObj.IsCompleted() {
		if m.ended == nil {
			now := time.Now()
			m.ended = &now
		}
		return m.ended.Sub(m.started)
	}

	return time.Since(m.started)
}

var devCommand = cli.Command{
	Name:  "dev",
	Usage: "Development mode with interactive build monitoring",
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
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return runDevMode(ctx, cmd)
	},
}

func runDevMode(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)

	gitUser, err := getGitUsername()
	if err != nil {
		return fmt.Errorf("failed to get git username: %v", err)
	}

	var selectedBranch string

	now := time.Now()
	randomBytes := make([]byte, 3)
	rand.Read(randomBytes)
	randomSuffix := base64.RawURLEncoding.EncodeToString(randomBytes)
	randomBranch := fmt.Sprintf("%s/%d%02d%02d-%s", gitUser, now.Year(), now.Month(), now.Day(), randomSuffix)

	branchOptions := []huh.Option[string]{}
	if currentBranch, err := getCurrentGitBranch(); err == nil && currentBranch != "main" && currentBranch != "master" {
		branchOptions = append(branchOptions,
			huh.NewOption(currentBranch, currentBranch),
		)
	}
	branchOptions = append(branchOptions,
		huh.NewOption(fmt.Sprintf("%s/dev", gitUser), fmt.Sprintf("%s/dev", gitUser)),
		huh.NewOption(fmt.Sprintf("%s/<random>", gitUser), randomBranch),
	)

	branchForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("branch").
				Description("Select a Stainless project branch to use for development").
				Options(branchOptions...).
				Value(&selectedBranch),
		),
	).WithTheme(GetFormTheme(0))

	if err := branchForm.Run(); err != nil {
		return fmt.Errorf("branch selection failed: %v", err)
	}

	Property("branch", selectedBranch)

	// Phase 2: Language selection
	var selectedTargets []string

	// Use cached workspace config for intelligent defaults
	config := cc.workspaceConfig

	targetInfo := getAvailableTargetInfo(ctx, cc.client, cmd.String("project"), config)
	targetOptions := targetInfoToOptions(targetInfo)

	targetForm := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("targets").
				Description("Select targets to generate (space to select, enter to confirm, select none to select all):").
				Options(targetOptions...).
				Value(&selectedTargets),
		),
	).WithTheme(GetFormTheme(0))

	if err := targetForm.Run(); err != nil {
		return fmt.Errorf("target selection failed: %v", err)
	}

	if len(selectedTargets) == 0 {
		return fmt.Errorf("no languages selected")
	}

	Property("targets", strings.Join(selectedTargets, ", "))

	// Convert string targets to stainless.Target
	targets := make([]stainless.Target, len(selectedTargets))
	for i, target := range selectedTargets {
		targets[i] = stainless.Target(target)
	}

	// Phase 3: Start build and monitor progress in a loop
	for {
		err := runDevBuild(ctx, cc, cmd, selectedBranch, targets)
		if err != nil {
			if errors.Is(err, ErrUserCancelled) {
				return nil
			}
			return err
		}
	}
}

func runDevBuild(ctx context.Context, cc *apiCommandContext, cmd *cli.Command, branch string, languages []stainless.Target) error {
	// Handle file flags by reading files and mutating JSON body
	if err := applyFileFlag(cmd, "openapi-spec", "revision.openapi\\.yml.content"); err != nil {
		return err
	}
	if err := applyFileFlag(cmd, "stainless-config", "revision.openapi\\.stainless\\.yml.content"); err != nil {
		return err
	}

	projectName := cmd.String("project")
	buildReq := stainless.BuildNewParams{
		Project:    stainless.String(projectName),
		Branch:     stainless.String(branch),
		Targets:    languages,
		AllowEmpty: stainless.Bool(true),
	}

	model := NewBuildModel(cc, ctx, branch, func() (*stainless.BuildObject, error) {
		build, err := cc.client.Builds.New(ctx, buildReq, option.WithMiddleware(cc.AsMiddleware()))
		if err != nil {
			return nil, fmt.Errorf("failed to create build: %v", err)
		}
		return build, err
	})

	p := tea.NewProgram(model)
	finalModel, err := p.Run()

	if err != nil {
		return fmt.Errorf("failed to run TUI: %v", err)
	}
	if buildModel, ok := finalModel.(BuildModel); ok {
		return buildModel.err
	}
	return nil
}

func getGitUsername() (string, error) {
	cmd := exec.Command("git", "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	username := strings.TrimSpace(string(output))
	if username == "" {
		return "", fmt.Errorf("git username not configured")
	}

	// Convert to lowercase and replace spaces with hyphens for branch name
	return strings.ToLower(strings.ReplaceAll(username, " ", "-")), nil
}

func getCurrentGitBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	branch := strings.TrimSpace(string(output))
	if branch == "" {
		return "", fmt.Errorf("could not determine current git branch")
	}

	return branch, nil
}








