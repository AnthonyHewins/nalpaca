package streaming

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

// The stocks, news and options streams each register their own error counters
// into one shared registry. Prometheus rejects two collectors that share a
// fully-qualified name, and conf.PrometheusHTTP turns that rejection into a
// startup failure — so a subsystem collision here means the whole app refuses to
// boot the moment metrics are enabled with more than one stream on.
//
// This is a regression test for exactly that: newMetric used to hardcode
// Subsystem: "stocks_stream", making all three sets of counters identical.
func TestMetricsForEveryStreamCanCoexist(t *testing.T) {
	reg := prometheus.NewRegistry()

	for _, subsystem := range []string{"stocks_stream", "news_stream", "options_stream"} {
		for _, c := range NewMetrics("nalpaca", subsystem).Collectors() {
			if err := reg.Register(c); err != nil {
				t.Fatalf("registering %s collectors failed: %v", subsystem, err)
			}
		}
	}

	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather failed: %v", err)
	}

	// 3 streams x 3 counters, each its own metric family.
	if len(families) != 9 {
		names := make([]string, 0, len(families))
		for _, f := range families {
			names = append(names, f.GetName())
		}
		t.Fatalf("want 9 distinct metric families, got %d: %v", len(families), names)
	}
}

// Registering the same subsystem twice must still fail. Without this, the test
// above could pass for the wrong reason (e.g. if registration errors were being
// swallowed somewhere).
func TestMetricsRejectDuplicateSubsystem(t *testing.T) {
	reg := prometheus.NewRegistry()

	for _, c := range NewMetrics("nalpaca", "stocks_stream").Collectors() {
		if err := reg.Register(c); err != nil {
			t.Fatalf("first registration should succeed: %v", err)
		}
	}

	err := reg.Register(NewMetrics("nalpaca", "stocks_stream").Collectors()[0])
	if err == nil {
		t.Fatal("expected a duplicate registration to be rejected")
	}
}

func TestMetricsCollectorsCoversEveryCounter(t *testing.T) {
	m := NewMetrics("nalpaca", "stocks_stream")

	if got := len(m.Collectors()); got != 3 {
		t.Fatalf("want 3 collectors, got %d", got)
	}

	// A nil counter here means Collectors() and the struct have drifted apart.
	for i, c := range m.Collectors() {
		if c == nil {
			t.Errorf("collector %d is nil", i)
		}
	}
}
