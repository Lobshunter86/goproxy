package main

import (
	"flag"
	"log"
	"os"
	"proxy/util"

	"github.com/armon/go-socks5"
)

var addr = flag.String("addr", "", "server addr to bind")
var caCrt = flag.String("cacert", "", "ca certificate file")
var serverCrt = flag.String("cert", "", "server sertificate file")
var serverKey = flag.String("key", "", "server private key file")

func main() {
	flag.Parse()
	logger := &log.Logger{}
	logger.SetOutput(os.Stdout)

	cfg, err := util.LoadServerCertificate(*caCrt, *serverCrt, *serverKey)
	if err != nil {
		logger.Printf("load certificate error: %v", err)
		return
	}

	server, err := InitServer(cfg, logger)
	if err != nil {
		logger.Printf("init server error: %v", err)
		return
	}

	conf := &socks5.Config{}
	handler, err := socks5.New(conf)
	if err != nil {
		logger.Printf("init socks5 server error: %v", err)
		return
	}

	err = server.ListenAndServe(*addr, handler)
	if err != nil {
		logger.Printf("listen and serve error: %v", err)
		return
	}
}
