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
	.SHELLFLAGS = -Command
	RACE_FLAG =
	RM_DIR = if (Test-Path "$(1)") { Remove-Item -Recurse -Force "$(1)" }
	MK_DIR = New-Item -ItemType Directory -Force -Path "$(1)" | Out-Null
	CP = Copy-Item "$(1)" "$(2)"
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

MATLAB_MCP_SERVER_BUILD_DIR ?= $(CURDIR)/.bin
MATLAB_MCP_SERVER_MLTBX_DIR ?= $(EMBEDDED_MLTBX_DIR)

# --- Build paths ---

TOOLS_BIN_DIR := $(MATLAB_MCP_SERVER_BUILD_DIR)/tools
SOURCEHASH_BIN := $(TOOLS_BIN_DIR)/sourcehash$(EXE_SUFFIX)
MCPB_GEN_BIN := $(TOOLS_BIN_DIR)/mcpb-gen$(EXE_SUFFIX)
REGISTRY_VALIDATE_BIN := $(TOOLS_BIN_DIR)/registry-validate$(EXE_SUFFIX)

WIN64_BIN_DIR := $(MATLAB_MCP_SERVER_BUILD_DIR)/win64
GLNXA64_BIN_DIR := $(MATLAB_MCP_SERVER_BUILD_DIR)/glnxa64
MACI64_BIN_DIR := $(MATLAB_MCP_SERVER_BUILD_DIR)/maci64
MACA64_BIN_DIR := $(MATLAB_MCP_SERVER_BUILD_DIR)/maca64
ALL_BIN_DIR := $(MATLAB_MCP_SERVER_BUILD_DIR)/all

MLTBX_DIR := $(MATLAB_MCP_SERVER_BUILD_DIR)/mltbx
SOURCES_HASH_FILE := $(EMBEDDED_MLTBX_DIR)/.sources-hash
MATLAB_TOOLBOX_DIR := $(CURDIR)/matlab/matlab_mcp_toolbox

TELEMETRY_ENDPOINT_FILE := $(CURDIR)/internal/adaptors/application/parameter/defaultparameters/endpoint_generated.go

MCPB_STAGING_DIR := $(MATLAB_MCP_SERVER_BUILD_DIR)/mcpb
MCPB_FILENAME := matlab-mcp-server.mcpb
registry-render: export VERSION := $(value VERSION)
registry-render: export FILE_SHA256 := $(value FILE_SHA256)
registry-render: export MCPB_FILENAME := $(value MCPB_FILENAME)

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

# Must stay equal to a substring of startupCode in localmatlabsession.go.
MATLAB_MCP_PROCESS_MARKER := matlab_mcp.initializeMCP

# Must stay equal to startupDirName in externalmatlab.go.
MATLAB_SYSTEM_TEST_PROCESS_MARKER := mcp-system-test-matlab

# Scoped exports prevent polluting other commands (e.g. npm uses HOST)
TEST_TARGETS := \
	unit-tests  matlab-unit-tests  ci-unit-tests          \
	matlab-integration-tests                              \
	matlab-integration-tests-dotnet-framework             \
	matlab-integration-tests-dotnet-core                  \
	integration-tests              ci-integration-tests    \
	functional-tests               ci-functional-tests     \
	system-tests                   ci-system-tests         \
	matlab-system-tests            ci-matlab-system-tests
$(TEST_TARGETS): export PATH := $(BIN_PATH)$(PATHSEP)$(PATH)
$(TEST_TARGETS): export MATLAB_MCP_SERVER_BUILD_DIR := $(MATLAB_MCP_SERVER_BUILD_DIR)
$(TEST_TARGETS): export MCP_MATLAB_PATH := $(MCP_MATLAB_PATH)
$(TEST_TARGETS): export MCPB_ARTIFACT_PATH := $(or $(MCPB_ARTIFACT_PATH),$(MCPB_STAGING_DIR)/$(MCPB_FILENAME))

# --- Build flags ---

# Make sure we do not have incompatible build flags set. 
# DEBUG and RELEASE are mutually exclusive.
ifneq ($(and $(filter true,$(DEBUG)), $(filter true,$(RELEASE))),)
$(error DEBUG and RELEASE cannot both be true)
endif

ifeq ($(DEBUG),true)
	# Used for debugging the MCP server - use something like
	# dlv attach $(pgrep -f "matlab-mcp-server") --headless --listen=:2345 --api-version=2
	BUILD_FLAGS := -gcflags="all=-N -l"
else
	BUILD_FLAGS := -trimpath
endif

ifeq ($(RELEASE),true)
	LDFLAGS_ARG := -ldflags "-s -w"
else
	LDFLAGS_ARG :=
endif

# Gating on RELEASE, not on the value, stops a developer's environment or .env baking in.
ifeq ($(RELEASE),true)
	TAGS_ARG := -tags release
	VERIFY_ENDPOINT = go run ./cmd/verify-telemetry-endpoint "$(1)"
else
	TAGS_ARG :=
	VERIFY_ENDPOINT =
endif

# The export lets a "make VAR=value" assignment reach go run's environment.
export TELEMETRY_COLLECTOR_ENDPOINT

all: wire mockery lint unit-tests integration-tests build mcpb-clean mcpb-dev functional-tests

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
	matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean unit-tests integration-tests;"

ci-matlab-unit-tests:
	@echo "ci-matlab-unit-tests is no longer needed. MATLAB unit tests are now included in ci-unit-tests."

ci-build-matlab-addon: build-matlab-addon

ci-build: build

ci-build-mcpb: build-mcpb-bundle mcpb-validate ci-validate-registry

ci-validate-registry: registry-validate

ci-render-registry: registry-render

ci-integration-tests: matlab-integration-tests
	go test $(RACE_FLAG) -json -count=1 $(INTEGRATION_TEST_PKGS)

ci-functional-tests: ensure-binary-executable
	go test $(RACE_FLAG) -json -count=1 $(FUNCTIONAL_TEST_PKGS)

ci-system-tests: ensure-binary-executable
	@$(call RUN_AND_REAP,go test $(RACE_FLAG) -timeout 120m -json -count=1 $(SYSTEM_TEST_PKGS),KILL_LEAKED_MATLABS)

ci-matlab-system-tests: matlab-system-tests

# =============================================================================
# 3. Code Generation
# =============================================================================

wire:
	go tool wire github.com/matlab/matlab-mcp-server/internal/wire

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
	@chmod +x "$(BIN_PATH)/matlab-mcp-server" 2>/dev/null || true
endif

# Windows .exe doesn't need execute bit; build-mcpb-bundle only runs on macOS/Linux
ensure-all-binaries-executable:
ifneq ($(OS),Windows_NT)
	@chmod +x "$(ALL_BIN_DIR)/matlab-mcp-server-linux-x64" "$(ALL_BIN_DIR)/matlab-mcp-server-macos-arm64" "$(ALL_BIN_DIR)/matlab-mcp-server-macos-x64" 2>/dev/null || true
endif

build: build-for-windows build-for-glnxa64 build-for-maci64 build-for-maca64
	@$(call MK_DIR,$(ALL_BIN_DIR))
	@$(call CP,$(GLNXA64_BIN_DIR)/matlab-mcp-server,$(ALL_BIN_DIR)/matlab-mcp-server-linux-x64)
	@$(call CP,$(MACA64_BIN_DIR)/matlab-mcp-server,$(ALL_BIN_DIR)/matlab-mcp-server-macos-arm64)
	@$(call CP,$(MACI64_BIN_DIR)/matlab-mcp-server,$(ALL_BIN_DIR)/matlab-mcp-server-macos-x64)
	@$(call CP,$(WIN64_BIN_DIR)/matlab-mcp-server.exe,$(ALL_BIN_DIR)/matlab-mcp-server-windows-x64.exe)

generate-telemetry-endpoint:
ifeq ($(RELEASE),true)
	go run ./cmd/generate-telemetry-endpoint "$(TELEMETRY_ENDPOINT_FILE)"
endif

build-for-windows: generate-telemetry-endpoint
ifeq ($(OS),Windows_NT)
	$$env:GOOS='windows'; $$env:GOARCH='amd64'; $$env:CGO_ENABLED='0'; go build $(BUILD_FLAGS) $(LDFLAGS_ARG) $(TAGS_ARG) -o $(WIN64_BIN_DIR)/matlab-mcp-server.exe ./cmd/matlab-mcp-server
else
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS_ARG) $(TAGS_ARG) -o "$(WIN64_BIN_DIR)/matlab-mcp-server.exe" ./cmd/matlab-mcp-server
endif
	$(call VERIFY_ENDPOINT,$(WIN64_BIN_DIR)/matlab-mcp-server.exe)

build-for-glnxa64: generate-telemetry-endpoint
ifeq ($(OS),Windows_NT)
	$$env:GOOS='linux'; $$env:GOARCH='amd64'; $$env:CGO_ENABLED='0'; go build $(BUILD_FLAGS) $(LDFLAGS_ARG) $(TAGS_ARG) -o $(GLNXA64_BIN_DIR)/matlab-mcp-server ./cmd/matlab-mcp-server
else
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS_ARG) $(TAGS_ARG) -o "$(GLNXA64_BIN_DIR)/matlab-mcp-server" ./cmd/matlab-mcp-server
endif
	$(call VERIFY_ENDPOINT,$(GLNXA64_BIN_DIR)/matlab-mcp-server)

build-for-maci64: generate-telemetry-endpoint
ifeq ($(OS),Windows_NT)
	$$env:GOOS='darwin'; $$env:GOARCH='amd64'; $$env:CGO_ENABLED='0'; go build $(BUILD_FLAGS) $(LDFLAGS_ARG) $(TAGS_ARG) -o $(MACI64_BIN_DIR)/matlab-mcp-server ./cmd/matlab-mcp-server
else
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS_ARG) $(TAGS_ARG) -o "$(MACI64_BIN_DIR)/matlab-mcp-server" ./cmd/matlab-mcp-server
endif
	$(call VERIFY_ENDPOINT,$(MACI64_BIN_DIR)/matlab-mcp-server)

build-for-maca64: generate-telemetry-endpoint
ifeq ($(OS),Windows_NT)
	$$env:GOOS='darwin'; $$env:GOARCH='arm64'; $$env:CGO_ENABLED='0'; go build $(BUILD_FLAGS) $(LDFLAGS_ARG) $(TAGS_ARG) -o $(MACA64_BIN_DIR)/matlab-mcp-server ./cmd/matlab-mcp-server
else
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) $(LDFLAGS_ARG) $(TAGS_ARG) -o "$(MACA64_BIN_DIR)/matlab-mcp-server" ./cmd/matlab-mcp-server
endif
	$(call VERIFY_ENDPOINT,$(MACA64_BIN_DIR)/matlab-mcp-server)

build-all:
	@$(call MK_DIR,$(ALL_BIN_DIR))
	@$(call CP,$(GLNXA64_BIN_DIR)/matlab-mcp-server,$(ALL_BIN_DIR)/matlab-mcp-server-linux-x64)
	@$(call CP,$(MACA64_BIN_DIR)/matlab-mcp-server,$(ALL_BIN_DIR)/matlab-mcp-server-macos-arm64)
	@$(call CP,$(MACI64_BIN_DIR)/matlab-mcp-server,$(ALL_BIN_DIR)/matlab-mcp-server-macos-x64)
	@$(call CP,$(WIN64_BIN_DIR)/matlab-mcp-server.exe,$(ALL_BIN_DIR)/matlab-mcp-server-windows-x64.exe)

build-tools: build-sourcehash build-mcpb-gen build-registry-validate

build-sourcehash:
	@$(call MK_DIR,$(TOOLS_BIN_DIR))
	go build -o "$(SOURCEHASH_BIN)" ./cmd/sourcehash

build-mcpb-gen:
	@$(call MK_DIR,$(TOOLS_BIN_DIR))
	go build -o "$(MCPB_GEN_BIN)" ./cmd/mcpb-gen

build-registry-validate:
	@$(call MK_DIR,$(TOOLS_BIN_DIR))
	go build -o "$(REGISTRY_VALIDATE_BIN)" ./cmd/registry-validate

print-mcpb-filename:
	@echo "$(MCPB_FILENAME)"

build-matlab-addon: sync-matlab-mcp
ifeq ($(OS),Windows_NT)
	$$env:MATLAB_MCP_SERVER_MLTBX_DIR='$(MLTBX_DIR)'; matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean package;"
else
	MATLAB_MCP_SERVER_MLTBX_DIR="$(MLTBX_DIR)" matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean package;"
endif

update-embedded-matlab-addon: build-tools sync-matlab-mcp
ifeq ($(OS),Windows_NT)
	$$env:MATLAB_MCP_SERVER_MLTBX_DIR='$(EMBEDDED_MLTBX_DIR)'; matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean package;"
	& "$(SOURCEHASH_BIN)" write "$(SOURCES_HASH_FILE)" "$(MATLAB_TOOLBOX_DIR)"
else
	MATLAB_MCP_SERVER_MLTBX_DIR="$(EMBEDDED_MLTBX_DIR)" matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean package;"
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
	npx @modelcontextprotocol/inspector matlab-mcp-server

unit-tests:
	go tool gotestsum --packages="$(UNIT_TEST_PKGS)" -- -race -coverprofile cover.out

matlab-unit-tests:
	matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean unit-tests integration-tests;"

ifeq ($(OS),Windows_NT)
matlab-integration-tests: matlab-integration-tests-dotnet-framework matlab-integration-tests-dotnet-core

matlab-integration-tests-dotnet-framework:
	matlab -batch "setenv('MATLAB_MCP_TEST_RESULTS_SUFFIX','dotnet-framework'); dotnetenv('framework'); cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean integration-tests;"

matlab-integration-tests-dotnet-core:
	matlab -batch "setenv('MATLAB_MCP_TEST_RESULTS_SUFFIX','dotnet-core'); dotnetenv('core'); cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean integration-tests;"
else
matlab-integration-tests:
	matlab -batch "cd(fullfile('$(CURDIR)', 'matlab', 'matlab_mcp_toolbox')); buildtool clean integration-tests;"
endif

integration-tests:
	go tool gotestsum --packages="$(INTEGRATION_TEST_PKGS)" -- -race

functional-tests:
	go tool gotestsum --packages="$(FUNCTIONAL_TEST_PKGS)" -- -race

system-tests:
	@$(call RUN_AND_REAP,go tool gotestsum --packages="$(SYSTEM_TEST_PKGS)" -- -race -count=1 -timeout 30m,REPORT_LEAKED_MATLABS)

matlab-system-tests: export MATLAB_MCP_SERVER_MLTBX_DIR := $(MATLAB_MCP_SERVER_MLTBX_DIR)
matlab-system-tests:
	matlab -batch "cd(fullfile('$(CURDIR)', 'tests', 'system', 'matlab')); buildtool;"

# Tests should clean up all MATLAB sessions; this catches leaks
check-matlab-leaks:
	@$(call REAP_MATLAB_LEAKS,REPORT_LEAKED_MATLABS)

# =============================================================================
# 8. MCPB
# =============================================================================

mcpb-stage: build-mcpb-gen
ifeq ($(OS),Windows_NT)
	$$env:MCPB_STAGING_DIR='$(MCPB_STAGING_DIR)'; & "$(MCPB_GEN_BIN)"
else
	MCPB_STAGING_DIR="$(MCPB_STAGING_DIR)" "$(MCPB_GEN_BIN)"
endif

pack-mcpb-bundle: mcpb-stage ensure-all-binaries-executable
ifeq ($(OS),Windows_NT)
	@Copy-Item "$(ALL_BIN_DIR)/matlab-mcp-server-*" "$(MCPB_STAGING_DIR)/bundle/bin/" -ErrorAction Stop
	@Push-Location "$(MCPB_STAGING_DIR)"; \
		npm ci; \
		if ($$LASTEXITCODE -ne 0) { Pop-Location; exit 1 }; \
		npm run mcpb-pack -- "$(MCPB_FILENAME)"; \
		$$ec = $$LASTEXITCODE; \
		Pop-Location; \
		exit $$ec
	@Write-Host ""
	@Write-Host "Created: $(MCPB_STAGING_DIR)/$(MCPB_FILENAME)"
else
	@cp "$(ALL_BIN_DIR)"/matlab-mcp-server-* "$(MCPB_STAGING_DIR)/bundle/bin/"
	@cd "$(MCPB_STAGING_DIR)" && npm ci && npm run mcpb-pack -- "$(MCPB_FILENAME)"
	@echo ""
	@echo "Created: $(MCPB_STAGING_DIR)/$(MCPB_FILENAME)"
endif

build-mcpb-bundle:
ifeq ($(OS),Windows_NT)
	@Write-Error "build-mcpb-bundle cannot be called directly on Windows. Use 'make mcpb-dev' instead."; exit 1
else
	@if [ ! -f "$(ALL_BIN_DIR)/matlab-mcp-server-linux-x64" ] || \
		[ ! -f "$(ALL_BIN_DIR)/matlab-mcp-server-macos-arm64" ] || \
		[ ! -f "$(ALL_BIN_DIR)/matlab-mcp-server-macos-x64" ] || \
		[ ! -f "$(ALL_BIN_DIR)/matlab-mcp-server-windows-x64.exe" ]; then \
		echo "Error: Missing binaries in $(ALL_BIN_DIR)."; \
		echo "Run 'make mcpb-dev' for local builds, or populate $(ALL_BIN_DIR) with signed binaries."; \
		exit 1; \
	fi
	@$(MAKE) pack-mcpb-bundle
endif

mcpb-clean:
	@$(call RM_DIR,$(MCPB_STAGING_DIR))
	@echo "Removed $(MCPB_STAGING_DIR)"

mcpb-dev: mcpb-clean build pack-mcpb-bundle

mcpb-validate:
ifeq ($(OS),Windows_NT)
	@Push-Location "$(MCPB_STAGING_DIR)"; npm run mcpb-validate; $$ec = $$LASTEXITCODE; Pop-Location; exit $$ec
else
	cd "$(MCPB_STAGING_DIR)"; \
	npm run mcpb-validate
endif

registry-validate: build-registry-validate
ifeq ($(OS),Windows_NT)
	& "$(REGISTRY_VALIDATE_BIN)" .registry/server.template.json
else
	"$(REGISTRY_VALIDATE_BIN)" .registry/server.template.json
endif

registry-render: build-registry-validate
ifeq ($(OS),Windows_NT)
	& "$(REGISTRY_VALIDATE_BIN)" --render .registry/server.json --version "$$env:VERSION" --file-sha256 "$$env:FILE_SHA256" --mcpb-filename "$$env:MCPB_FILENAME" .registry/server.template.json
else
	"$(REGISTRY_VALIDATE_BIN)" --render .registry/server.json --version "$${VERSION}" --file-sha256 "$${FILE_SHA256}" --mcpb-filename "$${MCPB_FILENAME}" .registry/server.template.json
endif

# =============================================================================
# 9. Internal
# =============================================================================
# define/endef for readability; $(strip ...) flattens for single-line execution.

ifeq ($(OS),Windows_NT)

define REAP_MATLAB_LEAKS_CMD
powershell -NoProfile -ExecutionPolicy Bypass -Command "& {
    Write-Host 'Waiting for processes to settle...';
    Start-Sleep -Seconds 5;
    Write-Host 'Checking for leaked MATLAB processes...';
    `$$leaked = Get-CimInstance Win32_Process |
        Where-Object { `$$_.Name -eq 'MATLAB.exe' -and (`$$_.CommandLine -like '*$(MATLAB_MCP_PROCESS_MARKER)*' -or `$$_.CommandLine -like '*$(MATLAB_SYSTEM_TEST_PROCESS_MARKER)*') };
    if (-not `$$leaked) {
        Write-Host 'No leaked MATLAB processes found.';
        exit 0
    };
    Write-Host 'ERROR: Found leaked MATLAB processes:';
    `$$leaked | ForEach-Object { Write-Host ('  ' + `$$_.ProcessId + '  ' + `$$_.CommandLine) };
    $($(1))
    exit 1
}"
endef

define KILL_LEAKED_MATLABS
`$$leaked | ForEach-Object { taskkill /T /F /PID `$$_.ProcessId 2>&1 | Out-Null };
Write-Host 'Killed the leaked MATLAB processes. A leak changes the result of later runs, so this fails the job.';
endef

define REPORT_LEAKED_MATLABS
Write-Host 'Stop these MATLAB processes before the next run. A leak changes the result of later runs, so this fails the run.';
endef

define RUN_AND_REAP_CMD
$$status = 0;
$(1);
if (-not $$?) { $$status = if ($$LASTEXITCODE) { $$LASTEXITCODE } else { 1 } };
$(call REAP_MATLAB_LEAKS,$(2));
if ($$LASTEXITCODE -ne 0 -and $$status -eq 0) { $$status = 1 };
exit $$status
endef

else

define REAP_MATLAB_LEAKS_CMD
echo "Waiting for processes to settle...";
sleep 5;
echo "Checking for leaked MATLAB processes...";
leaked=$$(pgrep -f '$(MATLAB_MCP_PROCESS_MARKER)|$(MATLAB_SYSTEM_TEST_PROCESS_MARKER)' 2>/dev/null |
    while read -r pid; do
        if ps -o comm= -p "$$pid" 2>/dev/null | sed 's|.*/||' | grep -qix matlab; then echo "$$pid"; fi;
    done);
if [ -z "$$leaked" ]; then
    echo "No leaked MATLAB processes found.";
else
    echo "ERROR: Found leaked MATLAB processes:";
    for pid in $$leaked; do ps -o pid= -o args= -p "$$pid"; done;
    $($(1))
    false;
fi
endef

define KILL_LEAKED_MATLABS
kill -9 $$leaked 2>/dev/null || true;
echo "Killed the leaked MATLAB processes. A leak changes the result of later runs, so this fails the job.";
endef

define REPORT_LEAKED_MATLABS
echo "Stop these MATLAB processes before the next run. A leak changes the result of later runs, so this fails the run.";
endef

define RUN_AND_REAP_CMD
status=0;
$(1) || status=$$?;
$(call REAP_MATLAB_LEAKS,$(2)) || { [ "$$status" -eq 0 ] && status=1; };
exit "$$status"
endef

endif

# Only CI kills a leaked MATLAB. On a developer machine, the markers also match the
# MATLABs that the developer's own MCP servers started.
# Usage: $(call REAP_MATLAB_LEAKS,KILL_LEAKED_MATLABS|REPORT_LEAKED_MATLABS)

REAP_MATLAB_LEAKS = $(strip $(REAP_MATLAB_LEAKS_CMD))

# Runs the tests and the reaper in one recipe line, so make cannot skip the reaper
# when the tests fail.
# Usage: $(call RUN_AND_REAP,<test command>,KILL_LEAKED_MATLABS|REPORT_LEAKED_MATLABS)

RUN_AND_REAP = $(strip $(RUN_AND_REAP_CMD))

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

REPO_URL = $(shell git remote get-url origin | sed 's|git@\(.*\):\(.*\)\.git|https://\1/\2|')

compare-commits:
ifndef BASE
	$(error Usage: make compare-commits BASE=branch-or-tag HEAD=branch-or-tag)
endif
ifndef HEAD
	$(error Usage: make compare-commits BASE=branch-or-tag HEAD=branch-or-tag)
endif
	@echo "$(REPO_URL)/compare/$$(git rev-parse $(BASE))...$$(git rev-parse $(HEAD))"
