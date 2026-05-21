% Copyright 2026 The MathWorks, Inc.

info = dir(tmpFile); % tmpFile is created in ./upload_mltbx.m
assert(~isempty(info), 'File not found: %s', tmpFile);
assert(info.bytes == EXPECTED_BYTE_SIZE, 'File size mismatch: expected EXPECTED_BYTE_SIZE, got %d', info.bytes);

clear info;
