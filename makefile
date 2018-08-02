all:
	go build -o bin/kache cmd/kache

dep:
	go get -d ./...