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

type Options struct {
	quotes, trades *symbolList

	logger  *slog.Logger
	s       *stream.OptionClient
	metrics Metrics
	js      publisher
	prefix  string
}

func NewOptions(logger *slog.Logger, metrics Metrics, js jetstream.JetStream, prefix, key, secret string, d *Stream) (*Options, error) {
	if d == nil {
		return nil, fmt.Errorf("missing stream opts")
	}

	switch d.Feed {
	case "opra", "indicative":
	default:
		logger.Error("invalid option feed", "feed", d.Feed)
		return nil, fmt.Errorf("invalid option feed %s: want opra or indicative", d.Feed)
	}

	// Alpaca's options websocket publishes trades and quotes only. Silently
	// ignoring a Bars request would leave someone waiting for data that is never
	// going to arrive.
	if d.Bars {
		logger.Error("bars requested on the options stream, which alpaca does not publish")
		return nil, fmt.Errorf("options stream does not support bars: alpaca publishes only option trades and quotes")
	}

	if !d.Quotes && !d.Trades {
		logger.Error("options stream enabled with no message types")
		return nil, fmt.Errorf("options stream enabled but quotes and trades are both disabled")
	}

	url := d.baseURL(optionsURL)

	s := &Options{
		quotes:  newSymbolList(d.quoteSymbols()...),
		trades:  newSymbolList(d.tradeSymbols()...),
		logger:  logger,
		metrics: metrics,
		js:      js,
		prefix:  prefix,
	}

	oo := []stream.OptionOption{}
	for _, v := range streamOpts(key, secret, url, logger, d) {
		oo = append(oo, v)
	}

	if d.Quotes {
		oo = append(oo, stream.WithOptionQuotes(s.quoteHandler, s.quotes.list()...))
	}
	if d.Trades {
		oo = append(oo, stream.WithOptionTrades(s.tradeHandler, s.trades.list()...))
	}

	logger.Info("creating options stream client",
		"conf", d,
		"key", key,
		"len(secret)>0", len(secret) > 0,
		"prefix", prefix,
		"url", url,
	)

	s.s = stream.NewOptionClient(d.Feed, oo...)
	return s, nil
}

func (c *Options) quoteHandler(q stream.OptionQuote) {
	c.publish(fmt.Sprintf("%s.quotes.%s", c.prefix, q.Symbol), q, &protoStream.OptionQuote{
		Symbol:      q.Symbol,
		BidExchange: q.BidExchange,
		BidPrice:    q.BidPrice,
		BidSize:     q.BidSize,
		AskExchange: q.AskExchange,
		AskPrice:    q.AskPrice,
		AskSize:     q.AskSize,
		Timestamp:   timestamppb.New(q.Timestamp),
		Condition:   q.Condition,
	})
}

func (c *Options) tradeHandler(t stream.OptionTrade) {
	c.publish(fmt.Sprintf("%s.trades.%s", c.prefix, t.Symbol), t, &protoStream.OptionTrade{
		Symbol:    t.Symbol,
		Exchange:  t.Exchange,
		Price:     t.Price,
		Size:      t.Size,
		Timestamp: timestamppb.New(t.Timestamp),
		Condition: t.Condition,
	})
}

func (c *Options) publish(subject string, raw any, msg proto.Message) {
	publish(c.logger, c.metrics, c.js, subject, raw, msg)
}

func (c *Options) ListQuoteSubscriptions() []string { return c.quotes.list() }
func (c *Options) ListTradeSubscriptions() []string { return c.trades.list() }

func (c *Options) AddQuoteSubscriptions(x ...string) error {
	return resubscribe(c.logger, "option quotes", c.quotes, c.s.SubscribeToQuotes, c.quoteHandler, add, x)
}

func (c *Options) DeleteQuoteSubscriptions(x ...string) error {
	return resubscribe(c.logger, "option quotes", c.quotes, c.s.SubscribeToQuotes, c.quoteHandler, del, x)
}

func (c *Options) AddTradeSubscriptions(x ...string) error {
	return resubscribe(c.logger, "option trades", c.trades, c.s.SubscribeToTrades, c.tradeHandler, add, x)
}

func (c *Options) DeleteTradeSubscriptions(x ...string) error {
	return resubscribe(c.logger, "option trades", c.trades, c.s.SubscribeToTrades, c.tradeHandler, del, x)
}

// Begin consuming data. Cancel context to initiate a shutdown?
// Unsure the underlying implementation, doesnt say in the alpaca docs
func (c *Options) Stream(ctx context.Context) error {
	if err := c.s.Connect(ctx); err != nil {
		c.logger.ErrorContext(ctx, "failed establishing options connection", "err", err)
		return err
	}

	if err := <-c.s.Terminated(); err != nil {
		c.logger.Error("options connection terminated with error", "err", err)
		return err
	}

	c.logger.Warn("options connection terminated gracefully")
	return nil
}
