package benchmark

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestRunResultTotal(t *testing.T) {
	r := RunResult{
		PhaseA: PhaseTiming{Duration: 4 * time.Second},
		PhaseB: PhaseTiming{Duration: 6 * time.Second},
		PhaseC: PhaseTiming{Duration: 2 * time.Second},
	}
	got := r.Total()
	want := 12 * time.Second
	if got != want {
		t.Errorf("Total() = %v, want %v", got, want)
	}
}

func TestRunResultTotalZeroPhaseC(t *testing.T) {
	r := RunResult{
		PhaseA: PhaseTiming{Duration: 3 * time.Second},
		PhaseB: PhaseTiming{Duration: 5 * time.Second},
	}
	got := r.Total()
	want := 8 * time.Second
	if got != want {
		t.Errorf("Total() = %v, want %v", got, want)
	}
}

func TestAggregateRuns(t *testing.T) {
	runs := []RunResult{
		{
			PhaseA: PhaseTiming{Duration: 4 * time.Second},
			PhaseB: PhaseTiming{Duration: 6 * time.Second, OpsCount: 40},
			PhaseC: PhaseTiming{Duration: 2 * time.Second},
		},
		{
			PhaseA: PhaseTiming{Duration: 4 * time.Second},
			PhaseB: PhaseTiming{Duration: 6 * time.Second, OpsCount: 44},
			PhaseC: PhaseTiming{Duration: 2 * time.Second},
		},
		{
			PhaseA: PhaseTiming{Duration: 4 * time.Second},
			PhaseB: PhaseTiming{Duration: 6 * time.Second, OpsCount: 42},
			PhaseC: PhaseTiming{Duration: 2 * time.Second},
		},
	}

	agg := Aggregate(runs)

	if agg.Runs != 3 {
		t.Errorf("Runs = %d, want 3", agg.Runs)
	}
	if agg.AvgA != 4*time.Second {
		t.Errorf("AvgA = %v, want 4s", agg.AvgA)
	}
	if agg.AvgB != 6*time.Second {
		t.Errorf("AvgB = %v, want 6s", agg.AvgB)
	}
	if agg.AvgC != 2*time.Second {
		t.Errorf("AvgC = %v, want 2s", agg.AvgC)
	}
	if agg.AvgTotal != 12*time.Second {
		t.Errorf("AvgTotal = %v, want 12s", agg.AvgTotal)
	}
	if agg.Errors != 0 {
		t.Errorf("Errors = %d, want 0", agg.Errors)
	}
	// Average ops: (40+44+42)/3 = 42
	if agg.AvgOps != 42 {
		t.Errorf("AvgOps = %d, want 42", agg.AvgOps)
	}
	// StdDev should be 0 since all totals are identical (12s)
	if agg.StdDev != 0 {
		t.Errorf("StdDev = %v, want 0", agg.StdDev)
	}
}

func TestAggregateRunsWithVariance(t *testing.T) {
	runs := []RunResult{
		{
			PhaseA: PhaseTiming{Duration: 3 * time.Second},
			PhaseB: PhaseTiming{Duration: 5 * time.Second, OpsCount: 20},
		},
		{
			PhaseA: PhaseTiming{Duration: 5 * time.Second},
			PhaseB: PhaseTiming{Duration: 7 * time.Second, OpsCount: 20},
		},
	}

	agg := Aggregate(runs)

	if agg.Runs != 2 {
		t.Errorf("Runs = %d, want 2", agg.Runs)
	}
	// Totals are 8s and 12s, avg 10s, stddev = sqrt(((−2)^2 + 2^2)/2) = 2s
	if agg.AvgTotal != 10*time.Second {
		t.Errorf("AvgTotal = %v, want 10s", agg.AvgTotal)
	}
	if agg.StdDev != 2*time.Second {
		t.Errorf("StdDev = %v, want 2s", agg.StdDev)
	}
}

func TestAggregateEmpty(t *testing.T) {
	agg := Aggregate(nil)
	if agg.Runs != 0 {
		t.Errorf("Runs = %d, want 0", agg.Runs)
	}
}

func TestAggregateWithErrors(t *testing.T) {
	runs := []RunResult{
		{
			PhaseA: PhaseTiming{Duration: 1 * time.Second},
			PhaseB: PhaseTiming{Duration: 2 * time.Second, OpsCount: 10, Errors: 1},
		},
		{
			PhaseA: PhaseTiming{Duration: 1 * time.Second},
			PhaseB: PhaseTiming{Duration: 2 * time.Second, OpsCount: 10},
			Error:  fmt.Errorf("connection lost"),
		},
	}

	agg := Aggregate(runs)
	// 1 phase error + 1 run error = 2
	if agg.Errors != 2 {
		t.Errorf("Errors = %d, want 2", agg.Errors)
	}
}

func TestAggregateString(t *testing.T) {
	agg := AggregateResult{
		Runs:         3,
		AvgA:         4200 * time.Millisecond,
		AvgB:         6100 * time.Millisecond,
		AvgC:         1800 * time.Millisecond,
		AvgTotal:     12100 * time.Millisecond,
		StdDev:       600 * time.Millisecond,
		AvgOps:       42,
		AvgOpsPerSec: 6.9,
		Errors:       0,
	}

	s := agg.String()

	checks := []string{
		"Benchmark Results",
		"Runs: 3",
		"Phase A",
		"LLM Gen",
		"4.2s",
		"Phase B",
		"CLI Exec",
		"6.1s",
		"ops: 42",
		"6.9 ops/s",
		"Phase C",
		"Verify",
		"1.8s",
		"TOTAL",
		"12.1s",
		"600ms",
		"Errors: 0/3",
	}

	for _, check := range checks {
		if !strings.Contains(s, check) {
			t.Errorf("String() output missing %q\nGot:\n%s", check, s)
		}
	}
}

func TestAggregateStringNoPhaseA(t *testing.T) {
	agg := AggregateResult{
		Runs:         1,
		AvgB:         3 * time.Second,
		AvgTotal:     3 * time.Second,
		AvgOps:       10,
		AvgOpsPerSec: 3.3,
	}

	s := agg.String()

	// Phase A should be omitted when zero
	if strings.Contains(s, "Phase A") {
		t.Errorf("String() should omit Phase A when AvgA is zero\nGot:\n%s", s)
	}
	// Phase B should still be present
	if !strings.Contains(s, "Phase B") {
		t.Errorf("String() missing Phase B\nGot:\n%s", s)
	}
}

func TestAggregateStringNoPhaseC(t *testing.T) {
	agg := AggregateResult{
		Runs:         1,
		AvgA:         2 * time.Second,
		AvgB:         3 * time.Second,
		AvgTotal:     5 * time.Second,
		AvgOps:       15,
		AvgOpsPerSec: 5.0,
	}

	s := agg.String()

	// Phase C should be omitted when zero
	if strings.Contains(s, "Phase C") {
		t.Errorf("String() should omit Phase C when AvgC is zero\nGot:\n%s", s)
	}
}

func TestFmtDur(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{500 * time.Millisecond, "500ms"},
		{999 * time.Millisecond, "999ms"},
		{1 * time.Second, "1.0s"},
		{4200 * time.Millisecond, "4.2s"},
		{12100 * time.Millisecond, "12.1s"},
	}

	for _, tt := range tests {
		got := fmtDur(tt.d)
		if got != tt.want {
			t.Errorf("fmtDur(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}
