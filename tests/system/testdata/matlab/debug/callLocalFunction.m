function result = callLocalFunction(inputValue)
    % Copyright 2026 The MathWorks, Inc.

    outerValue = inputValue + 1;
    offsetValue = addOffset(outerValue);    % MCP_BREAK_CALL
    result = offsetValue * 2;
    fprintf("MCP_DEBUG_RESULT=%g\n", result);
end

function outputValue = addOffset(inputValue)
    outputValue = inputValue + 10;
end
