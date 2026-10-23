function stopAtKeyboard()
    % Copyright 2026 The MathWorks, Inc.

    keyboardState = 11;
    keyboard(); %#ok<KEYBOARDFUN>
    fprintf("MCP_KEYBOARD_COMPLETE=%d\n", keyboardState);
end
