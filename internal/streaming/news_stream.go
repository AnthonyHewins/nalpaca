package streaming

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const newsComponent = "news"

var (
	_ slog.LogValuer                  = (*news)(nil)
	_ transmitter[stream.News, *news] = (*News)(nil)
)

type news struct {
	protoStream.News
}

func (n *news) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.Uint64("id", n.Id),
		slog.String("headline", n.Headline),
		slog.String("author", n.Author),
		slog.String("summary", n.Summary),
		slog.String("url", n.Url),
		slog.Any("symbols", n.Symbols),
	}

	if n.CreatedAt != nil {
		attrs = append(attrs, slog.Time("created", n.CreatedAt.AsTime()))
	}

	if n.UpdatedAt != nil {
		attrs = append(attrs, slog.Time("updated", n.UpdatedAt.AsTime()))
	}

	return slog.GroupValue(attrs...)
}

type News struct {
	client *Client
	metrics
	sync.Pool
	symbolList

	prefix string
	t      time.Duration

	n *stream.NewsClient
}

func (c *Client) News(d *Stream) (*News, error) {
	if d == nil {
		return nil, fmt.Errorf("missing stream opts")
	}

	n := News{
		client:     c,
		metrics:    newMetrics(newsComponent),
		symbolList: newSymbolList(d.Symbols...),
		prefix:     newsComponent,
		t:          d.Timeout,
		Pool:       sync.Pool{New: func() any { return make([]byte, 0, d.BufSize) }},
	}

	so := []stream.NewsOption{}
	for _, v := range c.streamOpts(d) {
		so = append(so, v)
	}

	symbols := n.list()
	if len(symbols) == 0 {
		return nil, fmt.Errorf("missing symbols for news")
	}

	n.n = stream.NewNewsClient(append(so, stream.WithNews(n.handler, symbols...))...)
	return &n, nil
}

func (c *News) bytePool() *sync.Pool       { return &c.Pool }
func (c *News) component() string          { return newsComponent }
func (c *News) componentMetrics() *metrics { return &c.metrics }
func (c *News) subject(w *news) string     { return newsComponent }
func (c *News) timeout() time.Duration     { return c.t }
func (c *News) toWire(n stream.News) (*news, error) {
	return &news{
		News: protoStream.News{
			Id:        uint64(n.ID),
			Symbols:   n.Symbols,
			Headline:  n.Headline,
			Author:    n.Author,
			Summary:   n.Summary,
			Content:   n.Content,
			Url:       n.URL,
			CreatedAt: timestamppb.New(n.CreatedAt),
			UpdatedAt: timestamppb.New(n.UpdatedAt),
		},
	}, nil
}

// Begin consuming data. Cancel context to initiate a shutdown?
// Unsure the underlying implementation, doesnt say in the alpaca docs
func (c *News) Stream(ctx context.Context) error {
	if err := c.n.Connect(ctx); err != nil {
		c.client.l.ErrorContext(ctx, "failed establishing stocks connection", "err", err)
		return err
	}

	if err := <-c.n.Terminated(); err != nil {
		c.client.l.Error("connection terminated with error", "err", err)
		return err
	}

	c.client.l.Warn("connection terminated gracefully")
	return nil
}

func (c *News) handler(x stream.News) { wrap(c.client, c, x) }

func (c *News) AddSubscriptions(x ...string) error {
	if len(x) == 0 {
		return nil
	}

	c.add(x...)
	l := c.list()
	if err := c.n.SubscribeToNews(c.handler, c.list()...); err != nil {
		c.client.l.Error("failed adding new subscriptions from news stream", "err", err, "wanted", x)
		return err
	}

	c.client.l.Info("added news subscription", "delta", x, "final", l)
	return nil
}

func (c *News) ListSubscriptions() []string { return c.list() }
func (c *News) DeleteSubscriptions(x ...string) error {
	if len(x) == 0 {
		return nil
	}

	c.del(x...)
	l := c.list()
	if err := c.n.SubscribeToNews(c.handler, c.list()...); err != nil {
		c.client.l.Error("failed deleting new subscriptions to news stream", "err", err, "wanted", x)
		return err
	}

	c.client.l.Info("removed news subscription", "delta", x, "final", l)
	return nil
}
