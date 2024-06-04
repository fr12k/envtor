.PHONY: build
# Define the PATH variable
PATH := $(PATH):$(shell go env GOPATH)/bin

lint:
	golangci-lint run ./...

build:
	go build -o build/envtor main.go

test:
	@go clean -testcache
	@go test -v -cover -coverprofile coverage.txt -race ./...
	@go tool cover -func coverage.txt

example:
	printf "ENVIRONMENT1=hello\nENVIRONMENT2=world" | ./build/envtor | docker-compose -f - up
