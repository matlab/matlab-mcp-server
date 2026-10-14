function startExternalMATLABCommands(commandDir)
    %startExternalMATLABCommands Run the test commands from commandDir on a timer
    %   A timer, not a loop, so that MATLAB can also serve the MCP server.

    % Copyright 2026 The MathWorks, Inc.

    arguments
        commandDir(1, 1) string
    end

    addpath(commandDir);
    commandTimer = timer(ExecutionMode="fixedSpacing", Period=0.2, ...
        TimerFcn=@(~, ~) runExternalMATLABCommands(commandDir));
    start(commandTimer);
    fclose(fopen(fullfile(commandDir, "ready"), "w"));
end
