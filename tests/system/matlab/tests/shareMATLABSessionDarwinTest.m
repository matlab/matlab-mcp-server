classdef shareMATLABSessionDarwinTest < shareMATLABSessionTestBase
    %shareMATLABSessionDarwinTest System tests for shareMATLABSession on macOS
    %   This test class verifies that shareMATLABSession correctly creates
    %   the session details file with proper permissions on macOS systems.

    % Copyright 2026 The MathWorks, Inc.

    methods (TestClassSetup)
        function checkPlatform(testCase)
            testCase.assumeTrue(ismac(), "Test requires macOS");
        end
    end

    methods (Access = protected)
        function envVarName = getRootEnvVar(~)
            envVarName = "HOME";
        end

        function expectedFolder = getExpectedAppDataFolder(testCase)
            expectedFolder = fullfile(testCase.TestBaseFolder, "Library", "Application Support", "MathWorks", "MATLAB MCP Core Server");
        end

        function verifyFolderPermissions(testCase, path)
            [status, result] = system(sprintf('stat -f "%%Lp" "%s"', path));
            testCase.verifyEqual(status, 0, ...
                sprintf("stat command should succeed for folder: %s", path));
            actualPerms = string(strip(result));
            testCase.verifyEqual(actualPerms, "700", ...
                sprintf("Folder %s should have 700 permissions", path));
        end

        function verifyFilePermissions(testCase, path)
            [status, result] = system(sprintf('stat -f "%%Lp" "%s"', path));
            testCase.verifyEqual(status, 0, ...
                sprintf("stat command should succeed for file: %s", path));
            actualPerms = string(strip(result));
            testCase.verifyEqual(actualPerms, "600", ...
                sprintf("File %s should have 600 permissions", path));
        end

        function setupSecureFolders(~, appDataFolder, v1Folder)
            mkdir(v1Folder);
            system(sprintf('chmod 700 "%s"', appDataFolder));
            system(sprintf('chmod 700 "%s"', v1Folder));
        end

        function createInsecureFile(~, filePath)
            writelines("{}", filePath);
            system(sprintf('chmod 644 "%s"', filePath));
        end
    end

end
