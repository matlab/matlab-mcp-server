# Use Custom Tools with MATLAB MCP Core Server

This guide shows how to use custom tools with the MATLAB MCP Core Server. 

You can expose any MATLAB function as an MCP tool defined in a JSON file. The server loads your tool definitions at startup and registers them alongside the built-in tools. When your AI application calls a custom tool, the server executes the MATLAB function and returns the command window output. The MATLAB function must be on the MATLAB path. To update your tool definitions, edit the extension file and restart the server.

Custom tool arguments support `string`, `number`, `integer`, and `boolean` data types. 

## Table of Contents
- [Get Started](#get-started)
- [Extension File Format](#extension-file-format)
    - [Tools](#tools)
    - [Signatures](#signatures)
    - [inputSchema](#inputschema)
    - [Supported Property Types](#supported-property-types)
    - [Annotations](#annotations)

## Get Started

This example shows how to add a custom tool to the MCP server. First you define the MATLAB function you want the tool to use. Then you define the tool in a JSON file, along with a signature for your MATLAB function. Then you pass this JSON file to the MCP server using the `extension file` argument.

1. Create the MATLAB function `greet_user.m`:
    ```matlab
    function greet_user(name, age)
        disp("Hello " + name + ", you are " + age + " years old!");
    end
    ```
    and add it to the MATLAB path:
    ```
    addpath('/path/to/your/functions')
    ```

1. Create an extension file `my-tools.json`, containing your tool definition and function signature:

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

1. Start the server with the [`--extension-file` argument](../README.md#arguments):

    ```
    --extension-file=/path/to/my-tools.json
    ```

    Your tool is now available. Call it with the `name = "Alice"` and `age = 30` to execute `greet_user("Alice", 30)` in MATLAB.

## Extension File Format

The extension file has two top-level fields: `tools` (an array) and `signatures` (an object).

### Tools

Each entry in the `tools` array defines an MCP tool following the [Tool Schema (MCP)](https://modelcontextprotocol.io/specification/latest/schema#tool):

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Unique tool identifier |
| `title` | Yes | Human-readable title |
| `description` | Yes | Explains what the tool does to the AI model |
| `inputSchema` | Yes | JSON Schema defining the tool's input arguments |
| `annotations` | No | MCP tool annotations for the AI client |

### Signatures

The `signatures` object maps each tool name to its MATLAB function call details. The signatures are separate from the tool definitions because they are not part of the MCP specification.

Every tool must have a corresponding signature entry:

| Field | Required | Description |
|-------|----------|-------------|
| `function` | Yes | MATLAB function to call (must be on the MATLAB path) |
| `input.order` | Yes | Array specifying the order arguments are passed to the function |

The `input.order` array must contain exactly the same entries as the `inputSchema.properties` keys. This determines the positional order of arguments in the MATLAB function call.

### inputSchema

Each tool's inputSchema field defines its arguments using the [JSON Schema](https://json-schema.org/) format:

```json
{
  "type": "object",
  "properties": {
    "argName": { "type": "string", "description": "Description of argument" }
  },
  "required": ["argName"]
}
```

| Constraint | Detail |
|------------|--------|
| Top-level `type` | Must be `"object"` |
| Property types | `string`, `number`, `integer`, `boolean` |
| `required` | Array of required argument names |

### Supported Property Types

| Type | JSON Input | MATLAB Receives |
|------|------------|-----------------|
| `string` | `"hello"` | `"hello"` |
| `number` | `42` or `3.14` | `42` or `3.14` |
| `integer` | `42` | `42` |
| `boolean` | `true` / `false` | `true` / `false` |

Arrays and objects are not currently supported.

### Annotations

The optional `annotations` field provides hints to AI clients about the tool's behavior. For details, see the [ToolAnnotations Schema (MCP)](https://modelcontextprotocol.io/specification/latest/schema#toolannotations). For example:

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

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `readOnlyHint` | boolean | `false` | Tool does not modify its environment |
| `destructiveHint` | boolean | `true` | Tool may perform destructive updates |
| `idempotentHint` | boolean | `false` | Repeated calls with same arguments have no additional effect |
| `openWorldHint` | boolean | `true` | Tool may interact with external entities |

---

Copyright 2026 The MathWorks, Inc.

---
