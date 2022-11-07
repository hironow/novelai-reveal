EXECUTABLE=novelai-reveal
VERSION=$(shell git describe --tags --always --long --dirty)
WINDOWS=$(EXECUTABLE)_windows_amd64_$(VERSION).exe
LINUX=$(EXECUTABLE)_linux_amd64_$(VERSION)
DARWIN=$(EXECUTABLE)_darwin_amd64_$(VERSION)

.PHONY: all test clean

all: test build ## Build and run tests

test: ## Run unit tests
	go test ./...

build: windows linux darwin ## Build binaries
	@echo version: $(VERSION)

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -v -o bin/$(WINDOWS) -ldflags="-s -w -X main.version=$(VERSION)"  .

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -v -o bin/$(LINUX) -ldflags="-s -w -X main.version=$(VERSION)"  .

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -v -o bin/$(DARWIN) -ldflags="-s -w -X main.version=$(VERSION)"  .

clean: ## Remove previous build
	rm -f $(WINDOWS) $(LINUX) $(DARWIN)

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
