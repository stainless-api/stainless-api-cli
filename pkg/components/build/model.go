package build

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-cli/pkg/stainlessutils"
	"github.com/stainless-api/stainless-api-go"
)

var ErrUserCancelled = errors.New("user cancelled")

type Model struct {
	stainless.Build // Current build. This component will keep fetching it until the build is completed. A zero value is permitted.

	Client    stainless.Client
	Ctx       context.Context
	Downloads map[stainless.Target]DownloadStatus // When a BuildTarget has a commit available, this target will download it, if it has been specified in the initialization.
	Err       error                               // This will be populated if the model concludes with an error
}

type DownloadStatus struct {
	Status string
	Path   string
}

type TickMsg time.Time
type FetchBuildMsg stainless.Build
type ErrorMsg error
type DownloadMsg stainless.Target

func NewModel(client stainless.Client, ctx context.Context, build stainless.Build, downloadPaths map[stainless.Target]string) Model {
	downloads := map[stainless.Target]DownloadStatus{}
	for target, path := range downloadPaths {
		downloads[target] = DownloadStatus{
			Path: path,
		}
	}

	return Model{
		Build:     build,
		Client:    client,
		Ctx:       ctx,
		Downloads: downloads,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return TickMsg(time.Now())
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.Err = ErrUserCancelled
			cmds = append(cmds, tea.Quit)
		}

	case TickMsg:
		if m.Build.ID != "" {
			cmds = append(cmds, m.fetchBuildStatus())
		}
		cmds = append(cmds, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return TickMsg(t)
		}))

	case DownloadMsg:
		download := m.Downloads[stainless.Target(msg)]
		download.Status = "completed"
		m.Downloads[stainless.Target(msg)] = download

	case FetchBuildMsg:
		m.Build = stainless.Build(msg)
		buildObj := stainlessutils.NewBuild(m.Build)
		languages := buildObj.Languages()
		for _, target := range languages {
			buildTarget := buildObj.BuildTarget(target)
			if buildTarget == nil {
				continue
			}
			status, _, conclusion := buildTarget.StepInfo("commit")
			if status == "completed" && conclusion != "fatal" {
				if download, ok := m.Downloads[target]; ok && download.Status == "not started" {
					download.Status = "started"
					cmds = append(cmds, m.downloadTarget(target))
					m.Downloads[target] = download
				}
			}
		}

	case ErrorMsg:
		m.Err = msg
	}

	return m, tea.Batch(cmds...)
}

func (m Model) downloadTarget(target stainless.Target) tea.Cmd {
	return func() tea.Msg {
		params := stainless.BuildTargetOutputGetParams{
			BuildID: m.Build.ID,
			Target:  stainless.BuildTargetOutputGetParamsTarget(target),
			Type:    "source",
			Output:  "git",
		}
		outputRes, err := m.Client.Builds.TargetOutputs.Get(
			context.TODO(),
			params,
		)
		if err != nil {
			return ErrorMsg(err)
		}
		err = PullOutput(outputRes.Output, outputRes.URL, outputRes.Ref, m.Downloads[target].Path, &console.Group{})
		if err != nil {
			return ErrorMsg(err)
		}
		return DownloadMsg(target)
	}
}

func (m Model) fetchBuildStatus() tea.Cmd {
	return func() tea.Msg {
		build, err := m.Client.Builds.Get(m.Ctx, m.Build.ID)
		if err != nil {
			return ErrorMsg(fmt.Errorf("failed to get build status: %v", err))
		}
		return FetchBuildMsg(*build)
	}
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

// PullOutput handles downloading or cloning a build target output
func PullOutput(output, url, ref, targetDir string, targetGroup *console.Group) error {
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
