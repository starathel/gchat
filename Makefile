.PHONY: server tidy proto clean all build_server build lint client build_client
all: build

server: build_server
	./bin/server

client: build_client
	./bin/client

build: proto build_server build_client

build_server:
	mkdir -p bin
	go build -o bin/server cmd/server/server.go

build_client:
	mkdir -p bin
	go build -o bin/client cmd/client/client.go

proto:
	mkdir -p gen
	protoc --go_out=gen --go_opt=module=github.com/starathel/gchat/gen \
		--go-grpc_out=gen --go-grpc_opt=module=github.com/starathel/gchat/gen \
		proto/*.proto

tidy:
	go mod tidy
	go mod vendor

lint:
	go vet ./...
	staticcheck ./...

clean:
	rm -rf bin
