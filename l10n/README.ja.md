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
  日本語 •
  <a href="README.ko.md">한국어</a> •
  <a href="README.zh-cn.md">简体中文</a>
</p>

MathWorks® 公式の MATLAB MCP Server を使用して、AI アプリケーションから MATLAB® を実行できます。MATLAB MCP Core Server を使用すると、AI アプリケーションから以下の操作が可能になります。

- MATLAB の起動と終了。
- MATLAB コードの記述と実行。
- MATLAB コードのスタイルと正確性の評価。

エージェントが MATLAB および Simulink を活用できるようにするために、[MATLAB Agentic Toolkit (GitHub)](https://github.com/matlab/matlab-agentic-toolkit) および [Simulink Agentic Toolkit (GitHub)](https://github.com/matlab/simulink-agentic-toolkit) のスキルを使用できます。これらのツールキットを使用して、本 MCP サーバーをインストールすることもできます。

## 目次

- [セットアップ](#セットアップ)
  - [Claude Code](#claude-code)
  - [Claude Desktop](#claude-desktop)
  - [Visual Studio Code の GitHub Copilot](#visual-studio-code-の-github-copilot)
- [引数](#引数)
- [ツール](#ツール)
- [リソース](#リソース)
- [データ収集](#データ収集)
- [セキュリティに関する考慮事項](#セキュリティに関する考慮事項)
- [ライセンスと使用について](#ライセンスと使用について)
- [サポートへのお問い合わせ](#サポートへのお問い合わせ)

## セットアップ

1. [MATLAB (MathWorks)](https://www.mathworks.com/help/install/ug/install-products-with-internet-connection.html) R2021a 以降をインストールし、システム パスに追加します。MATLAB MCP Core Server は、過去 5 年間の MATLAB リリースをサポートしています。
1. Claude Desktop 用に MATLAB MCP Core Server をセットアップする場合は、「[Claude Desktop](#claude-desktop)」の手順に進んでください。その他のアプリケーション用にサーバーをセットアップする場合は、以下の手順に従います。
   
   - Windows または Linux の場合、[**最新リリースをダウンロード**](https://github.com/matlab/matlab-mcp-core-server/releases/latest)します (または、**ソースからビルド**することもできます。[Go](https://go.dev/doc/install) をインストールし、`go install github.com/matlab/matlab-mcp-core-server/cmd/matlab-mcp-core-server@latest` を使用してバイナリをビルドします)。
    
   - macOS の場合、まずターミナルで以下のコマンドを実行して最新リリースをダウンロードします。
     - Apple シリコン プロセッサの場合、以下を実行します。
          ```sh
          curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maca64
          ```
      - Intel プロセッサの場合、以下を実行します。
          ```sh
          curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maci64
          ```
      次に、ダウンロードしたバイナリに実行権限を付与して、MATLAB MCP Core Server を実行できるようにします。

      ```sh
      chmod +x ~/Downloads/matlab-mcp-core-server
      ```

1. お使いの AI アプリケーションに MATLAB MCP Core Server を追加します。MCP サーバーの追加手順については、使用する AI アプリケーションのドキュメンテーションを参照してください。Claude Code®、Claude Desktop®、および Visual Studio® Code の GitHub Copilot での手順の例については、以下を参照してください。オプションの[引数](#引数)を指定してサーバーをカスタマイズすることもできます。

### Claude Code

ターミナルで以下を実行します。必ずセットアップで取得したサーバー バイナリの絶対パスを指定してください。

```sh
claude mcp add --transport stdio matlab -- /fullpath/to/matlab-mcp-core-server-binary
```

オプションの[引数](#引数)を指定してサーバーをカスタマイズできます。Claude Code のオプションとサーバー引数の間には `--` 区切り文字が必要なことに注意してください。

```sh
claude mcp add --transport stdio matlab -- /fullpath/to/matlab-mcp-core-server-binary --initial-working-folder=/home/username/myproject
```

Claude Code で MCP サーバーを追加する方法の詳細については、「[ローカル stdio サーバーを追加する (Claude Code)](https://code.claude.com/docs/ja/mcp#%E3%82%AA%E3%83%97%E3%82%B7%E3%83%A7%E3%83%B3-3%EF%BC%9A%E3%83%AD%E3%83%BC%E3%82%AB%E3%83%AB-stdio-%E3%82%B5%E3%83%BC%E3%83%90%E3%83%BC%E3%82%92%E8%BF%BD%E5%8A%A0%E3%81%99%E3%82%8B)」を参照してください。後でサーバーを削除するには、以下を実行します。

```sh
claude mcp remove matlab
```

### Claude Desktop

MATLAB MCP Core Server バンドルを使用して、Claude Desktop に MATLAB MCP Core Server をインストールします。

1. Claude Desktop に Filesystem 拡張機能をインストールして、Claude がシステム上のファイルを読み書きできるようにします。Claude Desktop で **[Settings]、[Extensions]、[Browse extensions]** の順にクリックします。Anthropic が開発した Filesystem 拡張機能を検索して **[Install]** をクリックします。MCP サーバーにアクセスを許可するフォルダーを指定し、**[Disabled]** ボタンを切り替えて Filesystem 拡張機能を**有効にします**。
   
2. [最新リリース](https://github.com/matlab/matlab-mcp-core-server/releases/latest)のページから MATLAB MCP Core Server バンドル `matlab-mcp-core-server.mcpb` をダウンロードします。

3. MATLAB MCP Core Server バンドルをデスクトップ拡張機能としてインストールするには、ダウンロードした `matlab-mcp-core-server.mcpb` ファイルをダブルクリックし、Claude Desktop で **[Install]** をクリックします (または、Claude で **[File] メニュー、[Settings]、[Extensions]、[Advanced Settings]、[Install Extension]** の順に移動し、`matlab-mcp-core-server.mcpb` ファイルを選択して **[Install]** をクリックします)。

MATLAB MCP Core Server の動作をカスタマイズするには、**[Settings]、[Extensions]、[Configure]** の順に移動します。ここで、サーバーの[引数](#引数)を変更できます。
   
### Visual Studio Code の GitHub Copilot

VS Code ワークスペースに `.vscode/mcp.json` という名前のファイルを作成します。以下の JSON を挿入し、セットアップで取得したサーバー バイナリの絶対パスと[引数](#引数)を忘れずに指定します。その後、ファイルを保存します (Windows ではパスにエスケープ文字として追加のバックスラッシュが必要であることに注意してください)。

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
VS Code で MCP サーバーを使用する方法の詳細については、「[Add and Manage MCP servers in VS Code (VS Code)](https://code.visualstudio.com/docs/copilot/customization/mcp-servers#_configure-the-mcpjson-file)」を参照してください。



## 引数

以下のいずれかの方法で引数を指定して、サーバーの動作をカスタマイズできます。
- AI アプリケーションの構成設定 (通常は `.json` ファイル) に引数を挿入します。
- サーバーの起動時にコマンドライン インターフェイス (CLI) フラグとして引数を入力します。
- 環境変数を使用します。CLI またはアプリケーションの構成設定のいずれかで指定します。CLI フラグから環境変数名を導き出すには、接頭辞 `MW_MCP_SERVER_` を追加し、大文字に変換し、ハイフン (`-`) をアンダースコア (`_`) に置き換えます。たとえば、引数 `--matlab-root` は環境変数 `MW_MCP_SERVER_MATLAB_ROOT` になります。両方を使用した場合、CLI フラグが環境変数より優先されます。

| 引数 | 説明 | 例 |
| ------------- | ------------- | ------------- |
| help | すべての引数のヘルプ情報を表示します。 | `--help` |
| version | MATLAB MCP Core Server のバージョンを表示します。 | `--version` |
| matlab-root | 起動する MATLAB の絶対パス。パスに `/bin` は含めないでください。既定では、システム パス上で最初に見つかった MATLAB が使用されます。 | Windows: `--matlab-root=C:\\Program Files\\MATLAB\\R2026a` <br><br> Linux/macOS: `--matlab-root=/home/usr/MATLAB/R2026a`<br><br>環境変数の場合: `MW_MCP_SERVER_MATLAB_ROOT=/home/usr/MATLAB/R2026a` |
| initialize-matlab-on-startup | サーバーの起動と同時に MATLAB を初期化するには、この引数を `true` に設定します。既定では、最初のツールが呼び出されたときにのみ MATLAB が起動します。 | `--initialize-matlab-on-startup=true` |
| initial-working-folder | MATLAB の起動フォルダーを指定します。値を指定しない場合、MATLAB は AI アプリケーションの最初の [Root (MCP)](https://modelcontextprotocol.io/specification/latest/client/roots) のパスで起動します。Root が定義されていない場合、MATLAB は次の場所で起動します。 <br> <ul><li>Linux: `/home/username` </li><li> Windows: `C:\Users\username\Documents`</li><li>Mac: `/Users/username/Documents`</li></ul> | Windows: `--initial-working-folder=C:\\Users\\username\\MyProject` <br><br> Linux/macOS: `--initial-working-folder=/Users/username/MyProject` |
| matlab-display-mode | MATLAB デスクトップを表示するかどうかを指定します。`desktop` モード (既定) では MATLAB デスクトップが表示されます。`nodesktop` モードでは、MATLAB デスクトップを使用せずに AI アプリケーションでのみ MATLAB を使用します。ただし、`nodesktop` モードでも、グラフィカル インターフェイスを必要とするコマンド (`edit`、`open`、`open_system`、`uifigure`、`appdesigner` など) の場合はデスクトップ上に MATLAB ウィンドウが開きます。 | `--matlab-display-mode=nodesktop` |
| matlab-session-mode | MCP サーバーが新しい MATLAB を起動するか、既存の MATLAB セッションに接続 (MATLAB R2023a 以降でサポート) するかを指定します。既定は **`auto`** モードです。<br><br>**`new` モード:** MCP サーバーが新しい MATLAB セッションを起動します。<br><br>**`auto` モード (既定):** サーバーは既存の MATLAB セッションへの接続を試みます。既存のセッションを使用するには、以下の手順で `existing` モードを事前に設定しておく必要があります。既存の MATLAB セッションが見つからない場合は、サーバーは新しいセッションを起動します。<br><br>**`existing` モード:** サーバーは既存の MATLAB セッションへの接続を試みます。このモードを使用するには、以下の手順で事前に MATLAB セッションを設定しておく必要があります。<br><br><ol><li>初めて `existing` モードを使用する場合は、`./matlab-mcp-core-server --setup-matlab` を実行します。<br><br>このコマンドにより、MATLAB MCP Core Server Toolbox という名前のアドオンが MATLAB にインストールされます。このテーブルの他の引数を使用してコマンドをカスタマイズできます。たとえば、ツールボックスのインストールに使用する MATLAB を指定するには、`./matlab-mcp-core-server --setup-matlab --matlab-root=/home/usr/MATLAB/R2026a` を使用できます。<br><br>Claude Desktop の場合、`./matlab-mcp-core-server --setup-matlab` を実行する前に、「[セットアップ](#セットアップ)」の手順に従って MATLAB MCP Core Server バイナリをダウンロードする必要があります。<br><br></li><li>実行中の MATLAB セッションのコマンド ウィンドウで `shareMATLABSession()` を実行します。`--matlab-session-mode=existing` を指定してサーバーを起動すると、MCP サーバーはこの MATLAB に接続します。複数の MATLAB セッションを実行している場合、サーバーは最後に `shareMATLABSession()` コマンドを実行した MATLAB セッションに接続します。<br><br>`shareMATLABSession()` を手動で実行する代わりに、MATLAB の[スタートアップ スクリプト (MathWorks)](https://www.mathworks.com/help/matlab/ref/startup.html) にこのコマンドを追加できます。</li></ol> | `--matlab-session-mode=existing` |
| extension-file | カスタム MCP ツールを使用するには、ツールを定義する JSON ファイルのパスを指定します。複数の拡張ファイルを使用することもできます。カスタム ツールの使用の詳細については、「[MATLAB MCP Core Server でのカスタム ツールの使用](guides/custom-tools.ja.md)」を参照してください。 | <br><br>Windows: `--extension-file=C:\\Users\\name\\my-tools.json` <br><br> Linux/macOS: `--extension-file=/path/to/my-tools.json` <br><br> **複数の拡張ファイルを使用する場合:**<br><br>Windows:`--extension-file=C:\\path\\to\\tools-1.json --extension-file=C:\\path\\to\\tools-2.json`<br><br>Linux/macOS:`--extension-file=/path/to/tools1.json --extension-file=/path/to/tools2.json` <br><br> **環境変数を使用する場合:** <br><br> Windows: `MW_MCP_SERVER_EXTENSION_FILE=C:\Users\name\tools1.json;C:\Users\name\tools2.json` <br><br> Linux/macOS: `MW_MCP_SERVER_EXTENSION_FILE=/path/to/tools1.json:/path/to/tools2.json` |
| log-folder | MCP サーバーがログ ファイルを保存するフォルダーを指定します。指定しない場合、サーバーはオペレーティング システムの既定の一時フォルダーを使用します。 | Windows: `--log-folder=C:\\Users\\name\\AppData\\Local\\Temp` <br><br> Linux/macOS: `--log-folder=/tmp/my-logs` |
| log-level | MCP サーバーのログ レベル。有効な値は、詳細度の高い順に `debug`、`info`、`warn`、`error` です。 | `--log-level=debug` |
| disable-telemetry | 匿名データ収集を無効にするには、この引数を `true` に設定します。詳細については、「[データ収集](#データ収集)」を参照してください。 | `--disable-telemetry=true` |

**複数の拡張ファイル**

Windows:
```
--extension-file=C:\\path\\to\\my-tools.json --extension-file=C:\\path\\to\\my-other-tools.json
```

Linux および macOS:
```
--extension-file=/path/to/my-tools.json --extension-file=/path/to/my-other-tools.json
```

**環境変数**

Windows:
```
MW_MCP_SERVER_EXTENSION_FILE=C:\Users\name\my-tools.json;C:\Users\name\my-other-tools.json
```

Linux および macOS:
```
MW_MCP_SERVER_EXTENSION_FILE=/path/to/my-tools.json:/path/to/my-other-tools.json
```
## ツール

1. `detect_matlab_toolboxes`
    - インストールされている MATLAB とツールボックスに関する情報 (バージョン番号を含む) を返します。

1. `check_matlab_code`
    - MATLAB スクリプトの静的コード解析を実行します。コーディング スタイル、潜在的なエラー、非推奨の関数、パフォーマンスの問題、およびベスト プラクティス違反に関する警告を返します。スクリプトを実行せずにコードの品質の問題を特定する、非破壊かつ読み取り専用の操作です。
    - 入力:
        - `script_path` (string): 解析する MATLAB スクリプト ファイルの絶対パス。有効な `.m` ファイルである必要があります。解析中にファイルは変更されません。例: `C:\Users\username\matlab\myFunction.m` または `/home/user/scripts/analysis.m`。

1. `evaluate_matlab_code`
    - MATLAB コードの文字列を評価し、出力を返します。
    - 入力:
        - `code` (string): 評価する MATLAB コード。
        - `project_path` (string): プロジェクト ディレクトリの絶対パス。MATLAB はこのディレクトリを現在の作業フォルダーとして設定します。例: `C:\Users\username\matlab-project` または `/home/user/research`。

1. `run_matlab_file`
    - MATLAB スクリプトを実行し、出力を返します。スクリプトは有効な `.m` ファイルである必要があります。
    - 入力:
        - `script_path` (string): 実行する MATLAB スクリプト ファイルの絶対パス。有効な `.m` ファイルである必要があります。例: `C:\Users\username\projects\analysis.m` または `/home/user/matlab/simulation.m`。

1. `run_matlab_test_file`
    - MATLAB テスト スクリプトを実行し、包括的なテスト結果を返します。MATLAB テスト フレームワークの規則に従った MATLAB ユニット テスト ファイル用に設計されています。
    - 入力:
        - `script_path` (string): MATLAB テスト スクリプト ファイルの絶対パス。MATLAB ユニット テストを含む有効な `.m` ファイルである必要があります。例: `C:\Users\username\tests\testMyFunction.m` または `/home/user/matlab/tests/test_analysis.m`。

## リソース

MCP サーバーは、AI アプリケーションが MATLAB コードを記述する際に役立つ[リソース (MCP)](https://modelcontextprotocol.io/specification/latest/server/resources) を提供します。リソースの使用手順については、使用する AI アプリケーションのリソースに関するドキュメンテーションを参照してください。

1. `matlab_coding_guidelines`
    - コードの可読性、保守性、およびコラボレーションを向上させるための包括的な MATLAB コーディング規約を提供します。命名規則、書式、コメント、パフォーマンスの最適化、およびエラー処理についてのガイドラインが含まれています。
    - URI: `guidelines://coding`
    - MIME タイプ: `text/markdown`
    - ソース: [MATLAB Coding Standards (GitHub)](https://github.com/matlab/rules/blob/main/matlab-coding-standards.md)

1. `plain_text_live_code_guidelines`
    - プレーン テキストのライブ コードの `.m` ファイル形式を使用したライブ スクリプトの生成に関するルールとガイドラインを提供します。バージョン管理および AI 支援開発に適しています。プレーン テキストのライブ スクリプトを実行するには、MATLAB R2025a 以降が必要です。詳細については、「[ライブ コード ファイル形式 (MathWorks)](https://www.mathworks.com/help/matlab/matlab_prog/plain-text-file-format-for-live-scripts.html)」を参照してください。
    - URI: `guidelines://plain-text-live-code`
    - MIME タイプ: `text/markdown`
    - ソース: [Plain Text Live Code Generation (GitHub)](https://github.com/matlab/rules/blob/main/live-script-generation.md)

## データ収集

MATLAB MCP Core Server は、サーバーの使用状況に関して、完全に匿名化された情報を収集して MathWorks に送信する場合があります。このデータ収集は MathWorks の製品改善に役立てるためのもので、既定で有効になっています。データ収集をオプトアウトするには、引数 `--disable-telemetry` を `true` に設定します。

## セキュリティに関する考慮事項

MATLAB MCP Core Server を使用する際は、すべてのツール呼び出しを実行前に十分に確認および検証してください。重要な操作では必ず人間が介在するようにし、呼び出しが期待どおりに動作することを確認してから実行してください。詳細については、「[User Interaction Model (MCP)](https://modelcontextprotocol.io/specification/latest/server/tools#user-interaction-model)」および「[Security Considerations (MCP)](https://modelcontextprotocol.io/specification/latest/server/tools#security-considerations)」を参照してください。

## ライセンスと使用について

ライセンスについては、本 GitHub リポジトリの [LICENSE.md](../LICENSE.md) ファイルをご確認ください。

MCP サーバーは、MathWorks Software License Agreement に従って MATLAB と併用する場合にのみ使用が許可されており、複数のユーザーで共有してはなりません。共有サーバーまたは集中型サーバーを使用する必要がある場合は、MathWorks にお問い合わせください。

## サポートへのお問い合わせ

MathWorks は本リポジトリの活用およびフィードバックの提供をお勧めしています。テクニカル サポートのリクエストや機能強化のご要望は、[GitHub issue を作成](https://github.com/matlab/matlab-mcp-core-server/issues)するか、[MathWorks テクニカル サポート](https://www.mathworks.com/support/contact_us.html)にお問い合わせください。

---

Copyright 2025-2026 The MathWorks, Inc.

---
