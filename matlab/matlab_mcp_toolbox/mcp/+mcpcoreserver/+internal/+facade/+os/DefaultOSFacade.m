classdef DefaultOSFacade < mcpcoreserver.internal.facade.os.OSFacade
    %DefaultOSFacade Default implementation for operating system detection
    %   This class delegates to MATLAB's built-in ismac, ispc, and isunix
    %   functions.

    % Copyright 2026 The MathWorks, Inc.

    methods
        function tf = ismac(~)
            tf = ismac();
        end

        function tf = ispc(~)
            tf = ispc();
        end

        function tf = isunix(~)
            tf = isunix();
        end

        function [status, result] = system(~, command)
            [status, result] = system(command);
        end

        function varargout = feature(~, varargin)
            [varargout{1:nargout}] = feature(varargin{:});
        end

        function value = getenv(~, name)
            value = getenv(name);
        end
    end

end
