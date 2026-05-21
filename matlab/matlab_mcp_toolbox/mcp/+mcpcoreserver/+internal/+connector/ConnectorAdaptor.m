classdef (Abstract) ConnectorAdaptor
    %ConnectorAdaptor Adaptor for MATLAB connector operations
    %   This class provides higher-level operations built on top of the
    %   connector facade.

    % Copyright 2026 The MathWorks, Inc.

    methods (Abstract)
        details = getConnectionDetails(obj)
    end

end
