.PHONY: build build-go build-plugin sync-plugin run-mcp run-ws test clean install lint

FIGMA_BINARY=ahd-figma
VERSION=0.0.0-dev

build-plugin:
	cd plugin && npm ci && npm run check

sync-plugin:
	mkdir -p internal/plugin/files/dist
	cp plugin/manifest.json internal/plugin/files/manifest.json
	cp plugin/dist/code.js internal/plugin/files/dist/code.js
	cp plugin/dist/ui.html internal/plugin/files/dist/ui.html

build: build-plugin sync-plugin
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(FIGMA_BINARY) ./cmd/ahd-figma

build-go:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(FIGMA_BINARY) ./cmd/ahd-figma

run-mcp: build
	./bin/$(FIGMA_BINARY) mcp

run-ws: build
	./bin/$(FIGMA_BINARY) ws

test:
	go test ./...

clean:
	rm -rf bin/ internal/plugin/files/

install: build
	cp bin/$(FIGMA_BINARY) $(GOPATH)/bin/

# Full deploy: build plugin + Go, sign, install to ~/bin, restart relay
deploy: build-plugin sync-plugin
	go build -o /tmp/$(FIGMA_BINARY) ./cmd/ahd-figma/
	codesign -f -s - /tmp/$(FIGMA_BINARY)
	cp /tmp/$(FIGMA_BINARY) ~/bin/$(FIGMA_BINARY)
	@echo "Binary installed to ~/bin/$(FIGMA_BINARY)"
	@# Restart relay if running
	@-pkill -f "ahd-figma ws" 2>/dev/null && echo "Stopped old relay" || true
	@sleep 1
	@nohup ~/bin/$(FIGMA_BINARY) ws > /tmp/ahd-relay.log 2>&1 & echo "Relay restarted (PID: $$!)"
	@echo "Deploy complete. Reopen Figma plugin to pick up new code."

lint:
	golangci-lint run
