.PHONY: build build-go build-plugin sync-plugin run-mcp run-ws test clean install lint

BINARY=ai-happy-design
ALIAS_BINARY=ahd-figma
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
	cp bin/$(BINARY) bin/$(ALIAS_BINARY)

build-go:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(BINARY) ./cmd/ai-happy-design
	cp bin/$(BINARY) bin/$(ALIAS_BINARY)

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
	cp bin/$(ALIAS_BINARY) $(GOPATH)/bin/

# Full deploy: build plugin + Go, sign, install to ~/bin, restart relay
deploy: build-plugin sync-plugin
	go build -o /tmp/$(BINARY) ./cmd/ai-happy-design/
	codesign -f -s - /tmp/$(BINARY)
	rm -f ~/bin/$(BINARY) ~/bin/$(ALIAS_BINARY)
	cp /tmp/$(BINARY) ~/bin/$(BINARY)
	cp /tmp/$(BINARY) ~/bin/$(ALIAS_BINARY)
	@echo "Binary installed to ~/bin/$(BINARY) and ~/bin/$(ALIAS_BINARY)"
	@# Restart relay through the managed lifecycle so state stays accurate.
	@-~/bin/$(ALIAS_BINARY) relay stop >/dev/null 2>&1 || true
	@~/bin/$(ALIAS_BINARY) relay start >/dev/null
	@~/bin/$(ALIAS_BINARY) relay status >/dev/null
	@echo "Relay restarted"
	@echo "Deploy complete. Reopen Figma plugin to pick up new code."

lint:
	golangci-lint run
