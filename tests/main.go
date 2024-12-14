package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/AnthonyHewins/nalpaca/tests/testcontrol"
	"github.com/nats-io/nats.go"
)

type config struct {
	logger
	url     string
	timeout time.Duration
}

func main() {
	c, err := newTester()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer c.Shutdown()

	var errs []error
	for _, v := range c.Tests() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err = v.Fn(ctx); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", v.Name, err))
		}
	}

	for _, v := range errs {
		fmt.Printf("FAIL %s\n", v)
	}
}

func newTester() (*testcontrol.Controller, error) {
	c := config{}

	flag.BoolVar(&c.quiet, "s", false, "Silent mode: no logs")
	flag.StringVar(&c.fmt, "format", "text", "Log format")
	flag.StringVar(&c.level, "log-level", "error", "Log level")
	flag.BoolVar(&c.src, "log-src", false, "log source line")

	flag.StringVar(&c.url, "server", "localhost:4225", "NATS url to connect to")

	flag.DurationVar(&c.timeout, "timeout", time.Second, "timeout each test after this")

	flag.Parse()

	l, err := c.logger.slog()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nc, err := nats.Connect(c.url)
	if err != nil {
		l.Error("failed connecting to nats", "err", err, "url", c.url)
		return nil, err
	}
	l.Debug("connected to nats", "url", c.url)

	x, err := testcontrol.NewController(l, nc)
	if err != nil {
		nc.Close()
		return nil, err
	}

	return x, nil
}
