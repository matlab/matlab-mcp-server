<!--
Source English Markdown:
- File: ./guides/custom-tools.md
- Branch: main
- Commit: a74f0896e40f57ed6755bcf6ac6511dd8bdb6cd2
-->

# MATLAB MCP Core Server에서 사용자 지정 툴 사용하기

<p align="center">
  <a href="../../guides/custom-tools.md">English</a> •
  <a href="custom-tools.es.md">Español</a> •
  <a href="custom-tools.ja.md">日本語</a> •
  한국어 •
  <a href="custom-tools.zh-cn.md">简体中文</a>
</p>

이 가이드에서는 MATLAB MCP Core Server에서 사용자 지정 툴을 사용하는 방법을 설명합니다.

MATLAB 함수에 대해 JSON 파일로 툴 정의를 만들어 MCP 툴로 노출할 수 있습니다. 서버는 시작 시 이 정의들을 불러와 기본 제공 툴과 함께 등록합니다. AI 애플리케이션이 사용자 지정 툴을 호출하면 서버가 MATLAB 함수를 실행하고 명령 창 출력을 반환합니다. MATLAB 함수는 MATLAB 경로에 있어야 합니다. 툴 정의를 업데이트하려면 확장 파일을 편집한 후 서버를 다시 시작하십시오.

사용자 지정 툴 인수로 지원되는 데이터형은 `string`, `number`, `integer`, `boolean`입니다.

## 목차
- [시작하기](#시작하기)
- [확장 파일 형식](#확장-파일-형식)
    - [툴](#툴)
    - [시그니처](#시그니처)
    - [inputSchema](#inputschema)
    - [지원되는 속성 유형](#지원되는-속성-유형)
    - [주석](#주석)

## 시작하기

이 예제에서는 MCP 서버에 사용자 지정 툴을 추가하는 방법을 보여줍니다. 먼저 툴로 사용할 MATLAB 함수를 정의합니다. 그런 다음 JSON 파일에 툴과 MATLAB 함수의 시그니처를 정의합니다. 마지막으로 `extension file` 인수를 사용하여 이 JSON 파일을 MCP 서버에 전달합니다. 여러 확장 파일을 전달할 수도 있습니다.

1. MATLAB 함수 `greet_user.m`을 작성합니다.
    ```matlab
    function greet_user(name, age)
        disp("Hello " + name + ", you are " + age + " years old!");
    end
    ```
    MATLAB 경로에 추가합니다.
    ```
    addpath('/path/to/your/functions')
    ```

1. 툴 정의와 함수 시그니처를 포함하는 확장 파일 `my-tools.json`을 작성합니다.

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

1. [`--extension-file` 인수](../README.ko.md#인수)를 사용하여 서버를 시작합니다.

    Windows:
    ```
    --extension-file=C:\\path\\to\\my-tools.json
    ```

    Linux 및 macOS:
    ```
    --extension-file=/path/to/my-tools.json
    ```

    이제 툴을 사용할 수 있습니다. `name = "Alice"` 및 `age = 30`으로 호출하면 MATLAB에서 `greet_user("Alice", 30)`이 실행됩니다.

1. 여러 확장 파일을 사용할 수도 있으며, 환경 변수를 사용하여 MCP 서버에 확장 파일을 제공할 수도 있습니다.

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

## 확장 파일 형식

확장 파일은 `tools`(배열)와 `signatures`(객체)라는 2개의 최상위 필드를 가집니다.

### 툴

`tools` 배열의 각 항목은 [Tool Schema (MCP)](https://modelcontextprotocol.io/specification/latest/schema#tool)에 따라 MCP 툴을 정의합니다.

| 필드 | 필수 | 설명 |
|-------|----------|-------------|
| `name` | 예 | 고유한 툴 식별자 |
| `title` | 예 | 사람이 읽을 수 있는 제목 |
| `description` | 예 | AI 모델에 툴의 기능을 설명 |
| `inputSchema` | 예 | 툴의 입력 인수를 정의하는 JSON Schema |
| `annotations` | 아니요 | AI 클라이언트를 위한 MCP 툴 주석 |

### 시그니처

`signatures` 객체는 각 툴 이름을 MATLAB 함수 호출 세부 정보에 매핑합니다. 시그니처는 MCP 사양의 일부가 아니기 때문에 툴 정의와 분리되어 있습니다.

모든 툴은 해당하는 시그니처 항목을 가져야 합니다.

| 필드 | 필수 | 설명 |
|-------|----------|-------------|
| `function` | 예 | 호출할 MATLAB 함수(MATLAB 경로에 있어야 합니다) |
| `input.order` | 예 | 함수에 전달되는 인수 순서를 지정하는 배열 |

`input.order` 배열은 `inputSchema.properties` 키들과 정확히 동일한 항목들을 포함해야 합니다. 이는 MATLAB 함수 호출 시 인수의 위치 순서를 결정합니다.

### inputSchema

각 툴의 inputSchema 필드는 [JSON Schema](https://json-schema.org/) 형식을 사용하여 인수를 정의합니다.

```json
{
  "type": "object",
  "properties": {
    "argName": { "type": "string", "description": "Description of argument" }
  },
  "required": ["argName"]
}
```

| 제약 조건 | 세부 사항 |
|------------|--------|
| 최상위 `type` | `"object"`여야 합니다 |
| 속성 유형 | `string`, `number`, `integer`, `boolean` |
| `required` | 필수 인수 이름의 배열 |

### 지원되는 속성 유형

| 유형 | JSON 입력 | MATLAB이 수신하는 값 |
|------|------------|-----------------|
| `string` | `"hello"` | `"hello"` |
| `number` | `42` 또는 `3.14` | `42` 또는 `3.14` |
| `integer` | `42` | `42` |
| `boolean` | `true` / `false` | `true` / `false` |

배열과 객체는 현재 지원되지 않습니다.

### 주석

선택적 `annotations` 필드는 툴의 동작에 대한 힌트를 AI 클라이언트에 제공합니다. 자세한 내용은 [ToolAnnotations Schema (MCP)](https://modelcontextprotocol.io/specification/latest/schema#toolannotations)를 참조하십시오. 예:

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

| 필드 | 유형 | 기본값 | 설명 |
|-------|------|---------|-------------|
| `readOnlyHint` | boolean | `false` | 툴이 환경을 수정하지 않음 |
| `destructiveHint` | boolean | `true` | 툴이 파괴적인 업데이트를 수행할 수 있음 |
| `idempotentHint` | boolean | `false` | 동일한 인수로 반복 호출해도 추가 효과 없음 |
| `openWorldHint` | boolean | `true` | 툴이 외부 엔티티와 상호작용할 수 있음 |

---

Copyright 2026 The MathWorks, Inc.

---
