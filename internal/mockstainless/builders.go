package mockstainless

import "time"

// --- CheckStep builders ---

func CheckStepNotStarted() M {
	return M{"status": "not_started"}
}

func CheckStepInProgress() M {
	return M{"status": "in_progress", "url": ""}
}

func CheckStepCompleted(conclusion string, opts ...func(M)) M {
	m := M{"status": "completed", "conclusion": conclusion, "url": ""}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// --- Commit builders ---

func CommitNotStarted() M {
	return M{"status": "not_started"}
}

func CommitInProgress() M {
	return M{"status": "in_progress"}
}

func CommitCompleted(conclusion string, opts ...func(M)) M {
	m := M{
		"status":     "completed",
		"conclusion": conclusion,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func WithCommitData(owner, repo, sha string, additions, deletions int) func(M) {
	return func(m M) {
		m["commit"] = M{
			"sha":      sha,
			"tree_oid": "tree_" + sha[:7],
			"repo": M{
				"owner": owner,
				"name":  repo,
			},
			"stats": M{
				"additions": additions,
				"deletions": deletions,
				"total":     additions + deletions,
			},
		}
	}
}

func WithMergeConflictPR(owner, repo string, number int) func(M) {
	return func(m M) {
		m["merge_conflict_pr"] = M{
			"repo": M{
				"owner": owner,
				"name":  repo,
			},
			"number": number,
		}
	}
}

// --- BuildTarget builders ---

type TargetOption func(M)

// Target creates a build target with the given status and commit state.
// Lint, build, and test default to not_started.
func Target(status string, commit M, opts ...TargetOption) M {
	m := M{
		"object":      "build_target",
		"status":      status,
		"install_url": "",
		"commit":      commit,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func WithLint(step M) TargetOption  { return func(m M) { m["lint"] = step } }
func WithBuild(step M) TargetOption { return func(m M) { m["build"] = step } }
func WithTest(step M) TargetOption  { return func(m M) { m["test"] = step } }

// Convenience target constructors

func CompletedTarget(owner, repo, sha string, additions, deletions int) M {
	return Target("completed",
		CommitCompleted("success", WithCommitData(owner, repo, sha, additions, deletions)),
		WithLint(CheckStepCompleted("success")),
		WithBuild(CheckStepCompleted("success")),
		WithTest(CheckStepCompleted("success")),
	)
}

func WarningTarget(owner, repo, sha string, additions, deletions int) M {
	return Target("completed",
		CommitCompleted("warning", WithCommitData(owner, repo, sha, additions, deletions)),
		WithLint(CheckStepCompleted("success")),
		WithBuild(CheckStepCompleted("success")),
		WithTest(CheckStepCompleted("success")),
	)
}

func ErrorTarget(owner, repo, sha string, additions, deletions int) M {
	return Target("completed",
		CommitCompleted("error", WithCommitData(owner, repo, sha, additions, deletions)),
		WithLint(CheckStepCompleted("success")),
		WithBuild(CheckStepCompleted("success")),
		WithTest(CheckStepCompleted("failure")),
	)
}

func FatalTarget() M {
	return Target("completed", CommitCompleted("fatal"))
}

func MergeConflictTarget(owner, repo string, prNum int) M {
	return Target("completed",
		CommitCompleted("merge_conflict", WithMergeConflictPR(owner, repo, prNum)),
	)
}

func NotStartedTarget() M {
	return Target("not_started", CommitNotStarted())
}

func InProgressTarget() M {
	return Target("codegen", CommitInProgress())
}

// --- Build builders ---

type BuildOption func(M)

// Build creates a build with sensible defaults.
func Build(id string, opts ...BuildOption) M {
	m := M{
		"id":            id,
		"config_commit": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0",
		"created_at":    time.Now().Format(time.RFC3339),
		"org":           DefaultOrg,
		"project":       DefaultProject,
		"targets":       M{},
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func WithTarget(name string, target M) BuildOption {
	return func(m M) {
		targets := m["targets"].(M)
		targets[name] = target
	}
}

func WithCreatedAt(t time.Time) BuildOption {
	return func(m M) { m["created_at"] = t.Format(time.RFC3339) }
}

func WithConfigCommit(sha string) BuildOption {
	return func(m M) { m["config_commit"] = sha }
}

// --- Diagnostic builders ---

type DiagnosticOption func(M)

func Diagnostic(code, level, message string, opts ...DiagnosticOption) M {
	m := M{
		"code":    code,
		"level":   level,
		"message": message,
		"ignored": false,
		"more":    nil,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func WithOASRef(ref string) DiagnosticOption    { return func(m M) { m["oas_ref"] = ref } }
func WithConfigRef(ref string) DiagnosticOption { return func(m M) { m["config_ref"] = ref } }
func WithMore(markdown string) DiagnosticOption {
	return func(m M) { m["more"] = M{"type": "markdown", "markdown": markdown} }
}
