classdef (Abstract) shareMATLABSessionTestBase < matlab.unittest.TestCase
    %shareMATLABSessionTestBase Base class for shareMATLABSession system tests
    %   This abstract class provides common helper methods for platform-specific
    %   system tests.

    % Copyright 2026 The MathWorks, Inc.

    properties (GetAccess=protected, SetAccess=private)
        TestBaseFolder
    end

    methods (Abstract, TestClassSetup)
        checkPlatform(testCase)
    end

    properties (Constant)
        MLTBX_DIR_ENV_VAR_NAME = "MATLAB_MCP_CORE_SERVER_MLTBX_DIR"
    end

    methods (TestClassSetup)
        function installMLTBX(testCase)
            mltbxDir = getenv(testCase.MLTBX_DIR_ENV_VAR_NAME);
            testCase.fatalAssertNotEmpty(mltbxDir, sprintf("environment variable %s is not set", testCase.MLTBX_DIR_ENV_VAR_NAME))

            mltbxPath = fullfile(mltbxDir, "MATLABMCPCoreServerToolbox.mltbx");
            testCase.fatalAssertEqual(exist(mltbxPath, "file"), 2, sprintf("MLTBX file does not exist at path: %s", mltbxPath))

            tbx = matlab.addons.toolbox.installToolbox(mltbxPath, true);
            testCase.addTeardown(@() matlab.addons.toolbox.uninstallToolbox(tbx))
        end
    end

    methods (TestMethodSetup)
        function setupEnvironment(testCase)
            % Create temporary folder with automatic cleanup
            tempFixture = testCase.applyFixture( ...
                matlab.unittest.fixtures.TemporaryFolderFixture ...
            );
            testCase.TestBaseFolder = tempFixture.Folder;

            % Set HOME environment variable with automatic restoration
            testCase.applyFixture( ...
                matlab.unittest.fixtures.EnvironmentVariableFixture( ...
                    testCase.getRootEnvVar(), ...
                    testCase.TestBaseFolder ...
                ) ...
            );
        end
    end

    methods (Abstract, Access = protected)
        envVarName = getRootEnvVar(testCase)
        expectedFolder = getExpectedAppDataFolder(testCase)
        verifyFolderPermissions(testCase, path)
        verifyFilePermissions(testCase, path)
        setupSecureFolders(testCase, appDataFolder, v1Folder)
        createInsecureFile(testCase, filePath)
    end

    methods (Test)
        function testShareMATLABSession_CreatesSecureFiles(testCase)
            % Arrange
            expectedAppDataFolder = testCase.getExpectedAppDataFolder();
            expectedV1Folder = fullfile(expectedAppDataFolder, "v1");
            expectedSessionDetailsPath = fullfile(expectedV1Folder, "sessionDetails.json");

            % Act
            shareMATLABSession();

            % Assert - File exists
            testCase.verifyTrue(isfile(expectedSessionDetailsPath), ...
                "sessionDetails.json should exist");

            % Assert - Check permissions
            testCase.verifyFolderPermissions(expectedAppDataFolder);
            testCase.verifyFolderPermissions(expectedV1Folder);
            testCase.verifyFilePermissions(expectedSessionDetailsPath);

            % Assert - Check file content
            testCase.verifySessionDetailsContent(expectedSessionDetailsPath);
        end

        function testShareMATLABSession_CalledTwice_Succeeds(testCase)
            % Arrange
            expectedAppDataFolder = testCase.getExpectedAppDataFolder();
            expectedV1Folder = fullfile(expectedAppDataFolder, "v1");
            expectedSessionDetailsPath = fullfile(expectedV1Folder, "sessionDetails.json");

            % Act - call twice to verify idempotency
            shareMATLABSession();
            shareMATLABSession();

            % Assert - File exists after second call
            testCase.verifyTrue(isfile(expectedSessionDetailsPath), ...
                "sessionDetails.json should exist after second call");

            % Assert - Permissions still correct after second call
            testCase.verifyFolderPermissions(expectedAppDataFolder);
            testCase.verifyFolderPermissions(expectedV1Folder);
            testCase.verifyFilePermissions(expectedSessionDetailsPath);

            % Assert - Content valid after second call
            testCase.verifySessionDetailsContent(expectedSessionDetailsPath);
        end

        function testShareMATLABSession_SessionDetailsFileExists_InsecurePermissions(testCase)
            % Arrange
            expectedAppDataFolder = testCase.getExpectedAppDataFolder();
            expectedV1Folder = fullfile(expectedAppDataFolder, "v1");
            expectedSessionDetailsPath = fullfile(expectedV1Folder, "sessionDetails.json");

            % Pre-create folders with correct permissions
            testCase.setupSecureFolders(expectedAppDataFolder, expectedV1Folder);

            % Pre-create file with WRONG permissions
            testCase.createInsecureFile(expectedSessionDetailsPath);

            % Act & Assert
            expectedError = mcpcoreserver.internal.error.Errors.InsecurePermissions(expectedSessionDetailsPath);
            testCase.verifyError( ...
                @() shareMATLABSession(), ...
                expectedError.identifier, ...
                "Should throw InsecurePermissions when sessionDetails.json has wrong permissions" ...
            );
        end
    end

    methods (Access = protected)
        function verifySessionDetailsContent(testCase, path)
            %verifySessionDetailsContent Verify sessionDetails.json has valid content

            % Build expected struct using connector methods directly
            expectedDetails = struct( ...
                "port", connector.securePort(), ...
                "certificate", connector.getCertificateLocation(), ...
                "apiKey", connector.internal.getConfig("apiKey"), ...
                "pid", feature("getpid") ...
            );

            % Read actual values from file
            jsonText = fileread(path);
            actualDetails = jsondecode(jsonText);

            % Verify exact match
            testCase.verifyEqual(actualDetails, expectedDetails, ...
                "sessionDetails.json should contain exact connection details");
        end
    end

end
