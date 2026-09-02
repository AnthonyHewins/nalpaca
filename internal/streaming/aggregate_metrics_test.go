package streaming

import "testing"

// StockSubscriptionManagers.Metrics() and OptionSubscriptionManagers.Metrics()
// fan out to whichever sub-clients are enabled. Neither had coverage, so a
// wiring mistake (e.g. forgetting to append one sub-client's metrics) would
// go unnoticed.

func TestStockSubscriptionManagersMetricsNilReceiver(t *testing.T) {
	var m *StockSubscriptionManagers
	if got := m.Metrics(); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestStockSubscriptionManagersMetricsOnlyIncludesEnabledSubClients(t *testing.T) {
	c := newClient()
	bars := newBars(c, &StreamTypeConfig{BufSize: 1}, nil)

	m := &StockSubscriptionManagers{Bars: bars}
	got := m.Metrics()

	// Bars alone contributes 5 collectors (see TestMetricsOmitsTransformErr);
	// Quotes/Trades are nil and must contribute nothing.
	if len(got) != 5 {
		t.Errorf("want 5 collectors from Bars alone, got %d", len(got))
	}
}

func TestStockSubscriptionManagersMetricsCombinesAllEnabledSubClients(t *testing.T) {
	c := newClient()
	m := &StockSubscriptionManagers{
		Bars:   newBars(c, &StreamTypeConfig{BufSize: 1}, nil),
		Quotes: newQuotes(c, &StreamTypeConfig{BufSize: 1}, nil),
		Trades: newTrades(c, &StreamTypeConfig{BufSize: 1}, nil),
	}

	got := m.Metrics()
	if len(got) != 15 {
		t.Errorf("want 15 collectors (5 per sub-client x 3), got %d", len(got))
	}
}

func TestOptionSubscriptionManagersMetricsNilReceiver(t *testing.T) {
	var m *OptionSubscriptionManagers
	if got := m.Metrics(); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestOptionSubscriptionManagersMetricsCombinesAllEnabledSubClients(t *testing.T) {
	c := newClient()
	m := &OptionSubscriptionManagers{
		Quotes: newOptionQuotes(c, &StreamTypeConfig{BufSize: 1}, nil),
		Trades: newOptionTrades(c, &StreamTypeConfig{BufSize: 1}, nil),
	}

	got := m.Metrics()
	if len(got) != 10 {
		t.Errorf("want 10 collectors (5 per sub-client x 2), got %d", len(got))
	}
}
