classdef (Abstract) SecureGenerator
    %SecureGenerator Abstract generator for cryptographically secure keys
    %   This abstract class defines the interface for generating
    %   cryptographically secure random keys using OS-level CSPRNGs.

    % Copyright 2026 The MathWorks, Inc.

    methods (Abstract)
        key = generateKey(obj)
    end

end
