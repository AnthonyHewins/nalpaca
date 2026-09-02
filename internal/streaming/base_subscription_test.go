package streaming

import (
	"errors"
	"sort"
	"testing"
	"time"
)

// addSubscription/rmSubscription had no direct coverage: symbol normalization
// ("clean"), and — more importantly — whether the in-memory symbol list stays
// consistent with the remote subscription when the SDK call fails.

func newTestBaseClient(t *testing.T) baseClient[string] {
	t.Helper()
	return newBaseClient[string](SubscriptionStockTrades, newClient(), nil, 64, time.Second)
}

func sortedList(b *baseClient[string]) []string {
	l := b.List()
	sort.Strings(l)
	return l
}

func TestAddSubscriptionCleansSymbolsBeforeCallingAddFn(t *testing.T) {
	b := newTestBaseClient(t)

	var gotArgs []string
	addFn := func(_ func(string), s ...string) error {
		gotArgs = append(gotArgs, s...)
		return nil
	}

	if err := b.addSubscription(addFn, func(string) {}, " aapl ", "msft"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sort.Strings(gotArgs)
	assertStrings(t, "addFn args", gotArgs, []string{"AAPL", "MSFT"})
	assertStrings(t, "List()", sortedList(&b), []string{"AAPL", "MSFT"})
}

func TestAddSubscriptionErrorDoesNotUpdateList(t *testing.T) {
	b := newTestBaseClient(t)

	want := errors.New("sdk rejected subscription")
	addFn := func(_ func(string), s ...string) error { return want }

	err := b.addSubscription(addFn, func(string) {}, "AAPL")
	if !errors.Is(err, want) {
		t.Fatalf("expected addFn's error to propagate, got %v", err)
	}

	// The remote subscribe failed, so the local list must not claim AAPL is
	// subscribed - otherwise List() lies about what's actually flowing.
	if list := b.List(); len(list) != 0 {
		t.Errorf("expected empty list after failed subscribe, got %v", list)
	}
}

func TestRmSubscriptionRemovesOnSuccess(t *testing.T) {
	b := newBaseClient[string](SubscriptionStockTrades, newClient(), []string{"AAPL", "MSFT"}, 64, time.Second)

	rmFn := func(s ...string) error { return nil }
	if err := b.rmSubscription(rmFn, "aapl"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertStrings(t, "List()", sortedList(&b), []string{"MSFT"})
}

func TestRmSubscriptionErrorLeavesListUnchanged(t *testing.T) {
	b := newBaseClient[string](SubscriptionStockTrades, newClient(), []string{"AAPL"}, 64, time.Second)

	want := errors.New("sdk rejected unsubscribe")
	rmFn := func(s ...string) error { return want }

	err := b.rmSubscription(rmFn, "AAPL")
	if !errors.Is(err, want) {
		t.Fatalf("expected rmFn's error to propagate, got %v", err)
	}

	// If the remote unsubscribe failed, AAPL is (presumably) still flowing;
	// the local list dropping it anyway would be a false negative.
	assertStrings(t, "List()", sortedList(&b), []string{"AAPL"})
}

func TestRmSubscriptionCleansSymbolBeforeCallingFn(t *testing.T) {
	b := newBaseClient[string](SubscriptionStockTrades, newClient(), []string{"AAPL"}, 64, time.Second)

	var gotArgs []string
	rmFn := func(s ...string) error {
		gotArgs = append(gotArgs, s...)
		return nil
	}

	if err := b.rmSubscription(rmFn, " aapl "); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertStrings(t, "rmFn args", gotArgs, []string{"AAPL"})
}
