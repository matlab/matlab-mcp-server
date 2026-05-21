classdef DefaultWindowsACLManagerTest < matlab.mock.TestCase
%DefaultWindowsACLManagerTest Tests for DefaultWindowsACLManager
%   This test class verifies that DefaultWindowsACLManager correctly
%   manages Windows ACL operations using .NET Security APIs.
%   Tests require Windows and are skipped on other platforms.

    % Copyright 2026 The MathWorks, Inc.

    properties
        TestFolder
    end

    methods (TestClassSetup)
        function checkPlatform(testCase)
            testCase.assumeTrue(ispc(), "Test requires Windows");
        end
    end

    methods (TestMethodSetup)
        function setupTempFolder(testCase)
            tempFixture = testCase.applyFixture( ...
                matlab.unittest.fixtures.TemporaryFolderFixture ...
            );
            testCase.TestFolder = tempFixture.Folder;
        end
    end

    methods (Test)
        function testDefaultWindowsACLManager_GetCurrentUserSid_HappyPath(testCase)
            % Arrange
            manager = mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.DefaultWindowsACLManager();

            % Act
            sid = manager.getCurrentUserSid();

            % Assert
            testCase.verifyTrue(startsWith(sid, "S-1-5-"), ...
                "SID should be a valid Windows SID starting with S-1-5-");
        end

        function testDefaultWindowsACLManager_GetCurrentUserSid_ReturnsCachedValue(testCase)
            % Arrange
            manager = mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.DefaultWindowsACLManager();

            % Act
            sid1 = manager.getCurrentUserSid();
            sid2 = manager.getCurrentUserSid();

            % Assert
            testCase.verifyEqual(sid1, sid2, ...
                "Subsequent calls should return the same cached SID");
        end

        function testDefaultWindowsACLManager_GetAllowedSids_HappyPath(testCase)
            % Arrange
            manager = mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.DefaultWindowsACLManager();

            % Act
            sids = manager.getAllowedSids(testCase.TestFolder);

            % Assert
            testCase.verifyNotEmpty(sids, ...
                "Temp folder should have at least one allowed SID");
            for i = 1:length(sids)
                testCase.verifyTrue(startsWith(sids(i), "S-1-"), ...
                    "Each SID should be a full Windows SID (S-1-...)");
            end
        end

        function testDefaultWindowsACLManager_IsDaclProtected_DefaultFolder(testCase)
            % Arrange
            manager = mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.DefaultWindowsACLManager();

            % Act
            tf = manager.isDaclProtected(testCase.TestFolder);

            % Assert
            testCase.verifyFalse(tf, ...
                "Default temp folder should inherit permissions (DACL not protected)");
        end

        function testDefaultWindowsACLManager_SetProtectedAcl_Directory(testCase)
            % Arrange
            manager = mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.DefaultWindowsACLManager();
            userSid = manager.getCurrentUserSid();
            sddlSids = [userSid, "SY", "BA"];
            expectedSids = [userSid, "S-1-5-18", "S-1-5-32-544"];

            % Act
            manager.setProtectedAcl(testCase.TestFolder, sddlSids, true);

            % Assert
            testCase.verifyTrue(manager.isDaclProtected(testCase.TestFolder), ...
                "DACL should be protected after setProtectedAcl");
            allowedSids = manager.getAllowedSids(testCase.TestFolder);
            testCase.verifyEqual(sort(allowedSids), sort(expectedSids), ...
                "Only the specified SIDs should have access");
        end

        function testDefaultWindowsACLManager_SetProtectedAcl_File(testCase)
            % Arrange
            manager = mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.DefaultWindowsACLManager();
            userSid = manager.getCurrentUserSid();
            sddlSids = [userSid, "SY", "BA"];
            expectedSids = [userSid, "S-1-5-18", "S-1-5-32-544"];
            testFile = fullfile(testCase.TestFolder, "test.txt");
            writelines("test content", testFile);

            % Act
            manager.setProtectedAcl(testFile, sddlSids, false);

            % Assert
            testCase.verifyTrue(manager.isDaclProtected(testFile), ...
                "DACL should be protected after setProtectedAcl");
            allowedSids = manager.getAllowedSids(testFile);
            testCase.verifyEqual(sort(allowedSids), sort(expectedSids), ...
                "Only the specified SIDs should have access");
        end

        function testDefaultWindowsACLManager_SetProtectedAcl_Directory_ChildInheritsAcl(testCase)
            % Arrange
            manager = mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.DefaultWindowsACLManager();
            userSid = manager.getCurrentUserSid();
            sddlSids = [userSid, "SY", "BA"];
            expectedSids = [userSid, "S-1-5-18", "S-1-5-32-544"];

            manager.setProtectedAcl(testCase.TestFolder, sddlSids, true);

            % Act — create a child file after the folder's ACL is set
            childFile = fullfile(testCase.TestFolder, "child.txt");
            writelines("test content", childFile);

            % Assert — child should inherit the parent's SIDs via OICI flags
            childSids = manager.getAllowedSids(childFile);
            testCase.verifyEqual(sort(childSids), sort(expectedSids), ...
                "Child file should inherit the parent folder's SIDs via OICI inheritance");
        end

        function testDefaultWindowsACLManager_SetProtectedAcl_SingleSid(testCase)
            % Arrange
            manager = mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.DefaultWindowsACLManager();
            userSid = manager.getCurrentUserSid();

            % Act
            manager.setProtectedAcl(testCase.TestFolder, userSid, true);

            % Assert
            testCase.verifyTrue(manager.isDaclProtected(testCase.TestFolder), ...
                "DACL should be protected after setProtectedAcl");
            allowedSids = manager.getAllowedSids(testCase.TestFolder);
            testCase.verifyEqual(allowedSids, userSid, ...
                "Only the current user SID should have access");
        end
    end

end
