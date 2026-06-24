# Variables
GOLANGCI_LINT := $(shell which golangci-lint)

# Modules in this repo: the core, plus opt-in plugins under x/.
MODULES := . x/mock

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

# Release both modules at VERSION (e.g. make release VERSION=v1.0.0-rc.6).
#
# A submodule's `require` can only name a tag that already exists, so multi-module
# releases are inherently a two-step dance. This target automates it: tag the core
# at HEAD, point x/mock's core requirement at that tag (via `go mod edit`, a plain
# text rewrite that needs no network and so works before the tag is pushed), commit
# that bump only if it actually changed, then tag x/mock. When the require was
# already bumped in the feature commit, both modules tag the same commit and no
# extra commit is made. Tags are created locally only — the final push is printed
# for you to run after review, since published module versions are immutable.
release:
	@test -n "$(VERSION)" || { echo "VERSION is required, e.g. make release VERSION=v1.0.0-rc.6"; exit 1; }
	@test -z "$$(git status --porcelain)" || { echo "working tree not clean — commit or stash first"; exit 1; }
	@$(MAKE) test
	@echo "== tag core $(VERSION) =="
	git tag $(VERSION)
	@echo "== point x/mock at core $(VERSION) =="
	cd x/mock && go mod edit -require=github.com/expectto/be@$(VERSION)
	@if git diff --quiet -- x/mock/go.mod; then \
		echo "   (x/mock already requires $(VERSION) — no bump commit needed)"; \
	else \
		git commit -m "x/mock: bump core require to $(VERSION)" x/mock/go.mod; \
	fi
	git tag x/mock/$(VERSION)
	@echo ""
	@echo "Tagged locally. Review, then publish both with:"
	@echo "  git push origin main $(VERSION) x/mock/$(VERSION)"

# Phony targets
.PHONY: all tidy lint-install lint test release
