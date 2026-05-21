function shareMATLABSession(options)
    %shareMATLABSession Share the current MATLAB session via MCP server
    %   This function enables sharing of the MATLAB session through the
    %   Model Context Protocol (MCP) server.

    % Copyright 2026 The MathWorks, Inc.

    arguments
        options.AppDataLocator(1, 1) mcpcoreserver.internal.appdata.AppDataLocator = mcpcoreserver.internal.appdata.DefaultAppDataLocator()
        options.FSAdaptor(1, 1) mcpcoreserver.internal.fs.FSAdaptor = mcpcoreserver.internal.fs.DefaultFSAdaptor()
        options.FSFacade(1, 1) mcpcoreserver.internal.facade.fs.FSFacade = mcpcoreserver.internal.facade.fs.DefaultFSFacade()
        options.ConnectorAdaptor(1, 1) mcpcoreserver.internal.connector.ConnectorAdaptor = mcpcoreserver.internal.connector.DefaultConnectorAdaptor()
        options.ConnectorFacade(1, 1) mcpcoreserver.internal.facade.connector.ConnectorFacade = mcpcoreserver.internal.facade.connector.DefaultConnectorFacade()
    end

    options.ConnectorFacade.ensureServiceOn()

    appDataFolder = options.AppDataLocator.getAppDataFolder();
    options.FSAdaptor.ensureSecureFolder(appDataFolder);

    v1Folder = fullfile(appDataFolder, "v1");
    options.FSAdaptor.ensureSecureFolder(v1Folder);

    sessionDetailsPath = fullfile(v1Folder, "sessionDetails.json");
    options.FSAdaptor.ensureSecureFile(sessionDetailsPath);

    sessionDetails = options.ConnectorAdaptor.getConnectionDetails();
    jsonText = jsonencode(sessionDetails, PrettyPrint=true);
    options.FSFacade.writelines(jsonText, sessionDetailsPath);
end
