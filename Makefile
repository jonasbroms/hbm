# HBM

.PHONY: all build clean test vendor format lint shellcheck dockerlint doclint tag help

DOCKER := $(shell command -v docker 2>/dev/null || command -v podman 2>/dev/null)

# Default target
all: build

# Build the HBM binary with version info
build:
	@./scripts/build-target

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/ dist/
	rm -f hbm hbm-test
	find . -name '*.test' -delete

# Run tests
test:
	@echo "Running Go tests... Don't have any of these..."
	go test ./...

# Vendor and clean
vendor:
	@echo "Vendoring dependencies"
	go mod vendor
	@echo "Cleaning"
	go mod tidy
	@echo "Verify"
	go mod verify

# Format Go code
format:
	@echo "Formatting Go code..."
	@for file in $$(find . -name '*.go' -type f -not -path "./.git/*" -not -path "./vendor/*"); do \
		echo "Formatting: $$file"; \
		gofmt -l -s -w "$$file"; \
	done

# Build distribution packages for all supported platforms
packages: build
	@./scripts/build-packages

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	golangci-lint run

# Lint shell scripts
shellcheck:
	@echo "Linting shell scripts..."
	@for file in $$(find . -type f -name '*.sh' -not -path "./.git/*" -not -path "./vendor/*"); do \
		echo "Checking: $$file"; \
		$(DOCKER) run --rm -v "$$PWD:/mnt:ro" koalaman/shellcheck -e SC2086 -e SC2046 -e SC1090 "$$file" || true; \
	done

# Lint Dockerfiles
dockerlint:
	@echo "Linting Dockerfiles..."
	@for file in $$(find . -name 'Dockerfile*' -not -name '*.dapper'); do \
		echo "Checking: $$file"; \
		$(DOCKER) run -i --rm hadolint/hadolint hadolint --ignore DL3018 --ignore DL3013 - < "$$file" || true; \
	done

# Lint Markdown documentation
doclint:
	@echo "Linting Markdown docs..."
	$(DOCKER) run --rm -v "$$PWD:/workdir:ro" davidanson/markdownlint-cli2 "docs/**/*.md" "*.md"

# Create an annotated git tag: make tag VERSION=0.22.0
tag:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make tag VERSION=<version>  (e.g. 0.22.0)"; \
		echo "Last tag: $$(git describe --tags --abbrev=0 2>/dev/null || echo none)"; \
		exit 1; \
	fi
	@if git tag -l | grep -qx "$(VERSION)"; then \
		echo "Tag $(VERSION) already exists"; exit 1; \
	fi
	git tag -a "$(VERSION)" -m "$$(git log --pretty=format:'  * %s' $$(git describe --tags --abbrev=0)..HEAD)"
	@echo "Created tag $(VERSION) — push with: git push origin $(VERSION)"

# Show help
help:
	@echo "HBM"
	@echo ""
	@echo "Targets:"
	@echo "  make build       - Build HBM binary with version info"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make test        - Run Go unit tests"
	@echo "  make format      - Format Go code with gofmt"
	@echo "  make packages    - Build RPM/DEB packages for all platforms"
	@echo "  make lint        - Run golint"
	@echo "  make shellcheck  - Lint shell scripts"
	@echo "  make dockerlint  - Lint Dockerfiles"
	@echo "  make doclint     - Lint Markdown documentation"
	@echo "  make tag         - Create annotated git tag (VERSION=x.y.z required)"
	@echo "  make help        - Show this help"

.DEFAULT_GOAL := build
