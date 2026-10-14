function runExternalMATLABCommands(commandDir)
    %runExternalMATLABCommands Run each .request file in commandDir and write a .response file
    %   No error leaves this function, so that the timer keeps running.

    % Copyright 2026 The MathWorks, Inc.

    arguments
        commandDir(1, 1) string
    end

    try
        requests = dir(fullfile(commandDir, "*.request"));
        for k = 1:numel(requests)
            respondTo(commandDir, requests(k).name);
        end
    catch pollerErr
        fprintf(2, "runExternalMATLABCommands: %s\n", pollerErr.message);
    end
end

function respondTo(commandDir, requestName)
    requestPath = fullfile(commandDir, requestName);
    [~, commandName] = fileparts(requestName);
    code = fileread(requestPath);
    delete(requestPath);
    try
        evalin("base", code);
        status = "done";
    catch err
        status = "error: " + err.message + newline + describeEnvironment();
    end
    responsePath = fullfile(commandDir, commandName + ".response");
    responseFile = fopen(responsePath + ".partial", "w");
    fprintf(responseFile, "%s\n", status);
    fclose(responseFile);
    movefile(responsePath + ".partial", responsePath);
end

function report = describeEnvironment()
    report = strjoin([
        "prefdir: " + describe(@() prefdir)
        "userpath: " + describe(@() userpath)
        "installed add-ons: " + describe(@installedAddOnNames)
        "shareMATLABSession: " + describe(@() which("-all", "shareMATLABSession"))
        "HOME: " + describe(@() getenv("HOME"))
        "APPDATA: " + describe(@() getenv("APPDATA"))
    ], newline);
end

function names = installedAddOnNames()
    addOns = matlab.addons.installedAddons();
    names = addOns.Name;
end

function text = describe(query)
    try
        values = string(query());
        if isempty(values)
            text = "none";
        else
            text = strjoin(values(:).', ", ");
        end
    catch err
        text = "unavailable: " + err.message;
    end
end
