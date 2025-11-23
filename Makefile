# HBM

.PHONY: all build clean test vendor format lint shellcheck dockerlint help

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

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

# Lint shell scripts
shellcheck:
	@echo "Linting shell scripts..."
	@for file in $$(find . -type f -name '*.sh' -not -path "./.git/*" -not -path "./vendor/*"); do \
		echo "Checking: $$file"; \
		docker run --rm -v "$$PWD:/mnt:ro" koalaman/shellcheck -e SC2086 -e SC2046 -e SC1090 "$$file" || true; \
	done

# Lint Dockerfiles
dockerlint:
	@echo "Linting Dockerfiles..."
	@for file in $$(find . -name 'Dockerfile*' -not -name '*.dapper'); do \
		echo "Checking: $$file"; \
		docker run -i --rm hadolint/hadolint hadolint --ignore DL3018 --ignore DL3013 - < "$$file" || true; \
	done

# Show help
help:
	@echo "HBM"
	@echo ""
	@echo "Targets:"
	@echo "  make build       - Build HBM binary with version info"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make release     - Create GitHub release"
	@echo "  make test        - Run Go unit tests"
	@echo "  make format      - Format Go code with gofmt"
	@echo "  make lint        - Run golint"
	@echo "  make shellcheck  - Lint shell scripts"
	@echo "  make dockerlint  - Lint Dockerfiles"
	@echo "  make help        - Show this help"

.DEFAULT_GOAL := build
