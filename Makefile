.PHONY: server tidy proto clean all build_server build
all: build

server: build_server
	./bin/server

build: build_server

build_server:
	mkdir -p bin
	go build -o bin/server cmd/server/server.go

proto:
	mkdir -p gen
	protoc --go_out=gen --go_opt=module=github.com/starathel/gchat/gen \
		--go-grpc_out=gen --go-grpc_opt=module=github.com/starathel/gchat/gen \
		proto/*.proto

tidy:
	go mod tidy
	go mod vendor

clean:
	rm -rf bin
