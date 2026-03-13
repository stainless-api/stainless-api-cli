package cmd

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/urfave/cli/v3"
)

//go:embed skill
var embeddedSkill embed.FS

var installCommand = cli.Command{
	Name:  "install",
	Usage: "Install Stainless development tools",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "skills",
			Usage: "Install coding agent skills for the Stainless CLI",
		},
	},
	Action:          handleInstall,
	HideHelpCommand: true,
}

func handleInstall(ctx context.Context, cmd *cli.Command) error {
	if !cmd.Bool("skills") {
		return fmt.Errorf("specify what to install, e.g.: stl install --skills")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	claudeExists := dirExists(filepath.Join(cwd, ".claude"))
	agentsExists := dirExists(filepath.Join(cwd, ".agents"))

	var primaryBase string
	var symlink bool

	switch {
	case claudeExists && agentsExists:
		primaryBase = ".agents"
		symlink = true
	case claudeExists && !agentsExists:
		primaryBase = ".claude"
	default:
		primaryBase = ".agents"
	}

	destDir := filepath.Join(cwd, primaryBase, "skills")

	err = fs.WalkDir(embeddedSkill, "skill", func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel("skill", p)
		if err != nil {
			return err
		}
		target := filepath.Join(destDir, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		data, err := embeddedSkill.ReadFile(p)
		if err != nil {
			return err
		}

		return os.WriteFile(target, data, 0644)
	})
	if err != nil {
		return fmt.Errorf("failed to install skills: %w", err)
	}

	skillDir := filepath.Join(destDir, "stl-cli")
	rel, err := filepath.Rel(cwd, skillDir)
	if err != nil {
		return fmt.Errorf("failed to compute relative path: %w", err)
	}
	console.Success("Skills installed to `%s`.", rel)

	if symlink {
		symlinkPath := filepath.Join(cwd, ".claude", "skills", "stl-cli")

		// If .claude is already a symlink to .agents, both paths resolve to the
		// same place and creating a symlink would cause a loop.
		realSkillDir, err1 := filepath.EvalSymlinks(skillDir)
		realSymlinkParent, err2 := filepath.EvalSymlinks(filepath.Dir(symlinkPath))
		sameRealPath := err1 == nil && err2 == nil && filepath.Join(realSymlinkParent, "stl-cli") == realSkillDir
		if !sameRealPath {
			if err := os.MkdirAll(filepath.Dir(symlinkPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

			if err := os.RemoveAll(symlinkPath); err != nil {
				return fmt.Errorf("failed to remove existing skill directory: %w", err)
			}

			relTarget, err := filepath.Rel(filepath.Dir(symlinkPath), skillDir)
			if err != nil {
				return fmt.Errorf("failed to compute symlink path: %w", err)
			}

			if err := os.Symlink(relTarget, symlinkPath); err != nil {
				return fmt.Errorf("failed to create symlink: %w", err)
			}

			symRel, err := filepath.Rel(cwd, symlinkPath)
			if err != nil {
				return fmt.Errorf("failed to compute relative path: %w", err)
			}
			console.Success("Symlinked `%s` → `%s`.", symRel, rel)
		}
	}

	return nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
