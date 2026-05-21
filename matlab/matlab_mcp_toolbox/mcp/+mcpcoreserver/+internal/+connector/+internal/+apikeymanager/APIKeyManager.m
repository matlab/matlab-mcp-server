classdef (Abstract) APIKeyManager
    %APIKeyManager Manages MATLAB Connector API key with version-specific logic
    %   This class handles API key retrieval based on MATLAB version.

    % Copyright 2026 The MathWorks, Inc.

    methods (Abstract)
        apiKey = getAPIKey(obj)
    end

end
