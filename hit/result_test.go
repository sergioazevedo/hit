package hit

import (
	"slices"
	"testing"
	"time"
)

func TestSummarizeFastestResult(t *testing.T) {
	results := []Result{
		{Duration: 100 * time.Millisecond},
		{Duration: 200 * time.Millisecond},
	}

	summary := Summarize(Results(slices.Values(results)))

	if summary.Fastest != 100*time.Millisecond {
		t.Errorf("expected fastest to be 100ms, got %v", summary.Fastest)
	}
}

func TestSummarizeSlowestResult(t *testing.T) {
	results := []Result{
		{Duration: 100 * time.Millisecond},
		{Duration: 300 * time.Millisecond},
	}

	summary := Summarize(Results(slices.Values(results)))

	if summary.Slowest != 300*time.Millisecond {
		t.Errorf("expected slowest to be 300ms, got %v", summary.Slowest)
	}
}

func TestSummarizeNilResults(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("expected no panic, got %v", r)
		}
	}()

	Summarize(nil)
}
