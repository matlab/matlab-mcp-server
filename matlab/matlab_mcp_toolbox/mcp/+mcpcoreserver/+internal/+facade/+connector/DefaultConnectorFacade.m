classdef DefaultConnectorFacade < mcpcoreserver.internal.facade.connector.ConnectorFacade
    %DefaultConnectorFacade Default implementation for MATLAB connector operations
    %   This class delegates to MATLAB's connector functions.

    % Copyright 2026 The MathWorks, Inc.

    methods
        function ensureServiceOn(~)
            connector.ensureServiceOn();
        end

        function port = securePort(~)
            port = connector.securePort();
        end

        function location = getCertificateLocation(~)
            location = connector.getCertificateLocation();
        end

        function value = getConfig(~, key)
            value = connector.internal.getConfig(key);
        end

        function setConfig(~, key, value)
            connector.internal.setConfig(key, value);
        end
    end

end
