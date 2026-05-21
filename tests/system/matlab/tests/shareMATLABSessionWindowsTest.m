classdef shareMATLABSessionWindowsTest < shareMATLABSessionTestBase
    %shareMATLABSessionWindowsTest System tests for shareMATLABSession on Windows
    %   This test class verifies that shareMATLABSession correctly creates
    %   the session details file with proper permissions on Windows systems.
    %   Expects 3 ACEs: current user, SYSTEM, and Administrators (matching
    %   the C++ UnixDomainSocketFolderPermissionsImpl security model).

    % Copyright 2026 The MathWorks, Inc.

    methods (TestClassSetup)
        function checkPlatform(testCase)
            testCase.assumeTrue(ispc(), "Test requires Windows");
        end
    end

    methods (Access = protected)
        function envVarName = getRootEnvVar(~)
            envVarName = "APPDATA";
        end

        function expectedFolder = getExpectedAppDataFolder(testCase)
            expectedFolder = fullfile(testCase.TestBaseFolder, "MathWorks", "MATLAB MCP Core Server");
        end

        function verifyFolderPermissions(testCase, path)
            testCase.verifyWindowsSecurePermissions(path);
        end

        function verifyFilePermissions(testCase, path)
            testCase.verifyWindowsSecurePermissions(path);
        end

        function verifyWindowsSecurePermissions(testCase, path)
            %verifyWindowsSecurePermissions Verify path has 3 ACEs: user, SYSTEM, Administrators
            [identities, aceFlags] = testCase.parseIcaclsOutput(path);

            username = getenv("USERNAME");
            userdomain = getenv("USERDOMAIN");
            qualifiedUsername = userdomain + "\" + username;

            aceCount = length(identities);
            userHasFullControl = false;
            systemHasFullControl = false;
            adminsHasFullControl = false;
            hasInheritedPermissions = false;

            for i = 1:aceCount
                identity = identities(i);
                flags = aceFlags(i);

                if contains(flags, "(I)")
                    hasInheritedPermissions = true;
                end

                hasFullControl = contains(flags, "(F)");

                if (strcmpi(identity, username) || strcmpi(identity, qualifiedUsername)) && hasFullControl
                    userHasFullControl = true;
                elseif strcmpi(identity, "NT AUTHORITY\SYSTEM") && hasFullControl
                    systemHasFullControl = true;
                elseif strcmpi(identity, "BUILTIN\Administrators") && hasFullControl
                    adminsHasFullControl = true;
                end
            end

            testCase.verifyTrue(userHasFullControl, ...
                sprintf("Path %s should have full control for current user", path));
            testCase.verifyTrue(systemHasFullControl, ...
                sprintf("Path %s should have full control for SYSTEM", path));
            testCase.verifyTrue(adminsHasFullControl, ...
                sprintf("Path %s should have full control for Administrators", path));
            testCase.verifyEqual(aceCount, 3, ...
                sprintf("Path %s should have exactly 3 ACEs, found %d", path, aceCount));
            testCase.verifyFalse(hasInheritedPermissions, ...
                sprintf("Path %s should not have inherited permissions", path));
        end

        function [identities, aceFlags] = parseIcaclsOutput(testCase, path)
            %parseIcaclsOutput Run icacls and parse the output into identities and flags
            %   Returns two string arrays: identities (e.g., "NT AUTHORITY\SYSTEM")
            %   and aceFlags (e.g., ":(OI)(CI)(F)"). Strips the queried path from
            %   line 1 of icacls output to avoid ambiguous parsing.
            [status, result] = system(sprintf('icacls "%s"', path));
            testCase.verifyEqual(status, 0, ...
                sprintf("icacls command should succeed for path: %s", path));

            % Strip the path from output (icacls puts it on line 1 with the first ACE)
            result = strrep(result, path, '');

            lines = splitlines(strtrim(result));
            identities = string.empty;
            aceFlags = string.empty;
            for idx = 1:length(lines)
                line = strtrim(lines{idx});
                if isempty(line) || startsWith(line, "Successfully")
                    continue;
                end
                token = regexp(line, '^(.+?)(:\([^)]+\)(?:\([^)]+\))*)$', 'tokens', 'once');
                if ~isempty(token)
                    identities(end+1) = strtrim(string(token{1})); %#ok<AGROW>
                    aceFlags(end+1) = string(token{2}); %#ok<AGROW>
                end
            end
        end

        function setupSecureFolders(~, appDataFolder, v1Folder)
            mkdir(v1Folder);
            userdomain = getenv("USERDOMAIN");
            username = getenv("USERNAME");
            qualifiedUsername = sprintf('%s\\%s', userdomain, username);
            system(sprintf('icacls "%s" /inheritance:r /grant:r "%s:(OI)(CI)F" /grant:r "NT AUTHORITY\\SYSTEM:(OI)(CI)F" /grant:r "BUILTIN\\Administrators:(OI)(CI)F"', ...
                appDataFolder, qualifiedUsername));
            system(sprintf('icacls "%s" /inheritance:r /grant:r "%s:(OI)(CI)F" /grant:r "NT AUTHORITY\\SYSTEM:(OI)(CI)F" /grant:r "BUILTIN\\Administrators:(OI)(CI)F"', ...
                v1Folder, qualifiedUsername));
        end

        function createInsecureFile(~, filePath)
            userdomain = getenv("USERDOMAIN");
            username = getenv("USERNAME");
            qualifiedUsername = sprintf('%s\\%s', userdomain, username);
            fid = fopen(filePath, "w");
            fprintf(fid, "{}");
            fclose(fid);
            system(sprintf('icacls "%s" /inheritance:r /grant:r "%s:F"', filePath, qualifiedUsername));
            system(sprintf('icacls "%s" /grant "Everyone:R"', filePath));
        end
    end

end
