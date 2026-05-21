#!/usr/bin/env bash
# Copyright 2026 The MathWorks, Inc.
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"

# Detect platform
case "$(uname -s)" in
  Linux*)
    case "$(uname -m)" in
      x86_64) BIN="$DIR/matlab-mcp-core-server-glnxa64" ;;
      *)      echo "Unsupported Linux architecture: $(uname -m)" >&2; exit 1 ;;
    esac ;;
  Darwin*)
    case "$(uname -m)" in
      arm64)  BIN="$DIR/matlab-mcp-core-server-maca64" ;;
      x86_64) BIN="$DIR/matlab-mcp-core-server-maci64" ;;
      *)      echo "Unsupported macOS architecture: $(uname -m)" >&2; exit 1 ;;
    esac ;;
  *)        echo "Unsupported platform: $(uname -s)" >&2; exit 1 ;;
esac

# Verify binary exists
if [[ ! -x "$BIN" ]]; then
  echo "MATLAB MCP core server binary not found or not executable: $BIN" >&2
  exit 1
fi

export MW_CONTEXT_TAGS="${MW_CONTEXT_TAGS:+$MW_CONTEXT_TAGS,}MATLAB:MCPB:V1"

# Env var to CLI flag mappings (format: ENV_VAR:type:flag)
# Types: string = pass value if non-empty, bool = pass flag if "true"
MCPB_MAPPINGS=(
    "__MATLAB_MCP_CORE_SERVER_MCPB_MATLAB_ROOT:string:--matlab-root"
    "__MATLAB_MCP_CORE_SERVER_MCPB_INITIAL_WD:string:--initial-working-folder"
    "__MATLAB_MCP_CORE_SERVER_MCPB_LOG_DIR:string:--log-folder"
    "__MATLAB_MCP_CORE_SERVER_MCPB_INIT_ON_START:bool:--initialize-matlab-on-startup"
    "__MATLAB_MCP_CORE_SERVER_MCPB_DISABLE_TELEM:bool:--disable-telemetry"
    "__MATLAB_MCP_CORE_SERVER_MCPB_DISPLAY_MODE:string:--matlab-display-mode"
    "__MATLAB_MCP_CORE_SERVER_MCPB_MATLAB_SESSION_MODE:string:--matlab-session-mode"
    "__MATLAB_MCP_CORE_SERVER_MCPB_EXTENSION_FILE:string:--extension-file"
)

ARGS=()
for mapping in "${MCPB_MAPPINGS[@]}"; do
    IFS=: read -r env_var type flag <<< "$mapping"
    val="${!env_var:-}"
    case "$type" in
        string) [[ -n "$val" ]] && ARGS+=("$flag" "$val") ;;
        bool)   [[ "$val" == "true" ]] && ARGS+=("$flag") ;;
    esac
    unset "$env_var"
done

exec "$BIN" ${ARGS[@]+"${ARGS[@]}"}
