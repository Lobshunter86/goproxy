all: server-side client-side

server-side: 
	go build -o bin/server ./cmd/server/

client-side: 
	go build -o bin/client ./cmd/client/