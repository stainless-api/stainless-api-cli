package cmd

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var ErrUserCancelled = errors.New("user cancelled")

// BuildModel represents the bubbletea model for build monitoring
type BuildModel struct {
	start       func() (*stainless.BuildObject, error)
	started     time.Time
	ended       *time.Time
	build       *stainless.BuildObject
	diagnostics []stainless.BuildDiagnosticListResponse
	view        string

	cc          *apiCommandContext
	ctx         context.Context
	err         error
	isCompleted bool
}

type tickMsg time.Time
type fetchBuildMsg *stainless.BuildObject
type fetchDiagnosticsMsg []stainless.BuildDiagnosticListResponse
type errorMsg error
type triggerNewBuildMsg struct{}

func NewBuildModel(cc *apiCommandContext, ctx context.Context, fn func() (*stainless.BuildObject, error)) BuildModel {
	return BuildModel{
		start:   fn,
		started: time.Now(),
		cc:      cc,
		ctx:     ctx,
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
			cmds = append(cmds, m.updateView("header"))
		}

		m.build = msg
		if !m.isCompleted && isCommitStepsCompleted(m.build) {
			m.isCompleted = true
			cmds = append(cmds, m.fetchDiagnostics())
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

	if isBuildCompleted(m.build) {
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
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return runDevMode(ctx, cmd)
	},
}

func runDevMode(ctx context.Context, cmd *cli.Command) error {
	projectName := GetProjectName(cmd, "project")
	if projectName == "" {
		return fmt.Errorf("project name is required")
	}

	configuredTargets := []stainless.Target{
		"typescript", "python", "go", "java", "kotlin",
		"ruby", "terraform", "cli", "php", "csharp", "node",
	}

	gitUser, err := getGitUsername()
	if err != nil {
		return fmt.Errorf("failed to get git username: %v", err)
	}

	// Phase 1: Branch selection
	var selectedBranch string
	branchOptions := []huh.Option[string]{
		huh.NewOption(fmt.Sprintf("%s/dev", gitUser), fmt.Sprintf("%s/dev", gitUser)),
		huh.NewOption("Random", "random"),
	}

	branchForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("branch").
				Description("Select a branch to use for development:").
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
	targetOptions := buildTargetOptions(configuredTargets)

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

	// Get API command context for the build
	cc, err := getAPICommandContextWithWorkspaceDefaults(cmd)
	if err != nil {
		return err
	}

	// Phase 3: Start build and monitor progress in a loop
	for {
		err := runDevBuild(ctx, cc, projectName, selectedBranch, targets)
		if err != nil {
			if errors.Is(err, ErrUserCancelled) {
				return nil
			}
			return err
		}
	}
}

func runDevBuild(ctx context.Context, cc *apiCommandContext, projectName string, branch string, languages []stainless.Target) error {
	buildReq := stainless.BuildNewParams{
		Project:    stainless.String(projectName),
		Branch:     stainless.String(branch),
		Targets:    languages,
		AllowEmpty: stainless.Bool(true),
	}

	model := NewBuildModel(cc, ctx, func() (*stainless.BuildObject, error) {
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

func getBuildTarget(build *stainless.BuildObject, target stainless.Target) *stainless.BuildTarget {
	switch target {
	case "node":
		if build.Targets.JSON.Node.Valid() {
			return &build.Targets.Node
		}
	case "typescript":
		if build.Targets.JSON.Typescript.Valid() {
			return &build.Targets.Typescript
		}
	case "python":
		if build.Targets.JSON.Python.Valid() {
			return &build.Targets.Python
		}
	case "go":
		if build.Targets.JSON.Go.Valid() {
			return &build.Targets.Go
		}
	case "java":
		if build.Targets.JSON.Java.Valid() {
			return &build.Targets.Java
		}
	case "kotlin":
		if build.Targets.JSON.Kotlin.Valid() {
			return &build.Targets.Kotlin
		}
	case "ruby":
		if build.Targets.JSON.Ruby.Valid() {
			return &build.Targets.Ruby
		}
	case "terraform":
		if build.Targets.JSON.Terraform.Valid() {
			return &build.Targets.Terraform
		}
	case "cli":
		if build.Targets.JSON.Cli.Valid() {
			return &build.Targets.Cli
		}
	case "php":
		if build.Targets.JSON.Php.Valid() {
			return &build.Targets.Php
		}
	case "csharp":
		if build.Targets.JSON.Csharp.Valid() {
			return &build.Targets.Csharp
		}
	}
	return nil
}

func getStepUnion(target *stainless.BuildTarget, step string) any {
	switch step {
	case "commit":
		if target.JSON.Commit.Valid() {
			return target.Commit
		}
	case "lint":
		if target.JSON.Lint.Valid() {
			return target.Lint
		}
	case "build":
		if target.JSON.Build.Valid() {
			return target.Build
		}
	case "test":
		if target.JSON.Test.Valid() {
			return target.Test
		}
	}
	return nil
}

func extractStepInfo(stepUnion interface{}) (status, conclusion string) {
	if u, ok := stepUnion.(stainless.BuildTargetCommitUnion); ok {
		status = u.Status
		if u.Status == "completed" {
			conclusion = u.Completed.Conclusion
		}
	}
	if u, ok := stepUnion.(stainless.CheckStepUnion); ok {
		status = u.Status
		if u.Status == "completed" {
			conclusion = u.Completed.Conclusion
		}
	}
	return
}

func isBuildTargetCompleted(build *stainless.BuildObject, target stainless.Target) bool {
	buildTarget := getBuildTarget(build, target)
	if buildTarget == nil {
		return false
	}

	steps := []string{"commit", "lint", "build", "test"}
	for _, step := range steps {
		if !gjson.Get(buildTarget.RawJSON(), step).Exists() {
			continue
		}
		stepUnion := getStepUnion(buildTarget, step)
		if stepUnion == nil {
			continue
		}
		status, _ := extractStepInfo(stepUnion)
		if status != "completed" {
			return false
		}
	}
	return true
}

func isBuildTargetInProgress(build *stainless.BuildObject, target stainless.Target) bool {
	buildTarget := getBuildTarget(build, target)
	if buildTarget == nil {
		return false
	}

	steps := []string{"commit", "lint", "build", "test", "upload"}
	for _, step := range steps {
		stepUnion := getStepUnion(buildTarget, step)
		if stepUnion == nil {
			continue
		}
		status, _ := extractStepInfo(stepUnion)
		if status == "in_progress" {
			return true
		}
	}
	return false
}

func isCommitStepsCompleted(build *stainless.BuildObject) bool {
	languages := getBuildLanguages(build)

	for _, target := range languages {
		buildTarget := getBuildTarget(build, target)
		if buildTarget == nil {
			return false
		}

		// Check if commit step is completed
		commitUnion := getStepUnion(buildTarget, "commit")
		if commitUnion == nil {
			continue
		}
		status, _ := extractStepInfo(commitUnion)
		if status != "completed" {
			return false
		}
	}
	return true
}

func isBuildCompleted(build *stainless.BuildObject) bool {
	languages := getBuildLanguages(build)
	for _, target := range languages {
		if !isBuildTargetCompleted(build, target) {
			return false
		}
	}
	return true
}

func getBuildSteps(buildTarget *stainless.BuildTarget) []string {
	if buildTarget == nil {
		return []string{}
	}

	var steps []string

	if gjson.Get(buildTarget.RawJSON(), "commit").Exists() {
		steps = append(steps, "commit")
	}
	if gjson.Get(buildTarget.RawJSON(), "lint").Exists() {
		steps = append(steps, "lint")
	}
	if gjson.Get(buildTarget.RawJSON(), "build").Exists() {
		steps = append(steps, "build")
	}
	if gjson.Get(buildTarget.RawJSON(), "test").Exists() {
		steps = append(steps, "test")
	}

	return steps
}

func getBuildLanguages(build *stainless.BuildObject) []stainless.Target {
	if build == nil {
		return []stainless.Target{}
	}

	var languages []stainless.Target
	targets := build.Targets

	if targets.JSON.Node.Valid() {
		languages = append(languages, "node")
	}
	if targets.JSON.Typescript.Valid() {
		languages = append(languages, "typescript")
	}
	if targets.JSON.Python.Valid() {
		languages = append(languages, "python")
	}
	if targets.JSON.Go.Valid() {
		languages = append(languages, "go")
	}
	if targets.JSON.Java.Valid() {
		languages = append(languages, "java")
	}
	if targets.JSON.Kotlin.Valid() {
		languages = append(languages, "kotlin")
	}
	if targets.JSON.Ruby.Valid() {
		languages = append(languages, "ruby")
	}
	if targets.JSON.Terraform.Valid() {
		languages = append(languages, "terraform")
	}
	if targets.JSON.Cli.Valid() {
		languages = append(languages, "cli")
	}
	if targets.JSON.Php.Valid() {
		languages = append(languages, "php")
	}
	if targets.JSON.Csharp.Valid() {
		languages = append(languages, "csharp")
	}

	return languages

}
