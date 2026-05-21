function initializeMCP()
% initializeMCP intializes the MATLAB sessions for the MATLAB MCP Core Server.

% Copyright 2025-2026 The MathWorks, Inc.

sessionDir = getenv("MW_MCP_SESSION_DIR");

try
    matlab.graphics.internal.initialize(true);
    connector.ensureServiceOn();

    securePortFile = fullfile(sessionDir, "connector.securePort");
    securePort = connector.securePort();

    securePortFileID = fopen(securePortFile, "w");
    if securePortFileID == -1
        error("Failed to open secure port file: %s", securePortFile);
    end
    closeSecurePortFile = onCleanup(@() fclose(securePortFileID));
    fprintf(securePortFileID, "%d", securePort);
catch e
    writeStartupError(sessionDir, e);
    rethrow(e);
end
end

function writeStartupError(sessionDir, e)
if isempty(sessionDir)
    return;
end

errFile = fullfile(sessionDir, "mcp_startup_error.txt");
fid = fopen(errFile, "w");
if fid == -1
    return;
end
closeFile = onCleanup(@() fclose(fid));

fprintf(fid, "%s\n\n%s\n", e.message, getReport(e, "extended", "hyperlinks", "off"));
end
