classdef DefaultMATLABFacade < mcpcoreserver.internal.facade.matlab.MATLABFacade
    %DefaultMATLABFacade Default implementation delegating to MATLAB built-ins

    % Copyright 2026 The MathWorks, Inc.

    methods
        function tf = isMATLABReleaseOlderThan(~, release)
            tf = isMATLABReleaseOlderThan(release);
        end
    end

end
