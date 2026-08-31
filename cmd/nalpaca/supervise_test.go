package main

import (
	"context"
	"io"
	"log/slog"
	"testing"

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
