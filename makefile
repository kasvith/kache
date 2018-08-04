# Makefile for kache

all: dep test build-kache build-cli

build-kache:
	go build -o bin/kache ./cmd/kache

build-cli:
	go build -o bin/kache-cli ./cmd/kache-cli

dep:
	go get -d ./...

test:
	go test -v ./...

clean:
	rm -f bin/kache bin/kache-cli
