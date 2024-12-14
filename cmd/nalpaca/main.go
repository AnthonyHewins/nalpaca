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

	"github.com/AnthonyHewins/nalpaca/internal/conf"
	"golang.org/x/sync/errgroup"
)

const appName = "nalpaca"

var version string

type config struct {
	conf.BootstrapConf
	CancelConf
	OrdersConf
	TradeUpdaterConf

	StreamPrefix string `env:"STREAM_PREFIX" envDefault:"nalpaca"`
	Bucket       string `env:"NATS_KV_BUCKET" envDefault:"nalpaca"`

	ProcessingTimeout time.Duration `env:"PROCESSING_TIMEOUT" envDefault:"3s"`
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
		a.Logger.InfoContext(ctx,
			"Starting "+appName,
			"version", info.Main.Version,
			"path", info.Main.Path,
			"checksum", info.Main.Sum,
			"codeVersion", version,
		)
	}

	g, ctx := errgroup.WithContext(ctx)
	a.start(ctx, g)

	select { // watch for signal interruptions or context completion
	case sig := <-interrupt:
		a.Logger.Warn("kill signal received", "sig", sig.String())
		cancel()
		break
	case <-ctx.Done():
		a.Logger.Warn("context canceled", "err", ctx.Err())
		break
	}

	a.shutdown()

	if err = g.Wait(); err == nil || errors.Is(err, http.ErrServerClosed) {
		return
	}

	a.Logger.ErrorContext(ctx, "server goroutines stopped with error", "error", err)
	os.Exit(1)
}

func (a *app) start(ctx context.Context, g *errgroup.Group) {
	g.Go(func() error {
		var err error
		a.order.ctx, err = a.order.ingestor.Consume(a.trader.Consume)
		if err != nil {
			a.Logger.ErrorContext(ctx, "failed starting order consumer", "err", err)
		}

		return err
	})

	g.Go(func() error {
		var err error
		a.cancel.ctx, err = a.cancel.ingestor.Consume(a.canceler.EventLoop)
		if err != nil {
			a.Logger.ErrorContext(ctx, "failed starting order cancel consumer", "err", err)
		}

		return err
	})

	g.Go(func() error {
		a.Logger.InfoContext(ctx, "starting trade updater event loop")
		return a.updater.TradeUpdateLoop(ctx)
	})

	if a.Metrics != nil {
		g.Go(func() error {
			a.Logger.InfoContext(ctx, "starting metrics server")
			return a.Metrics.ListenAndServe()
		})
	}

	if a.Health != nil {
		g.Go(func() error {
			a.Logger.InfoContext(ctx, "starting health server")
			return a.Health.Start(ctx)
		})
	}
}
