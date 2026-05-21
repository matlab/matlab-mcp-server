classdef (Abstract) OSFacade
    %OSFacade Abstract facade for operating system detection
    %   This abstract class defines the interface for detecting the
    %   operating system.

    % Copyright 2026 The MathWorks, Inc.

    methods (Abstract)
        tf = ismac(obj)
        tf = ispc(obj)
        tf = isunix(obj)
        [status, result] = system(obj, command)
        varargout = feature(obj, varargin)
        value = getenv(obj, name)
    end

end
