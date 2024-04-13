VERSION := $(shell echo $(shell git describe --tags || git branch --show-current) | sed 's/^v//')
GO_BIN := $(shell echo $(shell which go || echo "/usr/local/go/bin/go" ))

###############################################################################
###                                Build flags                              ###
###############################################################################

LD_FLAGS = -X github.com/bcdevtools/devd/v2/constants.VERSION=$(VERSION)

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

###############################################################################
###                                  Build                                  ###
###############################################################################

build: go.sum
	@echo "Building devd binary..."
	@echo "Flags $(BUILD_FLAGS)"
	@go build -mod=readonly $(BUILD_FLAGS) -o build/devd ./cmd/devd
	@echo "Builded successfully"
.PHONY: build

###############################################################################
###                                 Install                                 ###
###############################################################################

install: go.sum
	@echo "Build flags: $(BUILD_FLAGS)"
	@echo "Installing devd binary..."
	@$(GO_BIN) install -mod=readonly $(BUILD_FLAGS) ./cmd/devd
	@echo "Installed successfully"
.PHONY: install