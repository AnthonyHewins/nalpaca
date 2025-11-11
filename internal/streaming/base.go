package streaming

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

type symbolList struct {
	mu      sync.RWMutex
	symbols map[string]struct{}
}

func newSymbolList(x ...string) *symbolList {
	m := make(map[string]struct{}, len(x))
	for _, v := range x {
		m[v] = struct{}{}
	}
	return &symbolList{symbols: m}
}

func (s *symbolList) add(x ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range x {
		s.symbols[v] = struct{}{}
	}
}

func (s *symbolList) del(x ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range x {
		delete(s.symbols, v)
	}
}

func (s *symbolList) list() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	x := make([]string, len(s.symbols))
	i := 0
	for v := range s.symbols {
		x[i] = v
		i++
	}
	return x
}

type Stream struct {
	Feed           string        `env:"FEED_TYPE" envDefault:"sip"`
	Symbols        []string      `env:"SYMBOLS"` // use ',' as delimiter
	BaseURL        string        `env:"BASE_URL" envDefault:"https://stream.data.sandbox.alpaca.markets/v2"`
	Processors     uint16        `env:"PROCESSORS" envDefault:"1"`
	Buffer         uint32        `env:"BUFFER_SIZE" envDefault:"100000"`    // default in lib
	ReconnectLimit uint16        `env:"RECONNECT_LIMIT" envDefault:"20"`    // default in lib
	ReconnectDelay time.Duration `env:"RECONNECT_DELAY" envDefault:"150ms"` // default in lib
}

func streamOpts(key, secret string, logger *slog.Logger, s *Stream) []stream.Option {
	return []stream.Option{
		stream.WithCredentials(key, secret),
		stream.WithProcessors(int(s.Processors)),
		stream.WithBaseURL(s.BaseURL),
		stream.WithBufferSize(int(s.Buffer)), // default value
		stream.WithReconnectSettings(int(s.ReconnectLimit), s.ReconnectDelay),
		stream.WithBufferFillCallback(func(msg []byte) {
			logger.Info("buffer has been filled, processing interrupted", "len(bufferWaiting)", len(msg))
		}),
		stream.WithDisconnectCallback(func() { logger.Warn("stream was disconnected", "url", s.BaseURL) }),
		stream.WithConnectCallback(func() { logger.Info("stream connected", "url", s.BaseURL) }),
		stream.WithLogger(streamLogger{logger.With("alpaca", true)}),
	}
}

type streamLogger struct {
	l *slog.Logger
}

func (l streamLogger) Infof(format string, v ...interface{})  { l.l.Info(fmt.Sprintf(format, v...)) }
func (l streamLogger) Warnf(format string, v ...interface{})  { l.l.Warn(fmt.Sprintf(format, v...)) }
func (l streamLogger) Errorf(format string, v ...interface{}) { l.l.Error(fmt.Sprintf(format, v...)) }
