classdef (Abstract, HandleCompatible) WindowsACLManager
%WindowsACLManager Abstract interface for Windows ACL operations using SIDs
%   This abstract class defines the interface for Windows access control
%   operations. All methods use Security Identifiers (SIDs) for reliable
%   cross-domain, locale-independent identity resolution.

    % Copyright 2026 The MathWorks, Inc.

    methods (Abstract)
        sid = getCurrentUserSid(obj)
        sids = getAllowedSids(obj, path)
        tf = isDaclProtected(obj, path)
        setProtectedAcl(obj, path, sids, isDirectory)
    end

end
