package hit

import (
	"context"
	"net/http"
	"sync"
	"time"
)

func produce(ctx context.Context, n int, req *http.Request) <-chan *http.Request {
	out := make(chan *http.Request)

	go func() {
		defer close(out)

		for range n {
			select {
			case out <- req.Clone(ctx):
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func throttle(
	ctx context.Context,
	in <-chan *http.Request,
	delay time.Duration,
) <-chan *http.Request {
	out := make(chan *http.Request)

	go func() {
		defer close(out)

		t := time.NewTicker(delay)
		defer t.Stop()

		for req := range in {
			select {
			case <-t.C:
				out <- req
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func dispatch(
	ctx context.Context,
	in <-chan *http.Request,
	concurrency int,
	send SendFunc,
) <-chan Result {
	var (
		out       chan Result
		waitGroup sync.WaitGroup
	)
	out = make(chan Result)
	waitGroup.Add(concurrency)

	for range concurrency {
		go func() {
			defer waitGroup.Done()
			for req := range in {
				select {
				case out <- send(req):
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		waitGroup.Wait()
		defer close(out)
	}()

	return out
}

func runPipeline(
	ctx context.Context,
	n int,
	req *http.Request,
	opts Options,
) <-chan Result {
	requests := produce(ctx, n, req)
	if opts.RPS > 0 {
		requests = throttle(
			ctx,
			requests,
			time.Second/time.Duration(opts.RPS),
		)
	}

	return dispatch(
		ctx,
		requests,
		opts.Concurrency,
		opts.Send,
	)
}
