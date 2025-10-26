package exposeport

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/proto"

	"github.com/Warashi/muscat/v2/pb"
)

// RemotePortPolicy determines how the manager reacts to remote port conflicts.
type RemotePortPolicy int

const (
	RemotePortPolicyFail RemotePortPolicy = iota
	RemotePortPolicyNextFree
	RemotePortPolicySkip
)

func (p RemotePortPolicy) String() string {
	switch p {
	case RemotePortPolicyFail:
		return "fail"
	case RemotePortPolicyNextFree:
		return "next-free"
	case RemotePortPolicySkip:
		return "skip"
	default:
		return fmt.Sprintf("remote-port-policy(%d)", int(p))
	}
}

// ParseRemotePortPolicy parses a string into a RemotePortPolicy.
func ParseRemotePortPolicy(s string) (RemotePortPolicy, error) {
	switch strings.ToLower(s) {
	case "", "fail":
		return RemotePortPolicyFail, nil
	case "next-free":
		return RemotePortPolicyNextFree, nil
	case "skip":
		return RemotePortPolicySkip, nil
	default:
		return RemotePortPolicyFail, fmt.Errorf("unknown remote port policy %q", s)
	}
}

// Binding represents a single exposed port pairing.
type Binding struct {
	LocalPort  uint16
	RemotePort uint16
}

// EventHandlers are invoked to emit lifecycle events.
type EventHandlers struct {
	Ready   func(Binding)
	Error   func(Binding, error)
	Stopped func(Binding, error)
}

// ManagerConfig configures the Manager.
type ManagerConfig struct {
	BindAddress string
	Public      bool
	LocalAddr   string
	ChunkSize   int
	Logger      *log.Logger
	Handlers    EventHandlers
	Policy      RemotePortPolicy
}

type exposePortClient interface {
	ExposePort(context.Context, *pb.ExposePortInit, Config) error
}

// Manager manages multiple ExposePort sessions.
type Manager struct {
	client exposePortClient
	cfg    managerConfig

	mu       sync.Mutex
	sessions map[uint16]*sessionHandle
}

type managerConfig struct {
	bindAddress string
	public      bool
	localAddr   string
	chunkSize   int
	logger      *log.Logger
	handlers    EventHandlers
	policy      RemotePortPolicy
}

// NewManager constructs a Manager with sane defaults.
func NewManager(client exposePortClient, cfg ManagerConfig) *Manager {
	mcfg := managerConfig{
		bindAddress: cfg.BindAddress,
		public:      cfg.Public,
		localAddr:   cfg.LocalAddr,
		chunkSize:   cfg.ChunkSize,
		logger:      cfg.Logger,
		handlers:    cfg.Handlers,
		policy:      cfg.Policy,
	}
	if mcfg.logger == nil {
		mcfg.logger = log.Default()
	}
	if mcfg.localAddr == "" {
		mcfg.localAddr = "127.0.0.1"
	}
	return &Manager{
		client:   client,
		cfg:      mcfg,
		sessions: make(map[uint16]*sessionHandle),
	}
}

// Expose starts managing a local port exposure. The provided context controls the lifetime
// of the session and should outlive the manager.
func (m *Manager) Expose(ctx context.Context, localPort uint16, remotePort uint16) error {
	if localPort == 0 {
		return errors.New("local port must be specified")
	}

	sessionCtx, cancel := context.WithCancel(ctx)
	handle := &sessionHandle{
		manager:      m,
		binding:      Binding{LocalPort: localPort, RemotePort: remotePort},
		cancel:       cancel,
		done:         make(chan struct{}),
		policy:       m.cfg.policy,
		logger:       m.cfg.logger,
		chunkSize:    m.cfg.chunkSize,
		localAddress: m.cfg.localAddr,
		bindAddress:  m.cfg.bindAddress,
		public:       m.cfg.public,
	}

	m.mu.Lock()
	if _, exists := m.sessions[localPort]; exists {
		m.mu.Unlock()
		cancel()
		return fmt.Errorf("local port %d is already exposed", localPort)
	}
	m.sessions[localPort] = handle
	m.mu.Unlock()

	go handle.run(sessionCtx)
	return nil
}

// Stop terminates an active exposure for the given local port.
func (m *Manager) Stop(localPort uint16) {
	handle := m.take(localPort)
	if handle == nil {
		return
	}
	handle.cancel()
	<-handle.done
	if m.cfg.handlers.Stopped != nil {
		m.cfg.handlers.Stopped(handle.binding, nil)
	}
}

// Shutdown stops all active sessions.
func (m *Manager) Shutdown() {
	m.mu.Lock()
	handles := make([]*sessionHandle, 0, len(m.sessions))
	for port, handle := range m.sessions {
		delete(m.sessions, port)
		handles = append(handles, handle)
	}
	m.mu.Unlock()

	for _, handle := range handles {
		handle.cancel()
		<-handle.done
	}
}

func (m *Manager) take(localPort uint16) *sessionHandle {
	m.mu.Lock()
	defer m.mu.Unlock()
	handle, ok := m.sessions[localPort]
	if ok {
		delete(m.sessions, localPort)
	}
	return handle
}

func (m *Manager) cleanup(localPort uint16, handle *sessionHandle, err error) {
	removed := false
	m.mu.Lock()
	if current, ok := m.sessions[localPort]; ok && current == handle {
		delete(m.sessions, localPort)
		removed = true
	}
	m.mu.Unlock()

	if removed {
		if m.cfg.handlers.Stopped != nil {
			m.cfg.handlers.Stopped(handle.binding, err)
		}
	}
}

type sessionHandle struct {
	manager *Manager

	binding      Binding
	cancel       context.CancelFunc
	done         chan struct{}
	policy       RemotePortPolicy
	logger       *log.Logger
	chunkSize    int
	localAddress string
	bindAddress  string
	public       bool

	readyOnce       sync.Once
	fallbackAttempt atomic.Bool
}

func (h *sessionHandle) run(ctx context.Context) {
	defer close(h.done)
	var lastErr error
	defer func() {
		if errors.Is(lastErr, context.Canceled) {
			lastErr = nil
		}
		if lastErr == nil && ctx.Err() != nil && !errors.Is(ctx.Err(), context.Canceled) {
			lastErr = ctx.Err()
		}
		h.manager.cleanup(h.binding.LocalPort, h, lastErr)
	}()

	for {
		select {
		case <-ctx.Done():
			lastErr = ctx.Err()
			return
		default:
		}

		initBuilder := pb.ExposePortInit_builder{
			LocalPort: proto.Uint32(uint32(h.binding.LocalPort)),
		}
		if h.binding.RemotePort != 0 {
			initBuilder.RemotePort = proto.Uint32(uint32(h.binding.RemotePort))
		}
		if h.bindAddress != "" {
			initBuilder.BindAddress = proto.String(h.bindAddress)
		}
		if h.public {
			initBuilder.Public = proto.Bool(true)
		}
		init := initBuilder.Build()

		err := h.manager.client.ExposePort(ctx, init, Config{
			LocalAddress: h.localAddress,
			ChunkSize:    h.chunkSize,
			OnOpen:       h.onOpen,
			OnError:      h.onError,
			OnClose:      nil,
		})
		if err == nil || errors.Is(err, context.Canceled) {
			lastErr = err
			return
		}
		if ctx.Err() != nil {
			lastErr = ctx.Err()
			return
		}
		if h.policy == RemotePortPolicyNextFree && !h.fallbackAttempt.Load() && isAddrInUse(err) {
			h.logger.Printf(
				"exposeport: remote port %d unavailable, retrying with next free port",
				h.binding.RemotePort,
			)
			h.fallbackAttempt.Store(true)
			h.binding.RemotePort = 0
			continue
		}
		if h.policy == RemotePortPolicySkip && isAddrInUse(err) {
			if h.manager.cfg.handlers.Error != nil {
				h.manager.cfg.handlers.Error(h.binding, err)
			}
			lastErr = err
			return
		}
		if h.manager.cfg.handlers.Error != nil {
			h.manager.cfg.handlers.Error(h.binding, err)
		}
		lastErr = err
		return
	}
}

func (h *sessionHandle) onOpen(open *pb.ExposePortConnectionOpen) {
	h.readyOnce.Do(func() {
		if port := open.GetRemotePort(); port != 0 {
			h.binding.RemotePort = uint16(port)
		}
		if h.manager.cfg.handlers.Ready != nil {
			h.manager.cfg.handlers.Ready(h.binding)
		} else {
			h.logger.Printf(
				"exposeport: local port %d exposed on remote port %d",
				h.binding.LocalPort,
				h.binding.RemotePort,
			)
		}
	})
}

func (h *sessionHandle) onError(err error) {
	if h.manager.cfg.handlers.Error != nil {
		h.manager.cfg.handlers.Error(h.binding, err)
	} else {
		h.logger.Printf(
			"exposeport: error on local port %d (remote %d): %v",
			h.binding.LocalPort,
			h.binding.RemotePort,
			err,
		)
	}
}

func isAddrInUse(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "address already in use")
}
