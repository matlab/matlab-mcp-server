# Copyright 2025-2026 The MathWorks, Inc.
#
# Makefile layout:
#   1. Configuration    - Variables, platform detection, env setup
#   2. CI Targets       - Targets called by CI (prefixed ci-*)
#   3. Code Generation  - wire, mockery, and other generated code
#   4. Linting          - Go and MATLAB linting
#   5. Resources        - Syncing and updating external resources
#   6. Building         - Cross-platform binary compilation and packaging
#   7. Testing          - Local development test targets
#   8. MCPB             - MCP Bundle packaging
#   9. Internal         - Reusable command definitions
#
# Adding a new check to CI:
#   1. Add a target in the appropriate section (e.g. Code Generation)
#   2. Add it as a dependency of the relevant ci-* aggregate target
#   3. No CI configuration changes are needed

# =============================================================================
# 1. Configuration
# =============================================================================

# --- Platform ---

ifeq ($(OS),Windows_NT)
	SHELL = powershell.exe
	RACE_FLAG =
	RM_DIR = if (Test-Path "$(1)") { Remove-Item -Recurse -Force "$(1)" }
	MK_DIR = New-Item -ItemType Directory -Force -Path "$(1)" | Out-Null
	CP = Copy-Item $(1) $(2)
	PATHSEP = ;
	BIN_PATH = $(CURDIR)/.bin/win64
	EXE_SUFFIX = .exe
else
	SHELL = sh
	RACE_FLAG = -race
	RM_DIR = rm -rf $(1)
	MK_DIR = mkdir -p "$(1)"
	CP = cp $(1) $(2)
	PATHSEP = :
	EXE_SUFFIX =
	UNAME_S := $(shell uname -s)
	UNAME_M := $(shell uname -m)
	ifeq ($(UNAME_S),Darwin)
		ifeq ($(UNAME_M),arm64)
			BIN_PATH = $(CURDIR)/.bin/maca64
		else
			BIN_PATH = $(CURDIR)/.bin/maci64
		endif
	else
		BIN_PATH = $(CURDIR)/.bin/glnxa64
	endif
endif

# --- Variable precedence: CLI > .env file > env > default ---

ifneq (,$(wildcard .env))
	include .env
endif

EMBEDDED_MLTBX_DIR := $(CURDIR)/internal/adaptors/matlabmanager/addonmanager/installationsteps/assets/mltbx

MATLAB_MCP_CORE_SERVER_BUILD_DIR ?= $(CURDIR)/.bin
MATLAB_MCP_CORE_SERVER_MLTBX_DIR ?= $(EMBEDDED_MLTBX_DIR)

# --- Build paths ---

TOOLS_BIN_DIR := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/tools
SOURCEHASH_BIN := $(TOOLS_BIN_DIR)/sourcehash$(EXE_SUFFIX)
MCPB_GEN_BIN := $(TOOLS_BIN_DIR)/mcpb-gen$(EXE_SUFFIX)

WIN64_BIN_DIR := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/win64
GLNXA64_BIN_DIR := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/glnxa64
MACI64_BIN_DIR := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/maci64
MACA64_BIN_DIR := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/maca64
ALL_BIN_DIR := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/all

MLTBX_DIR := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/mltbx
SOURCES_HASH_FILE := $(EMBEDDED_MLTBX_DIR)/.sources-hash
MATLAB_TOOLBOX_DIR := $(CURDIR)/matlab/matlab_mcp_toolbox

MCPB_STAGING_DIR := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)/mcpb
MCPB_FILENAME := matlab-mcp-core-server.mcpb

# --- Resource paths ---

MATLAB_MCP_EMBEDDED_SRC := $(CURDIR)/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory/matlabfiles/assets/+matlab_mcp
MATLAB_MCP_TOOLBOX_DST := $(CURDIR)/matlab/matlab_mcp_toolbox/mcp/+matlab_mcp

CODING_GUIDELINES_URL := https://raw.githubusercontent.com/matlab/rules/main/matlab-coding-standards.md
CODING_GUIDELINES_PATH := $(CURDIR)/internal/adaptors/mcp/resources/codingguidelines/assets/codingguidelines.md

LIVE_CODE_GUIDELINES_URL := https://raw.githubusercontent.com/matlab/rules/main/live-script-generation.md
LIVE_CODE_GUIDELINES_PATH := $(CURDIR)/internal/adaptors/mcp/resources/plaintextlivecodegeneration/assets/plaintextlivecodegeneration.md

# --- Test configuration ---

UNIT_TEST_PKGS := ./internal/... ./pkg/... ./tests/testutils/...
INTEGRATION_TEST_PKGS := ./tests/integration/...
FUNCTIONAL_TEST_PKGS := ./tests/functional/...
SYSTEM_TEST_PKGS := ./tests/system/...

# Scoped exports prevent polluting other commands (e.g. npm uses HOST)
TEST_TARGETS := \
	unit-tests  matlab-unit-tests  ci-unit-tests          \
	integration-tests              ci-integration-tests    \
	functional-tests               ci-functional-tests     \
	system-tests                   ci-system-tests         \
	matlab-system-tests            ci-matlab-system-tests
$(TEST_TARGETS): export PATH := $(BIN_PATH)$(PATHSEP)$(PATH)
$(TEST_TARGETS): export MATLAB_MCP_CORE_SERVER_BUILD_DIR := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)
$(TEST_TARGETS): export MCP_MATLAB_PATH := $(MCP_MATLAB_PATH)

# --- Build flags ---

BUILD_FLAGS := -trimpath

ifeq ($(RELEASE),true)
	LDFLAGS_ARG := -ldflags "-s -w"
else
	LDFLAGS_ARG :=
endif

all: wire mockery lint unit-tests integration-tests build functional-tests mcpb-clean mcpb-dev

# =============================================================================
# 2. CI Targets
# =============================================================================
# These are the targets called by CI jobs. Each CI job should call exactly one
# ci-* target. CI targets either orchestrate base targets (e.g. ci-lint calls
# lint + matlab-lint) or provide CI-specific output formats (e.g. go test -json).
#
# To include a new check in CI, add it as a dependency of a ci-* aggregate
# target (e.g. ci-check-generated). To add a new test suite, add a new ci-*
# target here and a corresponding CI job.

ci-lint: lint matlab-lint

ci-check-generated: check-mockery check-wire check-embedded-matlab-addon

ci-unit-tests:
	go test $(RACE_FLAG) -json -count=1 -coverprofile cover.out $(UNIT_TEST_PKGS)
	matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean unit-tests;"

ci-matlab-unit-tests:
	@echo "ci-matlab-unit-tests is no longer needed. MATLAB unit tests are now included in ci-unit-tests."

ci-build-matlab-addon: build-matlab-addon

ci-build: build

ci-build-mcpb: build-mcpb-bundle mcpb-validate

ci-integration-tests:
	go test $(RACE_FLAG) -json -count=1 $(INTEGRATION_TEST_PKGS)

ci-functional-tests: ensure-binary-executable
	go test $(RACE_FLAG) -json -count=1 $(FUNCTIONAL_TEST_PKGS)

ci-system-tests: ensure-binary-executable
	go test $(RACE_FLAG) -timeout 120m -json -count=1 $(SYSTEM_TEST_PKGS)
	@$(CHECK_MATLAB_LEAKS)

ci-matlab-system-tests: matlab-system-tests

# =============================================================================
# 3. Code Generation
# =============================================================================

wire:
	go tool wire github.com/matlab/matlab-mcp-core-server/internal/wire

check-wire: wire
	@$(call CHECK_GIT_CLEAN,Wire generated code)

mockery:
	@$(call RM_DIR,./mocks)
	@$(call RM_DIR,./tests/mocks)
	go tool mockery

check-mockery: mockery
	@$(call CHECK_GIT_CLEAN,Generated mocks)

# =============================================================================
# 4. Linting
# =============================================================================

lint:
	go tool golangci-lint run ./...

fix-lint:
	go tool golangci-lint run ./... --fix

matlab-lint:
	matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean lint;"

# =============================================================================
# 5. Resources
# =============================================================================

sync-matlab-mcp:
	@$(call MK_DIR,$(MATLAB_MCP_TOOLBOX_DST))
	@$(call CP,$(MATLAB_MCP_EMBEDDED_SRC)/*.m,$(MATLAB_MCP_TOOLBOX_DST)/)
	@echo "Synced embedded MATLAB files to toolbox"

update-coding-guidelines:
ifeq ($(OS),Windows_NT)
	Invoke-WebRequest -Uri "$(CODING_GUIDELINES_URL)" -OutFile "$(CODING_GUIDELINES_PATH)"
else
	curl -sSL "$(CODING_GUIDELINES_URL)" -o "$(CODING_GUIDELINES_PATH)"
endif

update-live-code-guidelines:
ifeq ($(OS),Windows_NT)
	Invoke-WebRequest -Uri "$(LIVE_CODE_GUIDELINES_URL)" -OutFile "$(LIVE_CODE_GUIDELINES_PATH)"
else
	curl -sSL "$(LIVE_CODE_GUIDELINES_URL)" -o "$(LIVE_CODE_GUIDELINES_PATH)"
endif

# =============================================================================
# 6. Building
# =============================================================================

ensure-binary-executable:
ifneq ($(OS),Windows_NT)
	@chmod +x "$(BIN_PATH)/matlab-mcp-core-server" 2>/dev/null || true
endif

# Windows .exe doesn't need execute bit; build-mcpb-bundle only runs on macOS/Linux
ensure-all-binaries-executable:
ifneq ($(OS),Windows_NT)
	@chmod +x "$(ALL_BIN_DIR)/matlab-mcp-core-server-glnxa64" "$(ALL_BIN_DIR)/matlab-mcp-core-server-maca64" "$(ALL_BIN_DIR)/matlab-mcp-core-server-maci64" 2>/dev/null || true
endif

build: build-for-windows build-for-glnxa64 build-for-maci64 build-for-maca64
	@$(call MK_DIR,$(ALL_BIN_DIR))
	@$(call CP,$(GLNXA64_BIN_DIR)/matlab-mcp-core-server,$(ALL_BIN_DIR)/matlab-mcp-core-server-glnxa64)
	@$(call CP,$(MACA64_BIN_DIR)/matlab-mcp-core-server,$(ALL_BIN_DIR)/matlab-mcp-core-server-maca64)
	@$(call CP,$(MACI64_BIN_DIR)/matlab-mcp-core-server,$(ALL_BIN_DIR)/matlab-mcp-core-server-maci64)
	@$(call CP,$(WIN64_BIN_DIR)/matlab-mcp-core-server.exe,$(ALL_BIN_DIR)/matlab-mcp-core-server-win64.exe)

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

build-all:
	@$(call MK_DIR,$(ALL_BIN_DIR))
	@$(call CP,$(GLNXA64_BIN_DIR)/matlab-mcp-core-server,$(ALL_BIN_DIR)/matlab-mcp-core-server-glnxa64)
	@$(call CP,$(MACA64_BIN_DIR)/matlab-mcp-core-server,$(ALL_BIN_DIR)/matlab-mcp-core-server-maca64)
	@$(call CP,$(MACI64_BIN_DIR)/matlab-mcp-core-server,$(ALL_BIN_DIR)/matlab-mcp-core-server-maci64)
	@$(call CP,$(WIN64_BIN_DIR)/matlab-mcp-core-server.exe,$(ALL_BIN_DIR)/matlab-mcp-core-server-win64.exe)

build-tools: build-sourcehash build-mcpb-gen

build-sourcehash:
	@$(call MK_DIR,$(TOOLS_BIN_DIR))
	go build -o "$(SOURCEHASH_BIN)" ./cmd/sourcehash

build-mcpb-gen:
	@$(call MK_DIR,$(TOOLS_BIN_DIR))
	go build -o "$(MCPB_GEN_BIN)" ./cmd/mcpb-gen

build-matlab-addon: sync-matlab-mcp
ifeq ($(OS),Windows_NT)
	$$env:MATLAB_MCP_CORE_SERVER_MLTBX_DIR='$(MLTBX_DIR)'; matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean package;"
else
	MATLAB_MCP_CORE_SERVER_MLTBX_DIR="$(MLTBX_DIR)" matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean package;"
endif

update-embedded-matlab-addon: build-tools sync-matlab-mcp
ifeq ($(OS),Windows_NT)
	$$env:MATLAB_MCP_CORE_SERVER_MLTBX_DIR='$(EMBEDDED_MLTBX_DIR)'; matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean package;"
	& "$(SOURCEHASH_BIN)" write "$(SOURCES_HASH_FILE)" "$(MATLAB_TOOLBOX_DIR)"
else
	MATLAB_MCP_CORE_SERVER_MLTBX_DIR="$(EMBEDDED_MLTBX_DIR)" matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean package;"
	"$(SOURCEHASH_BIN)" write "$(SOURCES_HASH_FILE)" "$(MATLAB_TOOLBOX_DIR)"
endif

check-embedded-matlab-addon: build-sourcehash
ifeq ($(OS),Windows_NT)
	& "$(SOURCEHASH_BIN)" check "$(SOURCES_HASH_FILE)" "$(MATLAB_TOOLBOX_DIR)"
else
	"$(SOURCEHASH_BIN)" check "$(SOURCES_HASH_FILE)" "$(MATLAB_TOOLBOX_DIR)"
endif

# =============================================================================
# 7. Testing
# =============================================================================

mcp-inspector: export PATH := $(BIN_PATH)$(PATHSEP)$(PATH)
mcp-inspector: export HOST := localhost
mcp-inspector:
	npx @modelcontextprotocol/inspector matlab-mcp-core-server

unit-tests:
	go tool gotestsum --packages="$(UNIT_TEST_PKGS)" -- -race -coverprofile cover.out

matlab-unit-tests:
	matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean unit-tests;"

integration-tests:
	go tool gotestsum --packages="$(INTEGRATION_TEST_PKGS)" -- -race

functional-tests:
	go tool gotestsum --packages="$(FUNCTIONAL_TEST_PKGS)" -- -race

system-tests:
	go tool gotestsum --packages="$(SYSTEM_TEST_PKGS)" -- -race -count=1 -timeout 30m
	@$(CHECK_MATLAB_LEAKS)

matlab-system-tests: export MATLAB_MCP_CORE_SERVER_MLTBX_DIR := $(MATLAB_MCP_CORE_SERVER_MLTBX_DIR)
matlab-system-tests:
	matlab -batch "cd(fullfile('$(CURDIR)', 'tests', 'system', 'matlab')); buildtool;"

# Tests should clean up all MATLAB sessions; this catches leaks
check-matlab-leaks:
	@$(CHECK_MATLAB_LEAKS)

# =============================================================================
# 8. MCPB
# =============================================================================

mcpb-stage: build-mcpb-gen
ifeq ($(OS),Windows_NT)
	@echo "Error: MCPB manifest generation is only supported on macOS/Linux"; exit 1
else
	MCPB_STAGING_DIR="$(MCPB_STAGING_DIR)" "$(MCPB_GEN_BIN)"
endif

# Requires all 4 platform binaries in $(ALL_BIN_DIR).
# Local dev: make mcpb-dev | CI/signed: populate all/ then make build-mcpb-bundle
build-mcpb-bundle: mcpb-stage ensure-all-binaries-executable
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
	@echo "Removed $(MCPB_STAGING_DIR)"

mcpb-dev: mcpb-clean build build-mcpb-bundle

mcpb-validate:
ifeq ($(OS),Windows_NT)
	@echo "Error: MCPB validation is only supported on macOS/Linux"; exit 1
else
	cd "$(MCPB_STAGING_DIR)"; \
	npm run mcpb-validate
endif

# =============================================================================
# 9. Internal
# =============================================================================
# define/endef for readability; $(strip ...) flattens for single-line execution.

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

# Fails CI if generated files are stale.
# Usage: $(call CHECK_GIT_CLEAN,<description>,<optional-path-scope>)

define CHECK_GIT_CLEAN_CMD
echo "Checking for uncommitted changes$(if $(2), in $(2))...";
git_status=$$(git status --porcelain $(2));
if [ -n "$$git_status" ]; then
	echo "";
	echo "ERROR: $(1) are out of date. Please regenerate and commit.";
	echo "Changed files:";
	echo "$$git_status";
	exit 1;
fi;
echo "OK: $(1) are up to date."
endef

CHECK_GIT_CLEAN = $(strip $(CHECK_GIT_CLEAN_CMD))
