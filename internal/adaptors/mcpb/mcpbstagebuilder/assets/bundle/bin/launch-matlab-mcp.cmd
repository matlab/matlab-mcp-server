@echo off
REM Copyright 2026 The MathWorks, Inc.
setlocal enabledelayedexpansion
set BIN=%~dp0matlab-mcp-core-server-win64.exe
if defined MW_CONTEXT_TAGS (
    set "MW_CONTEXT_TAGS=%MW_CONTEXT_TAGS%,MATLAB:MCPB:V1"
) else (
    set "MW_CONTEXT_TAGS=MATLAB:MCPB:V1"
)
REM Env var to CLI flag mappings (format: ENV_VAR:type:flag)
REM Types: string = pass value if non-empty, bool = pass flag if "true"
set ARGS=
for %%M in (
    "__MATLAB_MCP_CORE_SERVER_MCPB_MATLAB_ROOT:string:--matlab-root"
    "__MATLAB_MCP_CORE_SERVER_MCPB_INITIAL_WD:string:--initial-working-folder"
    "__MATLAB_MCP_CORE_SERVER_MCPB_LOG_DIR:string:--log-folder"
    "__MATLAB_MCP_CORE_SERVER_MCPB_INIT_ON_START:bool:--initialize-matlab-on-startup"
    "__MATLAB_MCP_CORE_SERVER_MCPB_DISABLE_TELEM:bool:--disable-telemetry"
    "__MATLAB_MCP_CORE_SERVER_MCPB_DISPLAY_MODE:string:--matlab-display-mode"
    "__MATLAB_MCP_CORE_SERVER_MCPB_MATLAB_SESSION_MODE:string:--matlab-session-mode"
    "__MATLAB_MCP_CORE_SERVER_MCPB_EXTENSION_FILE:string:--extension-file"
) do (
    for /f "tokens=1,2,3 delims=:" %%A in ("%%~M") do (
        set "val=!%%A!"
        if "%%B"=="string" if defined %%A if not "!val!"=="" set ARGS=!ARGS! %%C "!val!"
        if "%%B"=="bool" if "!val!"=="true" set ARGS=!ARGS! %%C
        set "%%A="
    )
)

"%BIN%" %ARGS%
