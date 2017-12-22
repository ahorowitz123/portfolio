all: client server

client: deps client/ lib/
	go build -o bin/client client/client.go client/util.go client/clientrpc.go

server: deps server/ lib/
	go build -o bin/server server/server.go

deps: bin-dir

bin-dir:
	mkdir -p bin

clean:
	rm -rf bin
