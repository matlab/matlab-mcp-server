classdef (Abstract) FSFacade
    %FSFacade Abstract facade for filesystem operations
    %   This abstract class defines the interface for filesystem operations.

    % Copyright 2026 The MathWorks, Inc.

    methods (Abstract)
        varargout = mkdir(obj, varargin)
        tf = isfolder(obj, path)
        tf = isfile(obj, path)
        varargout = fileattrib(obj, varargin)
        writelines(obj, text, filename)
        varargout = fopen(obj, varargin)
        fclose(obj, fileID)
        varargout = fread(obj, fileID, varargin)
    end

end
