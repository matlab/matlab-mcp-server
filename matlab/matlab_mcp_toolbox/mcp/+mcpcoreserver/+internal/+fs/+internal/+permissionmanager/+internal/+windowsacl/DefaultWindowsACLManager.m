classdef DefaultWindowsACLManager < handle & mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.WindowsACLManager
    %DefaultWindowsACLManager Default implementation using .NET Security APIs
    %   Uses .NET System.Security.AccessControl for Windows ACL operations.
    %   Works entirely with SIDs via SDDL format. All operations run
    %   in-process with no subprocess overhead.
    %
    %   The current user SID is cached (never changes during a MATLAB session).

    % Copyright 2026 The MathWorks, Inc.

    properties (Access = private)
        CachedUserSid string = string.empty
    end

    methods
        function sid = getCurrentUserSid(obj)
            %getCurrentUserSid Get the SID of the current process user (cached)
            if ~isempty(obj.CachedUserSid)
                sid = obj.CachedUserSid;
                return;
            end

            identity = System.Security.Principal.WindowsIdentity.GetCurrent();
            obj.CachedUserSid = string(identity.User.ToString());
            sid = obj.CachedUserSid;
        end

        function sids = getAllowedSids(~, path)
            %getAllowedSids Get SIDs of all Allow ACEs on a path
            %   Uses .NET AccessRules API to enumerate all access rules,
            %   returning full SID strings (S-1-...) for every Allow ACE.
            accessSections = System.Security.AccessControl.AccessControlSections.Access;

            if isfolder(path)
                security = System.Security.AccessControl.DirectorySecurity( ...
                    char(path), accessSections);
            else
                security = System.Security.AccessControl.FileSecurity( ...
                    char(path), accessSections);
            end

            sidType = System.Type.GetType('System.Security.Principal.SecurityIdentifier');
            rules = security.GetAccessRules(true, true, sidType);

            sids = strings(1, rules.Count);
            n = 0;
            for i = 0:rules.Count-1
                rule = rules.Item(i);
                if rule.AccessControlType == System.Security.AccessControl.AccessControlType.Allow
                    n = n + 1;
                    sids(n) = string(rule.IdentityReference.ToString());
                end
            end
            sids = sids(1:n);
        end

        function tf = isDaclProtected(~, path)
            %isDaclProtected Check if the DACL is protected (inheritance blocked)
            %   Returns true if the SDDL contains "D:P", meaning the DACL
            %   is protected and does not inherit ACEs from parent folders.
            accessSections = System.Security.AccessControl.AccessControlSections.Access;

            if isfolder(path)
                security = System.Security.AccessControl.DirectorySecurity( ...
                    char(path), accessSections);
            else
                security = System.Security.AccessControl.FileSecurity( ...
                    char(path), accessSections);
            end

            sddl = string(security.GetSecurityDescriptorSddlForm(accessSections));
            tf = startsWith(sddl, "D:P");
        end

        function setProtectedAcl(~, path, sids, isDirectory)
            %setProtectedAcl Set a protected ACL with FullControl for the given SIDs
            %   Builds an SDDL string and applies it. D:P = Protected DACL
            %   (blocks inheritance from parent).

            % Build SDDL
            % "P" = Protected: blocks ACL inheritance from parent, ensuring only the explicitly listed SIDs have access.
            sddl = "D:P";
            for i = 1:length(sids)
                if isDirectory
                    % A;OICI;FA = Allow, ObjectInherit+ContainerInherit, FullAccess
                    sddl = sddl + sprintf("(A;OICI;FA;;;%s)", sids(i));
                else
                    % A;;FA = Allow, no inheritance flags, FullAccess
                    sddl = sddl + sprintf("(A;;FA;;;%s)", sids(i));
                end
            end

            % Apply SDDL to the path
            accessSections = System.Security.AccessControl.AccessControlSections.Access;
            if isDirectory
                security = System.Security.AccessControl.DirectorySecurity();
                security.SetSecurityDescriptorSddlForm(char(sddl), accessSections);
                dirInfo = System.IO.DirectoryInfo(char(path));
                dirInfo.SetAccessControl(security);
            else
                security = System.Security.AccessControl.FileSecurity();
                security.SetSecurityDescriptorSddlForm(char(sddl), accessSections);
                fileInfo = System.IO.FileInfo(char(path));
                fileInfo.SetAccessControl(security);
            end
        end
    end

end
