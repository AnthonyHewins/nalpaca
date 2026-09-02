package streaming

import (
	"testing"
	"time"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

// wrap() is the core publish path shared by every stream type (bars, trades,
// quotes, option trades/quotes, news): transform -> marshal -> publish, with
// metrics and error handling along the way. Despite mock_test.go providing a
// mockPublisher scaffold specifically for this, nothing exercised it before
// these tests.

func TestWrapPublishesTradeOnSuccess(t *testing.T) {
	c := newClient()
	mp := &mockPublisher{}
	c.js = mp

	tr := newTrades(c, &StreamTypeConfig{BufSize: 64, Timeout: time.Second}, nil)

	ts := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	tr.handler(stream.Trade{
		ID:         42,
		Symbol:     "AAPL",
		Exchange:   "Q",
		Price:      189.5,
		Size:       100,
		Timestamp:  ts,
		Conditions: []string{"@"},
		Tape:       "C",
	})

	p := mp.only(t)
	assertSubject(t, p, "stock_trades.AAPL")

	var got protoStream.Trade
	unmarshal(t, p, &got)

	if got.Id != 42 || got.Symbol != "AAPL" || got.Exchange != "Q" || got.Price != 189.5 || got.Size != 100 || got.Tape != "C" {
		t.Errorf("unexpected wire trade: %+v", &got)
	}
	if !got.Timestamp.AsTime().Equal(ts) {
		t.Errorf("timestamp: want %v, got %v", ts, got.Timestamp.AsTime())
	}
	assertStrings(t, "conditions", got.Conditions, []string{"@"})

	if n := counter(t, tr.metrics.receiveCount); n != 1 {
		t.Errorf("receiveCount: want 1, got %v", n)
	}
	if n := counter(t, tr.metrics.publishCount); n != 1 {
		t.Errorf("publishCount: want 1, got %v", n)
	}
	if n := counter(t, tr.metrics.totalErr); n != 0 {
		t.Errorf("totalErr: want 0, got %v", n)
	}
	if n := counter(t, tr.metrics.pubErr); n != 0 {
		t.Errorf("pubErr: want 0, got %v", n)
	}
}

func TestWrapPublishErrorIncrementsErrCountersNotPublishCount(t *testing.T) {
	c := newClient()
	c.js = failingPublisher()

	tr := newTrades(c, &StreamTypeConfig{BufSize: 64, Timeout: time.Second}, nil)
	tr.handler(stream.Trade{Symbol: "AAPL"})

	if n := counter(t, tr.metrics.receiveCount); n != 1 {
		t.Errorf("receiveCount: want 1, got %v", n)
	}
	if n := counter(t, tr.metrics.totalErr); n != 1 {
		t.Errorf("totalErr: want 1, got %v", n)
	}
	if n := counter(t, tr.metrics.pubErr); n != 1 {
		t.Errorf("pubErr: want 1, got %v", n)
	}
	if n := counter(t, tr.metrics.publishCount); n != 0 {
		t.Errorf("publishCount: want 0 on publish failure, got %v", n)
	}
	if n := counter(t, tr.metrics.marshalErr); n != 0 {
		t.Errorf("marshalErr: want 0 (this wasn't a marshal failure), got %v", n)
	}
}

func TestWrapPublishesOptionTrade(t *testing.T) {
	c := newClient()
	mp := &mockPublisher{}
	c.js = mp

	ot := newOptionTrades(c, &StreamTypeConfig{BufSize: 64, Timeout: time.Second}, nil)

	ts := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	ot.handler(stream.OptionTrade{
		Symbol:    "AAPL240119C00190000",
		Exchange:  "C",
		Price:     1.23,
		Size:      5,
		Timestamp: ts,
		Condition: "regular",
	})

	p := mp.only(t)
	assertSubject(t, p, "option_trades.AAPL240119C00190000")

	var got protoStream.OptionTrade
	unmarshal(t, p, &got)
	if got.Symbol != "AAPL240119C00190000" || got.Exchange != "C" || got.Price != 1.23 || got.Size != 5 || got.Condition != "regular" {
		t.Errorf("unexpected wire option trade: %+v", &got)
	}
}

func TestWrapPublishesBar(t *testing.T) {
	c := newClient()
	mp := &mockPublisher{}
	c.js = mp

	b := newBars(c, &StreamTypeConfig{BufSize: 64, Timeout: time.Second}, nil)

	ts := time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)
	b.handler(stream.Bar{Symbol: "AAPL", Open: 1, High: 2, Low: 0.5, Close: 1.5, Volume: 10, Timestamp: ts, TradeCount: 3, VWAP: 1.4})

	p := mp.only(t)
	assertSubject(t, p, "stock_bars.AAPL")

	var got protoStream.Bar
	unmarshal(t, p, &got)
	if got.Symbol != "AAPL" || got.Open != 1 || got.High != 2 || got.Low != 0.5 || got.Close != 1.5 || got.Volume != 10 || got.TradeCount != 3 || got.Vwap != 1.4 {
		t.Errorf("unexpected wire bar: %+v", &got)
	}
}

func TestWrapPublishesQuote(t *testing.T) {
	c := newClient()
	mp := &mockPublisher{}
	c.js = mp

	q := newQuotes(c, &StreamTypeConfig{BufSize: 64, Timeout: time.Second}, nil)
	q.handler(stream.Quote{Symbol: "AAPL", BidPrice: 1, AskPrice: 2})

	p := mp.only(t)
	assertSubject(t, p, "stock_quotes.AAPL")
}

func TestWrapPublishesOptionQuote(t *testing.T) {
	c := newClient()
	mp := &mockPublisher{}
	c.js = mp

	q := newOptionQuotes(c, &StreamTypeConfig{BufSize: 64, Timeout: time.Second}, nil)
	q.handler(stream.OptionQuote{Symbol: "AAPL240119C00190000", BidPrice: 1, AskPrice: 2})

	p := mp.only(t)
	assertSubject(t, p, "option_quotes.AAPL240119C00190000")
}

func TestWrapPublishesNewsAlwaysUsesNewsSubject(t *testing.T) {
	c := newClient()
	mp := &mockPublisher{}
	c.js = mp

	n := &News{baseClient: newBaseClient[stream.News](SubscriptionNews, c, nil, 64, time.Second)}

	created := time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC)
	updated := created.Add(time.Hour)
	n.handler(stream.News{
		ID:        7,
		Symbols:   []string{"AAPL", "MSFT"},
		Headline:  "headline",
		Author:    "author",
		Summary:   "summary",
		Content:   "content",
		URL:       "https://example.test/news/7",
		CreatedAt: created,
		UpdatedAt: updated,
	})

	p := mp.only(t)
	// unlike the per-symbol subjects, news always publishes to the bare
	// "news" subject regardless of how many symbols the article mentions.
	assertSubject(t, p, "news")

	var got protoStream.News
	unmarshal(t, p, &got)
	if got.Id != 7 || got.Headline != "headline" || got.Author != "author" || got.Summary != "summary" || got.Content != "content" || got.Url != "https://example.test/news/7" {
		t.Errorf("unexpected wire news: %+v", &got)
	}
	assertStrings(t, "symbols", got.Symbols, []string{"AAPL", "MSFT"})
	if !got.CreatedAt.AsTime().Equal(created) {
		t.Errorf("createdAt: want %v, got %v", created, got.CreatedAt.AsTime())
	}
	if !got.UpdatedAt.AsTime().Equal(updated) {
		t.Errorf("updatedAt: want %v, got %v", updated, got.UpdatedAt.AsTime())
	}
}
