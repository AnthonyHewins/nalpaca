package streaming

import (
	"testing"
	"time"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

func newTestNews(js publisher, m Metrics) *News {
	return &News{
		list:    newSymbolList(),
		logger:  testLogger(),
		metrics: m,
		js:      js,
		prefix:  "nalpaca.data.news",
	}
}

func TestNewsHandlerPublishesEveryField(t *testing.T) {
	js := &mockPublisher{}
	n := newTestNews(js, testMetrics())

	created := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
	updated := time.Date(2026, 8, 2, 12, 5, 0, 0, time.UTC)

	n.news(stream.News{
		ID:        4242,
		Author:    "somebody",
		CreatedAt: created,
		UpdatedAt: updated,
		Headline:  "a headline",
		Summary:   "a summary",
		Content:   "some content",
		URL:       "https://example.test/story",
		Symbols:   []string{"AAPL", "MSFT"},
	})

	got := js.only(t)
	// News is keyed by article id, not symbol: one article can cover many symbols.
	assertSubject(t, got, "nalpaca.data.news.4242")

	var msg protoStream.News
	unmarshal(t, got, &msg)

	if msg.Id != 4242 {
		t.Errorf("id: want 4242, got %d", msg.Id)
	}
	if msg.Author != "somebody" || msg.Headline != "a headline" {
		t.Errorf("author/headline: %q / %q", msg.Author, msg.Headline)
	}
	if msg.Summary != "a summary" || msg.Content != "some content" {
		t.Errorf("summary/content: %q / %q", msg.Summary, msg.Content)
	}
	if msg.Url != "https://example.test/story" {
		t.Errorf("url: got %q", msg.Url)
	}
	if !msg.CreatedAt.AsTime().Equal(created) || !msg.UpdatedAt.AsTime().Equal(updated) {
		t.Errorf("timestamps: created=%v updated=%v", msg.CreatedAt.AsTime(), msg.UpdatedAt.AsTime())
	}
	assertStrings(t, "symbols", msg.Symbols, []string{"AAPL", "MSFT"})
}

func TestNewsPublishFailureIncrementsCounters(t *testing.T) {
	m := testMetrics()
	n := newTestNews(failingPublisher(), m)

	n.news(stream.News{ID: 1})

	if got := counter(t, m.PubErr); got != 1 {
		t.Errorf("PubErr: want 1, got %v", got)
	}
	if got := counter(t, m.TotalErr); got != 1 {
		t.Errorf("TotalErr: want 1, got %v", got)
	}
}

func TestNewNewsNilConf(t *testing.T) {
	if _, err := NewNews(testLogger(), testMetrics(), nil, "p", "k", "s", nil); err == nil {
		t.Fatal("expected an error for nil stream opts")
	}
}

// News lives on its own endpoint. Before this fix the shared Stream.BaseURL
// defaulted to the stocks v2 URL, so an unset NEWS_STREAM_BASE_URL would have
// quietly pointed the news client at the stocks endpoint.
func TestNewsAppliesDefaultURL(t *testing.T) {
	d := &Stream{}
	if got := d.baseURL(newsURL); got != newsURL {
		t.Errorf("want %q, got %q", newsURL, got)
	}
	if newsURL == stocksURL {
		t.Error("news and stocks must not share an endpoint")
	}
}

func TestNewsSubscriptions(t *testing.T) {
	n := newTestNews(&mockPublisher{}, testMetrics())

	n.list.add("AAPL", "MSFT")
	if got := len(n.ListSubscriptions()); got != 2 {
		t.Errorf("want 2 subscriptions, got %d", got)
	}

	// Empty deltas are a no-op and must not touch alpaca.
	if err := n.AddSubscriptions(); err != nil {
		t.Errorf("empty add should be a no-op, got %v", err)
	}
	if err := n.DeleteSubscriptions(); err != nil {
		t.Errorf("empty delete should be a no-op, got %v", err)
	}
}
