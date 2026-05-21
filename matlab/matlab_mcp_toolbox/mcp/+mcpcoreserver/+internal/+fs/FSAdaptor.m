classdef (Abstract) FSAdaptor
    %FSAdaptor Adaptor for filesystem operations with security features
    %   This class provides higher-level operations built on top of the
    %   filesystem facade and permission manager.

    % Copyright 2026 The MathWorks, Inc.

    methods (Abstract)
        ensureSecureFolder(obj, folderPath)
        ensureSecureFile(obj, filePath)
    end

end
