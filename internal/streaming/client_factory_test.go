package streaming

import (
	"errors"
	"log/slog"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace/noop"
)

func newClient() *ClientFactory {
	return &ClientFactory{
		prefix:  "",
		key:     "",
		secret:  "",
		baseURL: "",
		js:      nil,
		l:       slog.New(slog.DiscardHandler),
		t:       noop.NewTracerProvider().Tracer(""),
	}
}

func TestNewClientFactoryStoresConstructorArgs(t *testing.T) {
	l := slog.New(slog.DiscardHandler)
	tr := noop.NewTracerProvider().Tracer("")

	c := New("prefix", "https://example.test", "key", "secret", l, nil, tr)

	if c.prefix != "prefix" || c.key != "key" || c.secret != "secret" || c.baseURL != "https://example.test" {
		t.Errorf("unexpected ClientFactory: %+v", c)
	}
}

// prepare() and streamOpts() sit underneath Stocks()/Options()/News() and
// were previously only exercised indirectly through those. Testing them
// directly pins down their contract (nil check, validate() error
// propagation, setDefaults() only running when enabled) independent of any
// particular stream type's config.

type fakeConfig struct {
	enabled        bool
	err            error
	defaultsCalled bool
}

func (f *fakeConfig) validate() (bool, error) { return f.enabled, f.err }
func (f *fakeConfig) setDefaults()            { f.defaultsCalled = true }

func TestPrepareNilConfigReturnsErrMissingOptions(t *testing.T) {
	c := newClient()

	enabled, err := c.prepare(nil)
	if enabled {
		t.Error("expected enabled=false for a nil config")
	}
	if !errors.Is(err, errMissingOptions) {
		t.Errorf("expected errMissingOptions, got %v", err)
	}
}

func TestPrepareValidateErrorPropagates(t *testing.T) {
	c := newClient()
	want := errors.New("boom")

	_, err := c.prepare(&fakeConfig{err: want})
	if !errors.Is(err, want) {
		t.Errorf("expected validate()'s error to propagate, got %v", err)
	}
}

func TestPrepareDisabledSkipsSetDefaults(t *testing.T) {
	c := newClient()
	f := &fakeConfig{enabled: false}

	enabled, err := c.prepare(f)
	if err != nil || enabled {
		t.Fatalf("expected (false, nil), got (%v, %v)", enabled, err)
	}
	if f.defaultsCalled {
		t.Error("setDefaults() should not run when the config is disabled")
	}
}

func TestPrepareEnabledRunsSetDefaults(t *testing.T) {
	c := newClient()
	f := &fakeConfig{enabled: true}

	enabled, err := c.prepare(f)
	if err != nil || !enabled {
		t.Fatalf("expected (true, nil), got (%v, %v)", enabled, err)
	}
	if !f.defaultsCalled {
		t.Error("setDefaults() should run when the config is enabled")
	}
}

func TestStreamOptsBaseOptionCount(t *testing.T) {
	c := newClient()
	opts := c.streamOpts(&StreamConfig{})

	// credentials, reconnect settings, buffer-fill callback, disconnect
	// callback, connect callback, logger - always present regardless of
	// config.
	if len(opts) != 6 {
		t.Errorf("want 6 base options, got %d", len(opts))
	}
}

func TestStreamOptsAddsProcessorsWhenSet(t *testing.T) {
	c := newClient()
	opts := c.streamOpts(&StreamConfig{Processors: 4})

	if len(opts) != 7 {
		t.Errorf("want 7 options with Processors set, got %d", len(opts))
	}
}

func TestStreamOptsAddsSocketBufSizeWhenSet(t *testing.T) {
	c := newClient()
	opts := c.streamOpts(&StreamConfig{SocketBufSize: 1024})

	if len(opts) != 7 {
		t.Errorf("want 7 options with SocketBufSize set, got %d", len(opts))
	}
}

func TestStreamOptsAddsURLWhenSet(t *testing.T) {
	c := newClient()
	opts := c.streamOpts(&StreamConfig{URL: "https://example.test"})

	if len(opts) != 7 {
		t.Errorf("want 7 options with URL set, got %d", len(opts))
	}
}

func TestStreamOptsAddsAllThreeOptionalOptions(t *testing.T) {
	c := newClient()
	opts := c.streamOpts(&StreamConfig{
		Processors:     4,
		SocketBufSize:  1024,
		URL:            "https://example.test",
		ReconnectDelay: time.Second,
	})

	if len(opts) != 9 {
		t.Errorf("want 9 options with all three optional fields set, got %d", len(opts))
	}
}
