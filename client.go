package proxy

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"

	"github.com/lucas-clemente/quic-go"
)

type LocalServer struct {
	remoteAddr string
	logger     *log.Logger
	tlsCfg     *tls.Config
}

func NewLocalServer(tlsCfg *tls.Config, logger *log.Logger, remoteAddr string) (*LocalServer, error) {
	return &LocalServer{
		tlsCfg:     tlsCfg,
		logger:     logger,
		remoteAddr: remoteAddr,
	}, nil
}

func (s *LocalServer) ListenAndServe(addr string) (err error) {
	defer func() {
		if err != nil {
			s.logger.Printf("[FATAL] Server %s exit on error: %v", addr, err)
		}
	}()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		println("yes is here:", addr)
		return err
	}

	s.logger.Println("proxy server started: ", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Printf("accept error: %v", err)
			continue
		} else {
			s.logger.Println("accept connection")
		}

		go s.ServeConn(conn)
	}
}

func (s *LocalServer) ServeConn(conn net.Conn) error {
	sess, err := quic.DialAddr(s.remoteAddr, s.tlsCfg, nil)
	if err != nil {
		s.logger.Printf("ServeConn dial error: %v", err)
		return err
	}

	// TODO: handle error properly, golang use syscall.Errno for this
	stream, err := sess.OpenStreamSync(context.Background())
	if err != nil {
		s.logger.Printf("ServeConn openstream error: %v", err)
		return err
	}

	done := make(chan struct{}, 1)
	go func() {
		io.Copy(stream, conn)
		done <- struct{}{}
	}()

	io.Copy(conn, stream)
	<-done

	conn.Close()
	stream.Close()
	sess.CloseWithError(0, "")
	s.logger.Printf("closed connection: %s -> %s", sess.LocalAddr().String(), sess.RemoteAddr().String())

	return nil
}
