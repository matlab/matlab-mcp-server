% IMPORTANT NOTICE:
% This file may contain calls to MathWorks internal APIs which are subject to
% change without any prior notice. Usage of these undocumented APIs outside of
% these files is not supported.

function results = mcpEval(code)
    % mcpEval A helper function for handling execution of MATLAB code and post-processing
    % the outputs. The MATLAB MCP Server will then convert those to the appropriate MCP Server Tool Content, see:
    % 
    % https://modelcontextprotocol.io/specification/2025-06-18/server/tools#tool-result
    % 
    % We use the Live Editor API for majority of the work.
    %
    % The entire MATLAB code given by user is treated as code within a single cell
    % of a unique Live Script. Hence, each execution request can be considered as
    % creating and running a new Live Script file.
        
    % This is largely a re-use of:
    % https://github.com/mathworks/jupyter-matlab-proxy/blob/057564dccb7de37f052e709f5380e3ece0b2c4a1/src/jupyter_matlab_kernel/matlab/%2Bjupyter/execute.m#L1

    % Copyright 2025-2026 The MathWorks, Inc.

    % Embed user MATLAB code in a try-catch block for MATLAB versions less than R2022b.
    % This is will disable inbuilt ErrorRecovery mechanism. Any exceptions created in
    % user code would be handled by +matlab_mcp/getOrStashExceptions.m
    if isMATLABReleaseOlderThan("R2022b")
        code = sprintf(['try\n'...
            '%s\n'...
            'catch MCPME\n'...
            'matlab_mcp.getOrStashExceptions(MCPME)\n'...
            'clear MCPME\n'...
            'end'], code);
    end

    fileToShowErrors = 'matlab_mcp_server';

    request = struct( ...
        'requestId', 'matlab_mcp_server',...
        'editorId', 'matlab_mcp_server',...
        'fullText', code,...
        'fullFilePath', fileToShowErrors);

    request = updateRequest(request, code);

    hotlinksPreviousState = feature('hotlinks','off');
    hotlinksCleanupObj = onCleanup(@() feature('hotlinks', hotlinksPreviousState));

    resp = jsondecode(matlab.internal.editor.evaluateSynchronousRequest(request));

    results = jsonencode(processOutputs(resp.outputs));
end

% Helper function to update fields in the request based on MATLAB and LiveEditor
% API version.
function request = updateRequest(request, code)
    % Support for MATLAB version <= R2023a.
    if isMATLABReleaseOlderThan("R2023b")
        request = updateRequestFromBefore23b(request, code);
    else
        % Support for MATLAB version >= R2023b.

        % To maintain backwards compatibility, each case in the switch
        % encodes conversion from the version number in the case
        % to the current version.
        switch matlab.internal.editor.getApiVersion('synchronous')
            case 1
                request = updateRequestFromVersion1(request, code);
            case 2
                request = updateRequestFromVersion2(request, code);
            otherwise
                error("Invalid API version. Create an issue at https://github.com/matlab/matlab-mcp-server for further support.");
        end
    end

    % Helper function to update fields in the request for MATLAB versions less than
    % R2023b
    function request = updateRequestFromBefore23b(request, code)
        jsonedRegionList = jsonencode(struct(...
            'regionLineNumber',1,...
            'regionString',code,...
            'regionNumber',0,...
            'endOfSection',true,...
            'sectionNumber',1));
        request.regionArray = jsonedRegionList;
    end

    % Helper function to update fields in the request when LiveEditor API version is 1.
    function request = updateRequestFromVersion1(request, code)
        request.sectionBoundaries = [];
        request.startLine = 1;
        request.endLine = builtin('count', code, newline) + 1;
    end

    % Helper function to update fields in the request when LiveEditor API version is 2.
    function request = updateRequestFromVersion2(request, code)
        request = updateRequestFromVersion1(request, code);

        % Request MIME based outputs.
        request.preferBasicOutputs = true;
    end
end

function result = processOutputs(outputs)
    result = cell(1,length(outputs));
    figureTrackingMap = containers.Map;

    % Post process each captured output based on its type.
    for ii = 1:length(outputs)
        out = outputs(ii);
        outputData = out.outputData;
        switch out.type
            case 'matrix'
                result{ii} = processMatrix(outputData);
            case 'variable'
                result{ii} = processVariable(outputData);
            case 'variableString'
                result{ii} = processVariableString(outputData);
            case 'symbolic'
                result{ii} = processSymbolic(outputData);
            case 'error'
                result{ii} = processError(outputData.text, getErrorIdentifier(outputData));
            case 'warning'
                result{ii} = processStream('stderr', outputData.text);
            case 'text'
                result{ii} = processStream('stdout', outputData.text);
            case 'stderr'
                result{ii} = processStream('stderr', outputData.text);
            case 'figure'
                % 'figure' outputType may not necessarily contain the actual image.
                % Hence, if the 'figure' is a placeholder, we store its position in
                % a map to preserve the ordering. In a later 'figure' output, if the
                % actual image data is present, we store the image in the corresponding
                % placeholder position if it exists, else the current position.
                if isfield(outputData, 'figurePlaceHolderId')
                    id = outputData.figurePlaceHolderId;
                    if ~figureTrackingMap.isKey(id)
                        figureTrackingMap(id) = ii;
                    end
                elseif isfield(outputData, 'figureImage')
                    id = outputData.figureId;
                    if figureTrackingMap.isKey(id)
                        idx = figureTrackingMap(id);
                    else
                        idx = ii;
                    end
                    result{idx} = processFigure(outputData.figureImage);
                end
            case 'text/html'
                result{ii} = processHtml(outputData);
        end
    end

    ME = matlab_mcp.getOrStashExceptions([], true);
    if ~isempty(ME)
        result{end+1} = processError(ME.message, ME.identifier);
    end

    % Helper functions to post process output of type 'matrix', 'variable' and
    % 'variableString'. These outputs are of HTML type due to various HTML tags
    % used in MATLAB outputs such as the <strong> tag in tables.
    function result = processText(text)
        result.type = 'execute_result';
        result.mimetype = {"text/html", "text/plain"};
        result.value = [sprintf("<html><body><pre>%s</pre></body></html>",text), text];
    end

    function result = processMatrix(output)
        text = sprintf("%s = %s %s\n%s", output.name, output.header, output.type, output.value);
        if output.rows > 10 || output.columns > 30
            text = strcat(text, "...");
        end
        result = processText(text);
    end

    function result = processVariable(output)
        if isempty(output.header)
            indentation = '';
        else
            indentation = sprintf('\n    ');
        end
        text = sprintf("%s = %s%s%s", output.name, output.header, indentation, output.value);
        result = processText(text);
    end

    function result = processVariableString(output)
        indentation = '';
        useSingleLineDisplay = ~(builtin('contains', output.value, newline));
        if useSingleLineDisplay
            if ~isempty(output.header)
                indentation = sprintf(newline);
            end
        else
            indentation = sprintf(newline);
        end
        text = sprintf("%s = %s%s%s", output.name, output.header, indentation, output.value);
        result = processText(text);
    end

    % Helper function for post-processing symbolic outputs. The captured output
    % contains MathML representation of symbolic expressions. We use EquationRenderer JS API to
    % convert the MathML to LaTeX values.
    function result = processSymbolic(output)
        % Use persistent variables to avoid loading multiple webwindows.
        persistent webwindow;
        persistent idler;

        if isempty(webwindow)
            url = 'toolbox/matlab/codetools/liveeditor/index.html';

            % MATLAB versions R2020b and R2021a requires specifying the base url.
            % Not doing so results in the URL not being loaded with the error
            %"Not found. Request outside of context root".
            if isMATLABReleaseOlderThan("R2021b")
                url = strcat(getenv("MWI_BASE_URL"), '/', url);
            end
            webwindow = matlab.internal.cef.webwindow(connector.getUrl(url));
            idler = jupyter.Idler;
            webwindow.PageLoadFinishedCallback = @(a,b) pageLoadCallback(a,b,idler);
        end

        % This will block the thread until stop loading is called. The values are logical
        pageLoaded = idler.startIdling(10);

        % If page is not loaded succesfully. We fallback to embedding MathML inside HTML.
        if ~pageLoaded
            result = processText(output.value);
            return
        end

        %  Use the EquationRenderer JS API to convert MathML to LaTeX.
        webwindow.executeJS('eq = require("equationrenderercore/EquationRenderer")');
        latexcode = jsondecode(webwindow.executeJS(sprintf('eq.convertMathMLToLaTeX(%s)', jsonencode(output.value))));
        if isempty(output.name)
            % If there is no variable name captured, then we only display the symbolic equation.
            % This happens in cases such as "disp(exp(b))".
            latexcode = strcat('$',latexcode,'$');
        else
            latexcode = strcat('$',output.name,' = ',latexcode,'$');
        end

        result.type = 'execute_result';
        result.mimetype = {"text/latex"};
        result.value = {latexcode};
    end

    % Helper function for processing outputs of stream type such as 'stdout' and 'stderr'
    function result = processStream(stream, text)
        result.type = 'stream';
        result.content.name = stream;
        result.content.text = text;
    end

    % Errors stay 'stderr' streams so servers that ignore isError still show them
    function result = processError(text, identifier)
        result = processStream('stderr', text);
        result.isError = true;
        result.identifier = identifier;
    end

    function identifier = getErrorIdentifier(output)
        identifier = '';
        if isfield(output, 'mException') && isfield(output.mException, 'identifier')
            identifier = char(output.mException.identifier);
        end
    end

    % Helper function for processing figure outputs.
    % base64Data will be 'data:image/png;base64,<base64_value>'
    function result = processFigure(base64Data)
        pattern = "data:(?<mimetype>.*);base64,(?<value>.*)";
        result = builtin('regexp', base64Data, pattern, 'names');
        assert(builtin('startsWith', result.mimetype, 'image'), 'Error in processFigure. ''mimetype'' is not an image');
        assert(~isempty(result.value), 'Error in processFigure. ''value'' is empty');
        result.mimetype = {result.mimetype};
        result.value = {result.value};
        result.type = 'execute_result';
    end

    % Helper function for processing text/html mime-type outputs.
    function result = processHtml(text)
        result.type = 'execute_result';
        result.mimetype = {"text/html", "text/plain"};
        result.value = [sprintf("%s",text), text];
    end

    % Helper function to notify browser page load finished
    function pageLoadCallback(webwindow,~,idler)
        idler.stopIdling();
        % Disable alert box which is preventing running JS after certain period of time.
        webwindow.executeJS('window.alert = function(){}');
    end
end
