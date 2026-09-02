package streaming

import "testing"

func TestNewMetricsCreatesAllCounters(t *testing.T) {
	m := newMetrics(SubscriptionStockBars)

	if m.publishCount == nil || m.receiveCount == nil || m.transformErr == nil ||
		m.totalErr == nil || m.marshalErr == nil || m.pubErr == nil {
		t.Fatalf("expected all counters to be non-nil: %+v", m)
	}
}

func TestMetricsNilReceiverReturnsNil(t *testing.T) {
	var m *metrics
	if got := m.Metrics(); got != nil {
		t.Errorf("expected nil from a nil *metrics receiver, got %v", got)
	}
}

// Metrics() is what actually gets registered with Prometheus. transformErr
// is created by newMetrics() (and documented as "error count transforming
// alpaca data types to serializable data") but Metrics() never includes it,
// so it can never be scraped even though the field exists on the struct.
// This test pins the current (surprising) behavior rather than silently
// "fixing" it, since it's not clear whether the intent is to drop the field
// or to wire it into Metrics().
func TestMetricsOmitsTransformErr(t *testing.T) {
	m := newMetrics(SubscriptionNews)
	got := m.Metrics()

	if len(got) != 5 {
		t.Errorf("want 5 registered collectors (transformErr is not among them), got %d", len(got))
	}
	for _, c := range got {
		if c == m.transformErr {
			t.Errorf("transformErr should not currently be registered by Metrics(), but it was found in the result")
		}
	}
}
