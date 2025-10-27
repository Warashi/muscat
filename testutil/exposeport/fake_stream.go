package exposeport

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/Warashi/muscat/v2/pb"
)

const (
	defaultBufferSize = 16
	defaultTimeout    = time.Second
)

// FakeClientStream provides a controllable implementation of the client-side
// ExposePort stream interface for tests.
type FakeClientStream struct {
	ctx context.Context

	requestCh  chan *pb.ExposePortRequest
	responseCh chan *pb.ExposePortResponse

	closeReqOnce sync.Once
	closeResOnce sync.Once
}

// NewFakeClientStream constructs a FakeClientStream bound to ctx.
func NewFakeClientStream(ctx context.Context) *FakeClientStream {
	return &FakeClientStream{
		ctx:        ctx,
		requestCh:  make(chan *pb.ExposePortRequest, defaultBufferSize),
		responseCh: make(chan *pb.ExposePortResponse, defaultBufferSize),
	}
}

// Send enqueues a client request frame.
func (f *FakeClientStream) Send(req *pb.ExposePortRequest) error {
	select {
	case <-f.ctx.Done():
		return f.ctx.Err()
	case f.requestCh <- req:
		return nil
	}
}

// Receive dequeues the next server response frame.
func (f *FakeClientStream) Receive() (*pb.ExposePortResponse, error) {
	select {
	case <-f.ctx.Done():
		return nil, f.ctx.Err()
	case resp, ok := <-f.responseCh:
		if !ok {
			return nil, io.EOF
		}
		return resp, nil
	}
}

// CloseRequest closes the client request channel.
func (f *FakeClientStream) CloseRequest() error {
	f.closeReqOnce.Do(func() {
		close(f.requestCh)
	})
	return nil
}

// CloseResponse closes the server response channel.
func (f *FakeClientStream) CloseResponse() error {
	f.closeResOnce.Do(func() {
		close(f.responseCh)
	})
	return nil
}

// NextRequest returns the next outbound request, failing the test on timeout.
func (f *FakeClientStream) NextRequest(tb testing.TB) *pb.ExposePortRequest {
	tb.Helper()
	select {
	case req := <-f.requestCh:
		return req
	case <-time.After(defaultTimeout):
		tb.Fatalf("timeout waiting for ExposePort request")
		return nil
	}
}

// SendResponse simulates a server frame arriving at the client.
func (f *FakeClientStream) SendResponse(resp *pb.ExposePortResponse) {
	f.responseCh <- resp
}

// CloseResponses closes the response side, causing subsequent Receive calls to
// return io.EOF.
func (f *FakeClientStream) CloseResponses() {
	f.closeResOnce.Do(func() {
		close(f.responseCh)
	})
}

// FakeServerStream provides a controllable implementation of the server-side
// ExposePort stream interface for tests.
type FakeServerStream struct {
	ctx context.Context

	requestCh chan *pb.ExposePortRequest
	sendCh    chan *pb.ExposePortResponse

	closeReqOnce  sync.Once
	closeSendOnce sync.Once
}

// NewFakeServerStream constructs a FakeServerStream bound to ctx.
func NewFakeServerStream(ctx context.Context) *FakeServerStream {
	return &FakeServerStream{
		ctx:       ctx,
		requestCh: make(chan *pb.ExposePortRequest, defaultBufferSize),
		sendCh:    make(chan *pb.ExposePortResponse, defaultBufferSize),
	}
}

// Receive returns the next client request.
func (f *FakeServerStream) Receive() (*pb.ExposePortRequest, error) {
	select {
	case <-f.ctx.Done():
		return nil, f.ctx.Err()
	case req, ok := <-f.requestCh:
		if !ok {
			return nil, io.EOF
		}
		return req, nil
	}
}

// Send enqueues a server response.
func (f *FakeServerStream) Send(resp *pb.ExposePortResponse) error {
	select {
	case <-f.ctx.Done():
		return f.ctx.Err()
	case f.sendCh <- resp:
		return nil
	}
}

// SendRequest injects a client frame into the stream.
func (f *FakeServerStream) SendRequest(req *pb.ExposePortRequest) {
	f.requestCh <- req
}

// CloseRequests closes the inbound request frontier.
func (f *FakeServerStream) CloseRequests() {
	f.closeReqOnce.Do(func() {
		close(f.requestCh)
	})
}

// NextResponse retrieves the next server response or fails on timeout.
func (f *FakeServerStream) NextResponse(tb testing.TB) *pb.ExposePortResponse {
	tb.Helper()
	select {
	case resp := <-f.sendCh:
		return resp
	case <-time.After(defaultTimeout):
		tb.Fatalf("timeout waiting for ExposePort response")
		return nil
	}
}

// CloseResponses ensures the send channel is closed.
func (f *FakeServerStream) CloseResponses() error {
	f.closeSendOnce.Do(func() {
		close(f.sendCh)
	})
	return nil
}
