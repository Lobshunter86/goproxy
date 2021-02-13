package tests

import "fmt"

// func localSocks5Config() {}
// func localHTTPConfig()   {}

func localConfig(certPath string) string {
	return fmt.Sprintf(`
servers:
  - name: "http proxy" # "name" is just for human to read, useless to code
    protocol: "http" # use for protocol negotiation with server, by setting TLS.Config.NextProto field
    localAddr: "127.0.0.1:18123" # local address to listen to
    serverAddr: "127.0.0.1:8888" # remote server address, server side can handles multiple protocols on single port, separate by "protocol" configuration

    # certificates and keys for TLS
    caCert: "%s/client.cert"
    clientCert: "%s/client.cert"
    clientKey: "%s/client.key"

  - name: "socks5 proxy"
    protocol: "socks5"
    localAddr: "127.0.0.1:11080"
    serverAddr: "127.0.0.1:8888"
    caCert: "%s/client.cert"
    clientCert: "%s/client.cert"
    clientKey: "%s/client.key"
`, certPath, certPath, certPath, certPath, certPath, certPath)
}

func remoteConfig(certPath string) string {
	return fmt.Sprintf(`
addr: "127.0.0.1:8888"
caCert: "%s/client.cert"
serverCert: "%s/server.cert"
serverKey: "%s/server.key"
protocols: 
  - "http"
  - "socks5"
`, certPath, certPath, certPath)
}
