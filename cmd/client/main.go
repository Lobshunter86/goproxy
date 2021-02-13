package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"

	"github.com/lobshunter86/goproxy/pkg/proxy"
	"github.com/lobshunter86/goproxy/pkg/util"
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

	// parse
	cfgData, err := ioutil.ReadFile(*configFile)
	util.DoneOrDieWithMesg(err, fmt.Sprintf("read config file %v", err))

	localCfg, err := proxy.ParseLocalServerCfg(cfgData)
	util.DoneOrDieWithMesg(err, fmt.Sprintf("parse config file %v", err))

	if len(localCfg.Servers) == 0 {
		panic("no configuration read")
	}

	// new local servers
	servers := []*proxy.LocalServer{}

	logger := &log.Logger{}
	logger.SetOutput(os.Stdout)
	for _, cfg := range localCfg.Servers {
		cert, err := ioutil.ReadFile(cfg.CaCert)
		util.DoneOrDieWithMesg(err, "read certificate")

		certProvider, err := proxy.NewLocalProvider(cfg.ClientCert, cfg.ClientKey)
		util.DoneOrDieWithMesg(err, "load certificate")

		tlsCfg := proxy.NewClientTLSConfig(cert, certProvider, cfg.Protocol)

		server, err := proxy.NewLocalServer(tlsCfg, logger, cfg.Protocol, cfg.ServerAddr)
		if err != nil {
			logger.Printf("listen and serve error: %v", err)
			return
		}

		servers = append(servers, server)
	}

	errGroup, ctx := errgroup.WithContext(context.TODO())

	// metrics
	if len(localCfg.Global.MetricsAddr) != 0 {
		http.Handle("/metrics", promhttp.Handler())
		errGroup.Go(func() error {
			return http.ListenAndServe(localCfg.Global.MetricsAddr, nil)
		})
	}

	// starts serers
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
