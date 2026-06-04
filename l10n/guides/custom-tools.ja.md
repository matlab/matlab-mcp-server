<!--
Source English Markdown:
- File: ./guides/custom-tools.md
- Branch: main
- Commit: a74f0896e40f57ed6755bcf6ac6511dd8bdb6cd2
-->

# MATLAB MCP Core Server でのカスタム ツールの使用

<p align="center">
  <a href="../../guides/custom-tools.md">English</a> •
  <a href="custom-tools.es.md">Español</a> •
  日本語 •
  <a href="custom-tools.ko.md">한국어</a> •
  <a href="custom-tools.zh-cn.md">简体中文</a>
</p>

このガイドでは、MATLAB MCP Core Server でカスタム ツールを使用する方法を説明します。

任意の MATLAB 関数を、JSON ファイルで定義された MCP ツールとして公開できます。サーバーは起動時にツール定義を読み込み、組み込みツールとともに登録します。AI アプリケーションがカスタム ツールを呼び出すと、サーバーは MATLAB 関数を実行し、コマンド ウィンドウの出力を返します。MATLAB 関数は MATLAB パス上に存在する必要があります。ツール定義を更新するには、拡張ファイルを編集してサーバーを再起動してください。

カスタム ツールの引数は `string`、`number`、`integer`、`boolean` のデータ型をサポートしています。

## 目次
- [はじめに](#はじめに)
- [拡張ファイルのフォーマット](#拡張ファイルのフォーマット)
    - [ツール](#ツール)
    - [シグネチャ](#シグネチャ)
    - [inputSchema](#inputschema)
    - [サポートされるプロパティの型](#サポートされるプロパティの型)
    - [アノテーション](#アノテーション)

## はじめに

この例では、MCP サーバーにカスタム ツールを追加する方法を示します。まず、ツールで使用する MATLAB 関数を定義します。次に、JSON ファイルでツールと MATLAB 関数のシグネチャを定義します。最後に、`extension file` 引数を使用してこの JSON ファイルを MCP サーバーに渡します。複数の拡張ファイルを渡すこともできます。

1. MATLAB 関数 `greet_user.m` を作成します。
    ```matlab
    function greet_user(name, age)
        disp("Hello " + name + ", you are " + age + " years old!");
    end
    ```
    その後、MATLAB パスに追加します。
    ```
    addpath('/path/to/your/functions')
    ```

1. ツール定義と関数シグネチャを含む拡張ファイル `my-tools.json` を作成します。

    ```json
    {
      "tools": [
        {
          "name": "greet_user",
          "title": "Greet User",
          "description": "Displays a greeting for the given user",
          "inputSchema": {
            "type": "object",
            "properties": {
              "name": { "type": "string", "description": "Name of the user to greet" },
              "age": { "type": "number", "description": "Age of the user" }
            },
            "required": ["name", "age"]
          }
        }
      ],
      "signatures": {
        "greet_user": {
          "function": "greet_user",
          "input": {
            "order": ["name", "age"]
          }
        }
      }
    }
    ```

1. [`--extension-file` 引数](../README.ja.md#引数)を使用してサーバーを起動します。

   Windows:
    ```
    --extension-file=C:\\path\\to\\my-tools.json
    ```

    Linux および macOS:
    ```
    --extension-file=/path/to/my-tools.json
    ```

    これでツールが利用可能になります。`name = "Alice"` と `age = 30` を指定して呼び出すと、MATLAB で `greet_user("Alice", 30)` が実行されます。

1. 複数の拡張ファイルを使用したり、環境変数を使用して MCP サーバーに拡張ファイルを提供したりすることもできます。

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

## 拡張ファイルのフォーマット

拡張ファイルには、`tools` (配列) と `signatures` (オブジェクト) の 2 つのトップレベル フィールドがあります。

### ツール

`tools` 配列の各エントリは、[Tool スキーマ (MCP)](https://modelcontextprotocol.io/specification/latest/schema#tool) に従って MCP ツールを定義します。

| フィールド | 必須 | 説明 |
|-------|----------|-------------|
| `name` | はい | 一意のツール識別子 |
| `title` | はい | 人間が読めるタイトル |
| `description` | はい | AI モデルに対するツール機能の説明 |
| `inputSchema` | はい | ツールの入力引数を定義する JSON Schema |
| `annotations` | いいえ | AI クライアント向けの MCP ツール アノテーション |

### シグネチャ

`signatures` オブジェクトは、各ツール名と、対応する MATLAB 関数の呼び出し方法をマッピングします。シグネチャは MCP の仕様の一部ではないため、ツール定義とは別にします。

すべてのツールに、対応するシグネチャ エントリが必要です。

| フィールド | 必須 | 説明 |
|-------|----------|-------------|
| `function` | はい | 呼び出す MATLAB 関数 (MATLAB パス上に存在する必要があります) |
| `input.order` | はい | 関数に渡される引数の順序を指定する配列 |

`input.order` 配列には、`inputSchema.properties` のキーと完全に同じエントリを含める必要があります。これにより、MATLAB 関数呼び出し時の引数の位置順が決定されます。

### inputSchema

各ツールの inputSchema フィールドは、[JSON Schema](https://json-schema.org/) フォーマットを使用して引数を定義します。

```json
{
  "type": "object",
  "properties": {
    "argName": { "type": "string", "description": "Description of argument" }
  },
  "required": ["argName"]
}
```

| 制約 | 詳細 |
|------------|--------|
| トップレベルの `type` | `"object"` である必要があります |
| プロパティの型 | `string`、`number`、`integer`、`boolean` |
| `required` | 必須引数名の配列 |

### サポートされるプロパティの型

| 型 | JSON 入力 | MATLAB が受け取る値 |
|------|------------|-----------------|
| `string` | `"hello"` | `"hello"` |
| `number` | `42` または `3.14` | `42` または `3.14` |
| `integer` | `42` | `42` |
| `boolean` | `true` / `false` | `true` / `false` |

配列とオブジェクトは現在サポートされていません。

### アノテーション

オプションの `annotations` フィールドは、ツールの動作に関するヒントを AI クライアントに提供します。詳細については、[ToolAnnotations スキーマ (MCP)](https://modelcontextprotocol.io/specification/latest/schema#toolannotations) を参照してください。以下に例を示します。

```json
{
  "annotations": {
    "readOnlyHint": false,
    "destructiveHint": true,
    "idempotentHint": false,
    "openWorldHint": true
  }
}
```

| フィールド | 型 | 既定 | 説明 |
|-------|------|---------|-------------|
| `readOnlyHint` | boolean | `false` | ツールは環境を変更しません |
| `destructiveHint` | boolean | `true` | ツールは破壊的な更新を行う可能性があります |
| `idempotentHint` | boolean | `false` | 同じ引数で繰り返し呼び出しても追加の効果はありません |
| `openWorldHint` | boolean | `true` | ツールは外部エンティティと対話する可能性があります |

---

Copyright 2026 The MathWorks, Inc.

---
