package streaming

import (
	"log/slog"

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
