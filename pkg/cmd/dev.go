package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var ErrUserCancelled = errors.New("user cancelled")

// BuildModel represents the bubbletea model for build monitoring
type BuildModel struct {
	start       func() (*stainless.Build, error)
	started     time.Time
	ended       *time.Time
	build       *stainless.Build
	branch      string
	help        help.Model
	diagnostics []stainless.BuildDiagnostic
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
type fetchBuildMsg *stainless.Build
type fetchDiagnosticsMsg []stainless.BuildDiagnostic
type errorMsg error
type downloadMsg stainless.Target
type fileChangeMsg struct{}

func NewBuildModel(cc *apiCommandContext, ctx context.Context, branch string, fn func() (*stainless.Build, error)) BuildModel {
	return BuildModel{
		start:   fn,
		started: time.Now(),
		cc:      cc,
		ctx:     ctx,
		branch:  branch,
		help:    help.New(),
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
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.err = ErrUserCancelled
			cmds = append(cmds, tea.Quit)
		case "enter":
			if m.cc.cmd.Bool("watch") {
				cmds = append(cmds, tea.Quit)
			}
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
		buildObj := NewBuild(m.build)
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

	case fileChangeMsg:
		// File change detected, exit with success
		cmds = append(cmds, tea.Quit)
	}
	return m, tea.Sequence(cmds...)
}

func (m BuildModel) ShortHelp() []key.Binding {
	if m.cc.cmd.Bool("watch") {
		return []key.Binding{
			key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl-c", "quit")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "rebuild")),
		}
	} else {
		return []key.Binding{key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl-c", "quit"))}
	}
}

func (m BuildModel) FullHelp() [][]key.Binding {
	if m.cc.cmd.Bool("watch") {
		return [][]key.Binding{{
			key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl-c", "quit")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "rebuild")),
		}}
	} else {
		return [][]key.Binding{{key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl-c", "quit"))}}
	}
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
		diags := []stainless.BuildDiagnostic{}
		diagnostics := m.cc.client.Builds.Diagnostics.ListAutoPaging(m.ctx, m.build.ID, stainless.BuildDiagnosticListParams{
			Limit: stainless.Float(100),
		})
		for diagnostics.Next() {
			diag := diagnostics.Current()
			if !diag.Ignored {
				diags = append(diags, diag)
			}
		}
		return fetchDiagnosticsMsg(diags)
	}
}

func (m *BuildModel) getBuildDuration() time.Duration {
	if m.build == nil {
		return time.Since(m.started)
	}

	buildObj := NewBuild(m.build)
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
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Which branch to use",
		},
		&cli.StringSliceFlag{
			Name:    "target",
			Aliases: []string{"t"},
			Usage:   "The target build language(s)",
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
		// Clear the screen and move the cursor to the top
		fmt.Print("\033[2J\033[H")
		os.Stdout.Sync()
	}

	cc := getAPICommandContext(cmd)

	gitUser, err := getGitUsername()
	if err != nil {
		Warn("Couldn't get a git user: %s", err)
		gitUser = "user"
	}

	var selectedBranch string
	if cmd.IsSet("branch") {
		selectedBranch = cmd.String("branch")
	} else {
		selectedBranch, err = chooseBranch(gitUser)
		if err != nil {
			return err
		}
	}
	Property("branch", selectedBranch)

	// Phase 2: Language selection
	var selectedTargets []string
	targetInfos := getAvailableTargetInfo(ctx, cc.client, cmd.String("project"), cc.workspaceConfig)
	if cmd.IsSet("target") {
		selectedTargets = cmd.StringSlice("target")
		for _, target := range selectedTargets {
			if !isValidTarget(targetInfos, target) {
				return fmt.Errorf("invalid language target: %s", target)
			}
		}
	} else {
		selectedTargets, err = chooseSelectedTargets(targetInfos)
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
		// Make the user get past linter errors
		if err := runLinter(ctx, cmd, true); err != nil {
			if errors.Is(err, ErrUserCancelled) {
				return nil
			}
			return err
		}

		// Start the build process
		if err := runDevBuild(ctx, cc, cmd, selectedBranch, targets); err != nil {
			if errors.Is(err, ErrUserCancelled) {
				return nil
			}
			return err
		}

		if !cmd.Bool("watch") {
			break
		}
	}
	return nil
}

func chooseBranch(gitUser string) (string, error) {
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

	var selectedBranch string
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
		return selectedBranch, fmt.Errorf("branch selection failed: %v", err)
	}

	return selectedBranch, nil
}

func chooseSelectedTargets(targetInfos []TargetInfo) ([]string, error) {
	targetOptions := targetInfoToOptions(targetInfos)

	var selectedTargets []string
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
		return nil, fmt.Errorf("target selection failed: %v", err)
	}
	return selectedTargets, nil
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

	model := NewBuildModel(cc, ctx, branch, func() (*stainless.Build, error) {
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

type GenerateSpecParams struct {
	Project string `json:"project"`
	Source  struct {
		Type            string `json:"type"`
		OpenAPISpec     string `json:"openapi_spec"`
		StainlessConfig string `json:"stainless_config"`
	} `json:"source"`
}

func getDiagnostics(ctx context.Context, cmd *cli.Command, cc *apiCommandContext) ([]stainless.BuildDiagnostic, error) {
	var specParams GenerateSpecParams
	if cmd.IsSet("project") {
		specParams.Project = cmd.String("project")
	} else {
		specParams.Project = cc.workspaceConfig.Project
	}
	specParams.Source.Type = "upload"

	configPath := cc.workspaceConfig.StainlessConfig
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

	oasPath := cc.workspaceConfig.OpenAPISpec
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

	var result []byte
	err = cc.client.Post(ctx, "api/generate/spec", specParams, &result, option.WithMiddleware(cc.AsMiddleware()))
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
