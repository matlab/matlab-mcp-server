 <!--
-Source English Markdown:
-- File: ./guides/custom-tools.md
-- Branch: main
-- Commit: a74f0896e40f57ed6755bcf6ac6511dd8bdb6cd2
 -->

# 在 MATLAB MCP Core Server 中使用自定义工具

<p align="center">
  <a href="../../guides/custom-tools.md">English</a> •
  <a href="custom-tools.es.md">Español</a> •
  <a href="custom-tools.ja.md">日本語</a> •
  <a href="custom-tools.ko.md">한국어</a> •
  简体中文
</p>

本指南介绍如何在 MATLAB MCP Core Server 中使用自定义工具。 

您可以将任意 MATLAB 函数公开为 MCP 工具，通过 JSON 文件进行定义。服务器在启动时加载您的工具定义，并将其与内置工具一起注册。当您的 AI 应用程序调用自定义工具时，服务器会执行 MATLAB 函数并返回命令行窗口输出。MATLAB 函数必须位于 MATLAB 路径上。要更新工具定义，请编辑扩展文件并重新启动服务器。

自定义工具参量支持 `string`、`number`、`integer` 和 `boolean` 数据类型。 

## 目录
- [快速入门](#快速入门)
- [扩展文件格式](#扩展文件格式)
    - [工具](#工具)
    - [签名](#签名)
    - [inputSchema](#inputschema)
    - [支持的属性类型](#支持的属性类型)
    - [注解](#注解)

## 快速入门

本示例介绍如何向 MCP 服务器添加自定义工具。首先定义您希望工具使用的 MATLAB 函数。然后在 JSON 文件中定义工具以及 MATLAB 函数的签名。然后使用 `extension file` 参量将此 JSON 文件传递给 MCP 服务器。您还可以传递多个扩展文件。

1. 创建 MATLAB 函数 `greet_user.m`：
    ```matlab
    function greet_user(name, age)
        disp("Hello " + name + ", you are " + age + " years old!");
    end
    ```
    and add it to the MATLAB path:
    ```
    addpath('/path/to/your/functions')
    ```

1. 创建扩展文件 `my-tools.json`，其中包含工具定义和函数签名：

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

1. 使用 [`--extension-file` 参量](../README.zh-cn.md#参量)启动服务器。

    Windows: 
    ```
    --extension-file=C:\\path\\to\\my-tools.json
    ```

    Linux 和 macOS: 
    ```
    --extension-file=/path/to/my-tools.json
    ```

    您的工具现已可用。使用 `name = "Alice"` 和 `age = 30` 进行调用，以在 MATLAB 中执行 `greet_user("Alice", 30)`。

1.  您还可以使用多个扩展文件，并使用环境变量向 MCP 服务器提供扩展文件。 

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

## 扩展文件格式

扩展文件有两个顶层字段: `tools`(一个数组) 和 `signatures` (一个对象)。

### 工具

`tools` 数组中的每个条目按照[工具架构 (MCP)](https://modelcontextprotocol.io/specification/latest/schema#tool) 定义一个 MCP 工具: 

| 字段 | 必需 | 说明 |
|-------|----------|-------------|
| `name` | 是 | 唯一工具标识符 |
| `title` | 是 | 人类可读标题 |
| `description` | 是 | 向 AI 模型说明工具的功能 |
| `inputSchema` | 是 | 定义工具输入参量的 JSON 架构 |
| `annotations` | 否 | 供 AI 客户端使用的 MCP 工具注解 |

### 签名

`signatures` 对象将每个工具名称映射到其 MATLAB 函数调用详细信息。签名与工具定义分开，因为它们不属于 MCP 规范。

每个工具都必须有对应的签名条目: 

| 字段 | 必需 | 说明 |
|-------|----------|-------------|
| `function` | 是 | 要调用的 MATLAB 函数 (必须位于 MATLAB 路径上) |
| `input.order` | 是 | 指定传递给函数的参量顺序的数组 |

`input.order` 数组必须包含与 `inputSchema.properties` 键完全相同的条目。这决定了 MATLAB 函数调用中参量的位置顺序。

### inputSchema

每个工具的 inputSchema 字段使用 [JSON 架构](https://json-schema.org/)格式定义其参量：

```json
{
  "type": "object",
  "properties": {
    "argName": { "type": "string", "description": "Description of argument" }
  },
  "required": ["argName"]
}
```

| 约束 | 详细信息 |
|------------|--------|
| 顶层 `type` | 必须为 `"object"` |
| 属性类型 | `string`、`number`、`integer`、`boolean` |
| `required` | 必需参量名称的数组 |

### 支持的属性类型

| 类型 | JSON 输入 | MATLAB 接收 |
|------|------------|-----------------| 
| `string` | `"hello"` | `"hello"` |
| `number` | `42` 或 `3.14` | `42` 或 `3.14` |
| `integer` | `42` | `42` |
| `boolean` | `true` / `false` | `true` / `false` |

当前不支持数组和对象。

### 注解

可选的 `annotations` 字段向 AI 客户端提供有关工具行为的提示。有关详细信息，请参阅 [ToolAnnotations Schema (MCP)](https://modelcontextprotocol.io/specification/latest/schema#toolannotations)。例如：

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

| 字段 | 类型 | 默认值 | 说明 |
|-------|------|---------|-------------|
| `readOnlyHint` | boolean | `false` | 工具不会修改其环境 |
| `destructiveHint` | boolean | `true` | 工具可能执行破坏性更新 |
| `idempotentHint` | boolean | `false` | 使用相同参量重复调用不会产生额外效果 |
| `openWorldHint` | boolean | `true` | 工具可能与外部实体交互 |

---

Copyright 2026 The MathWorks, Inc.

---
