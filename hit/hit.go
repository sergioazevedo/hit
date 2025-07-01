package hit

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func Send(client *http.Client, req *http.Request) Result {
	var (
		bytes     int64
		startedAt time.Time
		resp      *http.Response
		err       error
	)

	startedAt = time.Now()
	resp, err = client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		bytes, err = io.Copy(io.Discard, resp.Body)
	}

	return Result{
		Status:   resp.StatusCode,
		Bytes:    bytes,
		Duration: time.Since(startedAt),
		Error:    err,
	}
}

// Send N requests using [Send]
// it returns a single-use [Results] iterator that
// pushes a [Result] for each [http.Request] sent.
func SendN(ctx context.Context, n int, req *http.Request, opts Options) (Results, error) {
	opts = withDefaults(opts)
	if n <= 0 {
		return nil, fmt.Errorf("n must be greater than 0")
	}

	ctx, cancel := context.WithCancel(ctx)

	results := runPipeline(ctx, n, req, opts)

	return func(yield func(Result) bool) {
		defer cancel()
		for result := range results {
			if !yield(result) {
				return
			}
		}
	}, nil
}
