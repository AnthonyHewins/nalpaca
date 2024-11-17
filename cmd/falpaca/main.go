package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/AnthonyHewins/falpaca/internal/conf"
	"golang.org/x/sync/errgroup"
)

const appName = "falpaca"

type config struct {
	conf.Logger
	conf.Metrics
	conf.Health
	conf.Tracer
	conf.NATS
	conf.Alpaca

	HttpClientTimeout time.Duration `env:"HTTP_CLIENT_TIMEOUT" envDefault:"15s"`

	CacheSize uint16 `env:"CACHE_SIZE" envDefault:"100"`

	StreamPrefix       string          `env:"STREAM_PREFIX" envDefault:""`
	StreamMaxRedeliver uint8           `env:"STREAM_MAX_REDELIVER" envDefault:"3"`
	StreamBackoff      []time.Duration `env:"STREAM_BACKOFF" envDefault:""`
	ProcessingTimeout  time.Duration   `env:"PROCESSING_TIMEOUT" envDefault:"3s"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a, err := newApp(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	if info, ok := debug.ReadBuildInfo(); ok {
		a.logger.InfoContext(ctx,
			"Starting "+appName,
			"version", info.Main.Version,
			"path", info.Main.Path,
			"checksum", info.Main.Sum,
		)
	}

	g, ctx := errgroup.WithContext(ctx)
	a.start(ctx, g)

	select { // watch for signal interruptions or context completion
	case sig := <-interrupt:
		a.logger.Warn("kill signal received", "sig", sig.String())
		cancel()
		break
	case <-ctx.Done():
		a.logger.Warn("context canceled", "err", ctx.Err())
		break
	}

	a.shutdown()

	if err = g.Wait(); err == nil || errors.Is(err, http.ErrServerClosed) {
		return
	}

	a.logger.ErrorContext(ctx, "server goroutines stopped with error", "error", err)
	os.Exit(1)
}

func (a *app) start(ctx context.Context, g *errgroup.Group) {
	g.Go(func() error {
		var err error
		a.consumers.orderCtx, err = a.consumers.order.Consume(a.trader.Consume)
		if err != nil {
			a.logger.ErrorContext(ctx, "failed starting order consumer", "err", err)
		}

		return err
	})

	g.Go(func() error {
		a.logger.InfoContext(ctx, "starting metrics server")
		return a.metrics.ListenAndServe()
	})

	g.Go(func() error {
		a.logger.InfoContext(ctx, "starting health server")
		return a.health.Start(ctx)
	})
}
