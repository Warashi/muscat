package server

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"connectrpc.com/connect"
	"github.com/Warashi/go-swim"
	"github.com/skratchdot/open-golang/open"
	"google.golang.org/protobuf/proto"

	"github.com/Warashi/muscat/v2/consts"
	"github.com/Warashi/muscat/v2/pb"
	"github.com/Warashi/muscat/v2/pb/pbconnect"
	"github.com/Warashi/muscat/v2/server/internal/clipboard"
	"github.com/Warashi/muscat/v2/stream"
)

var _ pbconnect.MuscatServiceHandler = (*MuscatServer)(nil)

type MuscatServer struct {
	// mu clipboard用のmutex
	mu sync.Mutex
	// clipboard OSのクリップボードが使えないときにここに保持する
	clipboard []byte
}

// Health implements pbconnect.MuscatServiceHandler.
func (m *MuscatServer) Health(context.Context, *connect.Request[pb.HealthRequest]) (*connect.Response[pb.HealthResponse], error) {
	return connect.NewResponse(pb.HealthResponse_builder{Pid: proto.Int64(int64(os.Getpid()))}.Build()), nil
}

// Open implements pbconnect.MuscatServiceHandler.
func (m *MuscatServer) Open(_ context.Context, req *connect.Request[pb.OpenRequest]) (*connect.Response[pb.OpenResponse], error) {
	if err := open.Run(req.Msg.GetUri()); err != nil {
		return nil, fmt.Errorf("open.Run: %w", err)
	}
	return connect.NewResponse(new(pb.OpenResponse)), nil
}

// Copy implements pbconnect.MuscatServiceHandler.
func (m *MuscatServer) Copy(ctx context.Context, s *connect.ClientStream[pb.CopyRequest]) (*connect.Response[pb.CopyResponse], error) {
	buf, src := new(bytes.Buffer), bufio.NewReader(stream.NewReader(s))
	if _, err := io.Copy(buf, src); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("io.Copy: %w", err))
	}

	if clipboard.Unsupported() {
		// OSのクリップボードが使えないのでサーバーローカルに保持する
		m.mu.Lock()
		defer m.mu.Unlock()
		m.clipboard = buf.Bytes()
		return connect.NewResponse(new(pb.CopyResponse)), nil
	}
	if err := clipboard.Write(buf.Bytes()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("clipboard.WriteAll: %w", err))
	}
	return connect.NewResponse(new(pb.CopyResponse)), nil
}

// Paste implements pbconnect.MuscatServiceHandler.
func (m *MuscatServer) Paste(ctx context.Context, req *connect.Request[pb.PasteRequest], s *connect.ServerStream[pb.PasteResponse]) error {
	dst := stream.NewWriter(func(body []byte) *pb.PasteResponse { return pb.PasteResponse_builder{Body: body}.Build() }, s)
	var body []byte
	if clipboard.Unsupported() {
		m.mu.Lock()
		body = m.clipboard
		m.mu.Unlock()
	} else {
		var err error
		body, err = clipboard.Read()
		if err != nil {
			return fmt.Errorf("clipboard.ReadAll: %w", err)
		}
	}
	if _, err := io.Copy(dst, bytes.NewReader(body)); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	return nil
}

// GetInputMethod implements pbconnect.MuscatServiceHandler.
func (*MuscatServer) GetInputMethod(ctx context.Context, req *connect.Request[pb.GetInputMethodRequest]) (*connect.Response[pb.GetInputMethodResponse], error) {
	id, err := swim.Get()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("swim.Get: %w", err))
	}
	return connect.NewResponse(pb.GetInputMethodResponse_builder{Id: proto.String(id)}.Build()), nil
}

// SetInputMethod implements pbconnect.MuscatServiceHandler.
func (*MuscatServer) SetInputMethod(ctx context.Context, req *connect.Request[pb.SetInputMethodRequest]) (*connect.Response[pb.SetInputMethodResponse], error) {
	before, err := swim.Get()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("swim.Get: %w", err))
	}
	if err := swim.Set(req.Msg.GetId()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("swim.Set: %w", err))
	}
	return connect.NewResponse(pb.SetInputMethodResponse_builder{Before: proto.String(before)}.Build()), nil
}

// PortForward implements pbconnect.MuscatServiceHandler.
func (*MuscatServer) PortForward(ctx context.Context, s *connect.BidiStream[pb.PortForwardRequest, pb.PortForwardResponse]) error {
	port := s.RequestHeader().Get(consts.HeaderNameMuscatForwardedPort)
	if port == "" {
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("port is empty"))
	}
	conn, err := net.Dial("tcp", net.JoinHostPort("localhost", port))
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("net.Dial: %w", err))
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		dst := stream.NewWriter(func(body []byte) *pb.PortForwardResponse { return pb.PortForwardResponse_builder{Body: body}.Build() }, s)
		if _, err := io.Copy(dst, conn); err != nil {
			log.Printf("io.Copy: %v\n", err)
		}
	}()
	go func() {
		defer wg.Done()
		src := stream.NewBidiReader(s)
		if _, err := io.Copy(conn, src); err != nil {
			log.Printf("io.Copy: %v\n", err)
		}
	}()

	wg.Wait()

	return nil
}
