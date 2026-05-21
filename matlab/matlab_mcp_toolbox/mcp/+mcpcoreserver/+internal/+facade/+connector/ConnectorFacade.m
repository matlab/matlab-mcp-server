classdef (Abstract) ConnectorFacade
    %ConnectorFacade Abstract facade for MATLAB connector operations
    %   This abstract class defines the interface for MATLAB connector
    %   operations.

    % Copyright 2026 The MathWorks, Inc.

    methods (Abstract)
        ensureServiceOn(obj)
        port = securePort(obj)
        location = getCertificateLocation(obj)
        value = getConfig(obj, key)
        setConfig(obj, key, value)
    end

end
