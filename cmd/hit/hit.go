package main

import (
	"fmt"
	"os"
)

const logo = `

	██╗  ██╗██╗████████╗
	██║  ██║██║╚══██╔══╝
	███████║██║   ██║
	██╔══██║██║   ██║
	██║  ██║██║   ██║
	╚═╝  ╚═╝╚═╝   ╚═╝
`

const usage = `
Usage:
  hit [options]

Options:
  -url HTTP server url (required)

	-n number of requests

	-c concurrency level

	-rps request per second
`

func main() {
	config := config{
		n: 100,
		c: 2,
	}

	if err := parseArgs(&config, os.Args[1:]); err != nil {
		fmt.Printf("%s\n%s", err, usage)
		os.Exit(1)
	}

	fmt.Printf(
		"%s\n\n Sending %d requests to %q with %d concurrent requests\n",
		logo,
		config.n,
		config.url,
		config.c,
	)
}
