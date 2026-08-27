package streaming

import (
	"testing"
	"time"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

func newTestStocks(js publisher, m Metrics) *Stocks {
	return &Stocks{
		bars:    newSymbolList(),
		quotes:  newSymbolList(),
		trades:  newSymbolList(),
		logger:  testLogger(),
		metrics: m,
		js:      js,
		prefix:  "nalpaca.data.stocks",
	}
}

// Bars moved from <prefix>.<SYM> to <prefix>.bars.<SYM> so that quotes and trades
// could share the prefix unambiguously. This is a breaking change for existing
// consumers, so it gets an explicit test.
func TestStocksBarSubjectIncludesTypeToken(t *testing.T) {
	js := &mockPublisher{}
	s := newTestStocks(js, testMetrics())

	s.barHandler(stream.Bar{Symbol: "AAPL"})

	got := js.only(t)
	assertSubject(t, got, "nalpaca.data.stocks.bars.AAPL")

	if got.subject == "nalpaca.data.stocks.AAPL" {
		t.Error("bars are still on the old subject; the type token is missing")
	}
}

func TestStocksBarHandlerPublishesEveryField(t *testing.T) {
	js := &mockPublisher{}
	s := newTestStocks(js, testMetrics())

	ts := time.Date(2026, 8, 2, 9, 30, 0, 0, time.UTC)
	s.barHandler(stream.Bar{
		Symbol:     "AAPL",
		Open:       1,
		High:       2,
		Low:        0.5,
		Close:      1.5,
		Volume:     1000,
		Timestamp:  ts,
		TradeCount: 42,
		VWAP:       1.25,
	})

	var b protoStream.Bar
	unmarshal(t, js.only(t), &b)

	if b.Symbol != "AAPL" || b.Open != 1 || b.High != 2 || b.Low != 0.5 || b.Close != 1.5 {
		t.Errorf("ohlc mismatch: %+v", &b)
	}
	if b.Volume != 1000 || b.TradeCount != 42 || b.Vwap != 1.25 {
		t.Errorf("volume/count/vwap mismatch: %+v", &b)
	}
	if !b.Timestamp.AsTime().Equal(ts) {
		t.Errorf("timestamp: want %v, got %v", ts, b.Timestamp.AsTime())
	}
}

func TestStocksQuoteHandlerPublishesEveryField(t *testing.T) {
	js := &mockPublisher{}
	s := newTestStocks(js, testMetrics())

	ts := time.Date(2026, 8, 2, 9, 31, 0, 0, time.UTC)
	s.quoteHandler(stream.Quote{
		Symbol:      "MSFT",
		BidExchange: "A",
		BidPrice:    10.5,
		BidSize:     100,
		AskExchange: "B",
		AskPrice:    10.75,
		AskSize:     200,
		Timestamp:   ts,
		Conditions:  []string{"R", "S"},
		Tape:        "C",
	})

	got := js.only(t)
	assertSubject(t, got, "nalpaca.data.stocks.quotes.MSFT")

	var q protoStream.Quote
	unmarshal(t, got, &q)

	if q.BidPrice != 10.5 || q.AskPrice != 10.75 {
		t.Errorf("prices: bid=%v ask=%v", q.BidPrice, q.AskPrice)
	}
	if q.BidSize != 100 || q.AskSize != 200 {
		t.Errorf("sizes: bid=%d ask=%d", q.BidSize, q.AskSize)
	}
	if q.Tape != "C" {
		t.Errorf("tape: got %q", q.Tape)
	}
	// Conditions is repeated on stocks but a single string on options; make sure
	// every element survives the round trip.
	assertStrings(t, "conditions", q.Conditions, []string{"R", "S"})
}

func TestStocksTradeHandlerPublishesEveryField(t *testing.T) {
	js := &mockPublisher{}
	s := newTestStocks(js, testMetrics())

	ts := time.Date(2026, 8, 2, 9, 32, 0, 0, time.UTC)
	s.tradeHandler(stream.Trade{
		ID:         987654321,
		Symbol:     "SPY",
		Exchange:   "D",
		Price:      500.25,
		Size:       7,
		Timestamp:  ts,
		Conditions: []string{"@", "F", "T"},
		Tape:       "B",
	})

	got := js.only(t)
	assertSubject(t, got, "nalpaca.data.stocks.trades.SPY")

	var tr protoStream.Trade
	unmarshal(t, got, &tr)

	if tr.Id != 987654321 {
		t.Errorf("id: want 987654321, got %d", tr.Id)
	}
	if tr.Symbol != "SPY" || tr.Exchange != "D" || tr.Tape != "B" {
		t.Errorf("symbol/exchange/tape: %q / %q / %q", tr.Symbol, tr.Exchange, tr.Tape)
	}
	if tr.Price != 500.25 || tr.Size != 7 {
		t.Errorf("price/size: %v / %d", tr.Price, tr.Size)
	}
	if !tr.Timestamp.AsTime().Equal(ts) {
		t.Errorf("timestamp: want %v, got %v", ts, tr.Timestamp.AsTime())
	}
	assertStrings(t, "conditions", tr.Conditions, []string{"@", "F", "T"})
}

// A trade with no condition codes is normal; it must not panic or invent one.
func TestStocksTradeHandlerEmptyConditions(t *testing.T) {
	js := &mockPublisher{}
	s := newTestStocks(js, testMetrics())

	s.tradeHandler(stream.Trade{Symbol: "SPY", Conditions: nil})

	var tr protoStream.Trade
	unmarshal(t, js.only(t), &tr)

	if len(tr.Conditions) != 0 {
		t.Errorf("want no conditions, got %v", tr.Conditions)
	}
}

func TestStocksPublishFailureIncrementsCounters(t *testing.T) {
	m := testMetrics()
	s := newTestStocks(failingPublisher(), m)

	s.barHandler(stream.Bar{Symbol: "AAPL"})
	s.quoteHandler(stream.Quote{Symbol: "AAPL"})
	s.tradeHandler(stream.Trade{Symbol: "AAPL"})

	if got := counter(t, m.PubErr); got != 3 {
		t.Errorf("PubErr: want 3, got %v", got)
	}
	if got := counter(t, m.TotalErr); got != 3 {
		t.Errorf("TotalErr: want 3, got %v", got)
	}
}

func TestNewStocksFeedValidation(t *testing.T) {
	for _, tc := range []struct {
		feed    string
		wantErr bool
	}{
		{feed: "sip"},
		{feed: "iex"},
		{feed: "otc"},
		{feed: "delayed_sip"},
		{feed: "opra", wantErr: true},
		{feed: "", wantErr: true},
		{feed: "nonsense", wantErr: true},
	} {
		_, err := NewStocks(testLogger(), testMetrics(), nil, "p", "k", "s", &Stream{
			Feed: tc.feed,
			Bars: true,
		})

		if tc.wantErr && err == nil {
			t.Errorf("feed %q: expected an error, got none", tc.feed)
		}
		if !tc.wantErr && err != nil {
			t.Errorf("feed %q: unexpected error: %v", tc.feed, err)
		}
	}
}

// Enabling the stream without picking a message type would open a websocket that
// receives nothing, which is never what anyone means.
func TestNewStocksRequiresAMessageType(t *testing.T) {
	_, err := NewStocks(testLogger(), testMetrics(), nil, "p", "k", "s", &Stream{Feed: "sip"})
	if err == nil {
		t.Fatal("expected an error when bars, quotes and trades are all disabled")
	}
}

func TestNewStocksNilConf(t *testing.T) {
	if _, err := NewStocks(testLogger(), testMetrics(), nil, "p", "k", "s", nil); err == nil {
		t.Fatal("expected an error for nil stream opts")
	}
}

// Each message type keeps its own subscription list, but falls back to the shared
// Symbols list when it has none of its own.
func TestStocksSymbolListFallback(t *testing.T) {
	d := &Stream{
		Symbols:      []string{"A", "B"},
		QuoteSymbols: []string{"C"},
	}

	assertStrings(t, "bars", d.barSymbols(), []string{"A", "B"})
	assertStrings(t, "quotes", d.quoteSymbols(), []string{"C"})
	assertStrings(t, "trades", d.tradeSymbols(), []string{"A", "B"})
}

func TestStocksSubscriptionListsAreIndependent(t *testing.T) {
	s := newTestStocks(&mockPublisher{}, testMetrics())

	s.bars.add("A")
	s.quotes.add("B", "C")
	s.trades.add("D")

	if got := s.ListBarSubscriptions(); len(got) != 1 || got[0] != "A" {
		t.Errorf("bar subs: want [A], got %v", got)
	}
	if got := len(s.ListQuoteSubscriptions()); got != 2 {
		t.Errorf("quote subs: want 2, got %d", got)
	}
	if got := s.ListTradeSubscriptions(); len(got) != 1 || got[0] != "D" {
		t.Errorf("trade subs: want [D], got %v", got)
	}
}

func TestStocksAppliesDefaultURL(t *testing.T) {
	d := &Stream{Feed: "sip", Bars: true}
	if got := d.baseURL(stocksURL); got != stocksURL {
		t.Errorf("want %q, got %q", stocksURL, got)
	}
}
