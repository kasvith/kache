# Makefile for kache

all: dep vet test build-kache build-cli

build-kache:
	go build -o bin/kache ./cmd/kache

build-cli:
	go build -o bin/kache-cli ./cmd/kache-cli

dep:
	go get -d ./...

test:
	go test -v ./...

vet:
	go vet ./...

clean:
	rm -f bin/kache bin/kache-cli
