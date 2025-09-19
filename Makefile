# Makefile

APP      := matter
PKG      := ./cmd
BINDIR   := /usr/local/bin
BUILDDIR := build

GO       ?= go
GIT_TAG  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TS := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# 注入 main.version & main.buildTime，並加上 -s -w 與 PIE
LDFLAGS  := -s -w -X main.Version=$(GIT_TAG) -X main.BuildTime=$(BUILD_TS)

.PHONY: all build install clean run

all: build

build:
	@mkdir -p $(BUILDDIR)
	$(GO) build \
		-buildmode=pie \
		-ldflags "$(LDFLAGS)" \
		-o $(BUILDDIR)/$(APP) $(PKG)

install: build
	install -m 0755 $(BUILDDIR)/$(APP) $(BINDIR)/$(APP)

clean:
	rm -rf $(BUILDDIR)

run: build
	./$(BUILDDIR)/$(APP)
