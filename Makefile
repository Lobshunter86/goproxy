all: server-side client-side

server-side: 
	go build -o bin/server .

client-side: 
	go build -o bin/client ./client/