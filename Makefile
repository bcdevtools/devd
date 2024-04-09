VERSION := $(shell echo $(shell git describe --tags || git branch --show-current) | sed 's/^v//')
IS_SUDO_USER := $(shell if [ "$(shell whoami)" = "root" ] || [ "$(shell groups | grep -e 'sudo' -e 'admin' -e 'google-sudoers' | wc -l)" = "1" ]; then echo "1"; fi)
GO_BIN := $(shell echo $(shell which go || echo "/usr/local/go/bin/go" ))

###############################################################################
###                                Build flags                              ###
###############################################################################

LD_FLAGS = -X github.com/bcdevtools/devd/constants.VERSION=$(VERSION)

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

###############################################################################
###                                  Build                                  ###
###############################################################################

build: go.sum
	@echo "building devd binary..."
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
	@echo " [v] Installed in GOPATH/bin"
	@if [ "$(shell uname)" = "Linux" ] && [ "$(IS_SUDO_USER)" = "1" ]; then \
		sudo mv $(shell $(GO_BIN) env GOPATH)/bin/devd /usr/local/bin/devd; \
		echo " [v] Installed as global command"; \
	else \
		echo " [x] (Skipped) Install as global command"; \
	fi
.PHONY: install