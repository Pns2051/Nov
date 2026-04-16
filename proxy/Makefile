.PHONY: build-linux build-macos build-windows build-all clean

BINARY_NAME=adblock-proxy
DIST_DIR=../dist
VERSION ?= 1.0.0
LDFLAGS=-X github.com/Pns2051/Nov/proxy/internal/config.Version=$(VERSION)

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/adblock-proxy

build-macos:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/adblock-proxy

build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/adblock-proxy

build-all: build-linux build-macos build-windows

clean:
	rm -rf $(DIST_DIR)
