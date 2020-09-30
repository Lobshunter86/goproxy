package main

import (
	"context"
	"flag"
	"log"
	"os"
	"proxy"

	"golang.org/x/sync/errgroup"
)

func main() {
	configFile := flag.String("config", "./config.yaml", "config file path")
	flag.Parse()

	cfgs, err := proxy.ParseLocalServerCfg(*configFile)
	if err != nil {
		panic("parse config file error")
	}
	if len(cfgs) == 0 {
		panic("no configuration read")
	}

	servers := []*proxy.LocalServer{}
	for _, cfg := range cfgs {
		logger := &log.Logger{}
		logger.SetOutput(os.Stdout)

		tlsCfg, err := proxy.LoadClientCertificate(cfg.CaCert, cfg.ClientCert, cfg.ClientKey, cfg.Protocol)
		if err != nil {
			logger.Printf("init local server error: %v", err)
			return
		}

		server, err := proxy.NewLocalServer(tlsCfg, logger, cfg.ServerAddr)
		if err != nil {
			logger.Printf("listen and serve error: %v", err)
			return
		}

		servers = append(servers, server)
	}

	errGroup, ctx := errgroup.WithContext(context.TODO())

	for idx, server := range servers {
		// reassign to prevent data race
		// errgroup.Go calls "go" later on, make data race in this condiction
		s := server
		i := idx
		errGroup.Go(func() error {
			return s.ListenAndServe(cfgs[i].LocalAddr)
		})
	}

	<-ctx.Done()
	println("[FATAL] Local server exit on error")
}
