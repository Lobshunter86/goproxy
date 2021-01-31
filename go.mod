module proxy

go 1.15

require (
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5
	github.com/elazarl/goproxy v0.0.0-20200809112317-0581fc3aee2d
	github.com/elazarl/goproxy/ext v0.0.0-20200809112317-0581fc3aee2d // indirect
	github.com/lucas-clemente/quic-go v0.18.0
	github.com/prometheus/client_golang v1.9.0
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/lucas-clemente/quic-go v0.18.0 => /home/lob/go/src/github.com/lucas-clemente/quic-go
