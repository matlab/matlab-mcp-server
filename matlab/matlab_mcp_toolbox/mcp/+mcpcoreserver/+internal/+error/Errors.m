classdef Errors
    %Errors Error catalog for mcpcoreserver
    %   This class provides static methods that return MException objects
    %   for known error conditions.

    % Copyright 2026 The MathWorks, Inc.

    methods (Static)
        function ex = UnsupportedOS()
            %UnsupportedOS Create an error for unsupported operating systems
            ex = MException("mcpcoreserver:UnsupportedOS", ...
                "Unsupported operating system. Operating system must be Windows, Linux, or macOS.");
        end

        function ex = FailedToCreateDirectory(path, msg)
            %FailedToCreateDirectory Create an error for directory creation failure
            ex = MException("mcpcoreserver:FailedToCreateDirectory", ...
                "Failed to create directory '%s': %s", path, msg);
        end

        function ex = FailedToGetFileAttributes(path)
            %FailedToGetFileAttributes Create an error for file attribute retrieval failure
            ex = MException("mcpcoreserver:FailedToGetFileAttributes", ...
                "Failed to get file attributes for '%s'", path);
        end

        function ex = FailedToSetPermissions(path)
            %FailedToSetPermissions Create an error for permission setting failure
            ex = MException("mcpcoreserver:FailedToSetPermissions", ...
                "Failed to set permissions on '%s'", path);
        end

        function ex = InsecurePermissions(path)
            %InsecurePermissions Create an error for insecure permissions
            ex = MException("mcpcoreserver:InsecurePermissions", ...
                "Insecure permissions on '%s': Access must be restricted to the user only.", path);
        end

        function ex = UnsupportedMATLABVersion()
            %UnsupportedMATLABVersion Create an error for unsupported MATLAB version
            ex = MException("mcpcoreserver:UnsupportedMATLABVersion", ...
                "MCP server requires MATLAB R2023a or later.");
        end

        function ex = MissingEnvironmentVariable(varName)
            %MissingEnvironmentVariable Create an error for missing environment variable
            ex = MException("mcpcoreserver:MissingEnvironmentVariable", ...
                "Required environment variable '%s' is not set", varName);
        end

        function ex = SecureKeyGenerationFailed()
            %SecureKeyGenerationFailed Create an error for secure key generation failure
            ex = MException("mcpcoreserver:SecureKeyGenerationFailed", ...
                "Failed to generate secure key.");
        end

        function ex = FileExistsAtFolderPath(path)
            %FileExistsAtFolderPath Create an error when a file exists where a folder is expected
            ex = MException("mcpcoreserver:FileExistsAtFolderPath", ...
                "A file already exists at the expected folder path '%s'", path);
        end

        function ex = FolderExistsAtFilePath(path)
            %FolderExistsAtFilePath Create an error when a folder exists where a file is expected
            ex = MException("mcpcoreserver:FolderExistsAtFilePath", ...
                "A folder already exists at the expected file path '%s'", path);
        end
    end

end
