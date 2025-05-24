package client

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/proto"

	"github.com/Warashi/muscat/v2/consts"
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

func New(network, addr string) *MuscatClient {
	client := new(http.Client)
	client.Transport = &http2.Transport{
		AllowHTTP: true,
		DialTLSContext: func(ctx context.Context, _, _ string, cfg *tls.Config) (net.Conn, error) {
			return (new(net.Dialer)).DialContext(ctx, network, addr)
		},
	}
	return &MuscatClient{
		pb: pbconnect.NewMuscatServiceClient(client, "http://localhost", connect.WithSendGzip()),
	}
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
	if _, err := m.pb.Open(ctx, connect.NewRequest(pb.OpenRequest_builder{Uri: proto.String(uri)}.Build())); err != nil {
		return fmt.Errorf("m.pb.Open: %w", err)
	}
	return nil
}

func (m *MuscatClient) Copy(ctx context.Context, r io.Reader) (err error) {
	s := m.pb.Copy(ctx)
	defer func() {
		if _, closeErr := s.CloseAndReceive(); closeErr != nil && !errors.Is(closeErr, io.EOF) {
			err = errors.Join(err, fmt.Errorf("stream.CloseAndRecv: %w", closeErr))
		}
	}()
	dst, src := bufio.NewWriter(
		stream.NewWriter(
			func(body []byte) *pb.CopyRequest { return pb.CopyRequest_builder{Body: body}.Build() },
			s,
		),
	), bufio.NewReader(
		r,
	)
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
	res, err := m.pb.SetInputMethod(
		ctx,
		connect.NewRequest(pb.SetInputMethodRequest_builder{Id: proto.String(id)}.Build()),
	)
	if err != nil {
		return "", fmt.Errorf("m.pb.SetInputMethod: %w", err)
	}
	return res.Msg.GetBefore(), nil
}

func (m *MuscatClient) PortForward(ctx context.Context, port string) error {
	l, err := net.Listen("tcp", net.JoinHostPort("localhost", port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("l.Accept: %v", err)
			continue
		}
		go func() {
			defer conn.Close()
			if err := m.portForward(ctx, port, conn); err != nil {
				log.Printf("m.portForward: %v", err)
			}
		}()
	}
}

func (m *MuscatClient) portForward(ctx context.Context, port string, conn net.Conn) (err error) {
	s := m.pb.PortForward(ctx)

	s.RequestHeader().Set(consts.HeaderNameMuscatForwardedPort, port)
	s.Send(nil)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer s.CloseRequest()

		dst := stream.NewWriter(
			func(body []byte) *pb.PortForwardRequest { return pb.PortForwardRequest_builder{Body: body}.Build() },
			s,
		)
		if _, err := io.Copy(dst, conn); err != nil {
			log.Printf("io.Copy: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		defer s.CloseResponse()

		src := stream.NewBidiReader(s)
		if _, err := io.Copy(conn, src); err != nil {
			log.Printf("io.Copy: %v", err)
		}
	}()

	wg.Wait()

	return nil
}
