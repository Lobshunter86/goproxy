package util

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

func LoadServerCertificate(caCrt string, serverCrt string, serverKey string) (*tls.Config, error) {
	Cert, err := tls.LoadX509KeyPair(serverCrt, serverKey)
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
		Certificates: []tls.Certificate{Cert},
		ClientCAs:    caCrtPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		NextProtos:   []string{"http"},
	}, nil
}

func LoadClientCertificate(caCrt string, clientCrt string, clientKey string) (*tls.Config, error) {
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
		Certificates:       []tls.Certificate{Cert},
		RootCAs:            caCrtPool,
		InsecureSkipVerify: true,
		NextProtos:         []string{"http"},
	}, nil
}
