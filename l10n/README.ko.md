<!--
Source English Markdown:
- File: ./README.md
- Branch: main
- Commit: 8df303db49ea1c497eb0128257e81f16e1fcba34
-->

# MATLAB MCP Core Server

<p align="center">
  <a href="../README.md">English</a> •
  <a href="README.es.md">Español</a> •
  <a href="README.ja.md">日本語</a> •
  한국어 •
  <a href="README.zh-cn.md">简体中文</a>
</p>

> [!경고]
> 2026년 6월 18일 (v0.11.0)에서 MATLAB MCP Core Server의 이름이 MATLAB MCP Server로 변경됩니다. 이 날짜 이후에 최신 버전의 서버를 사용하려면 설정을 업데이트해야 합니다.
>
> | 변경 사항 | 필요한 조치 |
> |:-------:|:---------------:|
> | **리포지토리 URL**<br>`github.com/matlab/matlab-mcp-core-server` → **`github.com/matlab/matlab-mcp-server`** | 없음. GitHub가 자동으로 리디렉션합니다. |
> | **바이너리 이름**<br>새 형식: **`matlab-mcp-server-<os>-<arch>[.exe]`**<br>예: `matlab-mcp-server-windows-x64.exe` | AI 애플리케이션의 구성 설정(일반적으로 `.json` 파일)에서 바이너리 이름을 업데이트하십시오. |
> | **Go 모듈**<br>`github.com/matlab/matlab-mcp-core-server` → **`github.com/matlab/matlab-mcp-server`** | Go 프로젝트에서 MATLAB MCP Core Server 모듈을 사용하고 있는 경우 `go.mod`의 모듈 이름과 import 선언을 업데이트하십시오. |

MathWorks®의 공식 MATLAB MCP Server를 사용하여 AI 애플리케이션에서 MATLAB®을 실행하십시오. MATLAB MCP Core Server를 사용하면 AI 애플리케이션이 다음을 수행할 수 있습니다.

- MATLAB을 시작하고 종료합니다.
- MATLAB 코드를 작성하고 실행합니다.
- MATLAB 코드의 스타일과 정확성을 평가합니다.

[MATLAB Agentic Toolkit (GitHub)](https://github.com/matlab/matlab-agentic-toolkit)과 [Simulink Agentic Toolkit (GitHub)](https://github.com/matlab/simulink-agentic-toolkit)의 스킬을 사용하면 에이전트가 MATLAB과 Simulink를 사용하는 데 도움을 줄 수 있으며, 이들 툴킷은 이 MCP 서버를 대신 설치해 줄 수도 있습니다.

## 목차

- [설정](#설정)
  - [Claude Code](#claude-code)
  - [Claude Desktop](#claude-desktop)
  - [Visual Studio Code의 GitHub Copilot](#visual-studio-code의-github-copilot)
- [인수](#인수)
- [툴](#툴)
- [리소스](#리소스)
- [데이터 수집](#데이터-수집)
- [보안 고려 사항](#보안-고려-사항)
- [라이선스 및 사용](#라이선스-및-사용)
- [지원 문의](#지원-문의)

## 설정

1. [MATLAB (MathWorks)](https://www.mathworks.com/help/install/ug/install-products-with-internet-connection.html) 2021a 이상을 설치하고 시스템 PATH에 추가하십시오. MATLAB MCP Core Server는 최근 5년간의 MATLAB 릴리스를 지원합니다.
1. Claude Desktop용 MATLAB MCP Core Server를 설정하려면 [Claude Desktop](#claude-desktop)의 지침으로 건너뛰십시오. 다른 애플리케이션용 서버를 설정하려면 다음 지침을 따르십시오.
   
   - Windows 또는 Linux의 경우 [**최신 릴리스를 다운로드**](https://github.com/matlab/matlab-mcp-core-server/releases/latest)하십시오. (또는 **소스에서 빌드**할 수도 있습니다. [Go](https://go.dev/doc/install)를 설치하고 `go install github.com/matlab/matlab-mcp-core-server/cmd/matlab-mcp-core-server@latest` 명령을 사용하여 바이너리를 빌드하십시오.)
    
   - macOS의 경우 먼저 터미널에서 다음 명령을 실행하여 최신 릴리스를 다운로드하십시오.
     - Apple Silicon 프로세서의 경우:
          ```sh
          curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maca64
          ```
      - Intel 프로세서의 경우:
          ```sh
          curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maci64
          ```
      그런 다음 다운로드한 바이너리에 실행 권한을 부여하여 MATLAB MCP Core Server를 실행할 수 있도록 하십시오.

      ```sh
      chmod +x ~/Downloads/matlab-mcp-core-server
      ```

1. MATLAB MCP Core Server를 AI 애플리케이션에 추가하십시오. MCP 서버를 추가하는 방법은 AI 애플리케이션의 문서에서 확인할 수 있습니다. Claude Code®, Claude Desktop®, Visual Studio® Code의 GitHub Copilot 사용 예는 아래를 참조하십시오. 선택적 [인수](#인수)를 지정하여 서버를 사용자 지정할 수 있습니다.

### Claude Code

터미널에서 다음을 실행하십시오. 설정 단계에서 다운로드한 서버 바이너리의 전체 경로를 지정해야 합니다.

```sh
claude mcp add --transport stdio matlab -- /fullpath/to/matlab-mcp-core-server-binary
```

선택적 [인수](#인수)를 지정하여 서버를 사용자 지정할 수 있습니다. Claude Code의 옵션과 서버 인수 사이에 `--` 구분자가 필요합니다.

```sh
claude mcp add --transport stdio matlab -- /fullpath/to/matlab-mcp-core-server-binary --initial-working-folder=/home/username/myproject
```

Claude Code에서 MCP 서버를 추가하는 방법에 대한 자세한 내용은 [Add a local stdio server (Claude Code)](https://docs.claude.com/en/docs/claude-code/mcp#option-3%3A-add-a-local-stdio-server)를 참조하십시오. 서버를 나중에 제거하려면 다음을 실행하십시오.

```sh
claude mcp remove matlab
```

### Claude Desktop

MATLAB MCP Core Server 번들을 사용하여 Claude Desktop에 MATLAB MCP Core Server를 설치합니다.

1. Claude Desktop에 Filesystem 확장을 설치하여 Claude가 시스템의 파일을 읽고 쓸 수 있도록 하십시오. Claude Desktop에서 **Settings > Extensions > Browse extensions**를 클릭하십시오. Anthropic이 개발한 Filesystem 확장을 검색하고 **Install**을 클릭하십시오. MCP 서버에 액세스를 허용할 폴더를 지정한 다음 **Disabled** 버튼을 **Enable**로 전환하여 Filesystem 확장을 활성화하십시오.
   
2. [최신 릴리스](https://github.com/matlab/matlab-mcp-core-server/releases/latest) 페이지에서 MATLAB MCP Core Server 번들 `matlab-mcp-core-server.mcpb`를 다운로드하십시오.

3. MATLAB MCP Core Server 번들을 데스크탑 확장으로 설치하려면 다운로드한 `matlab-mcp-core-server.mcpb` 파일을 더블 클릭하고 Claude Desktop에서 **Install**을 클릭하십시오. (또는 Claude에서 **File 메뉴 > Settings > Extensions > Advanced Settings > Install Extension**으로 이동하여 `matlab-mcp-core-server.mcpb` 파일을 선택하십시오. **Install**을 클릭하십시오.)

MATLAB MCP Core Server의 동작을 사용자 지정하려면 **Settings > Extensions > Configure**로 이동하여 서버의 [인수](#인수)를 수정하십시오.
   
### Visual Studio Code의 GitHub Copilot

VS Code 작업 영역에서 `.vscode/mcp.json` 파일을 생성하십시오. 다음 JSON을 삽입하되 설정 단계에서 다운로드한 서버 바이너리의 전체 경로와 [인수](#인수)를 지정하십시오. 그런 다음 파일을 저장하십시오. (참고: Windows에서는 경로의 백슬래시를 이스케이프 문자로 사용하므로 추가 슬래시가 필요합니다.)

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
VS Code에서 MCP 서버를 사용하는 방법에 대한 자세한 내용은 [Add and Manage MCP servers in VS Code (VS Code)](https://code.visualstudio.com/docs/copilot/customization/mcp-servers#_configure-the-mcpjson-file)를 참조하십시오.

## 인수

다음 방법 중 하나로 인수를 지정하여 서버의 동작을 사용자 지정할 수 있습니다.
- AI 애플리케이션의 구성 설정(일반적으로 `.json` 파일)에 인수를 삽입합니다.
- 서버를 시작할 때 CLI(명령줄 인터페이스) 플래그로 인수를 입력합니다.
- CLI 또는 애플리케이션의 구성 설정에서 환경 변수를 사용합니다. CLI 플래그에서 환경 변수 이름을 유도하려면 접두사 `MW_MCP_SERVER_`를 추가하고, 대문자로 변환하고, 하이픈(`-`)을 밑줄(`_`)로 바꾸십시오. 예를 들어 인수 `--matlab-root`는 환경 변수 `MW_MCP_SERVER_MATLAB_ROOT`가 됩니다. CLI 플래그와 환경 변수를 모두 사용하는 경우 CLI 플래그가 우선합니다.

| 인수 | 설명 | 예 |
| ------------- | ------------- | ------------- |
| help | 모든 인수에 대한 도움말 정보를 표시합니다. | `--help` |
| version | MATLAB MCP Core Server의 버전을 표시합니다. | `--version` |
| matlab-root | 시작할 MATLAB을 지정하는 전체 경로입니다. 경로에 `/bin`을 포함하지 마십시오. 기본적으로 서버는 시스템 PATH에서 첫 번째 MATLAB을 찾습니다. | Windows: `--matlab-root=C:\\Program Files\\MATLAB\\R2026a` <br><br> Linux/macOS: `--matlab-root=/home/usr/MATLAB/R2026a`<br><br>환경 변수: `MW_MCP_SERVER_MATLAB_ROOT=/home/usr/MATLAB/R2026a` |
| initialize-matlab-on-startup | 서버를 시작하자마자 MATLAB을 초기화하려면 이 인수를 `true`로 설정하십시오. 기본적으로 MATLAB은 첫 번째 툴이 호출될 때만 시작됩니다. | `--initialize-matlab-on-startup=true` |
| initial-working-folder | MATLAB이 시작되는 폴더를 지정합니다. 값을 지정하지 않으면 MATLAB은 AI 애플리케이션의 첫 번째 [Root (MCP)](https://modelcontextprotocol.io/specification/latest/client/roots) 경로에서 시작됩니다. 루트를 정의하지 않은 경우 MATLAB은 다음 위치에서 시작됩니다. <br> <ul><li>Linux: `/home/username` </li><li> Windows: `C:\Users\username\Documents`</li><li>Mac: `/Users/username/Documents`</li></ul> | Windows: `--initial-working-folder=C:\\Users\\username\\MyProject` <br><br> Linux/macOS: `--initial-working-folder=/Users/username/MyProject` |
| matlab-display-mode | MATLAB 데스크탑 표시 여부를 지정합니다. MATLAB 데스크탑을 표시하려면 `desktop` 모드(기본값)를 사용하십시오. MATLAB 데스크탑 없이 AI 애플리케이션에서만 MATLAB을 사용하려면 `nodesktop` 모드를 사용하십시오. `nodesktop` 모드에서도 그래픽 인터페이스가 필요한 명령(예: `edit`, `open`, `open_system`, `uifigure`, `appdesigner`)은 데스크탑에 MATLAB 창을 엽니다. | `--matlab-display-mode=nodesktop` |
| matlab-session-mode | MCP 서버가 새 MATLAB을 시작할지 기존 MATLAB 세션에 연결할지 지정합니다(MATLAB R2023a 이상 지원). 기본값은 **`auto`** 모드입니다.<br><br> **`new` 모드:** MCP 서버가 새 MATLAB 세션을 시작합니다. <br><br>**`auto` 모드(기본값):** 서버가 아래 지침을 사용하여 `existing` 모드용으로 구성한 기존 MATLAB 세션에 연결을 시도합니다. 기존 MATLAB 세션을 찾을 수 없으면 새 세션을 시작합니다. <br><br>**`existing` 모드:** 서버가 기존 MATLAB 세션에 연결을 시도합니다. 이 모드를 사용하려면 다음 단계에 따라 MATLAB 세션을 사전에 구성해야 합니다.<br><br><ol><li>`existing` 모드를 처음 사용하는 경우 `./matlab-mcp-core-server --setup-matlab`을 실행하십시오.<br><br>이 명령은 MATLAB에 MATLAB MCP Core Server Toolbox라는 애드온을 설치합니다. 이 표의 다른 인수로 명령을 사용자 지정할 수 있습니다. 예를 들어 툴박스를 설치할 MATLAB을 지정하려면 `./matlab-mcp-core-server --setup-matlab --matlab-root=/home/usr/MATLAB/R2026a`를 사용할 수 있습니다.<br><br>Claude Desktop의 경우 `./matlab-mcp-core-server --setup-matlab`을 실행하기 전에 [설정](#설정)의 지침을 사용하여 MATLAB MCP Core Server 바이너리를 다운로드해야 합니다.<br><br></li><li>실행 중인 MATLAB 세션의 명령 창에서 `shareMATLABSession()`을 실행하십시오. `--matlab-session-mode=existing`으로 서버를 시작하면 MCP 서버가 이 MATLAB에 연결됩니다. 여러 MATLAB 세션을 실행 중인 경우 서버는 `shareMATLABSession()` 명령을 가장 최근에 실행한 MATLAB 세션에 연결됩니다.<br><br>`shareMATLABSession()`을 수동으로 실행하는 대신 MATLAB [시작 스크립트 (MathWorks)](https://www.mathworks.com/help/matlab/ref/startup.html)에 이 명령을 추가할 수 있습니다.</li></ol> | `--matlab-session-mode=existing` |
| extension-file | 사용자 지정 MCP 툴을 사용하려면 툴을 정의하는 JSON 파일의 경로를 제공하십시오. 여러 확장 파일을 사용할 수도 있습니다. 사용자 지정 툴 사용에 대한 자세한 내용은 [Use Custom Tools with the MATLAB MCP Core Server](guides/custom-tools.ko.md)를 참조하십시오. | <br><br>Windows: `--extension-file=C:\\Users\\name\\my-tools.json` <br><br> Linux/macOS: `--extension-file=/path/to/my-tools.json` <br><br> **여러 확장 파일 사용:**<br><br>Windows:`--extension-file=C:\\path\\to\\tools-1.json --extension-file=C:\\path\\to\\tools-2.json`<br><br>Linux/macOS:`--extension-file=/path/to/tools1.json --extension-file=/path/to/tools2.json` <br><br> **환경 변수 사용:** <br><br> Windows: `MW_MCP_SERVER_EXTENSION_FILE=C:\Users\name\tools1.json;C:\Users\name\tools2.json` <br><br> Linux/macOS: `MW_MCP_SERVER_EXTENSION_FILE=/path/to/tools1.json:/path/to/tools2.json` |
| log-folder | MCP 서버가 로그 파일을 저장하는 폴더를 지정합니다. 지정하지 않으면 서버는 운영 체제의 기본 임시 폴더를 사용합니다. | Windows: `--log-folder=C:\\Users\\name\\AppData\\Local\\Temp` <br><br> Linux/macOS: `--log-folder=/tmp/my-logs`  |
| log-level | MCP 서버의 로그 수준입니다. 유효한 값은 세부 정보가 자세한 순서대로 `debug`, `info`, `warn`, `error`입니다. | `--log-level=debug` |
| disable-telemetry | 익명화된 데이터 수집을 비활성화하려면 이 인수를 `true`로 설정하십시오. 자세한 내용은 [데이터 수집](#데이터-수집)을 참조하십시오. | `--disable-telemetry=true` |

**여러 확장 파일**

Windows:
```
--extension-file=C:\\path\\to\\my-tools.json --extension-file=C:\\path\\to\\my-other-tools.json
```

Linux 및 macOS:
```
--extension-file=/path/to/my-tools.json --extension-file=/path/to/my-other-tools.json
```

**환경 변수**

Windows:
```
MW_MCP_SERVER_EXTENSION_FILE=C:\Users\name\my-tools.json;C:\Users\name\my-other-tools.json
```

Linux 및 macOS:
```
MW_MCP_SERVER_EXTENSION_FILE=/path/to/my-tools.json:/path/to/my-other-tools.json
```
## 툴

1. `detect_matlab_toolboxes`
    - 설치된 MATLAB 및 툴박스에 대한 정보(버전 번호 포함)를 반환합니다.  

1. `check_matlab_code`
    - MATLAB 스크립트에 대한 정적 코드 분석을 수행합니다. 코딩 스타일, 잠재적 오류, 더 이상 사용되지 않는 함수, 성능 문제, 모범 사례 위반에 대한 경고를 반환합니다. 스크립트를 실행하지 않고 코드 품질 문제를 식별하는 비파괴적 읽기 전용 작업입니다.
    - 입력값:
        - `script_path` (string): 분석할 MATLAB 스크립트 파일의 절대 경로입니다. 유효한 `.m` 파일이어야 합니다. 분석 중에 파일이 수정되지 않습니다. 예: `C:\Users\username\matlab\myFunction.m` 또는 `/home/user/scripts/analysis.m`.

1. `evaluate_matlab_code`
    - MATLAB 코드 문자열을 실행하고 출력을 반환합니다.
    - 입력값:
        - `code` (string): 실행할 MATLAB 코드입니다.
        - `project_path` (string): 프로젝트 디렉터리의 절대 경로입니다. MATLAB이 이 디렉터리를 현재 작업 폴더로 설정합니다. 예: `C:\Users\username\matlab-project` 또는 `/home/user/research`.

1. `run_matlab_file`
    - MATLAB 스크립트를 실행하고 출력을 반환합니다. 스크립트는 유효한 `.m` 파일이어야 합니다.
    - 입력값:
        - `script_path` (string): 실행할 MATLAB 스크립트 파일의 절대 경로입니다. 유효한 `.m` 파일이어야 합니다. 예: `C:\Users\username\projects\analysis.m` 또는 `/home/user/matlab/simulation.m`.

1. `run_matlab_test_file`
    - MATLAB 테스트 스크립트를 실행하고 포괄적인 테스트 결과를 반환합니다. MATLAB 테스트 프레임워크 규칙을 따르는 MATLAB 단위 테스트 파일용으로 설계되었습니다.
    - 입력값:
        - `script_path` (string): MATLAB 테스트 스크립트 파일의 절대 경로입니다. MATLAB 단위 테스트가 포함된 유효한 `.m` 파일이어야 합니다. 예: `C:\Users\username\tests\testMyFunction.m` 또는 `/home/user/matlab/tests/test_analysis.m`.

## 리소스

MCP 서버는 AI 애플리케이션이 MATLAB 코드를 작성하는 데 도움이 되는 [리소스 (MCP)](https://modelcontextprotocol.io/specification/latest/server/resources)를 제공합니다. 이 리소스 사용 방법에 대한 지침은 리소스 사용 방법을 설명하는 AI 애플리케이션의 문서를 참조하십시오.

1. `matlab_coding_guidelines`
    - 코드 가독성, 유지 관리성, 협업을 개선하기 위한 포괄적인 MATLAB 코딩 표준을 제공합니다. 이 가이드라인은 명명 규칙, 서식 지정, 주석 달기, 성능 최적화, 오류 처리를 다룹니다.
    - URI: `guidelines://coding`
    - MIME Type: `text/markdown`
    - 출처: [MATLAB Coding Standards (GitHub)](https://github.com/matlab/rules/blob/main/matlab-coding-standards.md)

1. `plain_text_live_code_guidelines`
    - 버전 컨트롤과 AI 지원 개발에 적합한 일반 텍스트 라이브 코드 `.m` 파일 형식을 사용하여 라이브 스크립트를 생성하기 위한 규칙과 가이드라인을 제공합니다. 일반 텍스트 라이브 스크립트를 실행하려면 MATLAB R2025a 이상이 필요합니다. 자세한 내용은 [Live Code File Format (MathWorks)](https://www.mathworks.com/help/matlab/matlab_prog/plain-text-file-format-for-live-scripts.html)을 참조하십시오.
    - URI: `guidelines://plain-text-live-code`
    - MIME Type: `text/markdown`
    - 출처: [Plain Text Live Code Generation (GitHub)](https://github.com/matlab/rules/blob/main/live-script-generation.md)

## 데이터 수집

MATLAB MCP Core Server는 서버 사용에 대한 완전히 익명화된 정보를 수집하여 MathWorks에 전송할 수 있습니다. 이 데이터 수집은 MathWorks의 제품 개선에 도움이 되며 기본적으로 활성화되어 있습니다. 데이터 수집을 거부하려면 인수 `--disable-telemetry`를 `true`로 설정하십시오.

## 보안 고려 사항

MATLAB MCP Core Server를 사용할 때는 모든 툴 호출을 실행하기 전에 철저히 검토하고 유효성을 검사하십시오. 중요한 작업에는 항상 사람이 개입하도록 하고, 호출이 예상대로 정확히 수행될 것이라고 확신하는 경우에만 진행하십시오. 자세한 내용은 [User Interaction Model (MCP)](https://modelcontextprotocol.io/specification/latest/server/tools#user-interaction-model) 및 [Security Considerations (MCP)](https://modelcontextprotocol.io/specification/latest/server/tools#security-considerations)를 참조하십시오.

## 라이선스 및 사용

라이선스는 이 GitHub 리포지토리의 [LICENSE.md](../LICENSE.md) 파일에서 확인할 수 있습니다.

MathWorks 소프트웨어 라이선스 계약에 따라 MATLAB과 함께 사용하는 경우에만 MCP 서버 사용이 허용되며, 여러 사용자가 MCP 서버를 공유해서는 안 됩니다. 공유 또는 중앙 집중식 서버 사용을 지원해야 하는 경우 MathWorks에 문의하십시오.

## 지원 문의

MathWorks는 이 리포지토리를 사용하고 피드백을 제공해 주시기를 권장합니다. 기술 지원을 요청하거나 개선 사항을 제출하려면 [GitHub 이슈를 생성](https://github.com/matlab/matlab-mcp-core-server/issues)하거나 [MathWorks 기술 지원팀](https://www.mathworks.com/support/contact_us.html)에 문의하십시오.

---

Copyright 2025-2026 The MathWorks, Inc.

---
