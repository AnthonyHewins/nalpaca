package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/conf"
	"golang.org/x/sync/errgroup"
)

func testApp() *app {
	return &app{Server: &conf.Server{
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}}
}

// A stream that ends cleanly is also not an error for the group.
func TestSupervisedStreamGracefulExit(t *testing.T) {
	a := testApp()
	g, ctx := errgroup.WithContext(context.Background())

	a.superviseStream(ctx, g, "graceful", func(context.Context) error { return nil })

	if err := g.Wait(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSupervisedStreamPropagatesError(t *testing.T) {
	a := testApp()
	g, ctx := errgroup.WithContext(context.Background())

	wantErr := errors.New("boom")
	a.superviseStream(ctx, g, "failing", func(context.Context) error { return wantErr })

	if err := g.Wait(); !errors.Is(err, wantErr) {
		t.Fatalf("want %v, got %v", wantErr, err)
	}
}

// Regression test for the bug where a.optionStream was never initialized
// (initOptionStream was defined but never called from newApp()) combined
// with start() calling superviseStream unconditionally: a zero-value
// StockSubscriptionManagers/OptionSubscriptionManagers/nil News reaching
// Stream() would nil-dereference and crash the whole process at startup,
// even when streaming was intentionally left disabled.
//
// With every optional subsystem left unconfigured, start() must not spawn
// any goroutine that panics or hangs.
func TestStartIsANoopWithNoStreamsConfigured(t *testing.T) {
	a := testApp()
	g, ctx := errgroup.WithContext(context.Background())

	a.start(ctx, g)

	done := make(chan error, 1)
	go func() { done <- g.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("start() spawned a goroutine that never returned; expected a no-op with everything disabled")
	}
}
