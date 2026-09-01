package streaming

import (
	"context"
	"testing"
)

// Regression test: Options() used to swallow the error returned by
// prepare(), always returning nil regardless. A real misconfiguration
// (here, an invalid feed) was silently treated the same as "streaming
// disabled" instead of failing app startup.
func TestOptionsPrepareErrorIsPropagated(t *testing.T) {
	c := newClient()

	_, err := c.Options(&OptionsConfig{
		Feed:  "bogus",
		Quote: StreamTypeConfig{Enabled: true, Symbols: []string{"AAPL240119C00190000"}},
	})
	if err == nil {
		t.Fatal("expected an error for an invalid feed")
	}
}

func TestOptionsMissingSymbolsIsAnError(t *testing.T) {
	c := newClient()

	_, err := c.Options(&OptionsConfig{
		Feed:  "opra",
		Quote: StreamTypeConfig{Enabled: true}, // no symbols
	})
	if err == nil {
		t.Fatal("expected an error when quotes are enabled with no symbols")
	}
}

func TestOptionsDisabledReturnsZeroValueNoError(t *testing.T) {
	c := newClient()

	m, err := c.Options(&OptionsConfig{Feed: "opra"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Conn != nil {
		t.Errorf("expected nil Conn when disabled, got %v", m.Conn)
	}
}

func TestOptionsEnabledSetsConn(t *testing.T) {
	c := newClient()

	m, err := c.Options(&OptionsConfig{
		Feed:  "opra",
		Quote: StreamTypeConfig{Enabled: true, Symbols: []string{"AAPL240119C00190000"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Conn == nil {
		t.Fatal("expected Conn to be set when quotes are enabled")
	}
}

// Regression test: Stream() must not panic on a manager that was never
// wired up. This is exactly what happened when initOptionStream wasn't
// called from newApp() — a.optionStream stayed the zero value, and Stream()
// dereferenced a nil *stream.OptionClient with no guard.
func TestOptionSubscriptionManagersStreamNilConnIsNoop(t *testing.T) {
	m := &OptionSubscriptionManagers{}

	if err := m.Stream(context.Background()); err != nil {
		t.Errorf("expected nil error for an unconfigured manager, got %v", err)
	}
}
