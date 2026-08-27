package streaming

import (
	"context"
	"fmt"
	"log/slog"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type News struct {
	list    *symbolList
	logger  *slog.Logger
	n       *stream.NewsClient
	metrics Metrics
	js      publisher
	prefix  string
}

func NewNews(logger *slog.Logger, metrics Metrics, js jetstream.JetStream, prefix, key, secret string, d *Stream) (*News, error) {
	if d == nil {
		return nil, fmt.Errorf("missing stream opts")
	}

	url := d.baseURL(newsURL)

	s := &News{
		list:    newSymbolList(d.Symbols...),
		logger:  logger,
		metrics: metrics,
		js:      js,
		prefix:  prefix,
	}

	so := []stream.NewsOption{}
	for _, v := range streamOpts(key, secret, url, logger, d) {
		so = append(so, v)
	}

	logger.Info("creating news stream client",
		"conf", d,
		"key", key,
		"len(secret)>0", len(secret) > 0,
		"prefix", prefix,
		"url", url,
	)

	s.n = stream.NewNewsClient(append(so, stream.WithNews(s.news, s.list.list()...))...)
	return s, nil
}

func (c *News) news(n stream.News) {
	publish(c.logger, c.metrics, c.js, fmt.Sprintf("%s.%d", c.prefix, n.ID), n, &protoStream.News{
		Id:        uint64(n.ID),
		Symbols:   n.Symbols,
		Headline:  n.Headline,
		Author:    n.Author,
		Summary:   n.Summary,
		Content:   n.Content,
		Url:       n.URL,
		CreatedAt: timestamppb.New(n.CreatedAt),
		UpdatedAt: timestamppb.New(n.UpdatedAt),
	})
}

// Begin consuming data. Cancel context to initiate a shutdown?
// Unsure the underlying implementation, doesnt say in the alpaca docs
func (c *News) Stream(ctx context.Context) error {
	if err := c.n.Connect(ctx); err != nil {
		c.logger.ErrorContext(ctx, "failed establishing news connection", "err", err)
		return err
	}

	if err := <-c.n.Terminated(); err != nil {
		c.logger.Error("news connection terminated with error", "err", err)
		return err
	}

	c.logger.Warn("news connection terminated gracefully")
	return nil
}

func (c *News) ListSubscriptions() []string {
	return c.list.list()
}

func (c *News) AddSubscriptions(x ...string) error {
	return resubscribe(c.logger, "news", c.list, c.n.SubscribeToNews, c.news, add, x)
}

func (c *News) DeleteSubscriptions(x ...string) error {
	return resubscribe(c.logger, "news", c.list, c.n.SubscribeToNews, c.news, del, x)
}
