classdef (Abstract) MATLABFacade
    %MATLABFacade Abstract facade for MATLAB version detection
    %   This abstract class defines the interface for MATLAB version
    %   detection operations.

    % Copyright 2026 The MathWorks, Inc.

    methods (Abstract)
        tf = isMATLABReleaseOlderThan(obj, release)
    end

end
