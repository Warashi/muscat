package exposeport

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/Warashi/muscat/v2/pb"
	testexposeport "github.com/Warashi/muscat/v2/testutil/exposeport"
)

func TestSessionDataFlow(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := testexposeport.NewFakeServerStream(ctx)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen: %v", err)
	}
	var used int32
	cfg := Config{
		Listen: func(network, address string) (net.Listener, error) {
			if !atomic.CompareAndSwapInt32(&used, 0, 1) {
				return nil, errors.New("listener reused")
			}
			return listener, nil
		},
		ChunkSize: 4,
	}

	session := NewSession(ctx, stream, cfg)
	errCh := make(chan error, 1)
	go func() {
		errCh <- session.Run()
	}()

	stream.SendRequest(
		pb.ExposePortRequest_builder{
			Init: pb.ExposePortInit_builder{
				LocalPort: proto.Uint32(23456),
			}.Build(),
		}.Build(),
	)

	remoteConn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial: %v", err)
	}
	defer remoteConn.Close()

	resp := stream.NextResponse(t)
	open := resp.GetOpen()
	if open == nil {
		t.Fatalf("expected open frame, got %v", resp.WhichFrame())
	}
	connectionID := open.GetConnectionId()
	if connectionID == 0 {
		t.Fatalf("connection id is zero")
	}
	expectedPort := uint32(listener.Addr().(*net.TCPAddr).Port)
	if open.GetRemotePort() != expectedPort {
		t.Fatalf("remote port: got %d want %d", open.GetRemotePort(), expectedPort)
	}
	if open.GetRemoteAddress() == "" {
		t.Fatalf("remote address is empty")
	}

	message := []byte("abcdef")
	if _, err := remoteConn.Write(message); err != nil {
		t.Fatalf("remoteConn.Write: %v", err)
	}

	var received []byte
	for len(received) < len(message) {
		resp := stream.NextResponse(t)
		data := resp.GetData()
		if data == nil {
			t.Fatalf("expected data frame, got %v", resp.WhichFrame())
		}
		if data.GetConnectionId() != connectionID {
			t.Fatalf(
				"unexpected connection id: got %d want %d",
				data.GetConnectionId(),
				connectionID,
			)
		}
		received = append(received, data.GetPayload()...)
	}
	if !bytes.Equal(received, message) {
		t.Fatalf("payload mismatch: got %q want %q", received, message)
	}

	payload := []byte("xyz")
	stream.SendRequest(
		pb.ExposePortRequest_builder{
			Data: pb.ExposePortConnectionData_builder{
				ConnectionId: proto.Uint64(connectionID),
				Payload:      append([]byte(nil), payload...),
			}.Build(),
		}.Build(),
	)

	buf := make([]byte, len(payload))
	if _, err := io.ReadFull(remoteConn, buf); err != nil {
		t.Fatalf("io.ReadFull: %v", err)
	}
	if !bytes.Equal(buf, payload) {
		t.Fatalf("read mismatch: got %q want %q", buf, payload)
	}

	stream.SendRequest(
		pb.ExposePortRequest_builder{
			Close: pb.ExposePortConnectionClose_builder{
				ConnectionId: proto.Uint64(connectionID),
			}.Build(),
		}.Build(),
	)

	resp = stream.NextResponse(t)
	closeFrame := resp.GetClose()
	if closeFrame == nil {
		t.Fatalf("expected close frame, got %v", resp.WhichFrame())
	}
	if closeFrame.GetConnectionId() != connectionID {
		t.Fatalf("close id mismatch: got %d want %d", closeFrame.GetConnectionId(), connectionID)
	}

	_ = remoteConn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	n, readErr := remoteConn.Read(make([]byte, 1))
	if n != 0 || readErr == nil {
		t.Fatalf("remote connection should be closed, got n=%d err=%v", n, readErr)
	}

	stream.CloseRequests()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("session.Run: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for session completion")
	}
}

func TestSessionRequireInit(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := testexposeport.NewFakeServerStream(ctx)
	session := NewSession(ctx, stream, Config{})

	errCh := make(chan error, 1)
	go func() {
		errCh <- session.Run()
	}()

	stream.SendRequest(
		pb.ExposePortRequest_builder{
			Data: pb.ExposePortConnectionData_builder{
				ConnectionId: proto.Uint64(1),
				Payload:      []byte("x"),
			}.Build(),
		}.Build(),
	)
	stream.CloseRequests()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatalf("expected error but got nil")
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for session")
	}

	resp := stream.NextResponse(t)
	if resp.GetError() == nil {
		t.Fatalf("expected error frame, got %v", resp.WhichFrame())
	}
}

func TestSessionListenFallback(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := testexposeport.NewFakeServerStream(ctx)

	var (
		attempts  int
		fallback  net.Listener
		listenErr = errors.New("address already in use")
		readyCh   = make(chan struct{})
	)

	cfg := Config{
		Listen: func(network, address string) (net.Listener, error) {
			attempts++
			if attempts == 1 {
				if address != net.JoinHostPort("127.0.0.1", "12345") {
					t.Fatalf("unexpected first listen address: %s", address)
				}
				return nil, listenErr
			}
			if address != net.JoinHostPort("127.0.0.1", "0") {
				t.Fatalf("unexpected fallback listen address: %s", address)
			}
			l, err := net.Listen("tcp", address)
			if err != nil {
				return nil, err
			}
			fallback = l
			close(readyCh)
			return l, nil
		},
		ChunkSize: 4,
	}

	session := NewSession(ctx, stream, cfg)
	errCh := make(chan error, 1)
	go func() {
		errCh <- session.Run()
	}()

	stream.SendRequest(
		pb.ExposePortRequest_builder{
			Init: pb.ExposePortInit_builder{
				LocalPort: proto.Uint32(12345),
			}.Build(),
		}.Build(),
	)

	select {
	case <-readyCh:
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for fallback listener")
	}

	defer func() {
		if fallback != nil {
			_ = fallback.Close()
		}
	}()

	fallbackAddr := fallback.Addr().String()
	fallbackPort := uint32(fallback.Addr().(*net.TCPAddr).Port)

	remoteConn, err := net.Dial("tcp", fallbackAddr)
	if err != nil {
		t.Fatalf("net.Dial fallback: %v", err)
	}
	defer remoteConn.Close()

	resp := stream.NextResponse(t)
	open := resp.GetOpen()
	if open == nil {
		t.Fatalf("expected open frame, got %v", resp.WhichFrame())
	}
	if open.GetRemotePort() != fallbackPort {
		t.Fatalf("unexpected remote port: got %d want %d", open.GetRemotePort(), fallbackPort)
	}

	if _, err := remoteConn.Write([]byte("ping")); err != nil {
		t.Fatalf("remoteConn.Write: %v", err)
	}

	req := stream.NextResponse(t)
	if req.GetData() == nil {
		t.Fatalf("expected data frame after fallback open, got %v", req.WhichFrame())
	}

	stream.SendRequest(
		pb.ExposePortRequest_builder{
			Close: pb.ExposePortConnectionClose_builder{
				ConnectionId: proto.Uint64(open.GetConnectionId()),
			}.Build(),
		}.Build(),
	)

	stream.CloseRequests()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, context.Canceled) {
			t.Fatalf("unexpected run error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for session end")
	}
	if attempts != 2 {
		t.Fatalf("expected two listen attempts, got %d", attempts)
	}
}

func TestSessionListenError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := testexposeport.NewFakeServerStream(ctx)

	cfg := Config{
		Listen: func(network, address string) (net.Listener, error) {
			return nil, errors.New("failed to bind")
		},
	}

	session := NewSession(ctx, stream, cfg)
	errCh := make(chan error, 1)
	go func() {
		errCh <- session.Run()
	}()

	stream.SendRequest(
		pb.ExposePortRequest_builder{
			Init: pb.ExposePortInit_builder{
				LocalPort: proto.Uint32(4321),
			}.Build(),
		}.Build(),
	)

	var runErr error
	select {
	case runErr = <-errCh:
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for session error")
	}
	if runErr == nil {
		t.Fatalf("expected non-nil error from session")
	}
	if !strings.Contains(runErr.Error(), "listen") {
		t.Fatalf("unexpected error message: %v", runErr)
	}

	resp := stream.NextResponse(t)
	errFrame := resp.GetError()
	if errFrame == nil {
		t.Fatalf("expected error response, got %v", resp.WhichFrame())
	}
	if errFrame.GetMessage() == "" {
		t.Fatalf("error message must not be empty")
	}
}
