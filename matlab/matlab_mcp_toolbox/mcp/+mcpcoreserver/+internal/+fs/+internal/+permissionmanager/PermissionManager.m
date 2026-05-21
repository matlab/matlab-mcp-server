classdef (Abstract) PermissionManager
    %PermissionManager Abstract interface for permission management operations
    %   This abstract class defines the interface for permission management
    %   operations that involve logic beyond simple delegation to MATLAB built-ins.

    % Copyright 2026 The MathWorks, Inc.

    methods (Abstract)
        setPermissionsToUserOnly(obj, path)
        tf = checkPermissionsIsUserOnly(obj, path)
    end

end
