package streaming

import (
	"context"
	"testing"
)

// Regression test: News() used to skip validate() entirely by checking
// d.Enabled and calling d.setDefaults() directly instead of going through
// prepare(). That made NewsConfig.validate()'s "missing symbols" check dead
// code, so enabling news with no symbols configured silently started a
// client subscribed to nothing instead of failing fast.
func TestNewsMissingSymbolsIsAnError(t *testing.T) {
	c := newClient()

	_, err := c.News(&NewsConfig{StreamTypeConfig: StreamTypeConfig{Enabled: true}})
	if err == nil {
		t.Fatal("expected an error when news is enabled with no symbols")
	}
}

func TestNewsDisabledReturnsNilNoError(t *testing.T) {
	c := newClient()

	n, err := c.News(&NewsConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != nil {
		t.Errorf("expected nil News when disabled, got %v", n)
	}
}

func TestNewsEnabledAppliesDefaults(t *testing.T) {
	c := newClient()

	n, err := c.News(&NewsConfig{StreamTypeConfig: StreamTypeConfig{
		Enabled: true,
		Symbols: []string{"AAPL"},
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected a non-nil News client")
	}
}

func TestNewsStreamNilClientIsNoop(t *testing.T) {
	n := &News{}

	if err := n.Stream(context.Background()); err != nil {
		t.Errorf("expected nil error for an unconfigured client, got %v", err)
	}
}
