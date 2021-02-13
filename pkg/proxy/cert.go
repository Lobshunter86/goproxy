package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type ServerCertificateProvider interface {
	GetCert(*tls.ClientHelloInfo) (*tls.Certificate, error)
}

type ClientCertificateProvider interface {
	GetClientCert(*tls.CertificateRequestInfo) (*tls.Certificate, error)
}

type StaticProvider struct {
	cert tls.Certificate
}

func NewLocalProvider(certFile string, keyFile string) (*StaticProvider, error) {
	Cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	provider := new(StaticProvider)
	provider.cert = Cert
	return provider, nil
}

func (p *StaticProvider) GetClientCert(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
	return &p.cert, nil
}

func (p *StaticProvider) GetCert(helo *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return &p.cert, nil
}

type ACMEProvider struct {
	mgr autocert.Manager
}

// !!! ACMEProvider is untested
// not sure if it works
func NewACMEProvider(domains []string) *ACMEProvider {
	provider := new(ACMEProvider)
	provider.mgr = autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("certs"),
		HostPolicy: autocert.HostWhitelist(domains...),
	}

	return provider
}

// StartHTTP starts the ACME HTTP handler
func (p *ACMEProvider) StartHTTP() error {
	return http.ListenAndServe(":80", p.mgr.HTTPHandler(nil))
}

func (p *ACMEProvider) GetCert(helo *tls.ClientHelloInfo) (*tls.Certificate, error) {
	fmt.Println("got GetCert request")
	defer fmt.Println("done GetCert request")
	helo.SupportedProtos = append([]string{acme.ALPNProto}, helo.SupportedProtos...)
	defer func() { helo.SupportedProtos = helo.SupportedProtos[1:] }()
	return p.mgr.GetCertificate(helo)
}

// client authentication is process by TLS client certificate verification
func NewServerTLSConfig(clientCaCrt []byte, provider ServerCertificateProvider, protos []string) *tls.Config {
	caCrtPool := x509.NewCertPool()
	caCrtPool.AppendCertsFromPEM(clientCaCrt)

	return &tls.Config{
		MinVersion:     tls.VersionTLS13,
		GetCertificate: provider.GetCert,
		ClientCAs:      caCrtPool,
		ClientAuth:     tls.RequireAndVerifyClientCert,
		NextProtos:     protos,
	}
}

func NewClientTLSConfig(clientCaCrt []byte, provider ClientCertificateProvider, nextProto string) *tls.Config {
	caCrtPool := x509.NewCertPool()
	caCrtPool.AppendCertsFromPEM(clientCaCrt)

	return &tls.Config{
		MinVersion:           tls.VersionTLS13,
		GetClientCertificate: provider.GetClientCert,
		RootCAs:              caCrtPool,
		NextProtos:           []string{nextProto},
	}
}
