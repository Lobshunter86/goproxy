module proxy

go 1.14

require (
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5
	github.com/elazarl/goproxy v0.0.0-20200710112657-153946a5f232
	github.com/elazarl/goproxy/ext v0.0.0-20200710112657-153946a5f232 // indirect
	github.com/lucas-clemente/quic-go v0.17.3
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	gopkg.in/yaml.v2 v2.2.4
)

replace github.com/lucas-clemente/quic-go v0.17.3 => /home/lob/go/src/github.com/lucas-clemente/quic-go
