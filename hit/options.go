package hit

import (
	"net/http"
	"time"
)

// SendFunc is a type of function that sends a
// [http.Request] and returns a [Result]
type SendFunc func(*http.Request) Result

// Options defines the options for sending requests
// Uses default options for unset fields
type Options struct {
	// Concurrency is the number of concurrent requests to send
	// Defaults to 1
	Concurrency int
	// RPS is the number of requests per second to send
	// Defaults to 0 (unlimited)
	RPS int
	// Send is the function that sends a [http.Request] and returns a [Result]
	// Defaults to [Send]
	Send SendFunc
}

func Defaults() Options {
	return withDefaults(Options{})
}

func withDefaults(opts Options) Options {
	if opts.Concurrency == 0 {
		opts.Concurrency = 1
	}

	if opts.Send == nil {
		client := &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				MaxIdleConnsPerHost: opts.Concurrency,
			},
		}

		opts.Send = func(req *http.Request) Result {
			return Send(client, req)
		}
	}

	return opts
}
