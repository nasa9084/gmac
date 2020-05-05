PKG_NAME = github.com/nasa9084/gac
CMD_NAME = gac

GO ?= go
BIN_DIR = $(PWD)/bin
VERSION ?= $(shell git describe --tags --abbrev=0)
REVISION ?= $(shell git rev-parse --short HEAD)
LDFLAGS = -ldflags "-X $(PKG_NAME)/commands.Version=$(VERSION) -X $(PKG_NAME)/commands.Revision=$(REVISION)"
SRCS = $(foreach dir,$(shell find . -type d),$(wildcard $(dir)/*.go))

all: clean test build

.PHONY: build
build: $(BIN_DIR)/$(CMD_NAME)

.PHONY: test
test:
	@$(GO) test -v ./...

.PHONY: clean
clean:
	@-rm -fr $(BIN_DIR)

$(BIN_DIR)/$(CMD_NAME): $(SRCS)
	@mkdir -p $(BIN_DIR)
	@$(GO) build $(LDFLAGS) -o $(BIN_DIR)/$(CMD_NAME)
