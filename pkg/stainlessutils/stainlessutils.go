package stainlessutils

import (
	"fmt"

	"github.com/stainless-api/stainless-api-go"
	"github.com/tidwall/gjson"
)

// Build wraps stainless.Build to provide convenience methods
type Build struct {
	stainless.Build
}

// NewBuild creates a new Build wrapper
func NewBuild(build stainless.Build) *Build {
	return &Build{Build: build}
}

// BuildTarget returns the build target wrapper for a given target type, replacing getBuildTarget
func (b *Build) BuildTarget(target stainless.Target) *BuildTarget {
	switch target {
	case "node":
		if b.Targets.JSON.Node.Valid() {
			return NewBuildTarget(&b.Targets.Node, target)
		}
	case "typescript":
		if b.Targets.JSON.Typescript.Valid() {
			return NewBuildTarget(&b.Targets.Typescript, target)
		}
	case "python":
		if b.Targets.JSON.Python.Valid() {
			return NewBuildTarget(&b.Targets.Python, target)
		}
	case "go":
		if b.Targets.JSON.Go.Valid() {
			return NewBuildTarget(&b.Targets.Go, target)
		}
	case "java":
		if b.Targets.JSON.Java.Valid() {
			return NewBuildTarget(&b.Targets.Java, target)
		}
	case "kotlin":
		if b.Targets.JSON.Kotlin.Valid() {
			return NewBuildTarget(&b.Targets.Kotlin, target)
		}
	case "ruby":
		if b.Targets.JSON.Ruby.Valid() {
			return NewBuildTarget(&b.Targets.Ruby, target)
		}
	case "terraform":
		if b.Targets.JSON.Terraform.Valid() {
			return NewBuildTarget(&b.Targets.Terraform, target)
		}
	case "cli":
		if b.Targets.JSON.Cli.Valid() {
			return NewBuildTarget(&b.Targets.Cli, target)
		}
	case "php":
		if b.Targets.JSON.Php.Valid() {
			return NewBuildTarget(&b.Targets.Php, target)
		}
	case "csharp":
		if b.Targets.JSON.Csharp.Valid() {
			return NewBuildTarget(&b.Targets.Csharp, target)
		}
	}
	return nil
}

// Languages returns all available build languages/targets for this build
func (b *Build) Languages() []stainless.Target {
	var languages []stainless.Target
	targets := b.Targets

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

// IsCompleted checks if the entire build is completed (all targets)
func (b *Build) IsCompleted() bool {
	languages := b.Languages()
	for _, target := range languages {
		buildTarget := b.BuildTarget(target)
		if buildTarget == nil || !buildTarget.IsCompleted() {
			return false
		}
	}
	return true
}

// BuildTarget wraps stainless.BuildTarget to provide convenience methods
type BuildTarget struct {
	*stainless.BuildTarget
	target stainless.Target
}

// NewBuildTarget creates a new BuildTarget wrapper
func NewBuildTarget(buildTarget *stainless.BuildTarget, target stainless.Target) *BuildTarget {
	if buildTarget == nil {
		return nil
	}
	return &BuildTarget{
		BuildTarget: buildTarget,
		target:      target,
	}
}

// Target returns the target type (node, python, etc.)
func (bt *BuildTarget) Target() stainless.Target {
	return bt.target
}

// StepUnion returns the step union for a given step name
func (bt *BuildTarget) StepUnion(step string) any {
	if bt.BuildTarget == nil {
		return nil
	}

	switch step {
	case "commit":
		if bt.JSON.Commit.Valid() {
			return bt.Commit
		}
	case "lint":
		if bt.JSON.Lint.Valid() {
			return bt.Lint
		}
	case "build":
		if bt.JSON.Build.Valid() {
			return bt.Build
		}
	case "test":
		if bt.JSON.Test.Valid() {
			return bt.Test
		}
	}
	return nil
}

// StepInfo extracts status, url, and conclusion from a step union
func (bt *BuildTarget) StepInfo(step string) (status, url, conclusion string) {
	stepUnion := bt.StepUnion(step)
	if stepUnion == nil {
		return "", "", ""
	}

	if u, ok := stepUnion.(stainless.BuildTargetCommitUnion); ok {
		status = u.Status
		if u.Status == "completed" {
			conclusion = u.Completed.Conclusion
			url = fmt.Sprintf("https://github.com/%s/%s/commit/%s", u.Completed.Commit.Repo.Owner, u.Completed.Commit.Repo.Name, u.Completed.Commit.Sha)
		}
	}
	if u, ok := stepUnion.(stainless.CheckStepUnion); ok {
		status = u.Status
		if u.Status == "completed" {
			conclusion = u.Completed.Conclusion
			url = u.Completed.URL
		}
	}
	return
}

// Steps returns all available steps for this build target
func (bt *BuildTarget) Steps() []string {
	if bt.BuildTarget == nil {
		return []string{}
	}

	var steps []string

	if gjson.Get(bt.RawJSON(), "commit").Exists() {
		steps = append(steps, "commit")
	}
	if gjson.Get(bt.RawJSON(), "lint").Exists() {
		steps = append(steps, "lint")
	}
	if gjson.Get(bt.RawJSON(), "build").Exists() {
		steps = append(steps, "build")
	}
	if gjson.Get(bt.RawJSON(), "test").Exists() {
		steps = append(steps, "test")
	}

	return steps
}

func (bt *BuildTarget) IsCompleted() bool {
	steps := []string{"commit", "lint", "build", "test"}
	for _, step := range steps {
		if !gjson.Get(bt.RawJSON(), step).Exists() {
			continue
		}
		status, _, _ := bt.StepInfo(step)
		if status != "completed" {
			return false
		}
	}
	return true
}

func (bt *BuildTarget) IsInProgress() bool {
	steps := []string{"commit", "lint", "build", "test", "upload"}
	for _, step := range steps {
		status, _, _ := bt.StepInfo(step)
		if status == "in_progress" {
			return true
		}
	}
	return false
}

func (bt *BuildTarget) IsCommitCompleted() bool {
	status, _, _ := bt.StepInfo("commit")
	return status == "completed"
}

func (bt *BuildTarget) IsCommitFailed() bool {
	status, _, conclusion := bt.StepInfo("commit")
	return status == "completed" && conclusion == "fatal"
}
