.DEFAULT_GOAL := build

.PHONY:test fmt vet build clean

test:
	go test -v ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

build: test fmt vet
	go build -o bin/squire cmd/squire/main.go

get:
	go get ./...

clean:
	rm -rf bin/
