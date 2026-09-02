package streaming

import "testing"

// subscription_enumer.go is generated, but nothing pinned its round-trip
// behavior - a future manual edit to the Subscription const block without
// regenerating would silently desync String()/SubscriptionString().

func TestSubscriptionStringRoundTrip(t *testing.T) {
	for _, s := range SubscriptionValues() {
		str := s.String()

		got, err := SubscriptionString(str)
		if err != nil {
			t.Errorf("SubscriptionString(%q): unexpected error: %v", str, err)
			continue
		}
		if got != s {
			t.Errorf("round trip: %v -> %q -> %v", s, str, got)
		}
		if !s.IsASubscription() {
			t.Errorf("IsASubscription(%v): want true", s)
		}
	}
}

func TestSubscriptionStringUnknownValue(t *testing.T) {
	var s Subscription = 0
	if s.IsASubscription() {
		t.Error("Subscription(0) should not be a valid value")
	}

	if _, err := SubscriptionString("not_a_real_subscription"); err == nil {
		t.Error("expected an error for an unrecognized subscription name")
	}
}

func TestSubscriptionStringIsCaseInsensitive(t *testing.T) {
	got, err := SubscriptionString("STOCK_BARS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != SubscriptionStockBars {
		t.Errorf("want SubscriptionStockBars, got %v", got)
	}
}
