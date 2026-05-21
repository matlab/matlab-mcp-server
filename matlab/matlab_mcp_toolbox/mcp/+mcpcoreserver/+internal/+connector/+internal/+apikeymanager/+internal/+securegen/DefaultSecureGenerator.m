classdef DefaultSecureGenerator < mcpcoreserver.internal.connector.internal.apikeymanager.internal.securegen.SecureGenerator
    %DefaultSecureGenerator Default implementation for cryptographically secure key generation
    %   This class generates cryptographically secure random keys using
    %   OS-level CSPRNGs: /dev/urandom on Linux/macOS and
    %   .NET RandomNumberGenerator on Windows.

    % Copyright 2026 The MathWorks, Inc.

    properties (GetAccess = private, SetAccess = immutable)
        OSFacade(1, 1) mcpcoreserver.internal.facade.os.OSFacade = mcpcoreserver.internal.facade.os.DefaultOSFacade()
        FSFacade(1, 1) mcpcoreserver.internal.facade.fs.FSFacade = mcpcoreserver.internal.facade.fs.DefaultFSFacade()
    end

    properties (Constant, Access = private)
        KeySize(1,1) uint8 = 24
    end

    methods
        function obj = DefaultSecureGenerator(options)
            arguments
                options.?mcpcoreserver.internal.connector.internal.apikeymanager.internal.securegen.DefaultSecureGenerator
            end

            for prop = string(fieldnames(options).')
                obj.(prop) = options.(prop);
            end
        end

        function key = generateKey(obj)
            if obj.OSFacade.ispc()
                bytes = obj.generateKeyBytesOnWindows();
            elseif obj.OSFacade.ismac() || obj.OSFacade.isunix()
                bytes = obj.generateKeyBytesOnUnix();
            else
                throw(mcpcoreserver.internal.error.Errors.UnsupportedOS());
            end
            hex = string(dec2hex(bytes, 2));
            key = lower(strjoin(hex, ""));
        end
    end

    methods (Access=private)
        function bytes = generateKeyBytesOnUnix(obj)
            fid = obj.FSFacade.fopen(fullfile("/", "dev", "urandom"), "r");
            if fid == -1
                throw(mcpcoreserver.internal.error.Errors.SecureKeyGenerationFailed());
            end
            closeFile = onCleanup(@() obj.FSFacade.fclose(fid));
            bytes = obj.FSFacade.fread(fid, obj.KeySize, "*uint8");
        end

        function bytes = generateKeyBytesOnWindows(obj)
            rng = System.Security.Cryptography.RandomNumberGenerator.Create();
            disposeOfRng = onCleanup(@() rng.Dispose());
            bytes = NET.createArray('System.Byte', obj.KeySize);
            rng.GetBytes(bytes);
            bytes = uint8(bytes);
        end
    end

end
