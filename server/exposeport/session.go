package exposeport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/proto"

	"github.com/Warashi/muscat/v2/pb"
)

// Stream abstracts the server-side ExposePort bidi stream for testability.
type Stream interface {
	Receive() (*pb.ExposePortRequest, error)
	Send(*pb.ExposePortResponse) error
}

// ListenerFunc constructs a listener for the remote port exposure.
type ListenerFunc func(network, address string) (net.Listener, error)

// Config controls Session behaviour.
type Config struct {
	// Listen is used to create the listener on the server.
	Listen ListenerFunc
	// ChunkSize controls payload chunking when sending to the client.
	// When zero or negative, a default value is used.
	ChunkSize int
}

// Session handles a single ExposePort RPC invocation.
type Session struct {
	ctx    context.Context
	stream Stream
	cfg    sessionConfig

	connID   atomic.Uint64
	sendMu   sync.Mutex
	mu       sync.RWMutex
	conns    map[uint64]*serverConn
	wg       sync.WaitGroup
	failOnce sync.Once
}

type sessionConfig struct {
	listen    ListenerFunc
	chunkSize int
}

const defaultChunkSize = 32 * 1024

// NewSession creates a new Session with sane defaults.
func NewSession(ctx context.Context, stream Stream, cfg Config) *Session {
	sessionCfg := sessionConfig{
		listen:    net.Listen,
		chunkSize: cfg.ChunkSize,
	}
	if cfg.Listen != nil {
		sessionCfg.listen = cfg.Listen
	}
	if sessionCfg.chunkSize <= 0 {
		sessionCfg.chunkSize = defaultChunkSize
	}
	return &Session{
		ctx:    ctx,
		stream: stream,
		cfg:    sessionCfg,
		conns:  make(map[uint64]*serverConn),
	}
}

// Run serves the ExposePort session lifecycle.
func (s *Session) Run() error {
	init, err := s.receiveInit()
	if err != nil {
		s.fail(err)
		return err
	}

	listener, remotePort, err := s.openListener(init)
	if err != nil {
		s.fail(err)
		return err
	}
	defer listener.Close()

	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	errCh := make(chan error, 2)
	go func() {
		errCh <- s.acceptLoop(ctx, listener, remotePort)
	}()
	go func() {
		errCh <- s.receiveLoop(ctx)
	}()

	var runErr error
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil && !errors.Is(err, context.Canceled) {
			if runErr == nil {
				runErr = err
			}
		}
		cancel()
	}

	s.closeAllConnections("session ended", true)
	s.wg.Wait()

	if runErr != nil {
		s.fail(runErr)
	}
	return runErr
}

func (s *Session) receiveInit() (*pb.ExposePortInit, error) {
	req, err := s.stream.Receive()
	if err != nil {
		return nil, fmt.Errorf("stream.Receive: %w", err)
	}
	init := req.GetInit()
	if init == nil {
		return nil, fmt.Errorf("first message must be init frame")
	}
	if init.GetLocalPort() == 0 {
		return nil, fmt.Errorf("local_port must be set")
	}
	return init, nil
}

func (s *Session) openListener(init *pb.ExposePortInit) (net.Listener, uint32, error) {
	bindAddress := init.GetBindAddress()
	if bindAddress == "" {
		if init.GetPublic() {
			bindAddress = "0.0.0.0"
		} else {
			bindAddress = "127.0.0.1"
		}
	}

	remotePort := init.GetRemotePort()
	explicitRemotePort := remotePort != 0
	if remotePort == 0 {
		remotePort = init.GetLocalPort()
	}

	address := net.JoinHostPort(bindAddress, strconv.Itoa(int(remotePort)))
	listener, err := s.cfg.listen("tcp", address)
	if err != nil {
		if !explicitRemotePort {
			fallbackAddress := net.JoinHostPort(bindAddress, "0")
			listener, err = s.cfg.listen("tcp", fallbackAddress)
		}
		if err != nil {
			return nil, 0, fmt.Errorf("listen %s: %w", address, err)
		}
	}

	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return nil, 0, fmt.Errorf("listener.Addr: unexpected type %T", listener.Addr())
	}
	return listener, uint32(tcpAddr.Port), nil
}

func (s *Session) acceptLoop(ctx context.Context, listener net.Listener, remotePort uint32) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			return fmt.Errorf("listener.Accept: %w", err)
		}

		connectionID := s.connID.Add(1)
		serverConn := newServerConn(connectionID, conn)

		if err := s.storeConnection(serverConn); err != nil {
			conn.Close()
			return err
		}

		if err := s.sendOpen(serverConn, remotePort, conn.RemoteAddr().String()); err != nil {
			s.finishConnection(connectionID, fmt.Sprintf("send open: %v", err), true)
			return err
		}

		s.wg.Add(2)
		go func() {
			defer s.wg.Done()
			s.runRemoteToClient(ctx, serverConn)
		}()
		go func() {
			defer s.wg.Done()
			s.runClientToRemote(ctx, serverConn)
		}()
	}
}

func (s *Session) receiveLoop(ctx context.Context) error {
	for {
		req, err := s.stream.Receive()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("stream.Receive: %w", err)
		}

		if data := req.GetData(); data != nil {
			if len(data.GetPayload()) == 0 {
				continue
			}
			conn, ok := s.getConnection(data.GetConnectionId())
			if !ok {
				continue
			}
			payload := append([]byte(nil), data.GetPayload()...)
			if !conn.Enqueue(payload) {
				continue
			}
			continue
		}
		if close := req.GetClose(); close != nil {
			id := close.GetConnectionId()
			reset := close.GetReset()
			reason := close.GetError()
			s.finishConnection(id, reason, reset)
			continue
		}
		if req.GetInit() != nil {
			continue
		}

		select {
		case <-ctx.Done():
			return context.Canceled
		default:
		}
	}
}

func (s *Session) runRemoteToClient(ctx context.Context, conn *serverConn) {
	buf := make([]byte, s.cfg.chunkSize)
	for {
		select {
		case <-ctx.Done():
			s.finishConnection(conn.id, "context canceled", true)
			return
		case <-conn.Done():
			return
		default:
		}

		n, err := conn.conn.Read(buf)
		if n > 0 {
			payload := append([]byte(nil), buf[:n]...)
			if err := s.sendData(conn.id, payload); err != nil {
				s.finishConnection(conn.id, fmt.Sprintf("send data: %v", err), true)
				return
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				s.finishConnection(conn.id, "", false)
			} else {
				s.finishConnection(conn.id, err.Error(), true)
			}
			return
		}
	}
}

func (s *Session) runClientToRemote(ctx context.Context, conn *serverConn) {
	for {
		select {
		case <-ctx.Done():
			s.finishConnection(conn.id, "context canceled", true)
			return
		case <-conn.Done():
			return
		case payload := <-conn.Inbound():
			if len(payload) == 0 {
				continue
			}
			if _, err := conn.conn.Write(payload); err != nil {
				s.finishConnection(conn.id, err.Error(), true)
				return
			}
		}
	}
}

func (s *Session) storeConnection(conn *serverConn) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.conns[conn.id]; exists {
		return fmt.Errorf("connection %d already exists", conn.id)
	}
	s.conns[conn.id] = conn
	return nil
}

func (s *Session) getConnection(id uint64) (*serverConn, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	conn, ok := s.conns[id]
	return conn, ok
}

func (s *Session) takeConnection(id uint64) (*serverConn, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn, ok := s.conns[id]
	if ok {
		delete(s.conns, id)
	}
	return conn, ok
}

func (s *Session) finishConnection(id uint64, reason string, reset bool) {
	conn, ok := s.takeConnection(id)
	if !ok {
		return
	}
	_ = s.sendClose(conn, reason, reset)
	conn.Close()
}

func (s *Session) closeAllConnections(reason string, reset bool) {
	s.mu.Lock()
	conns := make([]*serverConn, 0, len(s.conns))
	for id, conn := range s.conns {
		delete(s.conns, id)
		conns = append(conns, conn)
	}
	s.mu.Unlock()
	for _, conn := range conns {
		_ = s.sendClose(conn, reason, reset)
		conn.Close()
	}
}

func (s *Session) sendOpen(conn *serverConn, remotePort uint32, remoteAddress string) error {
	resp := pb.ExposePortResponse_builder{
		Open: pb.ExposePortConnectionOpen_builder{
			ConnectionId:  proto.Uint64(conn.id),
			RemotePort:    proto.Uint32(remotePort),
			RemoteAddress: proto.String(remoteAddress),
		}.Build(),
	}.Build()
	return s.send(resp)
}

func (s *Session) sendData(id uint64, payload []byte) error {
	chunkSize := s.cfg.chunkSize
	for cursor := 0; cursor < len(payload); cursor += chunkSize {
		end := cursor + chunkSize
		if end > len(payload) {
			end = len(payload)
		}
		chunk := append([]byte(nil), payload[cursor:end]...)
		resp := pb.ExposePortResponse_builder{
			Data: pb.ExposePortConnectionData_builder{
				ConnectionId: proto.Uint64(id),
				Payload:      chunk,
			}.Build(),
		}.Build()
		if err := s.send(resp); err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) sendClose(conn *serverConn, reason string, reset bool) error {
	var sendErr error
	conn.notifyOnce.Do(func() {
		resp := pb.ExposePortResponse_builder{
			Close: pb.ExposePortConnectionClose_builder{
				ConnectionId: proto.Uint64(conn.id),
				Error:        proto.String(reason),
				Reset:        proto.Bool(reset),
			}.Build(),
		}.Build()
		sendErr = s.send(resp)
	})
	return sendErr
}

func (s *Session) sendError(message string) {
	if message == "" {
		return
	}
	_ = s.send(
		pb.ExposePortResponse_builder{
			Error: pb.ExposePortError_builder{
				Message: proto.String(message),
			}.Build(),
		}.Build(),
	)
}

func (s *Session) send(resp *pb.ExposePortResponse) error {
	s.sendMu.Lock()
	defer s.sendMu.Unlock()
	return s.stream.Send(resp)
}

func (s *Session) fail(err error) {
	if err == nil {
		return
	}
	s.failOnce.Do(func() {
		s.sendError(err.Error())
	})
}

type serverConn struct {
	id         uint64
	conn       net.Conn
	inbound    chan []byte
	done       chan struct{}
	closeOnce  sync.Once
	notifyOnce sync.Once
}

func newServerConn(id uint64, conn net.Conn) *serverConn {
	return &serverConn{
		id:      id,
		conn:    conn,
		inbound: make(chan []byte, 16),
		done:    make(chan struct{}),
	}
}

func (c *serverConn) Enqueue(payload []byte) bool {
	select {
	case <-c.done:
		return false
	default:
	}

	select {
	case c.inbound <- payload:
		return true
	case <-c.done:
		return false
	}
}

func (c *serverConn) Inbound() <-chan []byte {
	return c.inbound
}

func (c *serverConn) Done() <-chan struct{} {
	return c.done
}

func (c *serverConn) Close() {
	c.closeOnce.Do(func() {
		close(c.done)
		_ = c.conn.Close()
	})
}
