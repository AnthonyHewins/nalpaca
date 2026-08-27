package streaming

import (
	"context"
	"fmt"
	"log/slog"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Stocks struct {
	bars, quotes, trades *symbolList

	logger  *slog.Logger
	s       *stream.StocksClient
	metrics Metrics
	js      publisher
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

	if !d.Bars && !d.Quotes && !d.Trades {
		logger.Error("stock stream enabled with no message types")
		return nil, fmt.Errorf("stock stream enabled but bars, quotes and trades are all disabled")
	}

	url := d.baseURL(stocksURL)

	s := &Stocks{
		bars:    newSymbolList(d.barSymbols()...),
		quotes:  newSymbolList(d.quoteSymbols()...),
		trades:  newSymbolList(d.tradeSymbols()...),
		logger:  logger,
		metrics: metrics,
		js:      js,
		prefix:  prefix,
	}

	so := []stream.StockOption{}
	for _, v := range streamOpts(key, secret, url, logger, d) {
		so = append(so, v)
	}

	// Only subscribe to the types that are enabled. Registering a handler for a
	// disabled type would subscribe to the firehose and throw the messages away.
	if d.Bars {
		so = append(so, stream.WithBars(s.barHandler, s.bars.list()...))
	}
	if d.Quotes {
		so = append(so, stream.WithQuotes(s.quoteHandler, s.quotes.list()...))
	}
	if d.Trades {
		so = append(so, stream.WithTrades(s.tradeHandler, s.trades.list()...))
	}

	logger.Info("creating stocks stream client",
		"conf", d,
		"key", key,
		"len(secret)>0", len(secret) > 0,
		"prefix", prefix,
		"url", url,
	)

	s.s = stream.NewStocksClient(d.Feed, so...)
	return s, nil
}

func (c *Stocks) barHandler(b stream.Bar) {
	c.publish(fmt.Sprintf("%s.bars.%s", c.prefix, b.Symbol), b, &protoStream.Bar{
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
}

func (c *Stocks) quoteHandler(q stream.Quote) {
	c.publish(fmt.Sprintf("%s.quotes.%s", c.prefix, q.Symbol), q, &protoStream.Quote{
		Symbol:      q.Symbol,
		BidExchange: q.BidExchange,
		BidPrice:    q.BidPrice,
		BidSize:     q.BidSize,
		AskExchange: q.AskExchange,
		AskPrice:    q.AskPrice,
		AskSize:     q.AskSize,
		Timestamp:   timestamppb.New(q.Timestamp),
		Conditions:  q.Conditions,
		Tape:        q.Tape,
	})
}

func (c *Stocks) tradeHandler(t stream.Trade) {
	c.publish(fmt.Sprintf("%s.trades.%s", c.prefix, t.Symbol), t, &protoStream.Trade{
		Id:         t.ID,
		Symbol:     t.Symbol,
		Exchange:   t.Exchange,
		Price:      t.Price,
		Size:       t.Size,
		Timestamp:  timestamppb.New(t.Timestamp),
		Conditions: t.Conditions,
		Tape:       t.Tape,
	})
}

// publish marshals msg and pushes it to subject. raw is only used for logging.
func (c *Stocks) publish(subject string, raw any, msg proto.Message) {
	publish(c.logger, c.metrics, c.js, subject, raw, msg)
}

func (c *Stocks) ListBarSubscriptions() []string   { return c.bars.list() }
func (c *Stocks) ListQuoteSubscriptions() []string { return c.quotes.list() }
func (c *Stocks) ListTradeSubscriptions() []string { return c.trades.list() }

func (c *Stocks) AddBarSubscriptions(x ...string) error {
	return resubscribe(c.logger, "bars", c.bars, c.s.SubscribeToBars, c.barHandler, add, x)
}

func (c *Stocks) DeleteBarSubscriptions(x ...string) error {
	return resubscribe(c.logger, "bars", c.bars, c.s.SubscribeToBars, c.barHandler, del, x)
}

func (c *Stocks) AddQuoteSubscriptions(x ...string) error {
	return resubscribe(c.logger, "quotes", c.quotes, c.s.SubscribeToQuotes, c.quoteHandler, add, x)
}

func (c *Stocks) DeleteQuoteSubscriptions(x ...string) error {
	return resubscribe(c.logger, "quotes", c.quotes, c.s.SubscribeToQuotes, c.quoteHandler, del, x)
}

func (c *Stocks) AddTradeSubscriptions(x ...string) error {
	return resubscribe(c.logger, "trades", c.trades, c.s.SubscribeToTrades, c.tradeHandler, add, x)
}

func (c *Stocks) DeleteTradeSubscriptions(x ...string) error {
	return resubscribe(c.logger, "trades", c.trades, c.s.SubscribeToTrades, c.tradeHandler, del, x)
}

// Begin consuming data. Cancel context to initiate a shutdown?
// Unsure the underlying implementation, doesnt say in the alpaca docs
func (c *Stocks) Stream(ctx context.Context) error {
	if err := c.s.Connect(ctx); err != nil {
		c.logger.ErrorContext(ctx, "failed establishing stocks connection", "err", err)
		return err
	}

	if err := <-c.s.Terminated(); err != nil {
		c.logger.Error("stocks connection terminated with error", "err", err)
		return err
	}

	c.logger.Warn("stocks connection terminated gracefully")
	return nil
}
