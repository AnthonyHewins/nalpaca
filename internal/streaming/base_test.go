package streaming

import (
	"sort"
	"sync"
	"testing"
)

func TestNewSymbolList(t *testing.T) {
	sl := newSymbolList("AAPL", "GOOG")
	expected := []string{"AAPL", "GOOG"}
	got := sl.list()

	sort.Strings(got)
	sort.Strings(expected)

	if len(got) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
	for i := range got {
		if got[i] != expected[i] {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	}
}

func TestAddAndList(t *testing.T) {
	sl := newSymbolList()
	sl.add("MSFT", "TSLA")
	sl.add("NVDA") // add single
	sl.add("MSFT") // duplicate, should not increase size

	expected := []string{"MSFT", "NVDA", "TSLA"}
	got := sl.list()

	sort.Strings(got)
	sort.Strings(expected)

	if len(got) != len(expected) {
		t.Fatalf("length mismatch: want %d, got %d", len(expected), len(got))
	}
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("at index %d: want %s, got %s", i, expected[i], got[i])
		}
	}
}

func TestDelete(t *testing.T) {
	sl := newSymbolList("A", "B", "C")
	sl.del("B")
	sl.del("X") // delete non-existent, should be safe

	expected := []string{"A", "C"}
	got := sl.list()

	sort.Strings(got)
	sort.Strings(expected)

	if len(got) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("mismatch at %d: want %s, got %s", i, expected[i], got[i])
		}
	}
}

func TestDeleteAll(t *testing.T) {
	sl := newSymbolList("X", "Y")
	sl.del("X", "Y")

	if len(sl.list()) != 0 {
		t.Fatal("list should be empty after deleting all")
	}
}

func TestListEmpty(t *testing.T) {
	sl := newSymbolList()
	if len(sl.list()) != 0 {
		t.Fatal("new list should be empty")
	}
}

func TestConcurrentAccess(t *testing.T) {
	sl := newSymbolList()

	var wg sync.WaitGroup
	const goroutines = 100
	const itemsPerRoutine = 100

	// Writers: add symbols
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			for j := 0; j < itemsPerRoutine; j++ {
				symbol := string(rune('A'+offset%26)) + string(rune('0'+j%10))
				sl.add(symbol)
			}
		}(i)
	}

	// Readers: list symbols
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < itemsPerRoutine; j++ {
				_ = sl.list() // should not panic
			}
		}()
	}

	// Deleters
	for i := 0; i < goroutines/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < itemsPerRoutine/2; j++ {
				sl.del("NONEXISTENT") // safe no-op
			}
		}()
	}

	wg.Wait()

	// Final list should not panic and be consistent
	final := sl.list()
	if final == nil {
		t.Fatal("list returned nil")
	}
	// We don't check exact content due to concurrency, but ensure no crash
}

func TestAddDeleteRace(t *testing.T) {
	sl := newSymbolList("START")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			sl.add("ADD")
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			sl.del("ADD")
		}
	}()

	wg.Wait()

	// After race, list should be consistent (either has "START" or "START"+"ADD")
	list := sl.list()
	if len(list) < 1 {
		t.Fatal("list became empty unexpectedly")
	}
}

// Optional: Benchmark
func BenchmarkAdd(b *testing.B) {
	sl := newSymbolList()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.add("BENCH")
	}
}

func BenchmarkList(b *testing.B) {
	sl := newSymbolList("A", "B", "C", "D", "E")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sl.list()
	}
}
