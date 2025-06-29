package hit

import (
	"net/http"
	"sync"
	"time"
)

func produce(n int, req *http.Request) <-chan *http.Request {
	out := make(chan *http.Request)

	go func() {
		defer close(out)

		for range n {
			out <- req
		}
	}()

	return out
}

func runPipeline(n int, req *http.Request, opts Options) <-chan Result {
	requests := produce(n, req)
	if opts.RPS > 0 {
		requests = throttle(
			requests,
			time.Second/time.Duration(opts.RPS),
		)
	}

	return dispatch(
		requests,
		opts.Concurrency,
		opts.Send,
	)
}

func throttle(
	in <-chan *http.Request,
	delay time.Duration,
) <-chan *http.Request {
	out := make(chan *http.Request)

	go func() {
		defer close(out)

		t := time.NewTicker(delay)
		for req := range in {
			<-t.C
			out <- req
		}
	}()

	return out
}

func dispatch(
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
				out <- send(req)
			}
		}()
	}

	go func() {
		waitGroup.Wait()
		defer close(out)
	}()

	return out
}
