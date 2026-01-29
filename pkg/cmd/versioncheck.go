package cmd

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

// CheckForUpdate starts a background check for new versions and returns a channel
// that will contain an update message if one is available
func CheckForUpdate() <-chan string {
	updateMsg := make(chan string, 1)

	go func() {
		defer close(updateMsg)

		client := &http.Client{
			Timeout: 3 * time.Second,
		}

		resp, err := client.Get("https://api.github.com/repos/stainless-api/stainless-api-cli/releases/latest")
		if err != nil {
			return
		}
		defer resp.Body.Close()

		var release GitHubRelease
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			return
		}

		latest := release.TagName
		if !strings.HasPrefix(latest, "v") {
			latest = "v" + latest
		}
		current := Version
		if !strings.HasPrefix(current, "v") {
			current = "v" + current
		}

		if semver.Compare(latest, current) > 0 {
			updateMsg <- "New version available: " + latest + " (current: " + current + ")"
		}
	}()

	return updateMsg
}
