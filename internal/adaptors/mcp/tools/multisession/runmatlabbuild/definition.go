// Copyright 2025-2026 The MathWorks, Inc.

package runmatlabbuild

const (
	name        = "run_matlab_build_in_session"
	title       = "Run MATLAB Build in a MATLAB Session"
	description = "Run a MATLAB Build Tool build (`buildtool`) in an existing MATLAB session, given its session ID. Executes one or more build tasks defined in a `buildfile.m`. Optionally specify a working directory containing the `buildfile.m`, specific task names to run, and build options. Returns the full build log and whether the build succeeded."
)

type Args struct {
	SessionID         int      `json:"session_id"                    jsonschema:"The ID of the MATLAB session in which to run the build."`
	WorkingDirectory  string   `json:"working_directory,omitempty"   jsonschema:"(Optional) Absolute path to the folder containing buildfile.m. When omitted, uses MATLAB's current working folder. Example: /home/user/myproject or C:\\Users\\username\\myproject."`
	Tasks             []string `json:"tasks,omitempty"               jsonschema:"(Optional) Task names to run. When omitted, runs the default task(s) defined in buildfile.m. Example: [\"check\", \"test\"]."`
	ContinueOnFailure bool     `json:"continue_on_failure,omitempty" jsonschema:"(Optional) When true, continues running remaining tasks even if one fails. Equivalent to the -continueOnFailure flag."`
	Parallel          bool     `json:"parallel,omitempty"            jsonschema:"(Optional) When true, runs tests in parallel. Equivalent to the -parallel flag. Requires MATLAB R2024a or later."`
	Verbosity         string   `json:"verbosity,omitempty"           jsonschema:"(Optional) Controls the verbosity of build output. One of: terse, concise, detailed, verbose. Requires MATLAB R2024b or later."`
	Skip              []string `json:"skip,omitempty"                jsonschema:"(Optional) Task names to skip. Each named task is excluded from the build. Example: [\"clean\", \"check\"]."`
}

type ReturnArgs struct {
	Log     string `json:"log"     jsonschema:"The full console output produced by the build, including task status messages and any errors."`
	Success bool   `json:"success" jsonschema:"True if all executed tasks completed successfully, false if any task failed."`
}
