package exposeport

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/Warashi/muscat/v2/pb"
)

func TestManagerExposeReady(t *testing.T) {
	t.Parallel()

	client := &fakeExposePortClient{
		behaviours: []func(context.Context, *pb.ExposePortInit, Config) error{
			func(ctx context.Context, init *pb.ExposePortInit, cfg Config) error {
				if cfg.OnOpen == nil {
					t.Fatalf("OnOpen must be set")
				}
				cfg.OnOpen(
					pb.ExposePortConnectionOpen_builder{
						ConnectionId: proto.Uint64(1),
						RemotePort:   proto.Uint32(9090),
					}.Build(),
				)
				<-ctx.Done()
				return ctx.Err()
			},
		},
	}

	readyCh := make(chan Binding, 1)
	stoppedCh := make(chan error, 1)

	manager := NewManager(client, ManagerConfig{
		Handlers: EventHandlers{
			Ready: func(binding Binding) {
				readyCh <- binding
			},
			Stopped: func(_ Binding, err error) {
				stoppedCh <- err
			},
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := manager.Expose(ctx, 1234, 5678); err != nil {
		t.Fatalf("manager.Expose: %v", err)
	}

	select {
	case binding := <-readyCh:
		if binding.LocalPort != 1234 {
			t.Fatalf("unexpected binding local port: %d", binding.LocalPort)
		}
		if binding.RemotePort != 9090 {
			t.Fatalf("expected remote port 9090, got %d", binding.RemotePort)
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for ready event")
	}

	manager.Stop(1234)

	select {
	case err := <-stoppedCh:
		if err != nil {
			t.Fatalf("expected nil stop error, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for stop event")
	}
}

func TestManagerRemotePolicyNextFree(t *testing.T) {
	t.Parallel()

	var callMu sync.Mutex
	var remotePorts []uint32

	client := &fakeExposePortClient{
		behaviours: []func(context.Context, *pb.ExposePortInit, Config) error{
			func(context.Context, *pb.ExposePortInit, Config) error {
				return errors.New("listen 127.0.0.1:8080: bind: address already in use")
			},
			func(ctx context.Context, init *pb.ExposePortInit, cfg Config) error {
				callMu.Lock()
				remotePorts = append(remotePorts, init.GetRemotePort())
				callMu.Unlock()
				if cfg.OnOpen == nil {
					t.Fatalf("OnOpen must be set")
				}
				cfg.OnOpen(
					pb.ExposePortConnectionOpen_builder{
						ConnectionId: proto.Uint64(1),
						RemotePort:   proto.Uint32(4242),
					}.Build(),
				)
				<-ctx.Done()
				return ctx.Err()
			},
		},
	}

	readyCh := make(chan Binding, 1)

	manager := NewManager(client, ManagerConfig{
		Policy: RemotePortPolicyNextFree,
		Handlers: EventHandlers{
			Ready: func(binding Binding) {
				readyCh <- binding
			},
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := manager.Expose(ctx, 8080, 8080); err != nil {
		t.Fatalf("manager.Expose: %v", err)
	}

	select {
	case binding := <-readyCh:
		if binding.RemotePort != 4242 {
			t.Fatalf("expected remote port 4242, got %d", binding.RemotePort)
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for ready event")
	}

	manager.Stop(8080)

	callMu.Lock()
	defer callMu.Unlock()
	if len(client.calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(client.calls))
	}
	if client.calls[0].GetRemotePort() != 8080 {
		t.Fatalf("first call remote port: got %d want 8080", client.calls[0].GetRemotePort())
	}
	if client.calls[1].GetRemotePort() != 0 {
		t.Fatalf("second call should request port 0, got %d", client.calls[1].GetRemotePort())
	}
	if len(remotePorts) != 1 || remotePorts[0] != 0 {
		t.Fatalf("unexpected recorded remote ports: %#v", remotePorts)
	}
}

func TestManagerRemotePolicySkip(t *testing.T) {
	t.Parallel()

	client := &fakeExposePortClient{
		behaviours: []func(context.Context, *pb.ExposePortInit, Config) error{
			func(context.Context, *pb.ExposePortInit, Config) error {
				return errors.New("listen 127.0.0.1:8080: bind: address already in use")
			},
		},
	}

	errCh := make(chan error, 1)
	manager := NewManager(client, ManagerConfig{
		Policy: RemotePortPolicySkip,
		Handlers: EventHandlers{
			Error: func(_ Binding, err error) {
				errCh <- err
			},
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := manager.Expose(ctx, 8080, 8080); err != nil {
		t.Fatalf("manager.Expose: %v", err)
	}

	select {
	case err := <-errCh:
		if err == nil || !strings.Contains(err.Error(), "address already in use") {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for error event")
	}
}

type fakeExposePortClient struct {
	mu         sync.Mutex
	calls      []*pb.ExposePortInit
	behaviours []func(context.Context, *pb.ExposePortInit, Config) error
}

func (f *fakeExposePortClient) ExposePort(
	ctx context.Context,
	init *pb.ExposePortInit,
	cfg Config,
) error {
	f.mu.Lock()
	index := len(f.calls)
	f.calls = append(f.calls, proto.Clone(init).(*pb.ExposePortInit))
	var behaviour func(context.Context, *pb.ExposePortInit, Config) error
	switch {
	case index < len(f.behaviours):
		behaviour = f.behaviours[index]
	case len(f.behaviours) > 0:
		behaviour = f.behaviours[len(f.behaviours)-1]
	}
	f.mu.Unlock()

	if behaviour != nil {
		return behaviour(ctx, init, cfg)
	}

	<-ctx.Done()
	return ctx.Err()
}
