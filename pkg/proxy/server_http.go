package proxy

import (
	"net"
	"net/http"

	"github.com/elazarl/goproxy"
)

// This file is a wrapper for goproxy http proxy library
// Intends to make it satisfy Handler interface
// However, wrapping this implement isn't a good idea, because Handler interface doesn't work well with http.Handler
// it's better to impletement a http handler from scatch

type HttpHandler struct {
	*goproxy.ProxyHttpServer
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{goproxy.NewProxyHttpServer()}
}

func (h *HttpHandler) Serve(l net.Listener) error {
	return http.Serve(l, h.ProxyHttpServer)
}
