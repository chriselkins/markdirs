# Makefile for markdirs

# Name and version info
BINARY      := markdirs
VERSION     ?= $(shell git describe --tags --always --dirty)
COMMIT      := $(shell git rev-parse --short HEAD)
DATE        := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Target platforms
PLATFORMS := linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64

# Output directory for binaries
DIST        := dist

# ldflags for embedding version info
LDFLAGS := -s -w -buildid= -X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildDate=$(DATE)'

.PHONY: all clean release version

all: $(PLATFORMS:%=$(DIST)/$(BINARY)-%)

# Cross-compile for each platform
$(DIST)/$(BINARY)-%:
	@platform="$*"; \
	os=$${platform%-*}; arch=$${platform#*-}; \
	outfile="$(DIST)/$(BINARY)-$$os-$$arch"; \
	[ "$$os" = "windows" ] && outfile="$$outfile.exe"; \
	mkdir -p $(DIST); \
	GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 \
	go build -trimpath -ldflags="$(LDFLAGS)" -o "$$outfile" .

# Release target: build and zip artifacts
release: clean all
	@for platform in $(PLATFORMS); do \
		os=$${platform%-*}; arch=$${platform#*-}; \
		basename="$(DIST)/$(BINARY)-$$os-$$arch"; \
		outfile="$$basename"; \
		[ "$$os" = "windows" ] && outfile="$$basename.exe"; \
		zipfile="$$basename.zip"; \
		zip -j "$$zipfile" "$$outfile"; \
	done
	@cd $(DIST) && (command -v sha256sum >/dev/null 2>&1 && sha256sum * > SHA256SUMS || shasum -a 256 * > SHA256SUMS)

clean:
	rm -rf $(DIST)

# Print version info
version:
	@echo "Version:   $(VERSION)"
	@echo "Commit:    $(COMMIT)"
	@echo "BuildDate: $(DATE)"
