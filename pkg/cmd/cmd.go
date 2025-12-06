// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	docs "github.com/urfave/cli-docs/v3"
	"github.com/urfave/cli/v3"
)

var (
	Command       *cli.Command
	OutputFormats = []string{"auto", "explore", "json", "jsonl", "pretty", "raw", "yaml"}
)

func init() {
	Command = &cli.Command{
		Name:    "stl",
		Usage:   "CLI for the stainless API",
		Version: Version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug logging",
			},
			&cli.StringFlag{
				Name:        "base-url",
				DefaultText: "url",
				Usage:       "Override the base URL for API requests",
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "The format for displaying response data (one of: " + strings.Join(OutputFormats, ", ") + ")",
				Value: "auto",
				Validator: func(format string) error {
					if !slices.Contains(OutputFormats, strings.ToLower(format)) {
						return fmt.Errorf("format must be one of: %s", strings.Join(OutputFormats, ", "))
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "format-error",
				Usage: "The format for displaying error data (one of: " + strings.Join(OutputFormats, ", ") + ")",
				Value: "auto",
				Validator: func(format string) error {
					if !slices.Contains(OutputFormats, strings.ToLower(format)) {
						return fmt.Errorf("format must be one of: %s", strings.Join(OutputFormats, ", "))
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "transform",
				Usage: "The GJSON transformation for data output.",
			},
			&cli.StringFlag{
				Name:  "transform-error",
				Usage: "The GJSON transformation for errors.",
			},
			&cli.StringFlag{
				Name:  "environment",
				Usage: "Set the environment for API requests",
			},
		},
		Commands: []*cli.Command{
			{
				Name:     "projects",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&projectsCreate,
					&projectsRetrieve,
					&projectsUpdate,
					&projectsList,
				},
			},
			{
				Name:     "projects:branches",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&projectsBranchesCreate,
					&projectsBranchesRetrieve,
					&projectsBranchesList,
					&projectsBranchesDelete,
					&projectsBranchesRebase,
					&projectsBranchesReset,
				},
			},
			{
				Name:     "projects:configs",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&projectsConfigsRetrieve,
					&projectsConfigsGuess,
				},
			},
			{
				Name:     "builds",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&buildsCreate,
					&buildsRetrieve,
					&buildsList,
					&buildsCompare,
				},
			},
			{
				Name:     "builds:diagnostics",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&buildsDiagnosticsList,
				},
			},
			{
				Name:     "builds:target-outputs",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&buildsTargetOutputsRetrieve,
				},
			},
			{
				Name:     "orgs",
				Category: "API RESOURCE",
				Commands: []*cli.Command{
					&orgsRetrieve,
					&orgsList,
				},
			},
			{
				Name:            "@manpages",
				Usage:           "Generate documentation for 'man'",
				UsageText:       "stl @manpages [-o stl.1] [--gzip]",
				Hidden:          true,
				Action:          generateManpages,
				HideHelpCommand: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "write manpages to the given folder",
						Value:   "man",
					},
					&cli.BoolFlag{
						Name:    "gzip",
						Aliases: []string{"z"},
						Usage:   "output gzipped manpage files to .gz",
						Value:   true,
					},
					&cli.BoolFlag{
						Name:    "text",
						Aliases: []string{"z"},
						Usage:   "output uncompressed text files",
						Value:   false,
					},
				},
			},
		},
		EnableShellCompletion:      true,
		ShellCompletionCommandName: "@completion",
		HideHelpCommand:            true,
	}
}

func generateManpages(ctx context.Context, c *cli.Command) error {
	manpage, err := docs.ToManWithSection(Command, 1)
	if err != nil {
		return err
	}
	dir := c.String("output")
	err = os.MkdirAll(filepath.Join(dir, "man1"), 0755)
	if err != nil {
		// handle error
	}
	if c.Bool("text") {
		file, err := os.Create(filepath.Join(dir, "man1", "stl.1"))
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(manpage); err != nil {
			return err
		}
	}
	if c.Bool("gzip") {
		file, err := os.Create(filepath.Join(dir, "man1", "stl.1.gz"))
		if err != nil {
			return err
		}
		defer file.Close()
		gzWriter := gzip.NewWriter(file)
		defer gzWriter.Close()
		_, err = gzWriter.Write([]byte(manpage))
		if err != nil {
			return err
		}
	}
	fmt.Printf("Wrote manpages to %s\n", dir)
	return nil
}
