package streaming

import (
	"context"
	"errors"
	"log/slog"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultNewsBufSize = 100000
)

var (
	_ slog.LogValuer                       = (*newsProto)(nil)
	_ transmitter[stream.News, *newsProto] = (*News)(nil)
	_ Subscriber                           = (*News)(nil)
)

// simple wrapper of the protobuf type because we wanted to override
// logging behavior. printing out the content of the whole news article
// would be dumb, which is what it would do by default, so we implement LogValue()
type newsProto struct{ protoStream.News }

func (n *newsProto) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.Uint64("id", n.Id),
		slog.String("headline", n.Headline),
		slog.String("author", n.Author),
		slog.String("summary", n.Summary),
		slog.String("url", n.Url),
		slog.Any("symbols", n.Symbols),
		slog.Int("len(article)", len(n.Content)),
	}

	if n.CreatedAt != nil {
		attrs = append(attrs, slog.Time("created", n.CreatedAt.AsTime()))
	}

	if n.UpdatedAt != nil {
		attrs = append(attrs, slog.Time("updated", n.UpdatedAt.AsTime()))
	}

	return slog.GroupValue(attrs...)
}

var _ config = (*NewsConfig)(nil)

type NewsConfig struct {
	StreamConfig
	StreamTypeConfig
}

func (n *NewsConfig) validate() (bool, error) {
	if !n.Enabled {
		return false, nil
	}

	if len(n.Symbols) == 0 {
		return false, errors.New("missing symbols for news")
	}

	return true, nil
}

func (n *NewsConfig) setDefaults() {
	if n.BufSize == 0 {
		n.BufSize = defaultNewsBufSize
	}
}

type News struct {
	baseClient[stream.News]
	client *stream.NewsClient
}

func (c *ClientFactory) News(d *NewsConfig) (*News, error) {
	if !d.Enabled {
		return nil, nil
	}

	d.setDefaults()

	n := &News{
		baseClient: newBaseClient[stream.News](
			SubscriptionNews,
			c,
			d.Symbols,
			int(d.BufSize),
			d.Timeout,
		),
	}

	so := []stream.NewsOption{}
	for _, v := range c.streamOpts(&d.StreamConfig) {
		so = append(so, v)
	}

	c.l.Info("creating news stream client", "conf", d)
	n.client = stream.NewNewsClient(append(so, stream.WithNews(n.handler, n.List()...))...)
	return n, nil
}

func (n *News) subject(w *newsProto) string { return SubscriptionNews.String() }

func (n *News) toWire(x stream.News) *newsProto {
	return &newsProto{
		News: protoStream.News{
			Id:        uint64(x.ID),
			Symbols:   x.Symbols,
			Headline:  x.Headline,
			Author:    x.Author,
			Summary:   x.Summary,
			Content:   x.Content,
			Url:       x.URL,
			CreatedAt: timestamppb.New(x.CreatedAt),
			UpdatedAt: timestamppb.New(x.UpdatedAt),
		},
	}
}

func (n *News) Unsubscribe(x ...string) error {
	return n.rmSubscription(n.client.UnsubscribeFromNews, x...)
}

func (n *News) Subscribe(x ...string) error {
	return n.addSubscription(n.client.SubscribeToNews, n.handler, x...)
}

func (n *News) handler(x stream.News) { wrap(n.ClientFactory, n, x) }

// Begin consuming data. Cancel context to initiate a shutdown?
// Unsure the underlying implementation, doesnt say in the alpaca docs
func (n *News) Stream(ctx context.Context) error {
	if n.client == nil {
		return nil
	}

	if err := n.client.Connect(ctx); err != nil {
		return err
	}

	return <-n.client.Terminated()
}

func (n *News) Terminated() <-chan error { return n.client.Terminated() }
