package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
)

type config struct {
	url string
	n   int
	c   int
	rps int
}

func (c *config) validate() error {
	var errs error

	u, err := url.Parse(c.url)
	if err != nil {
		errs = errors.Join(errs, fmt.Errorf("invalid url: %q", c.url))
	}

	if c.url == "" || u.Scheme == "" || u.Host == "" {
		errs = errors.Join(errs, fmt.Errorf("invalid url: %q", c.url))
	}

	if c.n < c.c {
		errs = errors.Join(errs, fmt.Errorf("invalid number of requests: %d, it should be greater than or equal to concurrency level: %d", c.n, c.c))
	}

	if errs != nil {
		return errs
	}

	return nil
}

type positiveIntValue int

func asPositiveIntValue(p *int) *positiveIntValue {
	return (*positiveIntValue)(p)
}

func (i *positiveIntValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}

	if v < 0 {
		return fmt.Errorf("value must be positive")
	}

	*i = positiveIntValue(v)
	return nil
}

func (i *positiveIntValue) String() string {
	return strconv.Itoa(int(*i))
}

func parseArgs(c *config, args []string) error {
	flagSet := flag.NewFlagSet("hit", flag.ContinueOnError)

	flagSet.Usage = func() {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			"usage %s [options] url\n",
			flagSet.Name(),
		)
		flagSet.PrintDefaults()
	}

	flagSet.Var(asPositiveIntValue(&c.n), "n", "number of requests")
	flagSet.Var(asPositiveIntValue(&c.c), "c", "concurrency level")
	flagSet.Var(asPositiveIntValue(&c.rps), "rps", "request per second")

	if err := flagSet.Parse(args); err != nil {
		return err
	}
	c.url = flagSet.Arg(0)

	if err := c.validate(); err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "failed with:\n%v\n", err)
		flagSet.Usage()
		return err
	}

	return nil
}
