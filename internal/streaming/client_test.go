package streaming

import (
	"log/slog"

	"go.opentelemetry.io/otel/trace/noop"
)

func newClient() *Client {
	return &Client{
		prefix:  "",
		key:     "",
		secret:  "",
		baseURL: "",
		js:      nil,
		l:       slog.New(slog.DiscardHandler),
		t:       noop.NewTracerProvider().Tracer(""),
	}
}
