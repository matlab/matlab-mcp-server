function plan = buildfile
    %buildfile Build tasks for MATLAB MCP Core Server Toolbox

    % Copyright 2026 The MathWorks, Inc.

    import matlab.buildtool.Task
    import matlab.buildtool.tasks.CodeIssuesTask
    import matlab.buildtool.tasks.TestTask
    import matlab.buildtool.tasks.CleanTask

    plan = buildplan;

    % Get project root folder
    prj = plan.RootFolder;
    srcFiles = fullfile(prj, "mcp");
    testFiles = fullfile(prj, "tests");
    unitTestFiles = fullfile(testFiles, "unit");
    testResultsFolder = fullfile(prj, "test_results");

    distFolder = fullfile(prj, "dist");
    distFolderOverride = getenv("MATLAB_MCP_CORE_SERVER_MLTBX_DIR");
    if distFolderOverride ~= ""
        distFolder = distFolderOverride;
    end

    if ~isfolder(distFolder)
        mkdir(distFolder);
    end

    % Define clean task using built-in CleanTask
    plan("clean") = CleanTask;
    plan("clean").Description = "Clean up generated artifacts";

    % Define lint task using built-in CodeIssuesTask
    plan("lint") = CodeIssuesTask( ...
        [ ...
            srcFiles, ...
            testFiles ...
        ], ...
        WarningThreshold=0 ...
    );
    plan("lint").Description = "Check code for issues and style problems";

    % Define unit test task using built-in TestTask
    plan("unit-tests") = TestTask( ...
        unitTestFiles, ...
        SourceFiles=srcFiles, ...
        TestResults=fullfile(testResultsFolder, "unit_tests.xml"), ...
        CodeCoverageResults=[ ...
            fullfile(testResultsFolder, "cobertura.xml"), ...
            fullfile(testResultsFolder, "coverage.html") ...
        ] ...
    );
    plan("unit-tests").Description = "Run unit tests";

    % Define package task
    plan("package") = Task( ...
        Description="Package the toolbox", ...
        Actions=@(context) packageTask(context, srcFiles, distFolder), ...
        Inputs=srcFiles, ...
        Outputs=[ ...
            distFolder, ...
            fullfile(prj, "mcp", "Contents.m") ...
        ] ...
    );

    % Set default tasks
    plan.DefaultTasks = ["lint", "unit-tests", "package"];
end

function packageTask(context, toolboxFolder, distFolder)
    %packageTask Package the MATLAB MCP Core Server Toolbox
    %   This task creates an MLTBX file for distribution.

    % Get project root folder
    prj = context.Plan.RootFolder;

    % Create dist folder if it doesn't exist
    if ~isfolder(distFolder)
        mkdir(distFolder);
    end

    % Stable UUID for the toolbox (do not change)
    toolboxName = "MATLAB MCP Core Server Toolbox";
    toolboxVersion = "0.1.0";
    uuid = "fce99a56-1b63-4d49-9d70-c32c5605e4eb";

    % Create Contents.m
    createContents(toolboxFolder, toolboxName, toolboxVersion);

    % Create toolbox options
    opts = matlab.addons.toolbox.ToolboxOptions(toolboxFolder, uuid, ...
        ToolboxName=toolboxName, ...
        ToolboxVersion=toolboxVersion, ...
        Summary="MATLAB toolbox for the MATLAB MCP Core server", ...
        Description="Enables sharing of MATLAB sessions through the Model Context Protocol.", ...
        AuthorName="MathWorks", ...
        AuthorEmail="genai-support@mathworks.com", ...
        AuthorCompany="MathWorks", ...
        MinimumMatlabRelease="R2023a", ...
        MaximumMatlabRelease="", ...
        ToolboxGettingStartedGuide="", ...
        ToolboxImageFile=fullfile(prj, "icon.png"), ...
        OutputFile=fullfile(distFolder, "MATLABMCPCoreServerToolbox.mltbx") ...
    );

    % Package the toolbox
    matlab.addons.toolbox.packageToolbox(opts);
end

function createContents(toolboxFolder, toolboxName, toolboxVersion)
    %createContents Create Contents.m file for the toolbox
    %   Generates a Contents.m file in the mcp folder with the toolbox
    %   name and version.

    contentsFile = fullfile(toolboxFolder, "Contents.m");
    fid = fopen(contentsFile, "w");
    closeFile = onCleanup(@() fclose(fid));
    
    release = matlabRelease;
    now = datetime("now");
    now.Format = "dd-MMM-uuuu";

    fprintf(fid, "%% %s\n", toolboxName);
    fprintf(fid, "%% Version %s (%s) %s\n", toolboxVersion, release.Release, now);
    fprintf(fid, "%%\n");
    fprintf(fid, "%%   Copyright 2026 The MathWorks, %s\n", "Inc.");
end
