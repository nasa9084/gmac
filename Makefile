GO ?= go

.PHONY: test
test:
	@$(GO) test -v ./...
