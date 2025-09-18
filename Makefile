# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOFMT=gofmt
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod
PACKAGES=./...

# Benchmarking parameters
BENCH_TIME=10s
BENCH_COUNT=5
BENCH_OUTPUT=~/Documents/sourcemap-go-bench-$(shell date +%Y%m%d-%H%M%S).txt

# Color output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

.PHONY: all
all: fmt lint test build

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  ${GREEN}%-20s${NC} %s\n", $$1, $$2 } /^##@/ { printf "\n${YELLOW}%s${NC}\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Format code using gofmt
	@echo "${GREEN}Formatting code...${NC}"
	@$(GOFMT) -s -w .
	@echo "${GREEN}✓ Code formatted${NC}"

.PHONY: fmt-check
fmt-check: ## Check if code is formatted
	@echo "${GREEN}Checking code formatting...${NC}"
	@if [ -n "$$($(GOFMT) -s -l .)" ]; then \
		echo "${RED}✗ Code is not formatted. Run 'make fmt'${NC}"; \
		$(GOFMT) -s -d .; \
		exit 1; \
	else \
		echo "${GREEN}✓ Code is properly formatted${NC}"; \
	fi

.PHONY: lint
lint: fmt-check vet ## Run all linters

.PHONY: vet
vet: ## Run go vet
	@echo "${GREEN}Running go vet...${NC}"
	@$(GOVET) $(PACKAGES)
	@echo "${GREEN}✓ Go vet passed${NC}"

.PHONY: staticcheck
staticcheck: ## Run staticcheck (install with: go install honnef.co/go/tools/cmd/staticcheck@latest)
	@echo "${GREEN}Running staticcheck...${NC}"
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck $(PACKAGES); \
		echo "${GREEN}✓ Staticcheck passed${NC}"; \
	else \
		echo "${YELLOW}⚠ Staticcheck not installed. Install with: go install honnef.co/go/tools/cmd/staticcheck@latest${NC}"; \
	fi

##@ Testing

.PHONY: test
test: ## Run unit tests
	@echo "${GREEN}Running tests...${NC}"
	@$(GOTEST) -v -race -coverprofile=coverage.out $(PACKAGES)
	@echo "${GREEN}✓ Tests passed${NC}"

.PHONY: test-short
test-short: ## Run short tests only
	@echo "${GREEN}Running short tests...${NC}"
	@$(GOTEST) -short $(PACKAGES)
	@echo "${GREEN}✓ Short tests passed${NC}"

.PHONY: coverage
coverage: test ## Run tests with coverage report
	@echo "${GREEN}Generating coverage report...${NC}"
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}✓ Coverage report generated: coverage.html${NC}"

##@ Benchmarking

.PHONY: bench
bench: ## Run all benchmarks
	@echo "${GREEN}Running benchmarks...${NC}"
	@$(GOTEST) -bench=. -benchmem -benchtime=$(BENCH_TIME) $(PACKAGES)

.PHONY: bench-jsonv2
bench-jsonv2: ## Run benchmarks with jsonv2 tag (requires Go 1.25+ with GOEXPERIMENT=jsonv2)
	@echo "${GREEN}Running benchmarks with json/v2...${NC}"
	@echo "${YELLOW}Note: Requires Go 1.25+ with GOEXPERIMENT=jsonv2${NC}"
	@GOEXPERIMENT=jsonv2 $(GOTEST) -tags=jsonv2 -bench=. -benchmem -benchtime=$(BENCH_TIME) $(PACKAGES)

.PHONY: bench-parse
bench-parse: ## Run ParseSourceMap benchmarks only
	@echo "${GREEN}Running ParseSourceMap benchmarks...${NC}"
	@$(GOTEST) -bench=BenchmarkParseSourceMap -benchmem -benchtime=$(BENCH_TIME) github.com/redawl/go-sourcemap/spec

.PHONY: bench-save
bench-save: ## Run benchmarks and save results to file
	@echo "${GREEN}Running benchmarks and saving to $(BENCH_OUTPUT)...${NC}"
	@$(GOTEST) -bench=. -benchmem -count=$(BENCH_COUNT) -timeout=30m github.com/redawl/go-sourcemap/spec > $(BENCH_OUTPUT) 2>&1
	@echo "${GREEN}✓ Benchmark results saved to $(BENCH_OUTPUT)${NC}"

.PHONY: bench-compare
bench-compare: ## Compare two benchmark files (usage: make bench-compare OLD=file1.txt NEW=file2.txt)
	@if [ -z "$(OLD)" ] || [ -z "$(NEW)" ]; then \
		echo "${RED}Usage: make bench-compare OLD=file1.txt NEW=file2.txt${NC}"; \
		exit 1; \
	fi
	@echo "${GREEN}Comparing benchmarks...${NC}"
	@if command -v benchstat >/dev/null 2>&1; then \
		benchstat $(OLD) $(NEW); \
	else \
		echo "${YELLOW}Installing benchstat...${NC}"; \
		go install golang.org/x/perf/cmd/benchstat@latest; \
		benchstat $(OLD) $(NEW); \
	fi

##@ Build

.PHONY: build
build: ## Build the binary
	@echo "${GREEN}Building...${NC}"
	@$(GOBUILD) -v ./...
	@echo "${GREEN}✓ Build successful${NC}"

.PHONY: install
install: ## Install the binary
	@echo "${GREEN}Installing...${NC}"
	@$(GOCMD) install -v ./...
	@echo "${GREEN}✓ Installation successful${NC}"

##@ Dependencies

.PHONY: deps
deps: ## Download dependencies
	@echo "${GREEN}Downloading dependencies...${NC}"
	@$(GOMOD) download
	@echo "${GREEN}✓ Dependencies downloaded${NC}"

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "${GREEN}Updating dependencies...${NC}"
	@$(GOMOD) tidy
	@echo "${GREEN}✓ Dependencies updated${NC}"

##@ Cleanup

.PHONY: clean
clean: ## Clean build artifacts and test cache
	@echo "${GREEN}Cleaning...${NC}"
	@$(GOCMD) clean -testcache
	@rm -f coverage.out coverage.html
	@rm -f spec.test
	@echo "${GREEN}✓ Cleaned${NC}"

##@ CI/CD

.PHONY: ci
ci: deps lint test ## Run CI pipeline (deps, lint, test)
	@echo "${GREEN}✓ CI pipeline completed successfully${NC}"

.PHONY: pre-commit
pre-commit: fmt lint test ## Run pre-commit checks
	@echo "${GREEN}✓ Pre-commit checks passed${NC}"