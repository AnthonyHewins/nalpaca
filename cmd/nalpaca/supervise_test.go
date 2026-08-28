package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync/atomic"
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

// The whole point of superviseStream: a stream that dies must not take the
// others with it. Before this, every stream ran directly under the errgroup, so
// one returning an error cancelled the shared context and stopped all of them.
//
// This is the regression test for "news and stocks together stopped working".
func TestSupervisedStreamFailureDoesNotKillOthers(t *testing.T) {
	a := testApp()
	g, ctx := errgroup.WithContext(context.Background())

	var survivorSawCancel atomic.Bool

	// Fails immediately, the way a stream with a revoked entitlement or a
	// rejected connection would.
	a.superviseStream(ctx, g, "doomed", func(context.Context) error {
		return errors.New("connection limit exceeded")
	})

	// Healthy stream: blocks until its context is cancelled.
	healthy := make(chan struct{})
	a.superviseStream(ctx, g, "healthy", func(c context.Context) error {
		close(healthy)
		// Outlive the failing stream by enough to observe a cancellation if the
		// two are still coupled; the failure lands essentially instantly.
		select {
		case <-c.Done():
			survivorSawCancel.Store(true)
		case <-time.After(250 * time.Millisecond):
		}
		return nil
	})

	select {
	case <-healthy:
	case <-time.After(time.Second):
		t.Fatal("healthy stream never started")
	}

	if err := g.Wait(); err != nil {
		t.Fatalf("a failing stream must not surface an error from the errgroup, got %v", err)
	}

	if survivorSawCancel.Load() {
		t.Error("the failing stream cancelled the shared context and killed the healthy stream")
	}
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
