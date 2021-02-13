package tests

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/armon/go-socks5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goproxy "golang.org/x/net/proxy"
	"golang.org/x/sync/errgroup"

	"github.com/lobshunter86/goproxy/pkg/proxy"
	"github.com/lobshunter86/goproxy/pkg/util"
)

// TODO: maybe extract common code from main package

// FIXME: support shutting down servers, so that test can run multiple times
// 		  currently test server will occupy UDP/TCP port, second run will fail to bind port
// 		  shutdown server can be down by close net.Listener

func TestServe(t *testing.T) {
	curDir := getIntegrationDir()
	logger := new(log.Logger)
	logger.SetOutput(os.Stdout)

	localCfg := localConfig(curDir)
	remoteCfg := remoteConfig(curDir)

	clientCertProvider, err := proxy.NewLocalProvider(curDir+"/client.cert", curDir+"/client.key")
	require.Nil(t, err)

	serverCertProvider, err := proxy.NewLocalProvider(curDir+"/server.cert", curDir+"/server.key")
	require.Nil(t, err)

	go testRemoteServer([]byte(remoteCfg), serverCertProvider, logger)
	go testLocalServer([]byte(localCfg), clientCertProvider, logger)
	go http.ListenAndServe("127.0.0.1:18888", nil) // nolint:errcheck

	time.Sleep(time.Second) // FIXME: wait server started

	// test http
	proxyURL, err := url.Parse("http://127.0.0.1:18123")
	require.Nil(t, err)

	httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	resp, err := httpClient.Get("http://127.0.0.1:18888")
	assert.Nil(t, err)

	resp.Body.Close()
	t.Log("status code:", resp.StatusCode)

	// test socks5
	dialer, err := goproxy.SOCKS5("tcp", "127.0.0.1:11080", nil, goproxy.Direct)
	require.Nil(t, err)
	ctxDialer := func(ctx context.Context, network, addr string) (net.Conn, error) { return dialer.Dial(network, addr) }

	httpTransport := &http.Transport{}
	httpClient = &http.Client{Transport: httpTransport}
	httpTransport.DialContext = ctxDialer
	resp, err = httpClient.Get("http://127.0.0.1:18888")
	assert.Nil(t, err)

	resp.Body.Close()
	t.Log("status code:", resp.StatusCode)
}

func testLocalServer(configData []byte, certProvider proxy.ClientCertificateProvider, logger *log.Logger) {
	// parse
	localCfg, err := proxy.ParseLocalServerCfg(configData)
	util.DoneOrDieWithMesg(err, "parse config file error")

	servers := []*proxy.LocalServer{}

	for _, cfg := range localCfg.Servers {
		caCert, err := ioutil.ReadFile(cfg.CaCert)
		util.DoneOrDieWithMesg(err, "read certificate")

		tlsCfg := proxy.NewClientTLSConfig(caCert, certProvider, cfg.Protocol)
		tlsCfg.InsecureSkipVerify = true

		server, err := proxy.NewLocalServer(tlsCfg, logger, cfg.Protocol, cfg.ServerAddr)
		util.DoneOrDieWithMesg(err, fmt.Sprintf("listen and serve error: %v", err))

		servers = append(servers, server)
	}

	errGroup, ctx := errgroup.WithContext(context.TODO())

	// starts serers
	for idx, server := range servers {
		s := server
		i := idx
		errGroup.Go(func() error {
			return s.ListenAndServe(localCfg.Servers[i].LocalAddr)
		})
	}

	<-ctx.Done()
	logger.Fatal("[FATAL] Local server exit on error")
}

func testRemoteServer(configData []byte, certProvider proxy.ServerCertificateProvider, logger *log.Logger) {
	var defaultBufferSize = 4096

	cfg, err := proxy.ParseRemoteServerCfg(configData)
	util.DoneOrDieWithMesg(err, fmt.Sprintf("[FATAL] parse config file error: %v\n", err))

	clientCaCert, err := ioutil.ReadFile(cfg.CaCert)
	util.DoneOrDieWithMesg(err, fmt.Sprintf("read client ca cert %v", err))

	tlsCfg := proxy.NewServerTLSConfig(clientCaCert, certProvider, cfg.Protocols)

	server, err := proxy.NewProxyServer(tlsCfg, logger, cfg.Addr)
	util.DoneOrDieWithMesg(err, fmt.Sprintf("[FATAL] Setup ProxyServer error: %v\n", err))

	// register handlers
	socksCfg := &socks5.Config{}
	socksHandler, err := socks5.New(socksCfg)
	util.DoneOrDieWithMesg(err, fmt.Sprintf("[FATAL] Init socks5 proxy error: %v\n", err))

	err = server.Handle("socks5", socksHandler, defaultBufferSize)
	util.DoneOrDieWithMesg(err, fmt.Sprintf("[FATAL] register socks5 handler error: %v\n", err))

	httpHandler := proxy.NewHTTPHandler()
	err = server.Handle("http", httpHandler, defaultBufferSize)
	util.DoneOrDieWithMesg(err, fmt.Sprintf("[FATAL] register http handler error: %v\n", err))

	// starting server
	err = server.ListenAndServe()
	logger.Fatalf(fmt.Sprintf("[FATAL] server exit, err: %v", err))
}

func getIntegrationDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("unable to determine the current file path")
	}

	return path.Dir(filename)
}
