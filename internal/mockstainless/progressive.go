package mockstainless

import (
	"sync"
	"time"
)

// CallCounter is a thread-safe call counter.
type CallCounter struct {
	mu    sync.Mutex
	count int
}

func (c *CallCounter) Increment() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
	return c.count
}

// ProgressiveBuild returns a build that fills in incrementally over time.
// Each target progresses through 7 phases:
//
//	not_started → codegen → committed → lint → build → test → completed
//
// Target i begins its first phase at Delay*i, and each phase lasts Delay.
// So target i reaches "completed" at Delay*(i+6).
type ProgressiveBuild struct {
	ID            string
	ConfigCommit  string
	Targets       []string     // target names in completion order
	CompletedData map[string]M // final state per target
	Diagnostics   []M          // diagnostics returned for this build
	Delay         time.Duration
	StartTime     time.Time
	mu            sync.Mutex
}

// jitter returns a deterministic pseudo-random duration in [0, max) based on a seed.
// Uses a simple integer hash so the same seed always produces the same jitter.
func jitter(seed int, max time.Duration) time.Duration {
	// mix bits (based on splitmix / murmurhash finalizer)
	h := uint64(seed+1) * 0x9e3779b97f4a7c15
	h ^= h >> 30
	h *= 0xbf58476d1ce4e5b9
	return time.Duration(h % uint64(max))
}

// Snapshot returns the build state at the current time.
// If StartTime is zero the build has not been activated yet and all targets
// are returned as not_started.
func (p *ProgressiveBuild) Snapshot() M {
	p.mu.Lock()
	started := !p.StartTime.IsZero()
	var elapsed time.Duration
	if started {
		elapsed = time.Since(p.StartTime)
	}
	p.mu.Unlock()

	targets := M{}
	for i, name := range p.Targets {
		if !started {
			targets[name] = NotStartedTarget()
			continue
		}
		cd := p.CompletedData[name]
		// Stagger each target by ~1s with deterministic jitter
		offset := time.Duration(i)*time.Second + jitter(i*7, 400*time.Millisecond)
		targets[name] = targetSnapshot(i, elapsed-offset, cd)
	}

	build := Build(p.ID,
		WithConfigCommit(p.ConfigCommit),
		WithCreatedAt(p.StartTime),
	)
	build["targets"] = targets
	return build
}

// targetSnapshot returns the target state based on elapsed time since the target started.
// targetIdx is used to seed deterministic jitter for each step duration.
func targetSnapshot(targetIdx int, elapsed time.Duration, completed M) M {
	if elapsed <= 0 {
		return NotStartedTarget()
	}

	codegenDur := 2*time.Second + jitter(targetIdx*7+1, 600*time.Millisecond)
	if elapsed <= codegenDur {
		return InProgressTarget()
	}

	target := Target("completed", completed["commit"].(M))
	for i, name := range []string{"lint", "build", "test"} {
		step, ok := completed[name].(M)
		if !ok {
			continue
		}
		stepElapsed := elapsed - codegenDur
		stepQueueDelay := time.Second + jitter(targetIdx*7+2+i, 800*time.Millisecond)
		stepFinishDelay := 3*time.Second + jitter(targetIdx*7+2+i, 800*time.Millisecond)
		if stepElapsed < stepQueueDelay {
			target[name] = CheckStepNotStarted()
		} else if stepElapsed < stepQueueDelay+stepFinishDelay {
			target[name] = CheckStepInProgress()
		} else {
			target[name] = step
		}
	}

	return target
}

// Reset restarts the progression from now.
func (p *ProgressiveBuild) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.StartTime = time.Now()
}
