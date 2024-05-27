.PHONY: build
test:
				go test ./...

build:
				go build -v ./cmd/apiserver

build_amd64:
				GOOS=linux GOARCH=amd64 go build -v ./cmd/apiserver

run:
				./apiserver

.DEFAULT_GOAL := build
