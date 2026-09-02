package streaming

import (
	"context"
	"testing"
)

// Regression test for a bug where setDefaults() fell back to the bogus
// literal "v2" whenever URL was unset. streamOpts() then always passed that
// straight to stream.WithBaseURL, which the alpaca SDK parses as a host-less
// relative URL ("v2/iex") that can never be dialed. Leaving URL empty lets
// the SDK fall back to its own correct default base URL.
func TestStocksSetDefaultsLeavesURLEmpty(t *testing.T) {
	s := &StocksConfig{}
	s.setDefaults()

	if s.URL != "" {
		t.Errorf("expected URL to stay empty so the SDK default applies, got %q", s.URL)
	}
}

func TestStocksSetDefaultsPreservesExplicitURL(t *testing.T) {
	s := &StocksConfig{StreamConfig: StreamConfig{URL: "https://example.test/v9"}}
	s.setDefaults()

	if s.URL != "https://example.test/v9" {
		t.Errorf("explicit URL should be preserved, got %q", s.URL)
	}
}

func TestStocksSetDefaultsFillsBufSizes(t *testing.T) {
	s := &StocksConfig{}
	s.setDefaults()

	if s.Bar.BufSize != defaultBarBufSize {
		t.Errorf("Bar.BufSize: want %d, got %d", defaultBarBufSize, s.Bar.BufSize)
	}
	if s.Quote.BufSize != defaultQuoteBufSize {
		t.Errorf("Quote.BufSize: want %d, got %d", defaultQuoteBufSize, s.Quote.BufSize)
	}
	if s.Trade.BufSize != defaultTradeBufSize {
		t.Errorf("Trade.BufSize: want %d, got %d", defaultTradeBufSize, s.Trade.BufSize)
	}
}

func TestClientFactoryStocksDisabledReturnsZeroValue(t *testing.T) {
	c := newClient()

	m, err := c.Stocks(&StocksConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Conn != nil {
		t.Errorf("expected nil Conn when nothing is enabled, got %v", m.Conn)
	}
}

func TestClientFactoryStocksInvalidFeedReturnsError(t *testing.T) {
	c := newClient()

	_, err := c.Stocks(&StocksConfig{
		Feed: "bogus",
		Bar:  StreamTypeConfig{Enabled: true, Symbols: []string{"AAPL"}},
	})
	if err == nil {
		t.Fatal("expected an error for an invalid feed")
	}
}

func TestClientFactoryStocksEnabledSetsConn(t *testing.T) {
	c := newClient()

	m, err := c.Stocks(&StocksConfig{
		Feed: "iex",
		Bar:  StreamTypeConfig{Enabled: true, Symbols: []string{"AAPL"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Conn == nil {
		t.Fatal("expected Conn to be set when bars are enabled")
	}
}

// Regression test: Stream() must not panic on a manager that was never
// wired up (e.g. stock streaming disabled in config). Before the guard was
// added this dereferenced a nil *stream.StocksClient.
func TestStockSubscriptionManagersStreamNilConnIsNoop(t *testing.T) {
	m := &StockSubscriptionManagers{}

	if err := m.Stream(context.Background()); err != nil {
		t.Errorf("expected nil error for an unconfigured manager, got %v", err)
	}
}
