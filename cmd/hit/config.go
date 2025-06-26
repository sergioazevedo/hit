package main

import (
	"fmt"
	"strconv"
	"strings"
)

type config struct {
	url string
	n   int
	c   int
	rps int
}

func parseArgs(c *config, args []string) error {
	flagSet := map[string]parseFunc{
		"url": stringVar(&c.url),
		"n":   intVar(&c.n),
		"c":   intVar(&c.c),
		"rps": intVar(&c.rps),
	}

	for _, arg := range args {
		name, value, _ := strings.Cut(arg, "=")
		name = strings.TrimPrefix(name, "-")

		parseFunc, ok := flagSet[name]
		if !ok {
			return fmt.Errorf("unknown flag: %s", name)
		}
		if err := parseFunc(value); err != nil {
			return fmt.Errorf("invalid value %q for %s: %w", value, name, err)
		}
	}

	return nil
}

type parseFunc func(string) error

func stringVar(p *string) parseFunc {
	return func(s string) error {
		*p = s
		return nil
	}
}

func intVar(p *int) parseFunc {
	return func(s string) error {
		var err error
		*p, err = strconv.Atoi(s)
		return err
	}
}
