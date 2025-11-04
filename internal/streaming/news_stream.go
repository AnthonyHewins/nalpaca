package streaming

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type News struct {
	logger  *slog.Logger
	n       *stream.NewsClient
	metrics Metrics
	js      jetstream.JetStream
	prefix  string
}

func NewNews(logger *slog.Logger, metrics Metrics, js jetstream.JetStream, prefix, key, secret string, d *Stream) (*News, error) {
	if d == nil {
		return nil, fmt.Errorf("missing stream opts")
	}

	s := &News{
		logger:  logger,
		metrics: metrics,
		js:      js,
		prefix:  prefix,
	}

	so := []stream.NewsOption{}
	for _, v := range streamOpts(key, secret, logger, d) {
		so = append(so, v)
	}

	symbols := d.Symbols
	logger.Info("creating stocks stream client", "conf", d, "key", key, "len(secret)>0", len(secret) > 0, "prefix", prefix)
	s.n = stream.NewNewsClient(append(so, stream.WithNews(s.news, symbols...))...)
	return s, nil
}

func (c *News) news(n stream.News) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var err error
	defer func() {
		if err != nil {
			c.metrics.TotalErr.Inc()
		}
	}()

	buf, err := proto.Marshal(&protoStream.News{
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
	if err != nil {
		c.logger.ErrorContext(ctx, "failed marshal", "err", err, "raw", n)
		c.metrics.MarshalErr.Inc()
		return
	}

	if _, err = c.js.Publish(ctx, fmt.Sprintf("%s.%d", c.prefix, n.ID), buf); err != nil {
		c.logger.ErrorContext(ctx, "failed publishing", "err", err, "raw", n)
		c.metrics.PubErr.Inc()
	}
}

// Begin consuming data. Cancel context to initiate a shutdown?
// Unsure the underlying implementation, doesnt say in the alpaca docs
func (c *News) Stream(ctx context.Context) error {
	if err := c.n.Connect(ctx); err != nil {
		c.logger.ErrorContext(ctx, "failed establishing stocks connection", "err", err)
		return err
	}

	if err := <-c.n.Terminated(); err != nil {
		c.logger.Error("connection terminated with error", "err", err)
		return err
	}

	c.logger.Warn("connection terminated gracefully")
	return nil
}
