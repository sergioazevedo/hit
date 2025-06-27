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

func main() {
	config := config{
		n: 100,
		c: 2,
	}

	if err := parseArgs(&config, os.Args[1:]); err != nil {
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
