classdef DefaultFSFacade < mcpcoreserver.internal.facade.fs.FSFacade
    %DefaultFSFacade Default implementation for filesystem operations
    %   This class delegates to MATLAB's built-in filesystem functions.

    % Copyright 2026 The MathWorks, Inc.

    methods
        function varargout = mkdir(~, varargin)
            [varargout{1:nargout}] = mkdir(varargin{:});
        end

        function tf = isfolder(~, path)
            tf = isfolder(path);
        end

        function tf = isfile(~, path)
            tf = isfile(path);
        end

        function varargout = fileattrib(~, varargin)
            [varargout{1:nargout}] = fileattrib(varargin{:});
        end

        function writelines(~, text, filename)
            writelines(text, filename);
        end

        function varargout = fopen(~, varargin)
            [varargout{1:nargout}] = fopen(varargin{:});
        end

        function fclose(~, fileID)
            fclose(fileID);
        end

        function varargout = fread(~, fileID, varargin)
            [varargout{1:nargout}] = fread(fileID, varargin{:});
        end
    end

end
