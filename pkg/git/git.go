package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// RevParse runs git rev-parse and returns the SHA, or error if it fails
func RevParse(dir, ref string) (string, error) {
	cmd := exec.Command("git", "-C", dir, "rev-parse", ref)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

// Checkout runs git checkout with the given arguments
func Checkout(dir string, args ...string) error {
	fullArgs := append([]string{"-C", dir, "checkout"}, args...)
	cmd := exec.Command("git", fullArgs...)
	var stderr bytes.Buffer
	cmd.Stdout = nil
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git checkout failed: %v\nGit error: %s", err, stderr.String())
	}
	return nil
}

// Init initializes a git repository in the given directory
func Init(dir string) error {
	cmd := exec.Command("git", "-C", dir, "init")
	var stderr bytes.Buffer
	cmd.Stdout = nil
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git init failed: %v\nGit error: %s", err, stderr.String())
	}
	return nil
}

// RemoteAdd adds a remote to the repository
func RemoteAdd(dir, name, url string) error {
	cmd := exec.Command("git", "-C", dir, "remote", "add", name, url)
	var stderr bytes.Buffer
	cmd.Stdout = nil
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git remote add failed: %v\nGit error: %s", err, stderr.String())
	}
	return nil
}

// RemoteGetURL gets the URL of a remote
func RemoteGetURL(dir, name string) (string, error) {
	cmd := exec.Command("git", "-C", dir, "remote", "get-url", name)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

// Fetch fetches from a remote
func Fetch(dir, url string, refspecs ...string) error {
	args := append([]string{"-C", dir, "fetch", url}, refspecs...)
	cmd := exec.Command("git", args...)
	var stderr bytes.Buffer
	cmd.Stdout = nil
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git fetch failed: %v\nGit error: %s", err, stderr.String())
	}
	return nil
}
