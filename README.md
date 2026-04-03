# MATLAB MCP Core Server

Run MATLAB® using AI applications with the official MATLAB MCP Server from MathWorks®. The MATLAB MCP Core Server allows your AI applications to:

- Install MATLAB and toolboxes without requiring a pre-existing MATLAB installation.
- Start and quit MATLAB.
- Write and run MATLAB code.
- Assess your MATLAB code for style and correctness.

## Table of Contents

- [Setup](#setup)
  - [Claude Code](#claude-code)
  - [Claude Desktop](#claude-desktop)
  - [GitHub Copilot in Visual Studio Code](#github-copilot-in-visual-studio-code)
- [Arguments](#arguments)
- [Tools](#tools)
- [Resources](#resources)
- [Data Collection](#data-collection)
- [Contact Support](#contact-support)

## Setup

1. Install [MATLAB (MathWorks)](https://www.mathworks.com/help/install/ug/install-products-with-internet-connection.html) 2020b or later and add it to the system PATH. Alternatively, you can use the [`install_matlab`](#tools) tool to install MATLAB directly from your AI application after setting up the server.
1. To set up the MATLAB MCP Core Server for Claude Desktop, skip to the instructions for [Claude Desktop](#claude-desktop). To set up the server for other applications, follow these instructions:
   
   - For Windows or Linux, [**Download the Latest Release**](https://github.com/matlab/matlab-mcp-core-server/releases/latest). (Alternatively, you can **build from source**: install [Go](https://go.dev/doc/install) and build the binary using `go install github.com/matlab/matlab-mcp-core-server/cmd/matlab-mcp-core-server@latest`).
    
   - For macOS, first download the latest release by running the following command in your terminal:
     - For Apple silicon processors, run:
          ```sh
          curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maca64
          ```
      - For Intel processors, run:
          ```sh
          curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maci64
          ```
      Then grant executable permissions to the downloaded binary so you can run the MATLAB MCP Core Server:

      ```sh
      chmod +x ~/Downloads/matlab-mcp-core-server
      ```

1. Add the MATLAB MCP Core Server to your AI application. You can find instructions for adding MCP servers in the documentation of your AI application. For example instructions on using Claude Code®, Claude Desktop®, and GitHub Copilot in Visual Studio® Code, see below. Note that you can customize the server by specifying optional [Arguments](#arguments).

### Claude Code

In your terminal, run the following, remembering to insert the full path to the server binary you acquired in the setup:

```sh
claude mcp add --transport stdio matlab -- /fullpath/to/matlab-mcp-core-server-binary
```

You can customize the server by specifying optional [Arguments](#arguments). Note the `--` separator between Claude Code's options and the server arguments:

```sh
claude mcp add --transport stdio matlab -- /fullpath/to/matlab-mcp-core-server-binary --initial-working-folder=/home/username/myproject
```

For details on adding MCP servers in Claude Code, see [Add a local stdio server (Claude Code)](https://docs.claude.com/en/docs/claude-code/mcp#option-3%3A-add-a-local-stdio-server). To remove the server later, run:

```sh
claude mcp remove matlab
```

### Claude Desktop

You install the MATLAB MCP Core Server in Claude Desktop using the MATLAB MCP Core Server bundle.

1. Install the Filesystem extension in Claude Desktop to allow Claude to read and write files on your system. In Claude Desktop, click **Settings > Extensions > Browse extensions**. Search for the Filesystem extension developed by Anthropic and click **Install**. Specify the folders you want to allow the MCP server to access, then toggle the **Disabled** button to **Enable** the Filesystem extension.
   
2. Download the MATLAB MCP Core Server bundle `matlab-mcp-core-server.mcpb` from the [Latest Release](https://github.com/matlab/matlab-mcp-core-server/releases/latest) page. 

3. To install the MATLAB MCP Core Server bundle as a desktop extension, double click the downloaded `matlab-mcp-core-server.mcpb` file and click **Install** in Claude Desktop. (Alternatively, navigate in Claude to **File menu > Settings > Extensions > Advanced Settings > Install Extension** and select the `matlab-mcp-core-server.mcpb` file. Click **Install**).

   To customize the behaviour of the MATLAB MCP Core Server, navigate to **Settings > Extensions > Configure**, where you can modify the server's [Arguments](#arguments).

### GitHub Copilot in Visual Studio Code


In your VS Code workspace, create a file named `.vscode/mcp.json`. Insert the following JSON, remembering to specify the full path to the server binary you acquired in the setup, as well as any [Arguments](#arguments). Then save the file. (Note that on Windows, your paths require extra slashes as escape characters).

```json
{
    "servers": {
        "matlab": {
            "type": "stdio",
            "command": "C:\\fullpath\\to\\matlab-mcp-core-server-win64.exe",
            "args": []
        }
    }
}
```
For more information about using MCP servers in VS Code, see [Add and Manage MCP servers in VS Code (VS Code)](https://code.visualstudio.com/docs/copilot/customization/mcp-servers#_configure-the-mcpjson-file).



## Arguments

Customize the behavior of the server by providing arguments in the `args` array when configuring your AI application.

| Argument | Description | Example |
| ------------- | ------------- | ------------- |
| matlab-root | Full path specifying which MATLAB to start. Do not include `/bin` in the path. By default, the server tries to find the first MATLAB on the system PATH. | `"--matlab-root=/home/usr/MATLAB/R2025a"` |
| initialize-matlab-on-startup | To initialize MATLAB as soon as you start the server, set this argument to `true`. By default, MATLAB only starts when the first tool is called. | `"--initialize-matlab-on-startup=true"` |
| initial-working-folder | Specify the folder where MATLAB starts. If you do not specify a value, MATLAB starts at the path of your AI application's first [Root (MCP)](https://modelcontextprotocol.io/specification/2025-06-18/client/roots). If you have not defined a root, MATLAB starts in these locations: <br> <ul><li>Linux: `/home/username` </li><li> Windows: `C:\Users\username\Documents`</li><li>Mac: `/Users/username/Documents`</li></ul> | `"--initial-working-folder=C:\\Users\\name\\MyProject"` |
| matlab-display-mode | Specify whether to show the MATLAB desktop. Use `desktop` mode (default) to show the MATLAB desktop. Use `nodesktop` mode to use MATLAB only from your AI application, without the MATLAB desktop. Note that in `nodesktop` mode, commands requiring a graphical interface (such as `edit`, `open`, `open_system`, `uifigure`, and `appdesigner`) will still open MATLAB windows on your desktop.| `"--matlab-display-mode=nodesktop"` |
| disable-telemetry | To disable anonymized data collection, set this argument to `true`. For details, see [Data Collection](#data-collection). | `"--disable-telemetry=true"` |


## Tools

1. `install_matlab`
    - Downloads and installs MATLAB and/or MATLAB toolboxes using the MATLAB Package Manager (mpm). This tool does not require MATLAB to already be installed. It downloads the mpm CLI and uses it to install the specified products.
    - Inputs:
        - `release` (string): The MATLAB release to install (e.g., `R2025a` or `R2024b`).
        - `products` (string[]): List of product names to install. Use underscores for spaces (e.g., `MATLAB`, `Signal_Processing_Toolbox`, `Deep_Learning_Toolbox`). To install MATLAB itself, include `MATLAB` in the list.
        - `destination` (string, optional): The installation directory. If not specified, mpm uses its default location.

1. `detect_matlab_toolboxes`
    - Returns information about installed MATLAB and toolboxes, including version numbers.

1. `check_matlab_code`
    - Performs static code analysis on a MATLAB script. Returns warnings about coding style, potential errors, deprecated functions, performance issues, and best practice violations. This is a non-destructive, read-only operation that helps identify code quality issues without executing the script.
    - Inputs:
        - `script_path` (string): Absolute path to the MATLAB script file to analyze. Must be a valid `.m` file. The file is not modified during analysis. Example: `C:\Users\username\matlab\myFunction.m` or `/home/user/scripts/analysis.m`.

1. `evaluate_matlab_code`
    - Evaluates a string of MATLAB code and returns the output.
    - Inputs:
        - `code` (string): MATLAB code to evaluate.
        - `project_path` (string): Absolute path to your project directory. MATLAB sets this directory as the current working folder. Example: `C:\Users\username\matlab-project` or `/home/user/research`.

1. `run_matlab_file`
    - Executes a MATLAB script and returns the output. The script must be a valid `.m file`.
    - Inputs:
        - `script_path` (string): Absolute path to the MATLAB script file to execute. Must be a valid `.m` file. Example: `C:\Users\username\projects\analysis.m` or `/home/user/matlab/simulation.m`.

1. `run_matlab_test_file`
    - Executes a MATLAB test script and returns comprehensive test results. Designed specifically for MATLAB unit test files that follow MATLAB testing framework conventions.
    - Inputs:
        - `script_path` (string): Absolute path to the MATLAB test script file. Must be a valid `.m` file containing MATLAB unit tests. Example: `C:\Users\username\tests\testMyFunction.m` or `/home/user/matlab/tests/test_analysis.m`.

## Resources

The MCP server provides [Resources (MCP)](https://modelcontextprotocol.io/specification/2025-03-26/server/resources) to help your AI application write MATLAB code. To see instructions for using this resource, refer to the documentation of your AI application that explains how to use resources.

1. `matlab_coding_guidelines`
    - Provides comprehensive MATLAB coding standards for improving code readability, maintainability, and collaboration. The guidelines encompass naming conventions, formatting, commenting, performance optimization, and error handling.
    - URI: `guidelines://coding`
    - MIME Type: `text/markdown`
    - Source: [MATLAB Coding Standards (GitHub)](https://github.com/matlab/rules/blob/main/matlab-coding-standards.md)

1. `plain_text_live_code_guidelines`
    - Provides rules and guidelines for generating live scripts using the plain text Live Code `.m` file format, suitable for version control and AI-assisted development. Note that to run plain text live scripts you need MATLAB R2025a or newer. For details, see [Live Code File Format (MathWorks)](https://www.mathworks.com/help/matlab/matlab_prog/plain-text-file-format-for-live-scripts.html).
    - URI: `guidelines://plain-text-live-code`
    - MIME Type: `text/markdown`
    - Source: [Plain Text Live Code Generation (GitHub)](https://github.com/matlab/rules/blob/main/live-script-generation.md)

## Data Collection

The MATLAB MCP Core Server may collect fully anonymized information about your usage of the server and send it to MathWorks. This data collection helps MathWorks improve products and is on by default. To opt out of data collection, set the argument `--disable-telemetry` to `true`.

## Contact Support
MathWorks encourages you to use this repository and provide feedback. To request technical support or submit an enhancement request, [create a GitHub issue](https://github.com/matlab/matlab-mcp-core-server/issues) or email [genai-support@mathworks.com](mailto:genai-support@mathworks.com).

#

When using the MATLAB MCP Core Server, you should thoroughly review and validate all tool calls before you run them. Always keep a human in the loop for important actions and only proceed once you are confident the call will do exactly what you expect. For more information, see [User Interaction Model (MCP)](https://modelcontextprotocol.io/specification/2025-06-18/server/tools#user-interaction-model) and [Security Considerations (MCP)](https://modelcontextprotocol.io/specification/2025-06-18/server/tools#security-considerations).

The MATLAB MCP Core server may only be used with MATLAB installations that are used as a Personal Automation Server. Use with a central Automation Server is not allowed. Please contact MathWorks if Automation Server use is required. For more information see the [Program Offering Guide (MathWorks)](https://www.mathworks.com/help//pdf_doc/offering/offering.pdf).

---

Copyright 2025-2026 The MathWorks, Inc.

---
