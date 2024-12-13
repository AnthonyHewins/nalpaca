package main

import (
	"fmt"
	"os"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/conf"
	"github.com/AnthonyHewins/nalpaca/pkg/nalpaca"
	"github.com/caarlos0/env/v11"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type config struct {
	conf.Logger
	conf.NATS

	Timeout time.Duration `env:"TIMEOUT" envDefault:"1s"`
}

func main() {
	c := config{}
	if err := env.Parse(&c); err != nil {
		panic(err)
	}

	l, err := c.Logger.Slog()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nc, err := nats.Connect(c.NATS.URL)
	if err != nil {
		l.Error("failed connecting to nats", "err", err, "url", c.NATS.URL)
		os.Exit(1)
	}
	l.Debug("connected to nats", "url", c.NATS.URL)
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		l.Error("nats connection is not jetstream", "err", err)
		os.Exit(1)
	}

	x := controller{
		nc:      nc,
		timeout: c.Timeout,
		client:  nalpaca.NewClient(js, "nalpaca"),
	}

	for _, fn := range []func() error{
		x.cancel,
	} {
		if err = fn(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
