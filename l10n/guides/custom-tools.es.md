 <!--
-Source English Markdown:
-- File: ./guides/custom-tools.md
-- Branch: main
-- Commit: a74f0896e40f57ed6755bcf6ac6511dd8bdb6cd2
 -->

# Usar herramientas personalizadas con MATLAB MCP Core Server

<p align="center">
  <a href="../../guides/custom-tools.md">English</a> •
  Español •
  <a href="custom-tools.ja.md">日本語</a> •
  <a href="custom-tools.ko.md">한국어</a> •
  <a href="custom-tools.zh-cn.md">简体中文</a>
</p>

Esta guía muestra cómo usar herramientas personalizadas con MATLAB MCP Core Server. 

Puede exponer cualquier función de MATLAB como herramientas MCP definidas en archivos JSON. El servidor carga las definiciones de herramientas al inicio y las registra junto con las herramientas integradas. Cuando su aplicación de IA llama a una herramienta personalizada, el servidor ejecuta la función de MATLAB y devuelve la salida de la ventana de comandos. La función de MATLAB debe estar en el path de MATLAB. Para actualizar las definiciones de herramientas, edite los archivos de extensión y reinicie el servidor.

Los argumentos de herramientas personalizadas admiten los tipos de datos `string`, `number`, `integer` y `boolean`. 

## Tabla de contenido
- [Primeros pasos](#primeros-pasos)
- [Formato del archivo de extensión](#formato-del-archivo-de-extensión)
    - [Herramientas](#herramientas)
    - [Firmas](#firmas)
    - [inputSchema](#inputschema)
    - [Tipos de propiedades admitidos](#tipos-de-propiedades-admitidos)
    - [Anotaciones](#anotaciones)

## Primeros pasos

Este ejemplo muestra cómo agregar una herramienta personalizada al servidor MCP. Primero se define la función de MATLAB que la herramienta utilizará. Luego se define la herramienta en un archivo JSON, junto con una firma para la función de MATLAB. Después se pasa este archivo JSON al servidor MCP usando el argumento `extension file`. También se pueden pasar múltiples archivos de extensión.

1. Cree la función de MATLAB `greet_user.m`:
    ```matlab
    function greet_user(name, age)
        disp("Hello " + name + ", you are " + age + " years old!");
    end
    ```
    y agréguela al path de MATLAB:
    ```
    addpath('/path/a/sus/funciones')
    ```

1. Cree un archivo de extensión `my-tools.json`, que contenga la definición de la herramienta y la firma de la función:

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

1. Inicie el servidor con el [argumento `--extension-file`](../README.es.md#argumentos).

    Windows:
    ```
    --extension-file=C:\\path\\a\\mis-herramientas.json
    ```

    Linux y macOS:
    ```
    --extension-file=/path/a/mis-herramientas.json
    ```

    La herramienta ya está disponible. Llámela con `name = "Alice"` y `age = 30` para ejecutar `greet_user("Alice", 30)` en MATLAB.

1.  También puede usar múltiples archivos de extensión y proporcionar archivos de extensión al servidor MCP mediante variables de entorno. 

    **Múltiples archivos de extensión**

    Windows:
    ```
    --extension-file=C:\\path\\a\\mis-herramientas.json --extension-file=C:\\path\\a\\mis-otras-herramientas.json
    ```

    Linux y macOS:
    ```
    --extension-file=/path/a/mis-herramientas.json --extension-file=/path/a/mis-otras-herramientas.json
    ```
 
    **Variables de entorno**

    Windows:
    ```
    MW_MCP_SERVER_EXTENSION_FILE=C:\Users\name\mis-herramientas.json;C:\Users\name\mis-otras-herramientas.json
    ```

    Linux y macOS:
    ```
    MW_MCP_SERVER_EXTENSION_FILE=/path/a/mis-herramientas.json:/path/a/mis-otras-herramientas.json
    ```

## Formato del archivo de extensión

El archivo de extensión tiene dos campos de nivel superior: `tools` (un arreglo) y `signatures` (un objeto).

### Herramientas

Cada entrada en el arreglo `tools` define una herramienta MCP siguiendo el [Tool Schema (MCP)](https://modelcontextprotocol.io/specification/latest/schema#tool):

| Campo | Requerido | Descripción |
|-------|----------|-------------|
| `name` | Sí | Identificador único de la herramienta |
| `title` | Sí | Título legible para humanos |
| `description` | Sí | Explica al modelo de IA lo que hace la herramienta |
| `inputSchema` | Sí | JSON Schema que define los argumentos de entrada de la herramienta |
| `annotations` | No | Anotaciones de herramientas MCP para el cliente de IA |

### Firmas

El objeto `signatures` asigna cada nombre de herramienta a los detalles de la llamada a la función de MATLAB. Las firmas están separadas de las definiciones de herramientas porque no forman parte de la especificación MCP.

Cada herramienta debe tener una entrada de firma correspondiente:

| Campo | Requerido | Descripción |
|-------|----------|-------------|
| `function` | Sí | Función de MATLAB que se desea llamar (debe estar en el path de MATLAB) |
| `input.order` | Sí | Arreglo que especifica el orden en que los argumentos se pasan a la función |

El arreglo `input.order` debe contener exactamente las mismas entradas que las claves de `inputSchema.properties`. Esto determina el orden posicional de los argumentos en la llamada a la función de MATLAB.

### inputSchema

El campo inputSchema de cada herramienta define sus argumentos usando el formato [JSON Schema](https://json-schema.org/):

```json
{
  "type": "object",
  "properties": {
    "argName": { "type": "string", "description": "Description of argument" }
  },
  "required": ["argName"]
}
```

| Restricción | Detalle |
|------------|--------|
| `type` de nivel superior | Debe ser `"object"` |
| Tipos de propiedades | `string`, `number`, `integer`, `boolean` |
| `required` | Arreglo de nombres de argumentos requeridos |

### Tipos de propiedades admitidos

| Tipo | Entrada JSON | MATLAB recibe |
|------|------------|-----------------|
| `string` | `"hello"` | `"hello"` |
| `number` | `42` o `3.14` | `42` o `3.14` |
| `integer` | `42` | `42` |
| `boolean` | `true` / `false` | `true` / `false` |

Los arreglos y objetos no están admitidos actualmente.

### Anotaciones

El campo opcional `annotations` proporciona indicaciones a los clientes de IA sobre el comportamiento de la herramienta. Para obtener más información, consulte el [ToolAnnotations Schema (MCP)](https://modelcontextprotocol.io/specification/latest/schema#toolannotations). Por ejemplo:

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

| Campo | Tipo | Valor predeterminado | Descripción |
|-------|------|---------|-------------|
| `readOnlyHint` | boolean | `false` | La herramienta no modifica su entorno |
| `destructiveHint` | boolean | `true` | La herramienta puede realizar actualizaciones destructivas |
| `idempotentHint` | boolean | `false` | Las llamadas repetidas con los mismos argumentos no tienen efecto adicional |
| `openWorldHint` | boolean | `true` | La herramienta puede interactuar con entidades externas |

---

Copyright 2026 The MathWorks, Inc.

---
