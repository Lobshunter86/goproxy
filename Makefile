LDFLAGS := -w -s

default: vet server-side client-side

server-side: 
	go build -ldflags '$(LDFLAGS)' -o bin/server ./cmd/server/

client-side: 
	go build -ldflags '$(LDFLAGS)' -o bin/client ./cmd/client/

vet:
	go vet ./...