package exposeport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"

	"google.golang.org/protobuf/proto"

	"github.com/Warashi/muscat/v2/pb"
)

// Stream abstracts the client-side ExposePort bidi stream.
type Stream interface {
	Send(*pb.ExposePortRequest) error
	Receive() (*pb.ExposePortResponse, error)
	CloseRequest() error
	CloseResponse() error
}

// DialFunc dials the local endpoint for a connection.
type DialFunc func(ctx context.Context, network, address string) (net.Conn, error)

// Config controls client session behaviour.
type Config struct {
	Dial         DialFunc
	LocalAddress string
	ChunkSize    int
}

// Session handles the client side of ExposePort.
type Session struct {
	ctx    context.Context
	stream Stream
	init   *pb.ExposePortInit
	cfg    sessionConfig

	sendMu   sync.Mutex
	mu       sync.RWMutex
	conns    map[uint64]*clientConn
	wg       sync.WaitGroup
	failOnce sync.Once
}

type sessionConfig struct {
	dial         DialFunc
	localAddress string
	chunkSize    int
}

const (
	defaultLocalAddress = "127.0.0.1"
	defaultChunkSize    = 32 * 1024
)

// NewSession constructs a Session.
func NewSession(ctx context.Context, stream Stream, init *pb.ExposePortInit, cfg Config) *Session {
	sessionCfg := sessionConfig{
		dial:         (&net.Dialer{}).DialContext,
		localAddress: cfg.LocalAddress,
		chunkSize:    cfg.ChunkSize,
	}
	if sessionCfg.localAddress == "" {
		sessionCfg.localAddress = defaultLocalAddress
	}
	if cfg.Dial != nil {
		sessionCfg.dial = cfg.Dial
	}
	if sessionCfg.chunkSize <= 0 {
		sessionCfg.chunkSize = defaultChunkSize
	}
	var initCopy *pb.ExposePortInit
	if init != nil {
		initCopy = proto.Clone(init).(*pb.ExposePortInit)
	} else {
		initCopy = new(pb.ExposePortInit)
	}
	return &Session{
		ctx:    ctx,
		stream: stream,
		init:   initCopy,
		cfg:    sessionCfg,
		conns:  make(map[uint64]*clientConn),
	}
}

// Run executes the session lifecycle.
func (s *Session) Run() error {
	if err := s.sendInit(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.receiveLoop(ctx)
	}()

	err := <-errCh
	cancel()

	s.closeAllConnections("session ended", true, false)
	s.wg.Wait()

	_ = s.stream.CloseRequest()
	_ = s.stream.CloseResponse()

	return err
}

func (s *Session) sendInit() error {
	if s.init == nil {
		return errors.New("initial configuration is nil")
	}
	initReq := pb.ExposePortRequest_builder{
		Init: proto.Clone(s.init).(*pb.ExposePortInit),
	}.Build()
	return s.send(initReq)
}

func (s *Session) receiveLoop(ctx context.Context) error {
	for {
		resp, err := s.stream.Receive()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("stream.Receive: %w", err)
		}

		switch {
		case resp.GetOpen() != nil:
			if err := s.handleOpen(ctx, resp.GetOpen()); err != nil {
				return err
			}
		case resp.GetData() != nil:
			s.handleData(resp.GetData())
		case resp.GetClose() != nil:
			s.handleClose(resp.GetClose())
		case resp.GetError() != nil:
			return errors.New(resp.GetError().GetMessage())
		default:
			// ignore
		}
	}
}

func (s *Session) handleOpen(ctx context.Context, open *pb.ExposePortConnectionOpen) error {
	addr := net.JoinHostPort(s.cfg.localAddress, strconv.Itoa(int(s.init.GetLocalPort())))
	conn, err := s.cfg.dial(ctx, "tcp", addr)
	if err != nil {
		return s.sendCloseWithError(open.GetConnectionId(), err.Error(), true)
	}

	clientConn := newClientConn(open.GetConnectionId(), conn)
	if err := s.storeConnection(clientConn); err != nil {
		conn.Close()
		return s.sendCloseWithError(open.GetConnectionId(), err.Error(), true)
	}

	s.wg.Add(2)
	go func() {
		defer s.wg.Done()
		s.runLocalToRemote(ctx, clientConn)
	}()
	go func() {
		defer s.wg.Done()
		s.runRemoteToLocal(ctx, clientConn)
	}()

	return nil
}

func (s *Session) runLocalToRemote(ctx context.Context, conn *clientConn) {
	buf := make([]byte, s.cfg.chunkSize)
	for {
		select {
		case <-ctx.Done():
			s.finishConnection(conn.id, "context canceled", true, true)
			return
		case <-conn.Done():
			return
		default:
		}

		n, err := conn.conn.Read(buf)
		if n > 0 {
			payload := append([]byte(nil), buf[:n]...)
			if err := s.sendData(conn.id, payload); err != nil {
				s.finishConnection(conn.id, fmt.Sprintf("send data: %v", err), true, true)
				return
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				s.finishConnection(conn.id, "", false, true)
			} else {
				s.finishConnection(conn.id, err.Error(), true, true)
			}
			return
		}
	}
}

func (s *Session) runRemoteToLocal(ctx context.Context, conn *clientConn) {
	outbound := conn.Outbound()
	for {
		select {
		case <-ctx.Done():
			s.finishConnection(conn.id, "context canceled", true, true)
			return
		case <-conn.Done():
			return
		case payload := <-outbound:
			if len(payload) == 0 {
				continue
			}
			if _, err := conn.conn.Write(payload); err != nil {
				s.finishConnection(conn.id, err.Error(), true, true)
				return
			}
		}
	}
}

func (s *Session) handleData(data *pb.ExposePortConnectionData) {
	conn, ok := s.getConnection(data.GetConnectionId())
	if !ok {
		return
	}
	payload := append([]byte(nil), data.GetPayload()...)
	conn.Enqueue(payload)
}

func (s *Session) handleClose(close *pb.ExposePortConnectionClose) {
	s.finishConnection(close.GetConnectionId(), close.GetError(), close.GetReset(), false)
}

func (s *Session) storeConnection(conn *clientConn) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.conns[conn.id]; exists {
		return fmt.Errorf("connection %d already exists", conn.id)
	}
	s.conns[conn.id] = conn
	return nil
}

func (s *Session) getConnection(id uint64) (*clientConn, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	conn, ok := s.conns[id]
	return conn, ok
}

func (s *Session) takeConnection(id uint64) (*clientConn, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn, ok := s.conns[id]
	if ok {
		delete(s.conns, id)
	}
	return conn, ok
}

func (s *Session) finishConnection(id uint64, reason string, reset bool, notify bool) {
	conn, ok := s.takeConnection(id)
	if !ok {
		return
	}
	if notify {
		_ = s.sendClose(conn, reason, reset)
	}
	conn.Close()
}

func (s *Session) closeAllConnections(reason string, reset bool, notify bool) {
	s.mu.Lock()
	conns := make([]*clientConn, 0, len(s.conns))
	for id, conn := range s.conns {
		delete(s.conns, id)
		conns = append(conns, conn)
	}
	s.mu.Unlock()
	for _, conn := range conns {
		if notify {
			_ = s.sendClose(conn, reason, reset)
		}
		conn.Close()
	}
}

func (s *Session) sendData(id uint64, payload []byte) error {
	chunkSize := s.cfg.chunkSize
	for cursor := 0; cursor < len(payload); cursor += chunkSize {
		end := cursor + chunkSize
		if end > len(payload) {
			end = len(payload)
		}
		chunk := append([]byte(nil), payload[cursor:end]...)
		req := pb.ExposePortRequest_builder{
			Data: pb.ExposePortConnectionData_builder{
				ConnectionId: proto.Uint64(id),
				Payload:      chunk,
			}.Build(),
		}.Build()
		if err := s.send(req); err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) sendClose(conn *clientConn, reason string, reset bool) error {
	var sendErr error
	conn.notifyOnce.Do(func() {
		req := pb.ExposePortRequest_builder{
			Close: pb.ExposePortConnectionClose_builder{
				ConnectionId: proto.Uint64(conn.id),
				Error:        proto.String(reason),
				Reset:        proto.Bool(reset),
			}.Build(),
		}.Build()
		sendErr = s.send(req)
	})
	return sendErr
}

func (s *Session) sendCloseWithError(id uint64, reason string, reset bool) error {
	req := pb.ExposePortRequest_builder{
		Close: pb.ExposePortConnectionClose_builder{
			ConnectionId: proto.Uint64(id),
			Error:        proto.String(reason),
			Reset:        proto.Bool(reset),
		}.Build(),
	}.Build()
	return s.send(req)
}

func (s *Session) send(req *pb.ExposePortRequest) error {
	s.sendMu.Lock()
	defer s.sendMu.Unlock()
	return s.stream.Send(req)
}

type clientConn struct {
	id         uint64
	conn       net.Conn
	outbound   chan []byte
	done       chan struct{}
	closeOnce  sync.Once
	notifyOnce sync.Once
}

func newClientConn(id uint64, conn net.Conn) *clientConn {
	return &clientConn{
		id:       id,
		conn:     conn,
		outbound: make(chan []byte, 16),
		done:     make(chan struct{}),
	}
}

func (c *clientConn) Enqueue(payload []byte) bool {
	select {
	case <-c.done:
		return false
	default:
	}

	select {
	case c.outbound <- payload:
		return true
	case <-c.done:
		return false
	}
}

func (c *clientConn) Outbound() <-chan []byte {
	return c.outbound
}

func (c *clientConn) Done() <-chan struct{} {
	return c.done
}

func (c *clientConn) Close() {
	c.closeOnce.Do(func() {
		close(c.done)
		_ = c.conn.Close()
	})
}
