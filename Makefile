# Variables
GOLANGCI_LINT := $(shell which golangci-lint)

# Modules in this repo: the core, plus opt-in driver plugins under x/.
MODULES := . x/testify

# Default target
all: tidy

# Tidy: format and vet the code (across all modules)
tidy:
	@for m in $(MODULES); do \
		echo "== tidy $$m =="; \
		(cd $$m && go fmt ./... && go vet ./... && go mod tidy) || exit 1; \
	done

# Install golangci-lint only if it's not already installed
lint-install:
	@if ! [ -x "$(GOLANGCI_LINT)" ]; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi

# Lint the code using golangci-lint
# todo reuse var if possible
lint: lint-install
	$(shell which golangci-lint) run

test:
	@for m in $(MODULES); do \
		echo "== test $$m =="; \
		(cd $$m && go test ./...) || exit 1; \
	done

# Phony targets
.PHONY: all tidy lint-install lint test
