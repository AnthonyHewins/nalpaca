package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/conf"
	"github.com/AnthonyHewins/nalpaca/internal/system"
	"github.com/nats-io/nats.go/jetstream"
	"golang.org/x/sync/errgroup"
)

const appName = "nalpaca"

var version string

type config struct {
	conf.BootstrapConf
	conf.GrpcServerConfWithProxy

	Prefix string `env:"PREFIX" envDefault:"nalpaca"`

	ActionStream string `env:"ACTION_STREAM" envDefault:"nalpaca-action-stream"`
	DataStream   string `env:"DATA_STREAM" envDefault:"nalpaca-data-stream"`

	EnableCancel   bool   `env:"ENABLE_CANCELER" envDefault:"false"`
	CancelConsumer string `env:"CANCEL_CONSUMER" envDefault:"nalpaca-cancel-consumer"`

	EnableTradeUpdater bool `env:"ENABLE_TRADE_UPDATER" envDefault:"false"`

	EnableOrders      bool   `env:"ENABLE_ORDERS" envDefault:"false"`
	OrderConsumerName string `env:"ORDER_CONSUMER" envDefault:"nalpaca-orders-consumer"`

	Bucket string `env:"NATS_KV_BUCKET" envDefault:"nalpaca"`

	ProcessingTimeout time.Duration `env:"PROCESSING_TIMEOUT" envDefault:"3s"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	a, err := newApp(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		a.Logger.InfoContext(ctx,
			"Starting "+appName,
			"version", info.Main.Version,
			"path", info.Main.Path,
			"checksum", info.Main.Sum,
			"codeVersion", system.Version,
			"commit", system.Commit,
			"buildTime", system.BuildTime,
		)
	}

	g, ctx := errgroup.WithContext(ctx)
	a.start(ctx, g)

	<-ctx.Done()
	a.Logger.Warn("context canceled", "err", ctx.Err())
	a.shutdown()

	if err = g.Wait(); err == nil || errors.Is(err, http.ErrServerClosed) {
		return
	}

	a.Logger.ErrorContext(ctx, "server goroutines stopped with error", "err", err)
	switch {
	case errors.Is(err, context.Canceled):
		os.Exit(130)
	case errors.Is(err, io.EOF):
		os.Exit(2)
	case errors.Is(err, context.DeadlineExceeded):
		os.Exit(124)
	default:
		os.Exit(1)
	}
}

func (a *app) superviseStream(ctx context.Context, g *errgroup.Group, name string, run func(context.Context) error) {
	g.Go(func() error {
		a.Logger.InfoContext(ctx, "starting stream", "stream", name)

		if err := run(ctx); err != nil {
			a.Logger.ErrorContext(ctx, "stream_down: stream failed and will not be retried",
				"stream", name,
				"err", err,
			)
			return err
		}

		a.Logger.WarnContext(ctx, "stream_down: stream terminated gracefully", "stream", name)
		return nil
	})
}

func (a *app) start(ctx context.Context, g *errgroup.Group) {
	if a.grpc.Server != nil {
		g.Go(func() error {
			c := net.ListenConfig{KeepAlive: time.Minute * 5}
			ln, err := c.Listen(ctx, "tcp", fmt.Sprintf(":%d", a.grpc.Port))
			if err != nil {
				return err
			}

			a.Logger.InfoContext(ctx, "starting grpc server", "port", a.grpc.Port)
			return a.grpc.Server.Serve(ln)
		})
	}

	if a.grpcProxy != nil {
		a.Logger.InfoContext(ctx, "starting grpc proxy")
		g.Go(a.grpcProxy.ListenAndServe)
	}

	if a.order.ingestor != nil {
		g.Go(func() (err error) {
			a.Logger.InfoContext(ctx, "starting order consumer")
			if a.order.ctx, err = a.order.ingestor.Consume(a.trader.Consume); err != nil {
				a.Logger.ErrorContext(ctx, "failed starting order consumer", "err", err)
			}

			return err
		})
	}

	if a.cancel.ingestor != nil {
		g.Go(func() (err error) {
			a.Logger.InfoContext(ctx, "starting cancel consumer")
			if a.cancel.ctx, err = a.cancel.ingestor.Consume(a.canceler.EventLoop); err != nil {
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

	if a.stockStream.Conn != nil {
		a.superviseStream(ctx, g, "stocks", a.stockStream.Stream)
	}

	if a.news != nil {
		a.superviseStream(ctx, g, "news", a.news.Stream)
	}

	if a.optionStream.Conn != nil {
		a.superviseStream(ctx, g, "options", a.optionStream.Stream)
	}

	if a.Metrics.Server != nil {
		g.Go(func() error {
			a.Logger.InfoContext(ctx, "starting metrics server")
			return a.Metrics.Server.ListenAndServe()
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
	ctx, cancel := context.WithTimeout(context.Background(), a.ShutdownTimeout)
	defer cancel()

	if a.grpc.Server != nil {
		a.grpc.Server.GracefulStop()
	}

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
