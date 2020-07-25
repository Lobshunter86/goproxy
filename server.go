package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"

	"github.com/lucas-clemente/quic-go"
)

type Handler interface {
	ServeConn(net.Conn) error
}

type ProxyConn struct {
	quic.Session
	quic.Stream
}

type ProxyServer struct {
	tlsCfg *tls.Config
	logger *log.Logger
}

func InitServer(tlsCfg *tls.Config, logger *log.Logger) (*ProxyServer, error) {
	return &ProxyServer{tlsCfg: tlsCfg, logger: logger}, nil
}

func (s *ProxyServer) ListenAndServe(addr string, handler Handler) error {
	listener, err := quic.ListenAddr(addr, s.tlsCfg, nil)
	if err != nil {
		return err
	}

	s.logger.Println("server started")
	for {
		ctx := context.TODO()
		sess, err := listener.Accept(ctx)
		if err != nil {
			s.logger.Printf("accept session error: %v", err)
			continue
		} else {
			s.logger.Printf("accept session from: %v", sess.RemoteAddr())
		}

		go s.ServeSession(ctx, sess, handler)
	}
}

func (s *ProxyServer) ServeSession(ctx context.Context, sess quic.Session, handler Handler) error {
	for {
		stream, err := sess.AcceptStream(ctx)
		if err != nil {
			s.logger.Printf("accept stream error: %v", err)
			return err
		} else {
			s.logger.Printf("accept new stream from: %v", sess.RemoteAddr())
		}

		conn := ProxyConn{Session: sess, Stream: stream}
		go handler.ServeConn(conn)
	}
}
