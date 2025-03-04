.PHONY: server tidy
server:
	mkdir -p bin
	go build -o bin/server cmd/server/server.go

tidy:
	go mod tidy
	go mod vendor
