# Makefile for kache

all: dep vet fmt test build-kache build-cli

build-kache:
	go build -o bin/kache ./cmd/kache

build-cli:
	go build -o bin/kache-cli ./cmd/kache-cli

dep:
	dep ensure

test:
	go test -v ./...

vet:
	go vet ./...

fmt:
	gofmt -w `find . -name '*.go' | grep -v vendor`
	goimports -w -local='github.com/kasvith/kache' `find . -name '*.go' | grep -v vendor`

clean:
	rm -rf bin/*
