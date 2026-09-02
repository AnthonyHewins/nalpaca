package streaming

import (
	"strings"
	"sync"
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
		delete(s.symbols, s.clean(v))
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
