all: vendor build test

.PHONY: vendor
vendor:
	GOPROXY=direct GOSUMDB=off go mod vendor

.PHONY: build
build:
	GOOS=linux go build -ldflags="-s -w" -v -o bin/spacetrouble github.com/antonzhukov/spacetrouble/cmd/spacetrouble

.PHONY: test
test:
	go test ./...
