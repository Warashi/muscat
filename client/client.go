package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"

	"connectrpc.com/connect"

	"github.com/Warashi/muscat/v2/pb"
	"github.com/Warashi/muscat/v2/pb/pbconnect"
	"github.com/Warashi/muscat/v2/stream"
)

var osHostname string

func init() {
	if h, err := os.Hostname(); err == nil {
		osHostname = h
	}
}

func New(socketPath string) *MuscatClient {
	client := new(http.Client)
	client.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	}
	return &MuscatClient{pb: pbconnect.NewMuscatServiceClient(client, "http://localhost")}
}

type MuscatClient struct {
	pb pbconnect.MuscatServiceClient
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

func (m *MuscatClient) Health(ctx context.Context) (int, error) {
	res, err := m.pb.Health(ctx, connect.NewRequest(new(pb.HealthRequest)))
	if err != nil {
		return 0, fmt.Errorf("m.pb.Health: %w", err)
	}
	return int(res.Msg.GetPid()), nil
}

func (m *MuscatClient) Open(ctx context.Context, uri string) error {
	uri = replaceLoopback(uri)
	if _, err := m.pb.Open(ctx, connect.NewRequest(&pb.OpenRequest{Uri: uri})); err != nil {
		return fmt.Errorf("m.pb.Open: %w", err)
	}
	return nil
}

func (m *MuscatClient) Copy(ctx context.Context, r io.Reader) (err error) {
	s := m.pb.Copy(ctx)
	defer func() {
		if _, closeErr := s.CloseAndReceive(); closeErr != nil && !errors.Is(closeErr, io.EOF) {
			err = errors.Join(err, fmt.Errorf("stream.CloseAndRecv: %w", err))
		}
	}()
	dst, src := bufio.NewWriter(stream.NewWriter(func(body []byte) *pb.CopyRequest { return &pb.CopyRequest{Body: body} }, s)), bufio.NewReader(r)
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := dst.Flush(); err != nil {
		return fmt.Errorf("dst.Flush: %w", err)
	}
	return nil
}

func (m *MuscatClient) Paste(ctx context.Context) (io.Reader, error) {
	s, err := m.pb.Paste(ctx, connect.NewRequest(new(pb.PasteRequest)))
	if err != nil {
		return nil, fmt.Errorf("m.pb.Paste: %w", err)
	}
	return bufio.NewReader(stream.NewReader(s)), nil
}

func (m *MuscatClient) GetInputMethod(ctx context.Context) (string, error) {
	res, err := m.pb.GetInputMethod(ctx, connect.NewRequest(new(pb.GetInputMethodRequest)))
	if err != nil {
		return "", fmt.Errorf("m.pb.GetInputMethod: %w", err)
	}
	return res.Msg.GetId(), nil
}

func (m *MuscatClient) SetInputMethod(ctx context.Context, id string) (before string, err error) {
	res, err := m.pb.SetInputMethod(ctx, connect.NewRequest(&pb.SetInputMethodRequest{Id: id}))
	if err != nil {
		return "", fmt.Errorf("m.pb.SetInputMethod: %w", err)
	}
	return res.Msg.GetBefore(), nil
}
