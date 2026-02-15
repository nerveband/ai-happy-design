.PHONY: build build-go build-plugin sync-plugin run-mcp run-ws test clean install lint

BINARY=ai-happy-design
VERSION=0.0.0-dev

build-plugin:
	cd plugin && npm ci && npm run check

sync-plugin:
	mkdir -p internal/plugin/files/dist
	cp plugin/manifest.json internal/plugin/files/manifest.json
	cp plugin/dist/code.js internal/plugin/files/dist/code.js
	cp plugin/dist/ui.html internal/plugin/files/dist/ui.html

build: build-plugin sync-plugin
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

# Full deploy: build plugin + Go, sign, install to ~/bin, restart relay
deploy: build-plugin sync-plugin
	go build -o /tmp/$(BINARY) ./cmd/ai-happy-design/
	codesign -f -s - /tmp/$(BINARY)
	cp /tmp/$(BINARY) ~/bin/$(BINARY)
	@echo "Binary installed to ~/bin/$(BINARY)"
	@# Restart relay if running
	@-pkill -f "ai-happy-design ws" 2>/dev/null && echo "Stopped old relay" || true
	@sleep 1
	@nohup ~/bin/$(BINARY) ws > /tmp/ahd-relay.log 2>&1 & echo "Relay restarted (PID: $$!)"
	@echo "Deploy complete. Reopen Figma plugin to pick up new code."

lint:
	golangci-lint run
