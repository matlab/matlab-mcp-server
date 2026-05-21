classdef DefaultFSAdaptor < mcpcoreserver.internal.fs.FSAdaptor
    %DefaultFSAdaptor Default implementation for filesystem operations with security features
    %   This class provides higher-level operations built on top of the
    %   filesystem facade and permission manager.

    % Copyright 2026 The MathWorks, Inc.

    properties (GetAccess = private, SetAccess = immutable)
        FSFacade(1, 1) mcpcoreserver.internal.facade.fs.FSFacade = mcpcoreserver.internal.facade.fs.DefaultFSFacade()
        PermissionManager(1, 1) mcpcoreserver.internal.fs.internal.permissionmanager.PermissionManager = mcpcoreserver.internal.fs.internal.permissionmanager.DefaultPermissionManager()
    end

    methods
        function obj = DefaultFSAdaptor(options)
            arguments
                options.?mcpcoreserver.internal.fs.DefaultFSAdaptor
            end

            for prop = string(fieldnames(options).')
                obj.(prop) = options.(prop);
            end
        end

        function ensureSecureFolder(obj, folderPath)
            if obj.FSFacade.isfolder(folderPath)
                if ~obj.PermissionManager.checkPermissionsIsUserOnly(folderPath)
                    throw(mcpcoreserver.internal.error.Errors.InsecurePermissions(folderPath));
                end
            else
                if obj.FSFacade.isfile(folderPath)
                    throw(mcpcoreserver.internal.error.Errors.FileExistsAtFolderPath(folderPath));
                end
                [status, msg] = obj.FSFacade.mkdir(folderPath);
                if ~status
                    throw(mcpcoreserver.internal.error.Errors.FailedToCreateDirectory(folderPath, msg));
                end
                obj.PermissionManager.setPermissionsToUserOnly(folderPath);
            end
        end

        function ensureSecureFile(obj, filePath)
            if obj.FSFacade.isfile(filePath)
                if ~obj.PermissionManager.checkPermissionsIsUserOnly(filePath)
                    throw(mcpcoreserver.internal.error.Errors.InsecurePermissions(filePath));
                end
            else
                if obj.FSFacade.isfolder(filePath)
                    throw(mcpcoreserver.internal.error.Errors.FolderExistsAtFilePath(filePath));
                end
                obj.FSFacade.writelines("", filePath);
                obj.PermissionManager.setPermissionsToUserOnly(filePath);
            end
        end
    end

end
