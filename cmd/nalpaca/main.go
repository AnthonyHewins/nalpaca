package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/conf"
	"github.com/nats-io/nats.go/jetstream"
	"golang.org/x/sync/errgroup"
)

const appName = "nalpaca"

var version string

type config struct {
	conf.BootstrapConf

	GrpcPort   uint16 `env:"GRPC_PORT" envDefault:"9200"`
	EnableGrpc bool   `env:"ENABLE_GRPC" envDefault:"false"`

	Prefix string `env:"PREFIX" envDefault:"nalpaca"`

	ActionStream string `env:"ACTION_STREAM" envDefault:"nalpaca-action-stream-v0"`
	DataStream   string `env:"DATA_STREAM" envDefault:"nalpaca-data-stream-v0"`

	EnableCancel   bool   `env:"ENABLE_CANCELER" envDefault:"false"`
	CancelConsumer string `env:"CANCEL_CONSUMER" envDefault:"nalpaca-cancel-consumer-v0"`

	EnableTradeUpdater bool `env:"ENABLE_TRADE_UPDATER" envDefault:"false"`

	EnableOrders      bool   `env:"ENABLE_ORDERS" envDefault:"false"`
	OrderConsumerName string `env:"ORDER_CONSUMER" envDefault:"nalpaca-orders-consumer-v0"`

	EnableStockStream bool `env:"ENABLE_STOCK_STREAM" envDefault:"false"`

	EnableNewsStream bool `env:"ENABLE_NEWS_STREAM" envDefault:"false"`

	Bucket string `env:"NATS_KV_BUCKET" envDefault:"nalpaca"`

	ProcessingTimeout time.Duration `env:"PROCESSING_TIMEOUT" envDefault:"3s"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a, err := newApp(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
	if a.server != nil {
		g.Go(func() error {
			c := net.ListenConfig{KeepAlive: time.Minute * 5}
			ln, err := c.Listen(ctx, "tcp", fmt.Sprintf(":%d", a.grpcPort))
			if err != nil {
				return err
			}

			return a.server.Serve(ln)
		})
	}

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

	if a.newsStream != nil {
		g.Go(func() error {
			a.Logger.InfoContext(ctx, "starting news streaming")
			return a.newsStream.Stream(ctx)
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

func (a *app) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	a.server.GracefulStop()

	type consumers struct {
		name     string
		consumer jetstream.ConsumeContext
	}

	for _, v := range [...]consumers{
		{name: "order consumer", consumer: a.order.ctx},
		{name: "cancel consumer", consumer: a.cancel.ctx},
	} {
		if v.consumer == nil {
			continue
		}

		a.Logger.InfoContext(ctx, "draining "+v.name)
		v.consumer.Drain()
		a.Logger.InfoContext(ctx, "shut down "+v.name)
	}

	a.Server.Shutdown(ctx)
}
