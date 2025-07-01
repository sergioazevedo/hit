package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/sergioazevedo/hit/hit"
)

const logo = `

 ██╗  ██╗██╗████████╗
 ██║  ██║██║╚══██╔══╝
 ███████║██║   ██║
 ██╔══██║██║   ██║
 ██║  ██║██║   ██║
 ╚═╝  ╚═╝╚═╝   ╚═╝
`

func main() {
	config := config{
		n: 100,
		c: 2,
	}

	if err := parseArgs(&config, os.Args[1:]); err != nil {
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer stop()

	if err := runHit(ctx, config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runHit(ctx context.Context, config config) error {
	req, err := http.NewRequest(http.MethodGet, config.url, http.NoBody)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		os.Exit(1)
	}

	results, err := hit.SendN(
		ctx,
		config.n,
		req,
		hit.Options{
			Concurrency: config.c,
			RPS:         config.rps,
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending requests: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf(
		"%s\n\n Sending %d requests to %q with %d concurrent requests\n",
		logo,
		config.n,
		config.url,
		config.c,
	)

	summary := hit.Summarize(results)
	fmt.Printf(` Summary:
  Success Rate: %.1f%%
  RPS: %.1f
  Requests: %d
  Errors: %d
  Bytes: %d
  Duration: %s
  Fastest: %s
  Slowest: %s
`,
		summary.SuccessRate,
		summary.RPS,
		summary.Requests,
		summary.Errors,
		summary.Bytes,
		summary.Duration,
		summary.Fastest,
		summary.Slowest,
	)

	return ctx.Err()
}
