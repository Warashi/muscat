package server

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/atotto/clipboard"
	"github.com/skratchdot/open-golang/open"

	"github.com/Warashi/muscat/pb"
	"github.com/Warashi/muscat/stream"
)

type Muscat struct {
	pb.UnimplementedMuscatServer

	// mu clipboard用のmutex
	mu sync.Mutex
	// clipboard OSのクリップボードが使えないときにここに保持する
	clipboard string
}

func newPasteResponse(body []byte) *pb.PasteResponse {
	return &pb.PasteResponse{Body: body}
}

func (m *Muscat) Open(ctx context.Context, request *pb.OpenRequest) (*pb.OpenResponse, error) {
	if err := open.Run(request.Uri); err != nil {
		return nil, fmt.Errorf("open.Run: %w", err)
	}
	return new(pb.OpenResponse), nil
}

func (m *Muscat) Copy(s pb.Muscat_CopyServer) error {
	buf, src := new(bytes.Buffer), bufio.NewReader(stream.NewReader[*pb.CopyRequest](s))
	if _, err := io.Copy(buf, src); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if clipboard.Unsupported {
		// OSのクリップボードが使えないのでサーバーローカルに保持する
		m.mu.Lock()
		defer m.mu.Unlock()
		m.clipboard = buf.String()
		return nil
	}
	if err := clipboard.WriteAll(buf.String()); err != nil {
		return fmt.Errorf("clipboard.WriteAll: %w", err)
	}
	return nil
}

func (m *Muscat) Paste(_ *pb.PasteRequest, s pb.Muscat_PasteServer) error {
	dst := stream.NewWriter[*pb.PasteResponse](newPasteResponse, s)
	var body string
	if clipboard.Unsupported {
		m.mu.Lock()
		body = m.clipboard
		m.mu.Unlock()
	} else {
		var err error
		body, err = clipboard.ReadAll()
		if err != nil {
			return fmt.Errorf("clipboard.ReadAll: %w", err)
		}
	}
	if _, err := io.Copy(dst, strings.NewReader(body)); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	return nil
}

func (m *Muscat) Health(context.Context, *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Pid: int64(os.Getpid())}, nil
}
