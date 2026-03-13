package cmd

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func chdir(t *testing.T, dir string) {
	t.Helper()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { os.Chdir(oldWd) })
}

func TestSkills(t *testing.T) {
	t.Run("neither dir exists defaults to .agents", func(t *testing.T) {
		tmpDir := t.TempDir()
		chdir(t, tmpDir)

		err := Command.Run(context.Background(), []string{"stl", "skills"})
		require.NoError(t, err)
		require.FileExists(t,filepath.Join(tmpDir, ".agents", "skills", "stl-cli", "SKILL.md"))
		assert.NoDirExists(t, filepath.Join(tmpDir, ".claude"))
	})

	t.Run("only .claude exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		chdir(t, tmpDir)
		require.NoError(t, os.Mkdir(filepath.Join(tmpDir, ".claude"), 0755))

		err := Command.Run(context.Background(), []string{"stl", "skills"})
		require.NoError(t, err)
		require.FileExists(t,filepath.Join(tmpDir, ".claude", "skills", "stl-cli", "SKILL.md"))
		assert.NoDirExists(t, filepath.Join(tmpDir, ".agents"))
	})

	t.Run("only .agents exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		chdir(t, tmpDir)
		require.NoError(t, os.Mkdir(filepath.Join(tmpDir, ".agents"), 0755))

		err := Command.Run(context.Background(), []string{"stl", "skills"})
		require.NoError(t, err)
		require.FileExists(t,filepath.Join(tmpDir, ".agents", "skills", "stl-cli", "SKILL.md"))
		assert.NoDirExists(t, filepath.Join(tmpDir, ".claude"))
	})

	t.Run("both exist installs to .agents with symlink", func(t *testing.T) {
		tmpDir := t.TempDir()
		chdir(t, tmpDir)
		require.NoError(t, os.Mkdir(filepath.Join(tmpDir, ".claude"), 0755))
		require.NoError(t, os.Mkdir(filepath.Join(tmpDir, ".agents"), 0755))

		err := Command.Run(context.Background(), []string{"stl", "skills"})
		require.NoError(t, err)

		// Primary install goes to .agents
		require.FileExists(t,filepath.Join(tmpDir, ".agents", "skills", "stl-cli", "SKILL.md"))

		// .claude/skills/stl-cli should be a symlink
		symlinkPath := filepath.Join(tmpDir, ".claude", "skills", "stl-cli")
		target, err := os.Readlink(symlinkPath)
		require.NoError(t, err, ".claude/skills/stl-cli should be a symlink")
		assert.Equal(t, filepath.Join("..", "..", ".agents", "skills", "stl-cli"), target)

		// Symlink should resolve to a valid file
		require.FileExists(t,filepath.Join(symlinkPath, "SKILL.md"))
	})

	t.Run("idempotent re-run", func(t *testing.T) {
		tmpDir := t.TempDir()
		chdir(t, tmpDir)
		require.NoError(t, os.Mkdir(filepath.Join(tmpDir, ".claude"), 0755))

		err := Command.Run(context.Background(), []string{"stl", "skills"})
		require.NoError(t, err)
		err = Command.Run(context.Background(), []string{"stl", "skills"})
		require.NoError(t, err)
		require.FileExists(t,filepath.Join(tmpDir, ".claude", "skills", "stl-cli", "SKILL.md"))
	})

	t.Run("idempotent re-run with symlink", func(t *testing.T) {
		tmpDir := t.TempDir()
		chdir(t, tmpDir)
		require.NoError(t, os.Mkdir(filepath.Join(tmpDir, ".claude"), 0755))
		require.NoError(t, os.Mkdir(filepath.Join(tmpDir, ".agents"), 0755))

		err := Command.Run(context.Background(), []string{"stl", "skills"})
		require.NoError(t, err)
		err = Command.Run(context.Background(), []string{"stl", "skills"})
		require.NoError(t, err)

		require.FileExists(t,filepath.Join(tmpDir, ".agents", "skills", "stl-cli", "SKILL.md"))

		target, err := os.Readlink(filepath.Join(tmpDir, ".claude", "skills", "stl-cli"))
		require.NoError(t, err)
		assert.Equal(t, filepath.Join("..", "..", ".agents", "skills", "stl-cli"), target)
	})

	t.Run("upgrade from real directory to symlink", func(t *testing.T) {
		tmpDir := t.TempDir()
		chdir(t, tmpDir)
		require.NoError(t, os.Mkdir(filepath.Join(tmpDir, ".claude"), 0755))

		// First install: only .claude exists → real directory
		err := Command.Run(context.Background(), []string{"stl", "skills"})
		require.NoError(t, err)
		info, err := os.Lstat(filepath.Join(tmpDir, ".claude", "skills", "stl-cli"))
		require.NoError(t, err)
		assert.True(t, info.IsDir(), "should be a real directory before upgrade")

		// Now add .agents and re-run → should replace with symlink
		require.NoError(t, os.Mkdir(filepath.Join(tmpDir, ".agents"), 0755))
		err = Command.Run(context.Background(), []string{"stl", "skills"})
		require.NoError(t, err)

		info, err = os.Lstat(filepath.Join(tmpDir, ".claude", "skills", "stl-cli"))
		require.NoError(t, err)
		assert.True(t, info.Mode()&os.ModeSymlink != 0, "should be a symlink after upgrade")

		require.FileExists(t,filepath.Join(tmpDir, ".agents", "skills", "stl-cli", "SKILL.md"))
		require.FileExists(t,filepath.Join(tmpDir, ".claude", "skills", "stl-cli", "SKILL.md"))
	})

	t.Run(".claude symlinked to .agents skips symlink to avoid loop", func(t *testing.T) {
		tmpDir := t.TempDir()
		chdir(t, tmpDir)
		require.NoError(t, os.Mkdir(filepath.Join(tmpDir, ".agents"), 0755))
		require.NoError(t, os.Symlink(".agents", filepath.Join(tmpDir, ".claude")))

		err := Command.Run(context.Background(), []string{"stl", "skills"})
		require.NoError(t, err)

		// Files should be in .agents
		require.FileExists(t,filepath.Join(tmpDir, ".agents", "skills", "stl-cli", "SKILL.md"))

		// .claude/skills/stl-cli should NOT be a separate symlink (would cause a loop)
		symlinkPath := filepath.Join(tmpDir, ".claude", "skills", "stl-cli")
		_, err = os.Readlink(symlinkPath)
		assert.Error(t, err, ".claude/skills/stl-cli should not be a symlink itself")

		// But the file should still be accessible through .claude
		require.FileExists(t,filepath.Join(tmpDir, ".claude", "skills", "stl-cli", "SKILL.md"))
	})
}
