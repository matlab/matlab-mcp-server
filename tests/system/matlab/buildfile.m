function plan = buildfile
    %buildfile Build tasks for MATLAB MCP Core Server Toolbox system tests

    % Copyright 2026 The MathWorks, Inc.

    import matlab.buildtool.Task
    import matlab.buildtool.tasks.CodeIssuesTask
    import matlab.buildtool.tasks.TestTask
    import matlab.buildtool.tasks.CleanTask

    plan = buildplan;

    % Get project root folder
    prj = plan.RootFolder;
    testFiles = fullfile(prj, "tests");
    testResultsFolder = fullfile(prj, "test_results");

    % Define clean task using built-in CleanTask
    plan("clean") = CleanTask;
    plan("clean").Description = "Clean up generated artifacts";

    % Define lint task using built-in CodeIssuesTask
    plan("lint") = CodeIssuesTask( ...
        testFiles, ...
        WarningThreshold=0, ...
        Dependencies="clean" ...
    );
    plan("lint").Description = "Check code for issues and style problems";

    % Define system test task using built-in TestTask
    plan("system-tests") = TestTask( ...
        testFiles, ...
        TestResults=fullfile(testResultsFolder, "system_tests.xml"), ...
        Dependencies="lint" ...
    );
    plan("system-tests").Description = "Run system tests";

    % Set default tasks
    plan.DefaultTasks = "system-tests";
end