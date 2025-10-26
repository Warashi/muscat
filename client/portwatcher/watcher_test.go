package portwatcher

import (
	"context"
	"errors"
	"io"
	"log"
	"sync"
	"testing"
	"time"
)

type stubSnapshotProvider struct {
	mu       sync.Mutex
	results  []stubSnapshotResult
	position int
}

type stubSnapshotResult struct {
	snapshot Snapshot
	err      error
}

func (s *stubSnapshotProvider) Snapshot(context.Context) (Snapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.results) == 0 {
		return Snapshot{}, nil
	}
	result := s.results[s.position]
	if s.position < len(s.results)-1 {
		s.position++
	}
	return result.snapshot, result.err
}

func TestWatcherStart(t *testing.T) {
	provider := &stubSnapshotProvider{
		results: []stubSnapshotResult{
			{snapshot: Snapshot{}},
			{err: errors.New("transient failure")},
			{snapshot: Snapshot{Entries: []Entry{
				{Addr: "127.0.0.1", Port: 8080, Protocol: "tcp"},
			}}},
		},
	}
	watcher := New(
		Config{},
		WithInterval(10*time.Millisecond),
		WithLogger(log.New(io.Discard, "", 0)),
		WithSnapshotProvider(provider),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		mu    sync.Mutex
		diffs []Diff
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := watcher.Start(ctx, func(d Diff) {
			mu.Lock()
			diffs = append(diffs, d)
			mu.Unlock()
			cancel()
		}); !errors.Is(err, context.Canceled) {
			t.Errorf("Start returned error %v, want context canceled", err)
		}
	}()

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if len(diffs[0].Added) != 1 {
		t.Fatalf("expected 1 added entry, got %d", len(diffs[0].Added))
	}
	if diffs[0].Added[0].Port != 8080 {
		t.Fatalf("expected port 8080, got %d", diffs[0].Added[0].Port)
	}
}

func TestWatcherStartHandlerNil(t *testing.T) {
	watcher := New(Config{})
	if err := watcher.Start(context.Background(), nil); err == nil {
		t.Fatalf("expected error when handler is nil")
	}
}
