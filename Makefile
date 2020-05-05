PKG_NAME = github.com/nasa9084/gac
CMD_NAME = gac

GO ?= go
BIN_DIR = $(PWD)/bin
VERSION ?= $(shell git describe --tags --abbrev=0)
REVISION ?= $(shell git rev-parse --short HEAD)
LDFLAGS = -ldflags "-X $(PKG_NAME)/commands.Version=$(VERSION) -X $(PKG_NAME)/commands.Revision=$(REVISION)"

.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	@echo Version: $(VERSION), Revion: $(REVISION)
	@$(GO) build $(LDFLAGS) -o $(BIN_DIR)/$(CMD_NAME) main.go

.PHONY: test
test:
	@$(GO) test -v ./...

.PHONY: clean
clean:
	@-rm -fr $(BIN_DIR)
