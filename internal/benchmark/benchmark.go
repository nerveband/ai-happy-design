// Package benchmark provides types and formatting for provider-agnostic
// performance benchmarking of AI Happy Design workflows.
//
// The benchmark measures three phases:
//   - Phase A: LLM generation (user-reported via --phase-a-ms)
//   - Phase B: CLI batch execution against Figma (measured)
//   - Phase C: Export/verify (optional, measured)
package benchmark

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// PhaseTiming records the duration and metadata for a single benchmark phase.
type PhaseTiming struct {
	Label    string
	Duration time.Duration
	OpsCount int
	Errors   int
	Meta     map[string]string // e.g., "model", "tokens_in"
}

// RunResult holds the timing for a single benchmark run across all phases.
type RunResult struct {
	PhaseA PhaseTiming // LLM generation (user-reported via --phase-a-ms)
	PhaseB PhaseTiming // CLI batch execution (measured)
	PhaseC PhaseTiming // Export/verify (optional)
	Error  error
}

// Total returns the sum of all phase durations.
func (r *RunResult) Total() time.Duration {
	return r.PhaseA.Duration + r.PhaseB.Duration + r.PhaseC.Duration
}

// AggregateResult summarises multiple benchmark runs.
type AggregateResult struct {
	Runs                                   int
	AvgA, AvgB, AvgC, AvgTotal            time.Duration
	StdDev                                 time.Duration
	AvgOps                                 int
	AvgOpsPerSec                           float64
	Errors                                 int
}

// Aggregate computes summary statistics from a slice of RunResults.
func Aggregate(runs []RunResult) AggregateResult {
	n := len(runs)
	if n == 0 {
		return AggregateResult{}
	}

	var sumA, sumB, sumC, sumTotal time.Duration
	var totalOps int
	var totalOpsPerSec float64
	errors := 0

	for _, r := range runs {
		sumA += r.PhaseA.Duration
		sumB += r.PhaseB.Duration
		sumC += r.PhaseC.Duration
		sumTotal += r.Total()
		totalOps += r.PhaseB.OpsCount
		if r.PhaseB.Duration > 0 {
			totalOpsPerSec += float64(r.PhaseB.OpsCount) / r.PhaseB.Duration.Seconds()
		}
		if r.Error != nil {
			errors++
		}
		errors += r.PhaseA.Errors + r.PhaseB.Errors + r.PhaseC.Errors
	}

	avgTotal := sumTotal / time.Duration(n)

	// Standard deviation of total durations
	var varianceSum float64
	for _, r := range runs {
		diff := r.Total() - avgTotal
		varianceSum += float64(diff) * float64(diff)
	}
	stddev := time.Duration(math.Sqrt(varianceSum / float64(n)))

	return AggregateResult{
		Runs:         n,
		AvgA:         sumA / time.Duration(n),
		AvgB:         sumB / time.Duration(n),
		AvgC:         sumC / time.Duration(n),
		AvgTotal:     avgTotal,
		StdDev:       stddev,
		AvgOps:       totalOps / n,
		AvgOpsPerSec: totalOpsPerSec / float64(n),
		Errors:       errors,
	}
}

// fmtDur formats a duration as a human-friendly string (e.g. "4.2s", "120ms").
func fmtDur(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

// String returns a pretty-printed box with benchmark results.
func (a AggregateResult) String() string {
	const w = 48 // inner width (between box borders)
	var b strings.Builder

	hr := strings.Repeat("\u2500", w)
	pad := func(s string) string {
		if len(s) >= w {
			return s[:w]
		}
		return s + strings.Repeat(" ", w-len(s))
	}

	b.WriteString(fmt.Sprintf("\u256d%s\u256e\n", hr))
	b.WriteString(fmt.Sprintf("\u2502%s\u2502\n", pad("  AI Happy Design \u2014 Benchmark Results")))
	b.WriteString(fmt.Sprintf("\u251c%s\u2524\n", hr))
	b.WriteString(fmt.Sprintf("\u2502%s\u2502\n", pad(fmt.Sprintf("  Runs: %d", a.Runs))))

	if a.AvgA > 0 {
		b.WriteString(fmt.Sprintf("\u2502%s\u2502\n", pad(fmt.Sprintf("  Phase A (LLM Gen)     \u2502  avg %s", fmtDur(a.AvgA)))))
	}

	b.WriteString(fmt.Sprintf("\u2502%s\u2502\n", pad(fmt.Sprintf("  Phase B (CLI Exec)    \u2502  avg %s", fmtDur(a.AvgB)))))
	b.WriteString(fmt.Sprintf("\u2502%s\u2502\n", pad(fmt.Sprintf("    \u2514 ops: %d  \u2502  %.1f ops/s", a.AvgOps, a.AvgOpsPerSec))))

	if a.AvgC > 0 {
		b.WriteString(fmt.Sprintf("\u2502%s\u2502\n", pad(fmt.Sprintf("  Phase C (Verify)      \u2502  avg %s", fmtDur(a.AvgC)))))
	}

	b.WriteString(fmt.Sprintf("\u251c%s\u2524\n", hr))
	b.WriteString(fmt.Sprintf("\u2502%s\u2502\n", pad(fmt.Sprintf("  TOTAL                 \u2502  avg %s", fmtDur(a.AvgTotal)))))
	b.WriteString(fmt.Sprintf("\u2502%s\u2502\n", pad(fmt.Sprintf("                        \u2502  \u00b1 %s", fmtDur(a.StdDev)))))
	b.WriteString(fmt.Sprintf("\u2502%s\u2502\n", pad(fmt.Sprintf("  Errors: %d/%d runs", a.Errors, a.Runs))))
	b.WriteString(fmt.Sprintf("\u2570%s\u256f", hr))

	return b.String()
}
