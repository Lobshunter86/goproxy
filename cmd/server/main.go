package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/armon/go-socks5"

	"github.com/lobshunter86/goproxy/pkg/proxy"
	"github.com/lobshunter86/goproxy/pkg/util"
	"github.com/lobshunter86/goproxy/pkg/version"
)

const defaultBufferSize = 4096

// "protocols" field in yaml need to match these values
var supportedProtocols = map[string]bool{"http": true, "socks5": true}

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

	cfgData, err := ioutil.ReadFile(*configFile)
	util.DoneOrDieWithMesg(err, fmt.Sprintf("[FATAL] read config file error: %v\n", err))

	cfg, err := proxy.ParseRemoteServerCfg(cfgData)
	util.DoneOrDieWithMesg(err, fmt.Sprintf("[FATAL] parse config file error: %v\n", err))

OUTER:
	for _, proto := range cfg.Protocols {
		for supported := range supportedProtocols {
			if proto == supported {
				continue OUTER
			}
		}

		panic(fmt.Sprintf("unsupported protocol: %s", proto))
	}

	clientCaCert, err := ioutil.ReadFile(cfg.CaCert)
	util.DoneOrDieWithMesg(err, fmt.Sprintf("read client ca cert %v", err))

	certProvider, err := proxy.NewLocalProvider(cfg.ServerCert, cfg.ServerKey)
	util.DoneOrDieWithMesg(err, fmt.Sprintf("NewLocalProvider %v", err))

	tlsCfg := proxy.NewServerTLSConfig(clientCaCert, certProvider, cfg.Protocols)

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
