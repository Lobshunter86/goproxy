package main

import (
	"context"
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"proxy/util"

	"github.com/lucas-clemente/quic-go"
)

const (
	START_SESSION_RETRY = 3
)

var addr = flag.String("addr", "", "client addr to bind")
var serverAddr = flag.String("saddr", "", "server addr to dial to")
var caCrt = flag.String("cacert", "", "ca certificate file")
var clientCrt = flag.String("cert", "", "client sertificate file")
var clientKey = flag.String("key", "", "client private key file")

type LocalServer struct {
	remoteAddr string
	logger     *log.Logger
	tlsCfg     *tls.Config
}

func main() {
	flag.Parse()
	logger := &log.Logger{}
	logger.SetOutput(os.Stdout)

	cfg, err := util.LoadClientCertificate(*caCrt, *clientCrt, *clientKey)
	if err != nil {
		logger.Printf("load certificate error: %v", err)
		return
	}

	server, err := NewLocalServer(cfg, logger, *serverAddr)
	if err != nil {
		logger.Printf("init local server error: %v", err)
		return
	}

	err = server.ListenAndServe(*addr)
	if err != nil {
		logger.Printf("listen and serve error: %v", err)
		return
	}
}

func NewLocalServer(tlsCfg *tls.Config, logger *log.Logger, remoteAddr string) (*LocalServer, error) {
	return &LocalServer{
		tlsCfg:     tlsCfg,
		logger:     logger,
		remoteAddr: remoteAddr,
	}, nil
}

func (s *LocalServer) ListenAndServe(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.logger.Println("proxy server started")
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

	return nil
}
