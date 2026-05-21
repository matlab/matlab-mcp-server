classdef DefaultPermissionManager < mcpcoreserver.internal.fs.internal.permissionmanager.PermissionManager
    %DefaultPermissionManager Default implementation for permission management
    %   This class provides platform-aware permission management operations
    %   for files and folders.
    %
    %   On Windows, exactly 3 ACEs are set (current user, SYSTEM,
    %   Administrators), all with FullControl, and inheritance is blocked.
    %   SIDs are used for identity resolution to avoid domain-name ambiguity.
    %
    %   On Unix/macOS, chmod 700/600 is used for directories/files.

    % Copyright 2026 The MathWorks, Inc.

    properties (Constant, Access = private)
        CHMOD_FOLDER_PERMISSIONS = "700"
        CHMOD_FILE_PERMISSIONS = "600"
        % Well-known SIDs matching C++ WinLocalSystemSid and WinBuiltinAdministratorsSid
        SYSTEM_SID = "S-1-5-18"
        ADMINISTRATORS_SID = "S-1-5-32-544"
    end

    properties (GetAccess = private, SetAccess = immutable)
        FSFacade(1, 1) mcpcoreserver.internal.facade.fs.FSFacade = mcpcoreserver.internal.facade.fs.DefaultFSFacade()
        OSFacade(1, 1) mcpcoreserver.internal.facade.os.OSFacade = mcpcoreserver.internal.facade.os.DefaultOSFacade()
        WindowsACLManager(1, 1) mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.WindowsACLManager = mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.DefaultWindowsACLManager()
    end

    methods
        function obj = DefaultPermissionManager(options)
            arguments
                options.?mcpcoreserver.internal.fs.internal.permissionmanager.DefaultPermissionManager
            end

            for prop = string(fieldnames(options).')
                obj.(prop) = options.(prop);
            end
        end

        function setPermissionsToUserOnly(obj, path)
            if obj.OSFacade.ispc()
                obj.setPermissionsWindows(path);
            else
                obj.setPermissionsUnix(path);
            end
        end

        function tf = checkPermissionsIsUserOnly(obj, path)
            if obj.OSFacade.ispc()
                tf = obj.checkPermissionsWindows(path);
            else
                tf = obj.checkPermissionsUnix(path);
            end
        end
    end

    methods (Access = private)
        function setPermissionsUnix(obj, path)
            %setPermissionsUnix Set user-only permissions on Unix/macOS

            [status, attribs, ~] = obj.FSFacade.fileattrib(path);
            if ~status
                throw(mcpcoreserver.internal.error.Errors.FailedToGetFileAttributes(path));
            end

            % Use chmod for reliable permission setting on Unix
            if attribs.directory
                mode = obj.CHMOD_FOLDER_PERMISSIONS;  % rwx------
            else
                mode = obj.CHMOD_FILE_PERMISSIONS;  % rw-------
            end

            escapedPath = strrep(path, "'", "'\''");
            cmd = sprintf("chmod %s '%s'", mode, escapedPath);
            [status, ~] = obj.OSFacade.system(cmd);

            if status ~= 0
                throw(mcpcoreserver.internal.error.Errors.FailedToSetPermissions(path));
            end
        end

        function setPermissionsWindows(obj, path)
            %setPermissionsWindows Set user-only permissions on Windows using SID-based ACLs
            %   Grants FullControl to current user, SYSTEM, and Administrators,
            %   and protects the DACL from inheritance.

            isDir = obj.FSFacade.isfolder(path);

            try
                userSid = obj.WindowsACLManager.getCurrentUserSid();
                sids = [userSid, obj.SYSTEM_SID, obj.ADMINISTRATORS_SID];
                obj.WindowsACLManager.setProtectedAcl(path, sids, isDir);
            catch ME
                throw(addCause(mcpcoreserver.internal.error.Errors.FailedToSetPermissions(path), ME));
            end
        end

        function tf = checkPermissionsUnix(obj, path)
            %checkPermissionsUnix Check user-only permissions on Unix/macOS

            [status, attribs, ~] = obj.FSFacade.fileattrib(path);
            if ~status
                throw(mcpcoreserver.internal.error.Errors.FailedToGetFileAttributes(path));
            end

            groupHasAccess = attribs.GroupRead || attribs.GroupWrite || attribs.GroupExecute;
            otherHasAccess = attribs.OtherRead || attribs.OtherWrite || attribs.OtherExecute;
            userHasRead = attribs.UserRead;
            userHasWrite = attribs.UserWrite;
            userHasExecute = attribs.UserExecute;

            % Directories: rwx------ (700)
            % Files: rw------- (600)
            if attribs.directory
                tf = userHasRead && userHasWrite && userHasExecute && ~groupHasAccess && ~otherHasAccess;
            else
                tf = userHasRead && userHasWrite && ~userHasExecute && ~groupHasAccess && ~otherHasAccess;
            end
        end

        function tf = checkPermissionsWindows(obj, path)
            %checkPermissionsWindows Check user-only permissions on Windows using SIDs
            %   Verifies the DACL is protected (no inherited ACEs) and only
            %   current user, SYSTEM, and Administrators have access.
            %   Uses SIDs throughout for locale-independent comparison.

            try
                userSid = obj.WindowsACLManager.getCurrentUserSid();
                allowedSids = obj.WindowsACLManager.getAllowedSids(path);
                daclProtected = obj.WindowsACLManager.isDaclProtected(path);
            catch
                throw(mcpcoreserver.internal.error.Errors.FailedToGetFileAttributes(path));
            end

            if ~daclProtected
                tf = false;
                return;
            end

            trustedSids = [userSid, obj.SYSTEM_SID, obj.ADMINISTRATORS_SID];

            % Verify all 3 trusted SIDs have access
            allTrustedHaveAccess = true;
            for i = 1:length(trustedSids)
                if ~any(allowedSids == trustedSids(i))
                    allTrustedHaveAccess = false;
                    break;
                end
            end

            % Verify no unexpected SIDs have access
            hasUnexpectedSid = false;
            for i = 1:length(allowedSids)
                if ~any(allowedSids(i) == trustedSids)
                    hasUnexpectedSid = true;
                    break;
                end
            end

            tf = allTrustedHaveAccess && ~hasUnexpectedSid;
        end
    end

end
