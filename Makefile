

export PATH := ./bin:$(PATH)
export GO111MODULE := on
export BIN_NAME := ./bin/prometheus-exporter-scratch

# Initial development
.PHONY: init
init:
	go mod init github.com/mtulio/prometheus-exporter-scratch

# Install all the build and lint dependencies
.PHONY: setup
setup:
	go mod download

.PHONY: update
update:
	go mod tidy

# Build a beta version
.PHONY: build
build:
	@mkdir ./bin || true
	go build -o $(BIN_NAME) && strip $(BIN_NAME)

.PHONY: run
run:
	$(BIN_NAME)

.PHONY: version
version: build
	$(BIN_NAME) -v

.PHONY: clean
clean:
	@rm -f bin/$(BIN_NAME)

