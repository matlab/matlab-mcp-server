function throwAndCatchError()
    % Copyright 2026 The MathWorks, Inc.

    caughtState = 8;
    try
        error("MCP:ExpectedCaught", "MCP_EXPECTED_CAUGHT state=%d", caughtState);
    catch ME
        fprintf("MCP_CAUGHT_HANDLED=%s\n", ME.identifier);
    end
end
