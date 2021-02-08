package main

import (
	"flag"
	"log"
	"os"

	"github.com/armon/go-socks5"

	"github.com/lobshunter86/goproxy/pkg/proxy"
	"github.com/lobshunter86/goproxy/pkg/version"
)

const defaultBufferSize = 4096

func main() {
	if len(os.Args) > 1 &&
		(os.Args[1] == "-v" || os.Args[1] == "--version") {
		version.PrintVersion()
		return
	}

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

	if err = server.Handle("socks5", socksHandler, defaultBufferSize); err != nil {
		logger.Printf("[FATAL] register socks5 handler error: %v\n", err)
		return
	}

	httpHandler := proxy.NewHTTPHandler()
	if err = server.Handle("http", httpHandler, defaultBufferSize); err != nil {
		logger.Printf("[FATAL] register http handler error: %v\n", err)
		return
	}

	// starting server
	if server.ListenAndServe() != nil {
		logger.Printf("[FATAL] server exit, err: %v", err)
	}
}
