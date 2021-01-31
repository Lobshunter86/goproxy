package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"

	"github.com/lobshunter86/goproxy/pkg/proxy"
	"github.com/lobshunter86/goproxy/pkg/version"
)

func main() {
	if len(os.Args) > 1 &&
		(os.Args[1] == "-v" || os.Args[1] == "--version") {
		version.PrintVersion()
		return
	}

	configFile := flag.String("config", "./config.yaml", "config file path")
	flag.Parse()

	localCfg, err := proxy.ParseLocalServerCfg(*configFile)
	if err != nil {
		panic("parse config file error")
	}

	if len(localCfg.Servers) == 0 {
		panic("no configuration read")
	}

	servers := []*proxy.LocalServer{}

	for _, cfg := range localCfg.Servers {
		logger := &log.Logger{}
		logger.SetOutput(os.Stdout)

		tlsCfg, err := proxy.LoadClientCertificate(cfg.CaCert, cfg.ClientCert, cfg.ClientKey, cfg.Protocol)
		if err != nil {
			logger.Printf("init local server error: %v", err)
			return
		}

		server, err := proxy.NewLocalServer(tlsCfg, logger, cfg.Protocol, cfg.ServerAddr)
		if err != nil {
			logger.Printf("listen and serve error: %v", err)
			return
		}

		servers = append(servers, server)
	}

	errGroup, ctx := errgroup.WithContext(context.TODO())

	if len(localCfg.Global.MetricsAddr) != 0 {
		http.Handle("/metrics", promhttp.Handler())
		errGroup.Go(func() error {
			return http.ListenAndServe(localCfg.Global.MetricsAddr, nil)
		})
	}

	for idx, server := range servers {
		// reassign to prevent data race
		// errgroup.Go calls "go" later on, make data race in this condiction
		s := server
		i := idx
		errGroup.Go(func() error {
			return s.ListenAndServe(localCfg.Servers[i].LocalAddr)
		})
	}

	<-ctx.Done()
	println("[FATAL] Local server exit on error")
}
