all:
	go build -o bin/kache cmd/main.go

dep:
	go get -d ./...