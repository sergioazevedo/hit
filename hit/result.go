package hit

import (
	"iter"
	"time"
)

// Result is perfomance metrics of a single [http.Request]
type Result struct {
	Status   int
	Bytes    int64
	Duration time.Duration
	Error    error
}

type Results iter.Seq[Result]

// Summary is a summary of [Resultt values
type Summary struct {
	Requests    int
	Errors      int
	Bytes       int64
	RPS         float64
	Duration    time.Duration
	Fastest     time.Duration
	Slowest     time.Duration
	SuccessRate float64
}

func Summarize(results Results) Summary {
	if results == nil {
		return Summary{}
	}

	var summary Summary
	startedAt := time.Now()
	for result := range results {
		summary.Requests++
		summary.Bytes += result.Bytes

		if result.Error != nil {
			summary.Errors++
		}

		if summary.Fastest == 0 || result.Duration < summary.Fastest {
			summary.Fastest = result.Duration
		}

		if summary.Slowest == 0 || result.Duration > summary.Slowest {
			summary.Slowest = result.Duration
		}
	}
	summary.Duration = time.Since(startedAt)
	summary.RPS = float64(summary.Requests) / summary.Duration.Seconds()
	if summary.Requests > 0 {
		summary.SuccessRate = float64(summary.Requests-summary.Errors) / float64(summary.Requests) * 100
	}

	return summary
}
