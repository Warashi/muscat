package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Warashi/muscat/pb"
	"github.com/Warashi/muscat/stream"
)

var osHostname string

func init() {
	if h, err := os.Hostname(); err == nil {
		osHostname = h
	}
}

func socketDialer(ctx context.Context, addr string) (net.Conn, error) {
	var d net.Dialer
	return d.DialContext(ctx, "unix", addr)
}

func newCopyRequest(body []byte) *pb.CopyRequest {
	return &pb.CopyRequest{Body: body}
}

func New(socketPath string) (*Muscat, error) {
	conn, err := grpc.Dial(socketPath, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(socketDialer))
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %w", err)
	}
	return &Muscat{pb: pb.NewMuscatClient(conn), conn: conn}, nil
}

type Muscat struct {
	pb   pb.MuscatClient
	conn io.Closer
}

func replaceLoopback(uri string) (s string) {
	if osHostname == "" {
		return uri
	}
	u, err := url.Parse(uri)
	if err != nil || u.Host == "" {
		u, err = url.Parse("http://" + uri)
	}
	if err != nil {
		return uri
	}
	host, port := u.Host, ""
	if h, p, err := net.SplitHostPort(u.Host); err == nil {
		host, port = h, p
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		if port == "" {
			u.Host = osHostname
		} else {
			u.Host = net.JoinHostPort(osHostname, port)
		}
		return u.String()
	}
	if ip, err := net.LookupIP(host); err == nil && len(ip) > 0 && ip[0].IsLoopback() {
		if port == "" {
			u.Host = osHostname
		} else {
			u.Host = net.JoinHostPort(osHostname, port)
		}
		return u.String()
	}
	return uri
}

func (m *Muscat) Health(ctx context.Context) (int, error) {
	res, err := m.pb.Health(ctx, new(pb.HealthRequest))
	if err != nil {
		return 0, fmt.Errorf("m.pb.Health: %w", err)
	}
	return int(res.GetPid()), nil
}

func (m *Muscat) Close() error {
	if err := m.conn.Close(); err != nil {
		return fmt.Errorf("m.conn.Close: %w", err)
	}
	return nil
}

func (m *Muscat) Open(ctx context.Context, uri string) error {
	uri = replaceLoopback(uri)
	if _, err := m.pb.Open(ctx, &pb.OpenRequest{Uri: uri}); err != nil {
		return fmt.Errorf("m.pb.Open: %w", err)
	}
	return nil
}

func (m *Muscat) Copy(ctx context.Context, r io.Reader) error {
	s, err := m.pb.Copy(ctx)
	if err != nil {
		return fmt.Errorf("m.pb.Copy: %w", err)
	}

	dst, src := bufio.NewWriter(stream.NewWriter[*pb.CopyRequest](newCopyRequest, s)), bufio.NewReader(r)
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := dst.Flush(); err != nil {
		return fmt.Errorf("dst.Flush: %w", err)
	}
	if _, err := s.CloseAndRecv(); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("stream.CloseAndRecv: %w", err)
	}
	return nil
}

func (m *Muscat) Paste(ctx context.Context) (io.Reader, error) {
	s, err := m.pb.Paste(ctx, new(pb.PasteRequest))
	if err != nil {
		return nil, fmt.Errorf("m.pb.Paste: %w", err)
	}
	return bufio.NewReader(stream.NewReader[*pb.PasteResponse](s)), nil
}
