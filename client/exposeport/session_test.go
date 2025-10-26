package exposeport

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"sync"
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
	defer listener.Close()
	localPort := uint32(listener.Addr().(*net.TCPAddr).Port)

	init := pb.ExposePortInit_builder{
		LocalPort: proto.Uint32(localPort),
	}.Build()
	session := NewSession(ctx, stream, init, Config{ChunkSize: 4})

	errCh := make(chan error, 1)
	go func() {
		errCh <- session.Run()
	}()

	req := stream.NextRequest(t)
	if req.GetInit() == nil {
		t.Fatalf("expected init frame, got %v", req.WhichFrame())
	}

	connCh := make(chan net.Conn, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		connCh <- conn
	}()

	connectionID := uint64(1)
	stream.SendResponse(
		pb.ExposePortResponse_builder{
			Open: pb.ExposePortConnectionOpen_builder{
				ConnectionId:  proto.Uint64(connectionID),
				RemotePort:    proto.Uint32(5000),
				RemoteAddress: proto.String("remote:1234"),
			}.Build(),
		}.Build(),
	)

	var localConn net.Conn
	select {
	case localConn = <-connCh:
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for local accept")
	}
	defer localConn.Close()

	message := []byte("abcdef")
	if _, err := localConn.Write(message); err != nil {
		t.Fatalf("localConn.Write: %v", err)
	}

	var transmitted []byte
	for len(transmitted) < len(message) {
		req := stream.NextRequest(t)
		data := req.GetData()
		if data == nil {
			t.Fatalf("expected data frame, got %v", req.WhichFrame())
		}
		if data.GetConnectionId() != connectionID {
			t.Fatalf(
				"unexpected connection id: got %d want %d",
				data.GetConnectionId(),
				connectionID,
			)
		}
		transmitted = append(transmitted, data.GetPayload()...)
	}
	if !bytes.Equal(transmitted, message) {
		t.Fatalf("payload mismatch: got %q want %q", transmitted, message)
	}

	remotePayload := []byte("xyz")
	stream.SendResponse(
		pb.ExposePortResponse_builder{
			Data: pb.ExposePortConnectionData_builder{
				ConnectionId: proto.Uint64(connectionID),
				Payload:      append([]byte(nil), remotePayload...),
			}.Build(),
		}.Build(),
	)

	buf := make([]byte, len(remotePayload))
	if _, err := io.ReadFull(localConn, buf); err != nil {
		t.Fatalf("io.ReadFull: %v", err)
	}
	if !bytes.Equal(buf, remotePayload) {
		t.Fatalf("read mismatch: got %q want %q", buf, remotePayload)
	}

	stream.SendResponse(
		pb.ExposePortResponse_builder{
			Close: pb.ExposePortConnectionClose_builder{
				ConnectionId: proto.Uint64(connectionID),
			}.Build(),
		}.Build(),
	)

	_ = localConn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	if n, err := localConn.Read(make([]byte, 1)); err == nil || n != 0 {
		t.Fatalf("local connection should be closed, got n=%d err=%v", n, err)
	}

	stream.CloseResponses()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("session.Run: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for session")
	}
}

func TestSessionHandleDialError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := newTestStream(ctx)

	init := pb.ExposePortInit_builder{
		LocalPort: proto.Uint32(12345),
	}.Build()
	session := NewSession(
		ctx,
		stream,
		init,
		Config{
			Dial: func(context.Context, string, string) (net.Conn, error) {
				return nil, errors.New("dial failure")
			},
		},
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- session.Run()
	}()

	req := stream.NextRequest(t)
	if req.GetInit() == nil {
		t.Fatalf("expected init frame, got %v", req.WhichFrame())
	}

	stream.SendResponse(
		pb.ExposePortResponse_builder{
			Open: pb.ExposePortConnectionOpen_builder{
				ConnectionId: proto.Uint64(42),
			}.Build(),
		}.Build(),
	)

	closeReq := stream.NextRequest(t)
	if closeReq.GetClose() == nil {
		t.Fatalf("expected close frame, got %v", closeReq.WhichFrame())
	}
	if closeReq.GetClose().GetConnectionId() != 42 {
		t.Fatalf(
			"connection id mismatch: got %d want %d",
			closeReq.GetClose().GetConnectionId(),
			42,
		)
	}
	if closeReq.GetClose().GetError() == "" {
		t.Fatalf("expected error message")
	}

	stream.CloseResponses()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, io.EOF) {
			t.Fatalf("session Run error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for session")
	}
}

type testStream struct {
	ctx        context.Context
	requestCh  chan *pb.ExposePortRequest
	responseCh chan *pb.ExposePortResponse
}

func newTestStream(ctx context.Context) *testStream {
	return &testStream{
		ctx:        ctx,
		requestCh:  make(chan *pb.ExposePortRequest, 16),
		responseCh: make(chan *pb.ExposePortResponse, 16),
	}
}

func (t *testStream) Send(req *pb.ExposePortRequest) error {
	select {
	case <-t.ctx.Done():
		return t.ctx.Err()
	case t.requestCh <- req:
		return nil
	}
}

func (t *testStream) Receive() (*pb.ExposePortResponse, error) {
	select {
	case <-t.ctx.Done():
		return nil, t.ctx.Err()
	case resp, ok := <-t.responseCh:
		if !ok {
			return nil, io.EOF
		}
		return resp, nil
	}
}

func (t *testStream) CloseRequest() error {
	return nil
}

func (t *testStream) CloseResponse() error {
	return nil
}

func (t *testStream) NextRequest(tb testing.TB) *pb.ExposePortRequest {
	tb.Helper()
	select {
	case req := <-t.requestCh:
		return req
	case <-time.After(time.Second):
		tb.Fatalf("timeout waiting for request")
		return nil
	}
}

func (t *testStream) SendResponse(resp *pb.ExposePortResponse) {
	t.responseCh <- resp
}

func (t *testStream) CloseResponses() {
	close(t.responseCh)
}

func TestSessionEventHandlers(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := newTestStream(ctx)
	init := pb.ExposePortInit_builder{
		LocalPort: proto.Uint32(12345),
	}.Build()

	type event struct {
		kind string
		data any
	}
	eventCh := make(chan event, 3)

	var dialMu sync.Mutex
	var peers []net.Conn
	session := NewSession(
		ctx,
		stream,
		init,
		Config{
			Dial: func(context.Context, string, string) (net.Conn, error) {
				clientConn, serverConn := net.Pipe()
				dialMu.Lock()
				peers = append(peers, serverConn)
				dialMu.Unlock()
				return clientConn, nil
			},
			OnOpen: func(open *pb.ExposePortConnectionOpen) {
				eventCh <- event{kind: "open", data: open}
			},
			OnClose: func(close *pb.ExposePortConnectionClose) {
				eventCh <- event{kind: "close", data: close}
			},
			OnError: func(err error) {
				eventCh <- event{kind: "error", data: err}
			},
		},
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- session.Run()
	}()

	req := stream.NextRequest(t)
	if req.GetInit() == nil {
		t.Fatalf("expected init request, got %v", req.WhichFrame())
	}

	openFrame := pb.ExposePortResponse_builder{
		Open: pb.ExposePortConnectionOpen_builder{
			ConnectionId:  proto.Uint64(1),
			RemotePort:    proto.Uint32(8080),
			RemoteAddress: proto.String("remote:1234"),
		}.Build(),
	}.Build()
	stream.SendResponse(openFrame)

	select {
	case ev := <-eventCh:
		open := ev.data.(*pb.ExposePortConnectionOpen)
		if ev.kind != "open" {
			t.Fatalf("expected open event, got %s", ev.kind)
		}
		if open.GetConnectionId() != 1 || open.GetRemotePort() != 8080 {
			t.Fatalf("unexpected open payload: %#v", open)
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for open event")
	}

	closeFrame := pb.ExposePortResponse_builder{
		Close: pb.ExposePortConnectionClose_builder{
			ConnectionId: proto.Uint64(1),
			Error:        proto.String("done"),
			Reset:        proto.Bool(false),
		}.Build(),
	}.Build()
	stream.SendResponse(closeFrame)

	select {
	case ev := <-eventCh:
		closeMsg := ev.data.(*pb.ExposePortConnectionClose)
		if ev.kind != "close" {
			t.Fatalf("expected close event, got %s", ev.kind)
		}
		if closeMsg.GetError() != "done" {
			t.Fatalf("unexpected close payload: %#v", closeMsg)
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for close event")
	}

	errorFrame := pb.ExposePortResponse_builder{
		Error: pb.ExposePortError_builder{
			Message: proto.String("listener failed"),
		}.Build(),
	}.Build()
	stream.SendResponse(errorFrame)
	stream.CloseResponses()

	select {
	case ev := <-eventCh:
		err, ok := ev.data.(error)
		if ev.kind != "error" || !ok {
			t.Fatalf("expected error event, got %#v", ev)
		}
		if err.Error() != "listener failed" {
			t.Fatalf("unexpected error event: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for error event")
	}

	select {
	case err := <-errCh:
		if err == nil || err.Error() != "listener failed" {
			t.Fatalf("unexpected session result: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for session termination")
	}

	dialMu.Lock()
	for _, conn := range peers {
		_ = conn.Close()
	}
	dialMu.Unlock()
}
