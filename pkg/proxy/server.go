package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/lucas-clemente/quic-go"
	"golang.org/x/sync/errgroup"
)

// TODO: implement a http proxy fullfills ServeConn(net.Conn), it's way more effcient and elegant.

// Change Handler interface function from ServeConn(net.Conn) to Serve(net.Listener) to work with http proxy.
// But it's not a good approach because to work with ServeConn, server needs a buffered listener.
// Using buffered listener adds complexity, server will stop accepting new connection if one of buffers is full.
// also cause more CPU comsumption and worse performance because Go channel use mutex lock underneath.
type Handler interface {
	Serve(net.Listener) error
}

// Conn satisfies net.Conn interface
type Conn struct {
	quic.Session
	quic.Stream
}

func (c Conn) Close() error {
	return c.Session.CloseWithError(0, "")
}

type bufferedListener struct {
	addr net.Addr
	buf  chan net.Conn
}

func newBufferedListener(bufSize int, addr net.Addr) *bufferedListener {
	if bufSize < 0 {
		bufSize = 4096
	}
	return &bufferedListener{addr: addr, buf: make(chan net.Conn, bufSize)}
}

func (l *bufferedListener) Accept() (net.Conn, error) {
	conn, ok := <-l.buf
	if !ok {
		return nil, fmt.Errorf("Accept from closed buffer")
	}
	return conn, nil
}

func (l *bufferedListener) Close() error {
	return nil
}

func (l *bufferedListener) Addr() net.Addr {
	return l.addr
}

type Server struct {
	tlsCfg   *tls.Config
	logger   *log.Logger
	addr     string
	handlers map[string]Handler
	bufSize  map[string]int
}

func NewProxyServer(tlsCfg *tls.Config, logger *log.Logger, addr string) (*Server, error) {
	return &Server{
		tlsCfg:   tlsCfg,
		logger:   logger,
		addr:     addr,
		handlers: map[string]Handler{},
		bufSize:  map[string]int{}}, nil
}

// Handle register handler for given protocol and create a buffered listener for it
// Server will accept connections, then send connection to corresponding buffered listener according to connection's protocol
func (s *Server) Handle(protocol string, handler Handler, bufSize int) error {
	for _, proto := range s.tlsCfg.NextProtos {
		if protocol == proto {
			s.handlers[protocol] = handler
			s.bufSize[protocol] = bufSize
			return nil
		}
	}

	return fmt.Errorf("Protocol not supported or not configured: %s", protocol)
}

func (s *Server) ListenAndServe() error {
	if len(s.handlers) == 0 {
		return fmt.Errorf("no handler registered")
	}

	listener, err := quic.ListenAddr(s.addr, s.tlsCfg, nil)
	if err != nil {
		return err
	}
	defer listener.Close()

	errGroup, ctx := errgroup.WithContext(context.TODO())
	bufListeners := make(map[string]*bufferedListener)
	for proto, handler := range s.handlers {
		// reassign to prevent data race
		p := proto
		h := handler

		bufSize := s.bufSize[p]
		bufListener := newBufferedListener(bufSize, listener.Addr())
		bufListeners[p] = bufListener
		errGroup.Go(func() error {
			return h.Serve(bufListener)
		})
	}

	s.logger.Println("server started")
	go func() {
		for {
			ctx := context.TODO()
			sess, err := listener.Accept(ctx)
			if err != nil {
				s.logger.Printf("[ERROR] accept session err: %v\n", err)
				continue
			}
			s.logger.Printf("[INFO] Accept connection: %v", sess.RemoteAddr())

			stream, err := sess.AcceptStream(ctx)
			if err != nil {
				s.logger.Printf("[ERROR] accept stream err: %v\n", err)
				sess.CloseWithError(0, "") // nolint:errcheck
				continue
			}

			proto := sess.ConnectionState().NegotiatedProtocol
			l, ok := bufListeners[proto]
			if !ok {
				s.logger.Printf("[FATAL] NegotiatedProtocol not from config, proto: %s\n", proto)
				// code should not reach here
				// but quic-go mention NegotiatedProtocol is not guaranteed to be from tlsconfig.NextProtos
			}

			conn := Conn{Session: sess, Stream: stream}
			l.buf <- conn // buffer isn't big enough may block main thread
		}
	}()

	<-ctx.Done()
	return fmt.Errorf("[FATAL] proxy server exit on error")
}
