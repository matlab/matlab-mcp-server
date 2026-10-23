function throwUncaughtError()
    % Copyright 2026 The MathWorks, Inc.

    uncaughtState = 7;
    error("MCP:ExpectedUncaught", "MCP_EXPECTED_UNCAUGHT state=%d", uncaughtState);
end
