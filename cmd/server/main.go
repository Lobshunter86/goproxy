package main

import (
	"flag"
	"log"
	"os"
	"proxy"

	"github.com/armon/go-socks5"
)

// func main() {
// 	flag.Parse()
// 	logger := &log.Logger{}
// 	logger.SetOutput(os.Stdout)

// 	cfg, err := proxy.LoadServerCertificate(*caCrt, *serverCrt, *serverKey)
// 	if err != nil {
// 		logger.Printf("load certificate error: %v", err)
// 		return
// 	}

// 	server, err := proxy.InitServer(cfg, logger)
// 	if err != nil {
// 		logger.Printf("init server error: %v", err)
// 		return
// 	}

// 	conf := &socks5.Config{}
// 	handler, err := socks5.New(conf)
// 	if err != nil {
// 		logger.Printf("init socks5 server error: %v", err)
// 		return
// 	}

// 	err = server.ListenAndServe(*addr, handler)
// 	if err != nil {
// 		logger.Printf("listen and serve error: %v", err)
// 		return
// 	}
// }
const DEFAULT_BUFFER_SIZE = 4096

func main() {
	configFile := flag.String("config", "", "config file path")
	flag.Parse()

	logger := &log.Logger{}
	logger.SetOutput(os.Stdout)

	cfg, err := proxy.ParseRemoteServerCfg(*configFile)
	if err != nil {
		logger.Printf("[FATAL] parse config file error: %v\n", err)
		return
	}

	tlsCfg, err := proxy.LoadServerCertificate(cfg.CaCert, cfg.ServerCert, cfg.ServerKey, cfg.Protocols)
	if err != nil {
		logger.Printf("[FATAL] Load server Certificate error: %v\n", err)
		return
	}

	server, err := proxy.NewProxyServer(tlsCfg, logger, cfg.Addr)
	if err != nil {
		logger.Printf("[FATAL] Setup ProxyServer error: %v\n", err)
		return
	}

	// register handlers
	socksCfg := &socks5.Config{}
	socksHandler, err := socks5.New(socksCfg)
	if err != nil {
		logger.Printf("[FATAL] Init socks5 proxy error: %v\n", err)
		return
	}

	if err = server.Handle("socks5", socksHandler, DEFAULT_BUFFER_SIZE); err != nil {
		logger.Printf("[FATAL] register socks5 handler error: %v\n", err)
		return
	}

	httpHandler := proxy.NewHttpHandler()
	if err = server.Handle("http", httpHandler, DEFAULT_BUFFER_SIZE); err != nil {
		logger.Printf("[FATAL] register http handler error: %v\n", err)
		return
	}

	// starting server
	if server.ListenAndServe() != nil {
		logger.Printf("[FATAL] server exit, err: %v", err)
	}
}