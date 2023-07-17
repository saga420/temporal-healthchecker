# Constants
GREEN := $(shell tput setaf 2)
NORMAL := $(shell tput sgr0)

PKG := $(shell go list -m)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2> /dev/null || echo "unknown")
GIT_COMMIT_TIME := $(shell git show -s --format=%ct "$(GIT_COMMIT)" 2> /dev/null || echo "unknown")
APPNAME := $(shell basename "$(shell git rev-parse --show-toplevel)")
OS_ARCHS := darwin-amd64 darwin-arm64 linux-amd64 linux-arm64
TARGETS := $(OS_ARCHS)

# Flags
LDFLAGS := -s -w -X $(PKG)/version.GitRevision=$(GIT_COMMIT) -X $(PKG)/version.GitCommitAt=$(GIT_COMMIT_TIME)

# Targets
.PHONY: all
all: build

.PHONY: build
build: $(TARGETS)

$(TARGETS):
	$(eval GOOS := $(word 1,$(subst -, ,$@)))
	$(eval GOARCH := $(word 2,$(subst -, ,$@)))
	@echo "$(GREEN)Building for $(GOOS)/$(GOARCH)...$(NORMAL)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o "build/bin/$(APPNAME)_$(GOOS)_$(GOARCH)"

.PHONY: clean
clean:
	@echo "$(GREEN)Cleaning up...$(NORMAL)"
	rm -rf build/
