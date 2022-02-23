package client

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Warashi/muscat/pb"
)

func socketDialer(ctx context.Context, addr string) (net.Conn, error) {
	var d net.Dialer
	return d.DialContext(ctx, "unix", addr)
}

func New(socketPath string) (*Muscat, error) {
	conn, err := grpc.Dial(socketPath, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(socketDialer))
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %w", err)
	}
	return &Muscat{pb: pb.NewMuscatClient(conn)}, nil
}

type Muscat struct {
	pb pb.MuscatClient
}

type streamWriter struct {
	stream pb.Muscat_CopyClient
}

func (w *streamWriter) Write(p []byte) (n int, err error) {
	if err := w.stream.Send(&pb.CopyRequest{Body: p}); err != nil {
		return 0, fmt.Errorf("w.stream.Send: %w", err)
	}
	return len(p), nil
}

type streamReader struct {
	stream pb.Muscat_PasteClient
	buffer bytes.Buffer
}

func (r *streamReader) Read(p []byte) (n int, err error) {
	for r.buffer.Len() < len(p) {
		b, err := r.stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return 0, fmt.Errorf("r.stream.Recv: %w", err)
		}
		if _, err := r.buffer.Write(b.GetBody()); err != nil {
			return 0, fmt.Errorf("r.buffer.Write: %w", err)
		}
	}
	n, err = r.buffer.Read(p)
	if err != nil {
		return n, fmt.Errorf("r.buffer.Read: %w", err)
	}
	return n, nil
}

func (m *Muscat) Open(ctx context.Context, uri string) error {
	if _, err := m.pb.Open(ctx, &pb.OpenRequest{Uri: uri}); err != nil {
		return fmt.Errorf("m.pb.Open: %w", err)
	}
	return nil
}

func (m *Muscat) Copy(ctx context.Context, r io.Reader) error {
	stream, err := m.pb.Copy(ctx)
	if err != nil {
		return fmt.Errorf("m.pb.Copy: %w", err)
	}

	dst, src := bufio.NewWriter(&streamWriter{stream: stream}), bufio.NewReader(r)
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := dst.Flush(); err != nil {
		return fmt.Errorf("dst.Flush: %w", err)
	}
	return nil
}

func (m *Muscat) Paste(ctx context.Context) (io.Reader, error) {
	stream, err := m.pb.Paste(ctx, new(pb.PasteRequest))
	if err != nil {
		return nil, fmt.Errorf("m.pb.Paste: %w", err)
	}
	return bufio.NewReader(&streamReader{stream: stream}), nil
}
