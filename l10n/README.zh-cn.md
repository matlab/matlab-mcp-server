<!--
Source English Markdown:
- File: ./README.md
- Branch: main
- Commit: 9e25655af07ee633bcb5fe5bea5a0c9a844e7041
-->

# MATLAB MCP Core Server

<p align="center">
  <a href="../README.md">English</a> •
  <a href="README.es.md">Español</a> •
  <a href="README.ja.md">日本語</a> •
  <a href="README.ko.md">한국어</a> •
  简体中文
</p>

使用 MathWorks® 官方的 MATLAB MCP Server，通过 AI 应用程序运行 MATLAB®。MATLAB MCP Core Server 允许您的 AI 应用程序执行以下操作:

- 启动和退出 MATLAB。
- 编写和运行 MATLAB 代码。
- 评估 MATLAB 代码的风格和正确性。

[MATLAB Agentic Toolkit (GitHub)](https://github.com/matlab/matlab-agentic-toolkit) 和 [Simulink Agentic Toolkit (GitHub)](https://github.com/matlab/simulink-agentic-toolkit) 提供的技能 (skills) 可协助您的智能体使用 MATLAB 和 Simulink，还可以为您安装此 MCP 服务器。 

## 目录

- [设置](#设置)
  - [Claude Code](#claude-code)
  - [Claude Desktop](#claude-desktop)
  - [GitHub Copilot in Visual Studio Code](#github-copilot-in-visual-studio-code)
- [参量](#参量)
- [工具](#工具)
- [资源](#资源)
- [数据收集](#数据收集)
- [安全注意事项](#安全注意事项)
- [许可和使用](#许可和使用)
- [联系支持](#联系支持)

## 设置

1. 安装 [MATLAB (MathWorks)](https://www.mathworks.com/help/install/ug/install-products-with-internet-connection.html) 2021a 或更高版本，并将其添加到系统 PATH 环境变量中。MATLAB MCP Core Server 支持过去五年内的 MATLAB 版本。
1. 要为 Claude Desktop 设置 MATLAB MCP Core Server，请跳至 [Claude Desktop](#claude-desktop) 的说明。要为其他应用程序设置服务器，请按照以下说明操作:
   
   - 对于 Windows 或 Linux，请[**下载最新版本**](https://github.com/matlab/matlab-mcp-core-server/releases/latest)。(或者，您也可以**从源代码编译**: 安装 [Go](https://go.dev/doc/install) 并使用 `go install github.com/matlab/matlab-mcp-core-server/cmd/matlab-mcp-core-server@latest` 来编译二进制文件)。
   
   - 对于 macOS，请先在终端中运行以下命令以下载最新版本:
     - 对于 Apple silicon 处理器，请运行:
          ```sh
          curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maca64
          ```
      - 对于 Intel 处理器，请运行:
          ```sh
          curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maci64
          ```
           然后向下载的二进制文件授予可执行权限，以便您可以运行 MATLAB MCP Core Server:

      ```sh
      chmod +x ~/Downloads/matlab-mcp-core-server
      ```

1. 将 MATLAB MCP Core Server 添加到您的 AI 应用程序中。您可以在 AI 应用程序的文档中找到有关添加 MCP 服务器的说明。有关使用 Claude Code®、Claude Desktop® 和 GitHub Copilot in Visual Studio® Code 的示例说明，请参阅下文。请注意，您可以通过指定可选[参量](#参量)来自定义服务器。

### Claude Code

在终端中运行以下命令，请在其中插入您在设置过程中获取的服务器二进制文件的完整路径:

```sh
claude mcp add --transport stdio matlab -- /fullpath/to/matlab-mcp-core-server-binary
```

您可以通过指定可选[参量](#参量)来自定义服务器。请注意 Claude Code 选项与服务器参量之间需使用 `--` 分隔符:

```sh
claude mcp add --transport stdio matlab -- /fullpath/to/matlab-mcp-core-server-binary --initial-working-folder=/home/username/myproject
```

有关在 Claude Code 中添加 MCP 服务器的详细信息，请参阅 [Add a local stdio server (Claude Code)](https://docs.claude.com/en/docs/claude-code/mcp#option-3%3A-add-a-local-stdio-server)。稍后要删除服务器，请运行:

```sh
claude mcp remove matlab
```

### Claude Desktop

您可以使用 MATLAB MCP Core Server 捆绑包在 Claude Desktop 中安装 MATLAB MCP Core Server。

1. 在 Claude Desktop 中安装 Filesystem 扩展程序，以允许 Claude 在您的系统上读写文件。在 Claude Desktop 中，点击 **Settings > Extensions > Browse extensions**。搜索由 Anthropic 开发的 Filesystem 扩展程序，然后点击 **Install**。指定要允许 MCP 服务器访问的文件夹，然后将 **Disabled** 按钮切换为 **Enable** 以启用 Filesystem 扩展程序。
   
2. 从 [Latest Release](https://github.com/matlab/matlab-mcp-core-server/releases/latest) 页面下载 MATLAB MCP Core Server 捆绑包 `matlab-mcp-core-server.mcpb`。 

3. 要将 MATLAB MCP Core Server 捆绑包安装为桌面扩展程序，请双击下载的 `matlab-mcp-core-server.mcpb` 文件，然后在 Claude Desktop 中点击 **Install**。(或者，在 Claude 中导航到 **File 菜单 > Settings > Extensions > Advanced Settings > Install Extension**，然后选择 `matlab-mcp-core-server.mcpb` 文件。点击 **Install**)。

要自定义 MATLAB MCP Core Server 的行为，请导航到 **Settings > Extensions > Configure**，您可以在此修改服务器的[参量](#参量)。

### GitHub Copilot in Visual Studio Code

在您的 VS Code 工作区中，创建一个名为 `.vscode/mcp.json` 的文件。插入以下 JSON，请指定您在设置过程中获取的服务器二进制文件的完整路径，以及任何[参量](#参量)。然后保存文件。(请注意，在 Windows 平台上，您需要额外的斜杠作为转义符)。

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
有关在 VS Code 中使用 MCP 服务器的更多信息，请参阅 [Add and Manage MCP servers in VS Code (VS Code)](https://code.visualstudio.com/docs/copilot/customization/mcp-servers#_configure-the-mcpjson-file)。

## 参量

通过以下方式之一指定参量来自定义服务器的行为:
- 在 AI 应用程序的配置设置(通常为 `.json` 文件)中植入参量。
- 启动服务器时，以命令行接口 (CLI) 标志的形式输入参量。 
- 使用环境变量，可在 CLI 或应用程序的配置设置中进行指定。要从 CLI 标志派生环境变量名称，请添加前缀 `MW_MCP_SERVER_`，转换为大写，并将连字符 (`-`) 替换为下划线 (`_`)。例如，参量 `--matlab-root` 对应的环境变量为 `MW_MCP_SERVER_MATLAB_ROOT`。如果同时使用两者，CLI 标志将优先于环境变量。

| 参量 | 说明 | 示例 |
| ------------- | ------------- | ------------- |
| help | 显示所有参量的帮助信息。 | `--help` |
| version | 显示 MATLAB MCP Core Server 的版本。 | `--version` |
| matlab-root | 指定要启动的 MATLAB 的完整路径。路径中不要包含 `/bin`。默认情况下，服务器会尝试在系统的 PATH 环境变量中查找第一个 MATLAB。 | Windows: `--matlab-root=C:\Program Files\MATLAB\R2026a` <br><br> Linux/macOS: `--matlab-root=/home/usr/MATLAB/R2026a`<br><br>作为环境变量: <br>`MW_MCP_SERVER_MATLAB_ROOT=/home/usr/MATLAB/R2026a` |
| initialize-matlab-on-startup | 要在启动服务器后立即初始化 MATLAB，请将此参量设置为 `true`。默认情况下，MATLAB 仅在调用第一个工具时启动。 | `--initialize-matlab-on-startup=true` |
| initial-working-folder | 指定 MATLAB 启动时指向的工作文件夹。如果未指定值，MATLAB 会启动在 AI 应用程序的第一个 [Root (MCP)](https://modelcontextprotocol.io/specification/latest/client/roots) 路径处。如果您尚未定义 root，MATLAB 会启动在以下位置: <br> <ul><li>Linux: `/home/username` </li><li> Windows: `C:\Users\username\Documents`</li><li>Mac: `/Users/username/Documents`</li></ul> | Windows: `--initial-working-folder=C:\\Users\\username\\MyProject` <br><br> Linux/macOS: `--initial-working-folder=/Users/username/MyProject` |
| matlab-display-mode | 指定是否显示 MATLAB 桌面。使用 `desktop` 模式(默认)将显示 MATLAB 桌面。使用 `nodesktop` 模式仅从 AI 应用程序使用 MATLAB，而不显示 MATLAB 桌面。请注意，在 `nodesktop` 模式下，需要图形界面的命令(例如 `edit`、`open`、`open_system`、`uifigure` 和 `appdesigner`)仍将在桌面上打开 MATLAB 窗口。 | `--matlab-display-mode=nodesktop` |
| matlab-session-mode | 指定 MCP 服务器是启动新的 MATLAB 还是连接到现有的 MATLAB 会话 (支持 MATLAB R2023a 及更高版本)。默认为 **`auto`** 模式。<br><br> **`new` 模式: ** MCP 服务器启动新的 MATLAB 会话。<br><br>**`auto` 模式 (默认):** 服务器尝试连接到现有 MATLAB 会话，您必须已按照 `existing` 模式的说明配置该会话。如果服务器找不到现有 MATLAB 会话，则启动新会话。<br><br>**`existing` 模式:** 服务器尝试连接到现有 MATLAB 会话。您必须提前配置好 MATLAB 会话才能使用此模式，步骤如下:<br><br><ol><li>如果您首次使用 `existing` 模式，请运行 `./matlab-mcp-core-server --setup-matlab`。<br><br>此命令会在 MATLAB 中安装 一个名为 MATLAB MCP Core Server Toolbox 的附加功能。您可以使用此表中的其他参量来自定义命令。例如，要指定使用哪个 MATLAB 来安装工具箱，可以使用 `./matlab-mcp-core-server --setup-matlab --matlab-root=/home/usr/MATLAB/R2026a`。<br><br>对于 Claude Desktop，在运行 `./matlab-mcp-core-server --setup-matlab` 之前，您必须按照[设置](#设置)中的说明下载 MATLAB MCP Core Server 二进制文件。<br><br></li><li>在正在运行的 MATLAB 会话的命令行窗口中，运行 `shareMATLABSession()`。当您使用 `--matlab-session-mode=existing` 启动服务器时，MCP 服务器将连接到此 MATLAB 会话。如果您正在运行多个 MATLAB 会话，服务器将连接到您最近运行 `shareMATLABSession()` 命令的 MATLAB 会话。<br><br>作为手动运行 `shareMATLABSession()` 的替代方法，您可以将此命令添加到您的 MATLAB [Starup 脚本 (MathWorks)](https://www.mathworks.com/help/matlab/ref/startup.html) 中。</li></ol> | `--matlab-session-mode=existing` |
| extension-file | 要使用自定义 MCP 工具，请提供定义工具的 JSON 文件的路径。您还可以使用多个扩展文件。有关使用自定义工具的详细信息，请参阅[在 MATLAB MCP Core Server 中使用自定义工具](guides/custom-tools.zh-cn.md)。 | <br><br>Windows: `--extension-file=C:\\Users\\name\\my-tools.json` <br><br> Linux/macOS: `--extension-file=/path/to/my-tools.json` <br><br> **使用多个扩展文件:**<br><br>Windows:`--extension-file=C:\\path\\to\\tools-1.json --extension-file=C:\\path\\to\\tools-2.json`<br><br>Linux/macOS:`--extension-file=/path/to/tools1.json --extension-file=/path/to/tools2.json` <br><br> **使用环境变量:** <br><br> Windows: `MW_MCP_SERVER_EXTENSION_FILE=C:\Users\name\tools1.json;C:\Users\name\tools2.json` <br><br> Linux/macOS: `MW_MCP_SERVER_EXTENSION_FILE=/path/to/tools1.json:/path/to/tools2.json` |
| log-folder | 指定 MCP 服务器存储日志文件的文件夹。如果未指定，服务器将使用操作系统的默认临时文件夹。 | Windows: `--log-folder=C:\\Users\\name\\AppData\\Local\\Temp` <br><br> Linux/macOS: `--log-folder=/tmp/my-logs`  |
| log-level | MCP 服务器的日志级别。有效值按详细程度递减依次为 `debug`、`info`、`warn` 和 `error`。 | `--log-level=debug` |
| disable-telemetry | 要禁用匿名数据收集，请将此参量设置为 `true`。有关详细信息，请参阅[数据收集](#数据收集)。 | `--disable-telemetry=true` |

**多个扩展文件**

Windows:
```
--extension-file=C:\\path\\to\\my-tools.json --extension-file=C:\\path\\to\\my-other-tools.json
```

Linux 和 macOS:
```
--extension-file=/path/to/my-tools.json --extension-file=/path/to/my-other-tools.json
```

**环境变量**

Windows:
```
MW_MCP_SERVER_EXTENSION_FILE=C:\Users\name\my-tools.json;C:\Users\name\my-other-tools.json
```

Linux 和 macOS:
```
MW_MCP_SERVER_EXTENSION_FILE=/path/to/my-tools.json:/path/to/my-other-tools.json
```
## 工具

1. `detect_matlab_toolboxes`
    - 返回有关已安装的 MATLAB 和工具箱的信息，包括版本号。  

1. `check_matlab_code`
    - 对 MATLAB 脚本执行静态代码分析。返回有关编码风格、潜在错误、已弃用函数、性能问题和最佳实践违规的警告。这是一个非破坏性的只读操作，可在不执行脚本的情况下识别代码质量问题。
    - 输入:
        - `script_path` (字符串): 要分析的 MATLAB 脚本文件的绝对路径。必须是有效的 `.m` 文件。在分析过程中不会修改该文件。示例: `C:\Users\username\matlab\myFunction.m` 或 `/home/user/scripts/analysis.m`。

1. `evaluate_matlab_code`
    - 计算 MATLAB 代码字符串并返回输出。
    - 输入:
        - `code` (字符串): 要计算的 MATLAB 代码。
        - `project_path` (字符串): 工程目录的绝对路径。MATLAB 将此目录设置为当前工作文件夹。示例: `C:\Users\username\matlab-project` 或 `/home/user/research`。

1. `run_matlab_file`
    - 执行 MATLAB 脚本并返回输出。该脚本必须是有效的 `.m file`。
    - 输入:
        - `script_path` (字符串): 要执行的 MATLAB 脚本文件的绝对路径。必须是有效的 `.m` 文件。示例: `C:\Users\username\projects\analysis.m` 或 `/home/user/matlab/simulation.m`。

1. `run_matlab_test_file`
    - 执行 MATLAB 测试脚本并返回全面的测试结果。专为遵循 MATLAB 测试框架约定的 MATLAB 单元测试文件设计。
    - 输入:
        - `script_path` (字符串): MATLAB 测试脚本文件的绝对路径。必须是包含 MATLAB 单元测试的有效 `.m` 文件。示例: `C:\Users\username\tests\testMyFunction.m` 或 `/home/user/matlab/tests/test_analysis.m`。

## 资源

MCP 服务器提供了 [Resources (MCP)](https://modelcontextprotocol.io/specification/latest/server/resources)，以帮助您的 AI 应用程序编写 MATLAB 代码。要查看使用此资源的说明，请参阅 AI 应用程序中有关如何使用资源的文档。

1. `matlab_coding_guidelines`
    - 提供全面的 MATLAB 编码标准，以提高代码的可读性、可维护性和协作性。这些规范涵盖命名约定、格式、注释、性能优化和错误处理。
    - URI: `guidelines://coding`
    - MIME Type: `text/markdown`
    - 来源: [MATLAB Coding Standards (GitHub)](https://github.com/matlab/rules/blob/main/matlab-coding-standards.md)

1. `plain_text_live_code_guidelines`
    - 提供使用纯文本 Live Code `.m` 文件格式生成实时脚本的规则和规范，适用于版本控制和 AI 辅助开发。请注意，要运行纯文本实时脚本，您需要 MATLAB R2025a 或更高版本。有关详细信息，请参阅 [Live Code File Format (MathWorks)](https://www.mathworks.com/help/matlab/matlab_prog/plain-text-file-format-for-live-scripts.html)。
    - URI: `guidelines://plain-text-live-code`
    - MIME Type: `text/markdown`
    - 来源: [Plain Text Live Code Generation (GitHub)](https://github.com/matlab/rules/blob/main/live-script-generation.md)

## 数据收集

MATLAB MCP Core Server 可能会收集有关您使用服务器的完全匿名信息，并将其发送给 MathWorks。此数据收集有助于 MathWorks 改进产品，默认处于开启状态。要退出数据收集，请将参量 `--disable-telemetry` 设置为 `true`。

## 安全注意事项

使用 MATLAB MCP Core Server 时，在您运行任何工具调用之前，应先对其进行全面地审查和验证。对于重要操作，请始终保证有人参与其中，并且只有在您确信调用将完全按照您的预期执行时才继续操作。有关详细信息，请参阅 [User Interaction Model (MCP)](https://modelcontextprotocol.io/specification/latest/server/tools#user-interaction-model) 和 [Security Considerations (MCP)](https://modelcontextprotocol.io/specification/latest/server/tools#security-considerations)。

## 许可和使用

您可在此 GitHub 仓库的 [LICENSE.md](../LICENSE.md) 文件中查看许可证。

MCP 服务器仅可在遵守 MathWorks 软件许可协议的前提下与 MATLAB 一起使用，且不得由多个用户共享使用。如需支持共享式或集中式服务器使用，请联系 MathWorks。

## 联系支持

MathWorks 鼓励您使用此仓库并提供反馈。要请求技术支持或提交增强请求，请[创建 GitHub issue](https://github.com/matlab/matlab-mcp-core-server/issues) 或联系 [MathWorks Technical Support](https://www.mathworks.com/support/contact_us.html)。

---

Copyright 2025-2026 The MathWorks, Inc.

---
