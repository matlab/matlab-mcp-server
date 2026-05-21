classdef DefaultAPIKeyManager < mcpcoreserver.internal.connector.internal.apikeymanager.APIKeyManager
    %DefaultAPIKeyManager Default implementation for MATLAB Connector API key management
    %   This class handles API key retrieval based on MATLAB version:
    %   - Prior to R2023a: throws UnsupportedMATLABVersion error
    %   - R2023a to R2024b: generates secure key, calls setConfig, then getConfig
    %   - R2025a onwards: directly calls getConfig

    % Copyright 2026 The MathWorks, Inc.

    properties (GetAccess = private, SetAccess = immutable)
        ConnectorFacade(1, 1) mcpcoreserver.internal.facade.connector.ConnectorFacade = mcpcoreserver.internal.facade.connector.DefaultConnectorFacade()
        MATLABFacade(1, 1) mcpcoreserver.internal.facade.matlab.MATLABFacade = mcpcoreserver.internal.facade.matlab.DefaultMATLABFacade()
        SecureGenerator(1, 1) mcpcoreserver.internal.connector.internal.apikeymanager.internal.securegen.SecureGenerator = mcpcoreserver.internal.connector.internal.apikeymanager.internal.securegen.DefaultSecureGenerator()
    end

    methods
        function obj = DefaultAPIKeyManager(options)
            arguments
                options.?mcpcoreserver.internal.connector.internal.apikeymanager.DefaultAPIKeyManager
            end

            for prop = string(fieldnames(options).')
                obj.(prop) = options.(prop);
            end
        end

        function apiKey = getAPIKey(obj)
            % Prior to R2023a: not supported
            if obj.MATLABFacade.isMATLABReleaseOlderThan("R2023a")
                throw(mcpcoreserver.internal.error.Errors.UnsupportedMATLABVersion());
            end

            % Auto-generated key or key was previously set
            % This should true for R2025a onwards
            apiKey = obj.ConnectorFacade.getConfig("apiKey");
            if apiKey ~= ""
                return;
            end

            % R2023a to R2024b inclusive: generate and set secure key
            key = obj.SecureGenerator.generateKey();
            obj.ConnectorFacade.setConfig("apiKey", key);
            apiKey = obj.ConnectorFacade.getConfig("apiKey");
        end
    end

end
