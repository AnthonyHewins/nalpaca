package streaming

import (
	"testing"
	"time"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

// newTestOptions builds an Options wired to a mock publisher, skipping the
// constructor so no websocket client is involved.
func newTestOptions(js publisher, m Metrics) *Options {
	return &Options{
		quotes:  newSymbolList(),
		trades:  newSymbolList(),
		logger:  testLogger(),
		metrics: m,
		js:      js,
		prefix:  "nalpaca.data.options",
	}
}

func TestOptionsQuoteHandlerPublishesEveryField(t *testing.T) {
	js := &mockPublisher{}
	o := newTestOptions(js, testMetrics())

	ts := time.Date(2026, 8, 2, 14, 30, 0, 0, time.UTC)
	o.quoteHandler(stream.OptionQuote{
		Symbol:      "AAPL260116C00250000",
		BidExchange: "A",
		BidPrice:    1.23,
		BidSize:     10,
		AskExchange: "B",
		AskPrice:    4.56,
		AskSize:     20,
		Timestamp:   ts,
		Condition:   "R",
	})

	got := js.only(t)
	assertSubject(t, got, "nalpaca.data.options.quotes.AAPL260116C00250000")

	var q protoStream.OptionQuote
	unmarshal(t, got, &q)

	if q.Symbol != "AAPL260116C00250000" {
		t.Errorf("symbol: got %q", q.Symbol)
	}
	if q.BidExchange != "A" || q.AskExchange != "B" {
		t.Errorf("exchanges: got bid=%q ask=%q", q.BidExchange, q.AskExchange)
	}
	if q.BidPrice != 1.23 || q.AskPrice != 4.56 {
		t.Errorf("prices: got bid=%v ask=%v", q.BidPrice, q.AskPrice)
	}
	if q.BidSize != 10 || q.AskSize != 20 {
		t.Errorf("sizes: got bid=%d ask=%d", q.BidSize, q.AskSize)
	}
	if !q.Timestamp.AsTime().Equal(ts) {
		t.Errorf("timestamp: want %v, got %v", ts, q.Timestamp.AsTime())
	}
	if q.Condition != "R" {
		t.Errorf("condition: got %q", q.Condition)
	}
}

func TestOptionsTradeHandlerPublishesEveryField(t *testing.T) {
	js := &mockPublisher{}
	o := newTestOptions(js, testMetrics())

	ts := time.Date(2026, 8, 2, 14, 31, 0, 0, time.UTC)
	o.tradeHandler(stream.OptionTrade{
		Symbol:    "SPY260116P00400000",
		Exchange:  "C",
		Price:     7.89,
		Size:      3,
		Timestamp: ts,
		Condition: "I",
	})

	got := js.only(t)
	assertSubject(t, got, "nalpaca.data.options.trades.SPY260116P00400000")

	var tr protoStream.OptionTrade
	unmarshal(t, got, &tr)

	if tr.Symbol != "SPY260116P00400000" || tr.Exchange != "C" {
		t.Errorf("symbol/exchange: got %q / %q", tr.Symbol, tr.Exchange)
	}
	if tr.Price != 7.89 || tr.Size != 3 {
		t.Errorf("price/size: got %v / %d", tr.Price, tr.Size)
	}
	if !tr.Timestamp.AsTime().Equal(ts) {
		t.Errorf("timestamp: want %v, got %v", ts, tr.Timestamp.AsTime())
	}
	if tr.Condition != "I" {
		t.Errorf("condition: got %q", tr.Condition)
	}
}

func TestOptionsPublishFailureIncrementsCounters(t *testing.T) {
	m := testMetrics()
	o := newTestOptions(failingPublisher(), m)

	// Must not panic even though the publish fails.
	o.quoteHandler(stream.OptionQuote{Symbol: "AAPL260116C00250000"})

	if got := counter(t, m.PubErr); got != 1 {
		t.Errorf("PubErr: want 1, got %v", got)
	}
	if got := counter(t, m.TotalErr); got != 1 {
		t.Errorf("TotalErr: want 1, got %v", got)
	}
	if got := counter(t, m.MarshalErr); got != 0 {
		t.Errorf("MarshalErr: want 0, got %v", got)
	}
}

func TestNewOptionsFeedValidation(t *testing.T) {
	for _, tc := range []struct {
		feed    string
		wantErr bool
	}{
		{feed: "opra"},
		{feed: "indicative"},
		{feed: "sip", wantErr: true},
		{feed: "iex", wantErr: true},
		{feed: "", wantErr: true},
		{feed: "nonsense", wantErr: true},
	} {
		_, err := NewOptions(testLogger(), testMetrics(), nil, "p", "k", "s", &Stream{
			Feed:   tc.feed,
			Quotes: true,
		})

		if tc.wantErr && err == nil {
			t.Errorf("feed %q: expected an error, got none", tc.feed)
		}
		if !tc.wantErr && err != nil {
			t.Errorf("feed %q: unexpected error: %v", tc.feed, err)
		}
	}
}

// Alpaca publishes no option bars, so asking for them is a config mistake that
// should fail loudly rather than leaving someone waiting for data forever.
func TestNewOptionsRejectsBars(t *testing.T) {
	_, err := NewOptions(testLogger(), testMetrics(), nil, "p", "k", "s", &Stream{
		Feed:   "indicative",
		Bars:   true,
		Quotes: true,
	})

	if err == nil {
		t.Fatal("expected an error when bars are requested on the options stream")
	}
}

func TestNewOptionsRequiresAMessageType(t *testing.T) {
	_, err := NewOptions(testLogger(), testMetrics(), nil, "p", "k", "s", &Stream{Feed: "opra"})
	if err == nil {
		t.Fatal("expected an error when neither quotes nor trades are enabled")
	}
}

func TestNewOptionsNilConf(t *testing.T) {
	if _, err := NewOptions(testLogger(), testMetrics(), nil, "p", "k", "s", nil); err == nil {
		t.Fatal("expected an error for nil stream opts")
	}
}

func TestNewOptionsAppliesDefaultURL(t *testing.T) {
	d := &Stream{Feed: "opra", Quotes: true}
	if got := d.baseURL(optionsURL); got != optionsURL {
		t.Errorf("want %q, got %q", optionsURL, got)
	}

	d.BaseURL = "https://example.test/v9"
	if got := d.baseURL(optionsURL); got != "https://example.test/v9" {
		t.Errorf("explicit BaseURL should win, got %q", got)
	}
}

func TestOptionsSubscriptionListsAreIndependent(t *testing.T) {
	o := newTestOptions(&mockPublisher{}, testMetrics())

	o.quotes.add("A", "B")
	o.trades.add("C")

	if got := len(o.ListQuoteSubscriptions()); got != 2 {
		t.Errorf("quote subs: want 2, got %d", got)
	}
	if got := o.ListTradeSubscriptions(); len(got) != 1 || got[0] != "C" {
		t.Errorf("trade subs: want [C], got %v", got)
	}
}
