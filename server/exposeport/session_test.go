package exposeport

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/Warashi/muscat/v2/pb"
)

func TestSessionDataFlow(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := newTestStream(ctx)

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

	stream := newTestStream(ctx)
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

	select {
	case resp := <-stream.sendCh:
		if resp.GetError() == nil {
			t.Fatalf("expected error frame, got %v", resp.WhichFrame())
		}
	default:
		t.Fatalf("expected error frame to be sent")
	}
}

type testStream struct {
	ctx    context.Context
	recvCh chan *pb.ExposePortRequest
	sendCh chan *pb.ExposePortResponse
}

func newTestStream(ctx context.Context) *testStream {
	return &testStream{
		ctx:    ctx,
		recvCh: make(chan *pb.ExposePortRequest, 16),
		sendCh: make(chan *pb.ExposePortResponse, 16),
	}
}

func (t *testStream) Receive() (*pb.ExposePortRequest, error) {
	select {
	case <-t.ctx.Done():
		return nil, t.ctx.Err()
	case req, ok := <-t.recvCh:
		if !ok {
			return nil, io.EOF
		}
		return req, nil
	}
}

func (t *testStream) Send(resp *pb.ExposePortResponse) error {
	select {
	case <-t.ctx.Done():
		return t.ctx.Err()
	case t.sendCh <- resp:
		return nil
	}
}

func (t *testStream) SendRequest(req *pb.ExposePortRequest) {
	t.recvCh <- req
}

func (t *testStream) CloseRequests() {
	close(t.recvCh)
}

func (t *testStream) NextResponse(tb testing.TB) *pb.ExposePortResponse {
	tb.Helper()
	select {
	case resp := <-t.sendCh:
		return resp
	case <-time.After(time.Second):
		tb.Fatalf("timeout waiting for response")
		return nil
	}
}
