package portwatcher

import (
	"context"
	"errors"
	"log"
	"time"
)

const DefaultInterval = 5 * time.Second

type snapshotProvider interface {
	Snapshot(ctx context.Context) (Snapshot, error)
}

// Watcher periodically polls listening sockets and reports the difference to handlers.
type Watcher struct {
	cfg      Config
	provider snapshotProvider
	interval time.Duration
	logger   *log.Logger
}

// Option customises a watcher instance.
type Option func(*Watcher)

// WithInterval overrides the polling interval used by the watcher.
func WithInterval(interval time.Duration) Option {
	return func(w *Watcher) {
		if interval > 0 {
			w.interval = interval
		}
	}
}

// WithLogger injects a custom logger. When omitted the default logger is used.
func WithLogger(logger *log.Logger) Option {
	return func(w *Watcher) {
		if logger != nil {
			w.logger = logger
		}
	}
}

// WithSnapshotProvider replaces the snapshot provider, primarily for testing.
func WithSnapshotProvider(provider snapshotProvider) Option {
	return func(w *Watcher) {
		if provider != nil {
			w.provider = provider
		}
	}
}

// New constructs a Watcher using the provided configuration.
func New(cfg Config, opts ...Option) *Watcher {
	w := &Watcher{
		cfg:      cfg,
		interval: DefaultInterval,
		logger:   log.Default(),
	}
	for _, opt := range opts {
		opt(w)
	}
	if w.provider == nil {
		w.provider = newSystemSnapshotProvider(cfg, w.logger)
	}
	return w
}

// Start begins polling until the context is cancelled. Every detected diff is passed to handler.
func (w *Watcher) Start(ctx context.Context, handler func(Diff)) error {
	if handler == nil {
		return errors.New("portwatcher: handler cannot be nil")
	}
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	var prev Snapshot

	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		snapshot, err := w.provider.Snapshot(ctx)
		if err != nil {
			w.logger.Printf("portwatcher: failed to collect snapshot: %v", err)
		} else {
			diff := snapshot.Diff(prev)
			if !diff.Empty() {
				handler(diff)
			}
			prev = snapshot
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
