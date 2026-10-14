# System Tests

---

This README is intended for MathWorks® developers only.

---

System tests validate the MATLAB MCP Server end-to-end through realistic workflows. These tests use actual MATLAB installations and exercise the full server stack.

## Test Philosophy

System tests focus on **realistic user scenarios** rather than exhaustive feature testing:

- **Scenario-driven**: Tests simulate actual use cases that users encounter
- **Workflow-oriented**: Multiple features tested together as they would be used in practice

Think of system tests as answering: "Does this work for a real user trying to solve a real problem?"

## Running System Tests

### Prerequisites
```bash
# Build the server
make build
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `MCP_MATLAB_PATH` | Directory containing the MATLAB executable (e.g., `/usr/local/MATLAB/R2024b/bin`) | Yes |
| `MATLAB_MCP_SERVER_BUILD_DIR` | Base directory for built binaries. The test locates the server at `<base>/<os>/matlab-mcp-server` (e.g., `.bin/glnxa64/matlab-mcp-server`) | Only for `go test` |

The Makefile automatically sets `MATLAB_MCP_SERVER_BUILD_DIR` to `.bin/`, so you only need to set it when running `go test` directly.

### Run Tests
```bash
# Using make (recommended) - automatically sets binary path
export MCP_MATLAB_PATH=/path/to/matlab/bin
make system-tests

# Using go test directly - must set both variables
export MCP_MATLAB_PATH=/path/to/matlab/bin
export MATLAB_MCP_SERVER_BUILD_DIR=$(pwd)/.bin
go test -v ./tests/system/...

# Run specific suite
go test -v ./tests/system/... -run TestWorkflowSuite

# Run specific scenario
go test -v ./tests/system/... -run TestInteractiveDevelopmentWorkflow
go test -v ./tests/system/... -run TestParallelExperimentationWorkflow
```

### Existing Session Tests

`TestExistingSessionSuite` installs the MATLAB MCP Server Toolbox, starts MATLAB sessions, and connects the server to them.
The toolbox needs MATLAB R2023a or later.
With an earlier release, the suite skips.

The suite checks that:

- In `existing` mode, the server fails to attach when no MATLAB is shared.
- In `existing` and `auto` modes, the server attaches to the shared MATLAB.
- The server ignores a MATLAB that is not shared.
- When a second MATLAB is shared, the server attaches to it.
- In `auto` mode, the server starts a local MATLAB when no MATLAB is shared.
- Each shared MATLAB stays alive after the server disconnects.

The suite sets `HOME`, `APPDATA` and `MATLAB_PREFDIR` to folders in a temporary folder.
The server and the toolbox find the session details folder from `HOME`, or from `APPDATA` on Windows, so your shared MATLAB sessions do not change the result.
With `MATLAB_PREFDIR` set, the toolbox install does not change your real MATLAB settings or add-ons.

### Leaked MATLAB Check

After the system tests, the Makefile looks for MATLAB processes that the tests did not stop.
It finds a MATLAB process when its command line contains one of these markers:

- `MATLAB_MCP_PROCESS_MARKER`, for a MATLAB that the server started.
- `MATLAB_SYSTEM_TEST_PROCESS_MARKER`, for a MATLAB that the existing session tests started.

A leaked MATLAB makes the run fail.
`make system-tests` only reports a leaked MATLAB, because the markers also find MATLAB sessions that your own MCP servers started.
`make ci-system-tests` stops each leaked MATLAB.
`go test` does not do this check.
After you run `go test` directly, run `make check-matlab-leaks`.

## Adding New Tests

When adding new features, think in terms of **user scenarios** rather than individual tools:

**Ask yourself:**
- "What user problem does this solve?"
- "When would someone actually use this?"
- "Does this fit into an existing scenario or require a new one?"

**Guidelines:**
1. **Fits existing scenario**: Add to the appropriate workflow test
2. **New user scenario**: Create a new test method describing the scenario 

---

Copyright 2025-2026 The MathWorks, Inc.

---