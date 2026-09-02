package streaming

import (
	"testing"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

// Direct unit tests of each type's toWire() field mapping and subject()
// formatting, without going through the publish path (see wrap_test.go for
// that). None of these transforms had any coverage before.

func TestBarsToWire(t *testing.T) {
	b := &Bars{}
	ts := time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC)

	got := b.toWire(stream.Bar{
		Symbol:     "AAPL",
		Open:       100,
		High:       110,
		Low:        90,
		Close:      105,
		Volume:     1000,
		Timestamp:  ts,
		TradeCount: 42,
		VWAP:       103.5,
	})

	if got.Symbol != "AAPL" || got.Open != 100 || got.High != 110 || got.Low != 90 ||
		got.Close != 105 || got.Volume != 1000 || got.TradeCount != 42 || got.Vwap != 103.5 {
		t.Errorf("unexpected wire bar: %+v", got)
	}
	if !got.Timestamp.AsTime().Equal(ts) {
		t.Errorf("timestamp: want %v, got %v", ts, got.Timestamp.AsTime())
	}
	if got := b.subject(got); got != "stock_bars.AAPL" {
		t.Errorf("subject: want stock_bars.AAPL, got %q", got)
	}
}

func TestQuotesToWire(t *testing.T) {
	q := &Quotes{}
	ts := time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC)

	got := q.toWire(stream.Quote{
		Symbol:      "AAPL",
		BidExchange: "Q",
		BidPrice:    99.5,
		BidSize:     10,
		AskExchange: "K",
		AskPrice:    100.5,
		AskSize:     20,
		Timestamp:   ts,
		Conditions:  []string{"R"},
		Tape:        "C",
	})

	if got.Symbol != "AAPL" || got.BidExchange != "Q" || got.BidPrice != 99.5 || got.BidSize != 10 ||
		got.AskExchange != "K" || got.AskPrice != 100.5 || got.AskSize != 20 || got.Tape != "C" {
		t.Errorf("unexpected wire quote: %+v", got)
	}
	assertStrings(t, "conditions", got.Conditions, []string{"R"})
	if !got.Timestamp.AsTime().Equal(ts) {
		t.Errorf("timestamp: want %v, got %v", ts, got.Timestamp.AsTime())
	}
	if s := q.subject(got); s != "stock_quotes.AAPL" {
		t.Errorf("subject: want stock_quotes.AAPL, got %q", s)
	}
}

func TestOptionQuotesToWire(t *testing.T) {
	q := &optionQuotes{}
	ts := time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC)

	got := q.toWire(stream.OptionQuote{
		Symbol:      "AAPL240119C00190000",
		BidExchange: "Q",
		BidPrice:    1.1,
		BidSize:     2,
		AskExchange: "K",
		AskPrice:    1.2,
		AskSize:     3,
		Timestamp:   ts,
		Condition:   "regular",
	})

	if got.Symbol != "AAPL240119C00190000" || got.BidExchange != "Q" || got.BidPrice != 1.1 || got.BidSize != 2 ||
		got.AskExchange != "K" || got.AskPrice != 1.2 || got.AskSize != 3 || got.Condition != "regular" {
		t.Errorf("unexpected wire option quote: %+v", got)
	}
	if !got.Timestamp.AsTime().Equal(ts) {
		t.Errorf("timestamp: want %v, got %v", ts, got.Timestamp.AsTime())
	}
	if s := q.subject(got); s != "option_quotes.AAPL240119C00190000" {
		t.Errorf("subject: want option_quotes.AAPL240119C00190000, got %q", s)
	}
}

func TestOptionTradesSubjectUsesOptionTradesPrefix(t *testing.T) {
	tr := &optionTrades{}
	got := tr.toWire(stream.OptionTrade{Symbol: "AAPL240119C00190000"})
	if s := tr.subject(got); s != "option_trades.AAPL240119C00190000" {
		t.Errorf("subject: want option_trades.AAPL240119C00190000, got %q", s)
	}
}

func TestTradesSubjectUsesStockTradesPrefix(t *testing.T) {
	tr := &Trades{}
	got := tr.toWire(stream.Trade{Symbol: "AAPL"})
	if s := tr.subject(got); s != "stock_trades.AAPL" {
		t.Errorf("subject: want stock_trades.AAPL, got %q", s)
	}
}

func TestNewsSubjectIgnoresSymbolArgument(t *testing.T) {
	n := &News{}
	got := n.toWire(stream.News{Symbols: []string{"AAPL", "MSFT"}})
	if s := n.subject(got); s != "news" {
		t.Errorf("subject: want news, got %q", s)
	}
}

func TestNewsProtoLogValueOmitsTimestampsWhenNil(t *testing.T) {
	n := &newsProto{}
	n.Id = 7
	n.Headline = "headline"

	v := n.LogValue()
	attrs := v.Group()

	seen := map[string]bool{}
	for _, a := range attrs {
		seen[a.Key] = true
	}

	if seen["created"] || seen["updated"] {
		t.Errorf("expected no created/updated attrs when timestamps are nil, got %v", attrs)
	}
	if !seen["id"] || !seen["headline"] {
		t.Errorf("expected id and headline attrs, got %v", attrs)
	}
}

func TestNewsProtoLogValueIncludesTimestampsWhenSet(t *testing.T) {
	n := &News{}
	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	np := n.toWire(stream.News{CreatedAt: ts, UpdatedAt: ts})

	v := np.LogValue()
	seen := map[string]bool{}
	for _, a := range v.Group() {
		seen[a.Key] = true
	}
	if !seen["created"] || !seen["updated"] {
		t.Errorf("expected created/updated attrs when timestamps are set, got %v", v.Group())
	}
}
