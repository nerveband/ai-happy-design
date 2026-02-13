.PHONY: build run-mcp run-ws test clean install lint

BINARY=ai-happy-design
VERSION=1.0.0

build:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(BINARY) ./cmd/ai-happy-design

run-mcp: build
	./bin/$(BINARY) mcp

run-ws: build
	./bin/$(BINARY) ws

test:
	go test ./...

clean:
	rm -rf bin/

install: build
	cp bin/$(BINARY) $(GOPATH)/bin/

lint:
	golangci-lint run
