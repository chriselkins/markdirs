# Makefile for markdirs

# Name and version info
BINARY      := markdirs
VERSION     ?= $(shell git describe --tags --always --dirty)
COMMIT      := $(shell git rev-parse --short HEAD)
DATE        := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Target platforms (add or remove as needed)
PLATFORMS   := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Output directory for binaries
DIST        := dist

# ldflags for embedding version info
LDFLAGS := -s -w -X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildDate=$(DATE)'

.PHONY: all clean release

all: $(PLATFORMS:%=$(DIST)/$(BINARY)-%)

# Cross-compile for each platform
$(DIST)/$(BINARY)-%:
	@platform=$* ; \
	GOOS=$${platform%/*} GOARCH=$${platform#*/} \
	CGO_ENABLED=0 \
	go build -trimpath -ldflags="$(LDFLAGS)" \
	-o $(DIST)/$(BINARY)-$${platform%/*}-$${platform#*/}$(if $(findstring windows,$${platform%/*}),.exe,) \
	.

# Release target: build and zip artifacts
release: clean all
	@for platform in $(PLATFORMS); do \
		outfile=$(DIST)/$(BINARY)-$${platform%/*}-$${platform#*/} ; \
		[ "$${platform%/*}" = "windows" ] && outfile=$${outfile}.exe ; \
		zip -j "$${outfile}.zip" "$${outfile}" ; \
	done

clean:
	rm -rf $(DIST)

# Print version info
version:
	@echo "Version:   $(VERSION)"
	@echo "Commit:    $(COMMIT)"
	@echo "BuildDate: $(DATE)"
