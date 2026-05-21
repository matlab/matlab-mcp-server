% Copyright 2026 The MathWorks, Inc.

tmpFile = [tempname '.mltbx'];

fid = fopen(tmpFile, 'w');
assert(fid ~= -1, 'Failed to open file');

try
    fwrite(fid, matlab.net.base64decode('MTBX_BINARY_CONTENT_AS_BASE_ENCODED_STRING'), 'uint8');
    fclose(fid);
catch ME
    fclose(fid);
    rethrow(ME)
end

clear fid;
