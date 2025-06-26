package main

import "fmt"

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
	fmt.Printf("%s\n%s\n", logo, usage)
}
