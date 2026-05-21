classdef DefaultConnectorAdaptor < mcpcoreserver.internal.connector.ConnectorAdaptor
    %DefaultConnectorAdaptor Default implementation for MATLAB connector operations
    %   This class provides higher-level operations built on top of the
    %   connector facade.

    % Copyright 2026 The MathWorks, Inc.

    properties (GetAccess = private, SetAccess = immutable)
        ConnectorFacade(1, 1) mcpcoreserver.internal.facade.connector.ConnectorFacade = mcpcoreserver.internal.facade.connector.DefaultConnectorFacade()
        APIKeyManager(1, 1) mcpcoreserver.internal.connector.internal.apikeymanager.APIKeyManager = mcpcoreserver.internal.connector.internal.apikeymanager.DefaultAPIKeyManager()
        OSFacade(1, 1) mcpcoreserver.internal.facade.os.OSFacade = mcpcoreserver.internal.facade.os.DefaultOSFacade()
    end

    methods
        function obj = DefaultConnectorAdaptor(options)
            arguments
                options.?mcpcoreserver.internal.connector.DefaultConnectorAdaptor
            end

            for prop = string(fieldnames(options).')
                obj.(prop) = options.(prop);
            end
        end

        function details = getConnectionDetails(obj)
            details = struct( ...
                "port", obj.ConnectorFacade.securePort(), ...
                "certificate", obj.ConnectorFacade.getCertificateLocation(), ...
                "apiKey", obj.APIKeyManager.getAPIKey(), ...
                "pid", obj.OSFacade.feature("getpid") ...
            );
        end
    end

end
