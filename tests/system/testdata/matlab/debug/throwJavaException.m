function throwJavaException()
    % Copyright 2026 The MathWorks, Inc.

    javaState = 12; %#ok<NASGU>
    import java.net.ServerSocket
    import java.net.Socket

    % Close a port the OS just handed out, so nothing listens on it and the connection is refused.
    server = ServerSocket(0);
    closedPort = server.getLocalPort();
    server.close();
    socket = Socket('localhost', closedPort); %#ok<NASGU>
end
