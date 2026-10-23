function issueWarning()
    % Copyright 2026 The MathWorks, Inc.

    warningState = 9;
    warning("MCP:ExpectedWarning", "MCP_EXPECTED_WARNING state=%d", warningState);
    fprintf("MCP_WARNING_COMPLETE\n");
end
