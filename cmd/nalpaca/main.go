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

	Prefix string `env:"PREFIX" envDefault:"nalpaca"`

	ActionStream string `env:"ACTION_STREAM" envDefault:"nalpaca-action-stream-v0"`
	DataStream   string `env:"DATA_STREAM" envDefault:"nalpaca-data-stream-v0"`

	EnableCancel   bool   `env:"ENABLE_CANCELER" envDefault:"false"`
	CancelConsumer string `env:"CANCEL_CONSUMER" envDefault:"nalpaca-cancel-consumer-v0"`

	EnableTradeUpdater bool `env:"ENABLE_TRADE_UPDATER" envDefault:"false"`

	EnableOrders      bool   `env:"ENABLE_ORDERS" envDefault:"false"`
	OrderConsumerName string `env:"ORDER_CONSUMER" envDefault:"nalpaca-orders-consumer-v0"`

	EnableStockStream bool `env:"ENABLE_STOCK_STREAM" envDefault:"false"`

	Bucket string `env:"NATS_KV_BUCKET" envDefault:"nalpaca"`

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
	if a.order.ingestor != nil {
		g.Go(func() error {
			var err error
			a.order.ctx, err = a.order.ingestor.Consume(a.trader.Consume)
			if err != nil {
				a.Logger.ErrorContext(ctx, "failed starting order consumer", "err", err)
			}

			return err
		})
	}

	if a.cancel.ingestor != nil {
		g.Go(func() error {
			var err error
			a.cancel.ctx, err = a.cancel.ingestor.Consume(a.canceler.EventLoop)
			if err != nil {
				a.Logger.ErrorContext(ctx, "failed starting order cancel consumer", "err", err)
			}

			return err
		})
	}

	if a.updater != nil {
		g.Go(func() error {
			a.Logger.InfoContext(ctx, "starting trade updater event loop")
			if err := a.updater.UpdatePositionsKV(ctx); err != nil {
				return err
			}

			return a.updater.TradeUpdateLoop(ctx)
		})
	}

	if a.stockStream != nil {
		g.Go(func() error {
			a.Logger.InfoContext(ctx, "starting stock streaming")
			return a.stockStream.Stream(ctx)
		})
	}

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
