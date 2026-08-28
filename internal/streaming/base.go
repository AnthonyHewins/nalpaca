package streaming

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"
)

type symbolList struct {
	mu      sync.RWMutex
	symbols map[string]struct{}
}

func newSymbolList(x ...string) symbolList {
	m := make(map[string]struct{}, len(x))
	for _, v := range x {
		m[strings.ToUpper(strings.TrimSpace(v))] = struct{}{}
	}
	return symbolList{symbols: m}
}

func (s *symbolList) clean(x string) string { return strings.ToUpper(strings.TrimSpace(x)) }

func (s *symbolList) add(x ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range x {
		s.symbols[s.clean(v)] = struct{}{}
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

// Stream is the configuration used for particular clients, e.g. if you have a client for stocks or news,
// this config is available for you to use
type Stream struct {
	Version string `env:"VERSION"`

	Feed       string   `env:"FEED_TYPE"`
	Symbols    []string `env:"SYMBOLS"` // use ',' as delimiter
	Processors uint16   `env:"PROCESSORS" envDefault:"1"`

	// Below are options directly on the alpaca socket. These options are passed directly on to the
	// SDK
	SocketBufSize  uint32        `env:"SOCKET_BUFFER_SIZE" envDefault:"100000"` // default in lib
	ReconnectLimit uint16        `env:"RECONNECT_LIMIT" envDefault:"20"`        // default in lib
	ReconnectDelay time.Duration `env:"RECONNECT_DELAY" envDefault:"150ms"`     // default in lib

	// Timeout for nalpaca to send the message
	Timeout time.Duration `env:"TIMEOUT" envDefault:"1s"`
	// This is a separate buffer pool than the socket. This is nalpaca's configured buffer size for when
	// proto messages are serialized, and hence it is smaller and more optimized per message type
	BufSize uint32 `env:"BUFFER_SIZE"`
}

type streamLogger struct {
	l *slog.Logger
}

func (l streamLogger) Infof(format string, v ...interface{})  { l.l.Info(fmt.Sprintf(format, v...)) }
func (l streamLogger) Warnf(format string, v ...interface{})  { l.l.Warn(fmt.Sprintf(format, v...)) }
func (l streamLogger) Errorf(format string, v ...interface{}) { l.l.Error(fmt.Sprintf(format, v...)) }
