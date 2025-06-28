package hit

import (
	"fmt"
	"net/http"
	"time"
)

func Send(_ *http.Client, _ *http.Request) Result {
	const roundTripTime = 100 * time.Millisecond

	time.Sleep(roundTripTime)

	return Result{
		Status:   http.StatusOK,
		Bytes:    10,
		Duration: roundTripTime,
	}
}

// Send N requests using [Send]
// it returns a single-use [Results] iterator that
// pushes a [Result] for each [http.Request] sent.
func SendN(n int, req *http.Request, opts Options) (Results, error) {
	opts = withDefaults(opts)
	if n <= 0 {
		return nil, fmt.Errorf("n must be greater than 0")
	}

	return func(yield func(Result) bool) {
		for range n {
			result := opts.Send(req)
			if !yield(result) {
				return
			}
		}
	}, nil
}
