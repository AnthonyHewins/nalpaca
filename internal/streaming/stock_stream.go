package streaming

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newMetric(appName, name, help string) prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: appName,
		Subsystem: "stocks_stream",
		Name:      name,
		Help:      help,
	})
}

type Metrics struct {
	TotalErr, MarshalErr, PubErr prometheus.Counter
}

func NewMetrics(appName string) Metrics {
	return Metrics{
		TotalErr:   newMetric(appName, "total_err", "total error count"),
		MarshalErr: newMetric(appName, "marshal_err", "marshal error count"),
		PubErr:     newMetric(appName, "pub_err", "nats pub err count"),
	}
}

type Stocks struct {
	logger  *slog.Logger
	s       *stream.StocksClient
	metrics Metrics
	js      jetstream.JetStream
	prefix  string
}

func NewStocks(logger *slog.Logger, metrics Metrics, js jetstream.JetStream, prefix, key, secret string, d *Stream) (*Stocks, error) {
	if d == nil {
		return nil, fmt.Errorf("missing stream opts")
	}

	switch d.Feed {
	case "sip", "iex", "otc", "delayed_sip":
	default:
		logger.Error("invalid feed", "feed", d.Feed)
		return nil, fmt.Errorf("invalid feed %s", d.Feed)
	}

	s := &Stocks{
		logger:  logger,
		metrics: metrics,
		js:      js,
		prefix:  prefix,
	}

	so := []stream.StockOption{}
	for _, v := range streamOpts(key, secret, logger, d) {
		so = append(so, v)
	}

	symbols := d.Symbols
	logger.Info("creating stocks stream client", "conf", d, "key", key, "len(secret)>0", len(secret) > 0, "prefix", prefix)
	s.s = stream.NewStocksClient(d.Feed, append(so, stream.WithBars(s.bars, symbols...))...)
	return s, nil
}

func (c *Stocks) bars(b stream.Bar) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var err error
	defer func() {
		if err != nil {
			c.metrics.TotalErr.Inc()
		}
	}()

	buf, err := proto.Marshal(&protoStream.Bar{
		Symbol:     b.Symbol,
		Open:       b.Open,
		High:       b.High,
		Low:        b.Low,
		Close:      b.Close,
		Volume:     b.Volume,
		Timestamp:  timestamppb.New(b.Timestamp),
		TradeCount: b.TradeCount,
		Vwap:       b.VWAP,
	})
	if err != nil {
		c.logger.ErrorContext(ctx, "failed marshal", "err", err, "raw", b)
		c.metrics.MarshalErr.Inc()
		return
	}

	if _, err = c.js.Publish(ctx, fmt.Sprintf("%s.%s", c.prefix, b.Symbol), buf); err != nil {
		c.logger.ErrorContext(ctx, "failed publishing", "err", err, "raw", b)
		c.metrics.PubErr.Inc()
	}
}

// Begin consuming data. Cancel context to initiate a shutdown?
// Unsure the underlying implementation, doesnt say in the alpaca docs
func (c *Stocks) Stream(ctx context.Context) error {
	if err := c.s.Connect(ctx); err != nil {
		c.logger.ErrorContext(ctx, "failed establishing stocks connection", "err", err)
		return err
	}

	if err := <-c.s.Terminated(); err != nil {
		c.logger.Error("connection terminated with error", "err", err)
		return err
	}

	c.logger.Warn("connection terminated gracefully")
	return nil
}
