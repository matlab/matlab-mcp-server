# Copyright 2025-2026 The MathWorks, Inc.

# Set shell based on OS
ifeq ($(OS),Windows_NT)
	SHELL = powershell.exe
else
	SHELL = sh
endif

# Race detector flag
# Note: Disabled on Windows because CI agents don't have gcc available (required for -race)
ifeq ($(OS),Windows_NT)
    RACE_FLAG =
else
    RACE_FLAG = -race
endif

ifeq ($(OS),Windows_NT)
    RM_DIR = if (Test-Path "$(1)") { Remove-Item -Recurse -Force "$(1)" }
	PATHSEP = ;
	BIN_PATH = $(CURDIR)/.bin/win64
else
    RM_DIR = rm -rf $(1)
	PATHSEP = :
	BIN_PATH = $(CURDIR)/.bin/glnxa64
endif

# Capture CLI Environment variables
CLI_MATLAB_MCP_CORE_SERVER_BUILD_DIR := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)
CLI_MCP_MATLAB_PATH := $(MCP_MATLAB_PATH)

# Include .env file if it exists
ifneq (,$(wildcard .env))
    include .env
endif

# Set MATLAB_MCP_CORE_SERVER_BUILD_DIR with precendence CLI > .env > default
ifdef CLI_MATLAB_MCP_CORE_SERVER_BUILD_DIR
	MATLAB_MCP_CORE_SERVER_BUILD_DIR = $(CLI_MATLAB_MCP_CORE_SERVER_BUILD_DIR)
endif
ifndef MATLAB_MCP_CORE_SERVER_BUILD_DIR
	MATLAB_MCP_CORE_SERVER_BUILD_DIR = $(CURDIR)/.bin
endif
export MATLAB_MCP_CORE_SERVER_BUILD_DIR

# Set MCP_MATLAB_PATH with precendence CLI > .env > default (empty)
ifdef CLI_MCP_MATLAB_PATH
	MCP_MATLAB_PATH = $(CLI_MCP_MATLAB_PATH)
endif
export MCP_MATLAB_PATH

# Variables for MCP Inspector
export HOST = localhost
export PATH := $(BIN_PATH)$(PATHSEP)$(PATH)

# Go build flags
BUILD_FLAGS := -trimpath

# Strip symbol table and debug info for release builds only
ifeq ($(RELEASE),true)
	LDFLAGS_ARG := -ldflags "-s -w"
else
	LDFLAGS_ARG :=
endif

all: wire mockery lint unit-tests integration-tests build functional-tests mcpb-clean mcpb-dev

mcp-inspector: build
	npx @modelcontextprotocol/inspector matlab-mcp-core-server

# File checks

wire:
	go tool wire github.com/matlab/matlab-mcp-core-server/internal/wire

install:
	@echo "No longer needed"

mockery:
	@$(call RM_DIR,./mocks)
	@$(call RM_DIR,./tests/mocks)
	go tool mockery

lint:
	go tool golangci-lint run ./...

fix-lint:
	go tool golangci-lint run ./... --fix

matlab-lint:
ifeq ($(OS),Windows_NT)
	matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolkit')); buildtool clean lint;"
else
	matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolkit')); buildtool clean lint;"
endif

# Resources

CODING_GUIDELINES_URL := https://raw.githubusercontent.com/matlab/rules/main/matlab-coding-standards.md
CODING_GUIDELINES_PATH := $(CURDIR)/internal/adaptors/mcp/resources/codingguidelines/assets/codingguidelines.md

update-coding-guidelines:
ifeq ($(OS),Windows_NT)
	Invoke-WebRequest -Uri "$(CODING_GUIDELINES_URL)" -OutFile "$(CODING_GUIDELINES_PATH)"
else
	curl -sSL "$(CODING_GUIDELINES_URL)" -o "$(CODING_GUIDELINES_PATH)"
endif

LIVE_CODE_GUIDELINES_URL := https://raw.githubusercontent.com/matlab/rules/main/live-script-generation.md
LIVE_CODE_GUIDELINES_PATH := $(CURDIR)/internal/adaptors/mcp/resources/plaintextlivecodegeneration/assets/plaintextlivecodegeneration.md

update-live-code-guidelines:
ifeq ($(OS),Windows_NT)
	Invoke-WebRequest -Uri "$(LIVE_CODE_GUIDELINES_URL)" -OutFile "$(LIVE_CODE_GUIDELINES_PATH)"
else
	curl -sSL "$(LIVE_CODE_GUIDELINES_URL)" -o "$(LIVE_CODE_GUIDELINES_PATH)"
endif

# Building

WIN64_BIN_DIR :=$(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/win64
GLNXA64_BIN_DIR :=$(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/glnxa64
MACI64_BIN_DIR :=$(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/maci64
MACA64_BIN_DIR :=$(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/maca64
ALL_BIN_DIR := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/all
MLTBX_DIR :=$(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/mltbx

build: build-for-windows build-for-glnxa64 build-for-maci64 build-for-maca64
ifeq ($(OS),Windows_NT)
	@New-Item -ItemType Directory -Force -Path "$(ALL_BIN_DIR)" | Out-Null
	@Copy-Item "$(GLNXA64_BIN_DIR)/matlab-mcp-core-server" "$(ALL_BIN_DIR)/matlab-mcp-core-server-glnxa64"
	@Copy-Item "$(MACA64_BIN_DIR)/matlab-mcp-core-server" "$(ALL_BIN_DIR)/matlab-mcp-core-server-maca64"
	@Copy-Item "$(MACI64_BIN_DIR)/matlab-mcp-core-server" "$(ALL_BIN_DIR)/matlab-mcp-core-server-maci64"
	@Copy-Item "$(WIN64_BIN_DIR)/matlab-mcp-core-server.exe" "$(ALL_BIN_DIR)/matlab-mcp-core-server-win64.exe"
else
	@mkdir -p "$(ALL_BIN_DIR)"
	@cp "$(GLNXA64_BIN_DIR)/matlab-mcp-core-server" "$(ALL_BIN_DIR)/matlab-mcp-core-server-glnxa64"
	@cp "$(MACA64_BIN_DIR)/matlab-mcp-core-server" "$(ALL_BIN_DIR)/matlab-mcp-core-server-maca64"
	@cp "$(MACI64_BIN_DIR)/matlab-mcp-core-server" "$(ALL_BIN_DIR)/matlab-mcp-core-server-maci64"
	@cp "$(WIN64_BIN_DIR)/matlab-mcp-core-server.exe" "$(ALL_BIN_DIR)/matlab-mcp-core-server-win64.exe"
endif

build-for-windows:
ifeq ($(OS),Windows_NT)
	$$env:GOOS='windows'; $$env:GOARCH='amd64'; $$env:CGO_ENABLED='0'; go build $(BUILD_FLAGS) $(LDFLAGS_ARG) -o $(WIN64_BIN_DIR)/matlab-mcp-core-server.exe ./cmd/matlab-mcp-core-server
else
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS_ARG) -o "$(WIN64_BIN_DIR)/matlab-mcp-core-server.exe" ./cmd/matlab-mcp-core-server
endif

build-for-glnxa64:
ifeq ($(OS),Windows_NT)
	$$env:GOOS='linux'; $$env:GOARCH='amd64'; $$env:CGO_ENABLED='0'; go build $(BUILD_FLAGS) $(LDFLAGS_ARG) -o $(GLNXA64_BIN_DIR)/matlab-mcp-core-server ./cmd/matlab-mcp-core-server
else
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS_ARG) -o "$(GLNXA64_BIN_DIR)/matlab-mcp-core-server" ./cmd/matlab-mcp-core-server
endif

build-for-maci64:
ifeq ($(OS),Windows_NT)
	$$env:GOOS='darwin'; $$env:GOARCH='amd64'; $$env:CGO_ENABLED='0'; go build $(BUILD_FLAGS) $(LDFLAGS_ARG) -o $(MACI64_BIN_DIR)/matlab-mcp-core-server ./cmd/matlab-mcp-core-server
else
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS_ARG) -o "$(MACI64_BIN_DIR)/matlab-mcp-core-server" ./cmd/matlab-mcp-core-server
endif

build-for-maca64:
ifeq ($(OS),Windows_NT)
	$$env:GOOS='darwin'; $$env:GOARCH='arm64'; $$env:CGO_ENABLED='0'; go build $(BUILD_FLAGS) $(LDFLAGS_ARG) -o $(MACA64_BIN_DIR)/matlab-mcp-core-server ./cmd/matlab-mcp-core-server
else
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS_ARG) -o "$(MACA64_BIN_DIR)/matlab-mcp-core-server" ./cmd/matlab-mcp-core-server
endif

build-matlab-addon:
ifeq ($(OS),Windows_NT)
	$$env:MATLAB_MCP_CORE_SERVER_MLTBX_DIR='$(MLTBX_DIR)'; matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolkit')); buildtool clean package;"
else
	MATLAB_MCP_CORE_SERVER_MLTBX_DIR="$(MLTBX_DIR)" matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolkit')); buildtool clean package;"
endif

# Testing

unit-tests:
	go tool gotestsum --packages="./internal/... ./pkg/... ./tests/testutils/..." -- -race -coverprofile cover.out

matlab-unit-tests:
ifeq ($(OS),Windows_NT)
	matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolkit')); buildtool clean unit-tests;"
else
	matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolkit')); buildtool clean unit-tests;"
endif

integration-tests:
	go tool gotestsum --packages="./tests/integration/..." -- -race
	
functional-tests:
	go tool gotestsum --packages="./tests/functional/..." -- -race

system-tests:
	go tool gotestsum --packages="./tests/system/..." -- -race -count=1 -timeout 30m
	@$(CHECK_MATLAB_LEAKS)

matlab-system-tests: build-matlab-addon
ifeq ($(OS),Windows_NT)
	$$env:MATLAB_MCP_CORE_SERVER_MLTBX_PATH='$(MLTBX_DIR)/MATLABMCPCoreServerToolkit.mltbx'; matlab -batch "cd(fullfile('$(CURDIR)', 'tests', 'system', 'matlab')); buildtool;"
else
	MATLAB_MCP_CORE_SERVER_MLTBX_PATH="$(MLTBX_DIR)/MATLABMCPCoreServerToolkit.mltbx" matlab -batch "cd(fullfile('$(CURDIR)', 'tests', 'system', 'matlab')); buildtool;"
endif

ci-unit-tests:
	go test $(RACE_FLAG) -json -count=1 -coverprofile cover.out ./internal/... ./pkg/... ./tests/testutils/...

ci-matlab-unit-tests: matlab-unit-tests

ci-integration-tests:
	go test $(RACE_FLAG) -json -count=1 ./tests/integration/...
	
ci-functional-tests:
	go test $(RACE_FLAG) -json -count=1 ./tests/functional/...

ci-system-tests:
	go test $(RACE_FLAG) -timeout 120m -json -count=1 ./tests/system/...
	@$(CHECK_MATLAB_LEAKS)

# Requires MATLAB_MCP_CORE_SERVER_MLTBX_PATH to be set
ci-matlab-system-tests:
	matlab -batch "cd(fullfile('$(CURDIR)', 'tests', 'system', 'matlab')); buildtool;"

# Check for leaked MATLAB processes after system tests
# Tests should clean up all MATLAB sessions they create
check-matlab-leaks:
	@$(CHECK_MATLAB_LEAKS)

# =============================================================================
# Platform-specific multi-line command definitions
# =============================================================================
# These use define/endef for readability and $(strip ...) to flatten for execution

ifeq ($(OS),Windows_NT)

define CHECK_MATLAB_LEAKS_CMD
powershell -NoProfile -ExecutionPolicy Bypass -Command "& {
    Write-Host 'Waiting for processes to settle...';
    Start-Sleep -Seconds 5;
    Write-Host 'Checking for leaked MATLAB processes...';
    `$$p = Get-Process -Name MATLAB -ErrorAction SilentlyContinue |
        Where-Object { `$$_.CommandLine -like '*matlab-mcp-core-server*' };
    if (`$$p) {
        Write-Host 'WARNING: Found leaked MATLAB processes:';
        `$$p | Format-Table Id,ProcessName,StartTime;
        exit 1
    } else {
        Write-Host 'No leaked MATLAB processes found.'
    }
}"
endef

else

define CHECK_MATLAB_LEAKS_CMD
echo "Waiting for processes to settle...";
sleep 5;
echo "Checking for leaked MATLAB processes...";
leaked=$$(pgrep -a -f -l 'addpath\(sessionPath\);matlab_mcp\.initializeMCP\(\);clear sessionPath;' | grep -v 'make\|grep' || true);
if [ -n "$$leaked" ]; then
    echo "WARNING: Found leaked MATLAB processes:";
    echo "$$leaked";
    exit 1;
else
    echo "No leaked MATLAB processes found.";
fi
endef

endif

CHECK_MATLAB_LEAKS := $(strip $(CHECK_MATLAB_LEAKS_CMD))

# MCPB Bundle Configuration
MCPB_STAGING_DIR := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/mcpb
MCPB_FILENAME := matlab-mcp-core-server.mcpb
MCPB_GEN_BIN := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/mcpb-gen/mcpb-gen

# Generate mcpb sources for bundling
# Build mcpb-gen first (not go run) to get proper version from debug.ReadBuildInfo()
mcpb-stage:
ifeq ($(OS),Windows_NT)
	@echo "Error: MCPB manifest generation is only supported on macOS/Linux"; exit 1
else
	@mkdir -p "$(dir $(MCPB_GEN_BIN))"
	go build -o "$(MCPB_GEN_BIN)" ./cmd/mcpb-gen
	MCPB_STAGING_DIR="$(MCPB_STAGING_DIR)" "$(MCPB_GEN_BIN)"
endif

# Main mcpb target: generates manifest and packs bundle
# Requires all 4 platform binaries in $(ALL_BIN_DIR).
# For local dev: make mcpb-dev (builds + copies to all/ + packs)
# For CI/signed: populate all/ externally, then make mcpb
mcpb: mcpb-stage
ifeq ($(OS),Windows_NT)
	@echo "Error: MCPB packaging is only supported on macOS/Linux"; exit 1
else
	@if [ ! -f "$(ALL_BIN_DIR)/matlab-mcp-core-server-glnxa64" ] || \
	    [ ! -f "$(ALL_BIN_DIR)/matlab-mcp-core-server-maca64" ] || \
	    [ ! -f "$(ALL_BIN_DIR)/matlab-mcp-core-server-maci64" ] || \
	    [ ! -f "$(ALL_BIN_DIR)/matlab-mcp-core-server-win64.exe" ]; then \
		echo "Error: Missing binaries in $(ALL_BIN_DIR)."; \
		echo "Run 'make mcpb-dev' for local builds, or populate $(ALL_BIN_DIR) with signed binaries."; \
		exit 1; \
	fi
	@echo "Using binaries from $(ALL_BIN_DIR)"
	@cp "$(ALL_BIN_DIR)"/matlab-mcp-core-server-* "$(MCPB_STAGING_DIR)/bundle/bin/"
	@cd "$(MCPB_STAGING_DIR)" && npm i && npm run mcpb-pack -- "$(MCPB_FILENAME)"
	@echo ""
	@echo "Created: $(MCPB_STAGING_DIR)/$(MCPB_FILENAME)"
endif

mcpb-clean:
	@$(call RM_DIR,$(MCPB_STAGING_DIR))
	@$(call RM_DIR,$(dir $(MCPB_GEN_BIN)))
	@echo "Removed $(MCPB_STAGING_DIR) and $(dir $(MCPB_GEN_BIN))"

# Development workflow: build, copy to all/, and pack
mcpb-dev: mcpb-clean build mcpb

mcpb-validate:
ifeq ($(OS),Windows_NT)
	@echo "Error: MCPB validation is only supported on macOS/Linux"; exit 1
else
	cd "$(MCPB_STAGING_DIR)"; \
	npm run mcpb-validate
endif
