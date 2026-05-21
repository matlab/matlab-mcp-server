classdef (Abstract) AppDataLocator
    %AppDataLocator Abstract locator for accessing application data folder
    %   This abstract class defines the interface for retrieving the
    %   platform-specific application data folder path.

    % Copyright 2026 The MathWorks, Inc.

    methods (Abstract)
        path = getAppDataFolder(obj)
    end

end
