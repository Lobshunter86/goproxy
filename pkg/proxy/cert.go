package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

func LoadServerCertificate(caCrts []string, serverCrt string, serverKey string, protos []string) (*tls.Config, error) {
	Cert, err := tls.LoadX509KeyPair(serverCrt, serverKey)
	if err != nil {
		return nil, err
	}

	caCrtPool := x509.NewCertPool()

	for _, caCrt := range caCrts {
		caCert, err := ioutil.ReadFile(caCrt)
		if err != nil {
			return nil, err
		}

		caCrtPool.AppendCertsFromPEM(caCert)
	}

	for _, p := range protos {
		if !supportedProtocols[p] {
			return nil, fmt.Errorf("protocol %s not supported", p)
		}
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{Cert},
		ClientCAs:    caCrtPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		NextProtos:   protos,
	}, nil
}

func LoadClientCertificate(caCrt string, clientCrt string, clientKey string, nextProto string) (*tls.Config, error) {
	Cert, err := tls.LoadX509KeyPair(clientCrt, clientKey)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(caCrt)
	if err != nil {
		return nil, err
	}

	caCrtPool := x509.NewCertPool()
	caCrtPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{Cert},
		RootCAs:      caCrtPool,
		NextProtos:   []string{nextProto},
	}, nil
}
