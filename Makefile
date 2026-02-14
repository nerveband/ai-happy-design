.PHONY: build build-go build-plugin sync-plugin run-mcp run-ws test clean install lint

BINARY=ai-happy-design
VERSION=1.0.0

build-plugin:
	cd plugin && npm install && npm run build

sync-plugin:
	mkdir -p internal/plugin/files/dist
	cp plugin/manifest.json internal/plugin/files/manifest.json
	cp plugin/dist/code.js internal/plugin/files/dist/code.js
	cp plugin/dist/ui.html internal/plugin/files/dist/ui.html

build: sync-plugin
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(BINARY) ./cmd/ai-happy-design

build-go:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(BINARY) ./cmd/ai-happy-design

run-mcp: build
	./bin/$(BINARY) mcp

run-ws: build
	./bin/$(BINARY) ws

test:
	go test ./...

clean:
	rm -rf bin/ internal/plugin/files/

install: build
	cp bin/$(BINARY) $(GOPATH)/bin/

lint:
	golangci-lint run
