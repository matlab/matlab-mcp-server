<!--
Source English Markdown:
- File: ./README.md
- Branch: main
- Commit: 9e25655af07ee633bcb5fe5bea5a0c9a844e7041
-->

# MATLAB MCP Core Server

<p align="center">
  <a href="../README.md">English</a> •
  Español •
  <a href="README.ja.md">日本語</a> •
  <a href="README.ko.md">한국어</a> •
  <a href="README.zh-cn.md">简体中文</a>
</p>

Ejecute MATLAB® con aplicaciones de IA mediante el servidor oficial MATLAB MCP Server de MathWorks®. MATLAB MCP Core Server permite a las aplicaciones de IA:

- Iniciar y cerrar MATLAB.
- Escribir y ejecutar código de MATLAB.
- Evaluar el estilo y la corrección del código de MATLAB.

Para ayudar a su agente a utilizar MATLAB y Simulink, puede usar habilidades de [MATLAB Agentic Toolkit (GitHub)](https://github.com/matlab/matlab-agentic-toolkit) y [Simulink Agentic Toolkit (GitHub)](https://github.com/matlab/simulink-agentic-toolkit), que también pueden instalar este servidor MCP automáticamente. 

## Tabla de contenido

- [Configuración](#configuración)
  - [Claude Code](#claude-code)
  - [Claude Desktop](#claude-desktop)
  - [GitHub Copilot en Visual Studio Code](#github-copilot-en-visual-studio-code)
- [Argumentos](#argumentos)
- [Herramientas](#herramientas)
- [Recursos](#recursos)
- [Recopilación de datos](#recopilación-de-datos)
- [Consideraciones de seguridad](#consideraciones-de-seguridad)
- [Licencia y uso](#licencia-y-uso)
- [Contactar con soporte](#contactar-con-soporte)

## Configuración

1. Instale [MATLAB (MathWorks)](https://www.mathworks.com/help/install/ug/install-products-with-internet-connection.html) 2021a o posterior y agréguelo al PATH del sistema. MATLAB MCP Core Server es compatible con las versiones de MATLAB de los últimos cinco años.
1. Para configurar MATLAB MCP Core Server para Claude Desktop, vaya directamente a las instrucciones de [Claude Desktop](#claude-desktop). Para configurar el servidor para otras aplicaciones, siga estas instrucciones:
   
   - Para Windows o Linux, [**descargue la última versión**](https://github.com/matlab/matlab-mcp-core-server/releases/latest). (Alternativamente, puede **compilar desde el código fuente**: instale [Go](https://go.dev/doc/install) y compile el binario con `go install github.com/matlab/matlab-mcp-core-server/cmd/matlab-mcp-core-server@latest`).
    
   - Para macOS, primero descargue la última versión ejecutando el siguiente comando en su terminal:
     - Para procesadores Apple silicon, ejecute:
          ```sh
          curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maca64
          ```
      - Para procesadores Intel, ejecute:
          ```sh
          curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maci64
          ```
      A continuación, otorgue permisos de ejecución al binario descargado para poder ejecutar MATLAB MCP Core Server:

      ```sh
      chmod +x ~/Downloads/matlab-mcp-core-server
      ```

1. Agregue MATLAB MCP Core Server a su aplicación de IA. Puede encontrar instrucciones para agregar servidores MCP en la documentación de su aplicación de IA. Para instrucciones de ejemplo sobre el uso de Claude Code®, Claude Desktop® y GitHub Copilot en Visual Studio® Code, consulte la siguiente información. Tenga en cuenta que puede personalizar el servidor especificando [argumentos](#argumentos) opcionales.

### Claude Code

En su terminal, ejecute el siguiente comando, recordando insertar la ruta completa al binario del servidor que obtuvo en la configuración:

```sh
claude mcp add --transport stdio matlab -- /rutacompleta/al/binario-de-matlab-mcp-core-server
```

Puede personalizar el servidor especificando [argumentos](#argumentos) opcionales. Observe el separador `--` entre las opciones de Claude Code y los argumentos del servidor:

```sh
claude mcp add --transport stdio matlab -- /rutacompleta/al/binario-de-matlab-mcp-core-server --initial-working-folder=/home/username/myproject
```

Para obtener más información sobre cómo agregar servidores MCP en Claude Code, consulte [Add a local stdio server (Claude Code)](https://docs.claude.com/en/docs/claude-code/mcp#option-3%3A-add-a-local-stdio-server). Para eliminar el servidor posteriormente, ejecute:

```sh
claude mcp remove matlab
```

### Claude Desktop

Instale MATLAB MCP Core Server en Claude Desktop mediante el paquete MATLAB MCP Core Server.

1. Instale la extensión Filesystem en Claude Desktop para permitir que Claude lea y escriba archivos en su sistema. En Claude Desktop, haga clic en **Settings > Extensions > Browse extensions**. Busque la extensión Filesystem desarrollada por Anthropic y haga clic en **Install**. Especifique las carpetas a las que desea permitir el acceso del servidor MCP, luego cambie el botón **Disabled** a **Enable** para habilitar la extensión Filesystem.
   
2. Descargue el paquete MATLAB MCP Core Server `matlab-mcp-core-server.mcpb` desde la página de la [última versión](https://github.com/matlab/matlab-mcp-core-server/releases/latest). 

3. Para instalar el paquete MATLAB MCP Core Server como extensión de escritorio, haga doble clic en el archivo `matlab-mcp-core-server.mcpb` descargado y haga clic en **Install** en Claude Desktop. (Alternativamente, navegue en Claude a **File menu > Settings > Extensions > Advanced Settings > Install Extension** y seleccione el archivo `matlab-mcp-core-server.mcpb`. Haga clic en **Install**).

Para personalizar el comportamiento de MATLAB MCP Core Server, navegue a **Settings > Extensions > Configure**, donde puede modificar los [argumentos](#argumentos) del servidor.
   
### GitHub Copilot en Visual Studio Code

En su área de trabajo de VS Code, cree un archivo llamado `.vscode/mcp.json`. Inserte el siguiente código JSON, recordando especificar la ruta completa al binario del servidor que obtuvo en la configuración, así como cualquier [argumento](#argumentos). Luego guarde el archivo. (Tenga en cuenta que en Windows, las rutas requieren barras invertidas adicionales como caracteres de escape).

```json
{
    "servers": {
        "matlab": {
            "type": "stdio",
            "command": "C:\\rutacompleta\\a\\matlab-mcp-core-server-win64.exe",
            "args": []
        }
    }
}
```
Para obtener más información sobre el uso de servidores MCP en VS Code, consulte [Add and Manage MCP servers in VS Code (VS Code)](https://code.visualstudio.com/docs/copilot/customization/mcp-servers#_configure-the-mcpjson-file).

## Argumentos

Personalice el comportamiento del servidor especificando argumentos de una de estas formas:
- Inserte los argumentos en los ajustes de configuración de su aplicación de IA (generalmente un archivo `.json`).
- Introduzca los argumentos como indicadores de interfaz de línea de comandos (CLI) al iniciar el servidor. 
- Use variables de entorno, especificadas en su CLI o en los ajustes de configuración de la aplicación. Para derivar el nombre de la variable de entorno a partir de un indicador de CLI, agregue el prefijo `MW_MCP_SERVER_`, convierta a mayúsculas y reemplace los guiones (`-`) por guiones bajos (`_`). Por ejemplo, el argumento `--matlab-root` se convierte en la variable de entorno `MW_MCP_SERVER_MATLAB_ROOT`. Los indicadores de CLI tienen prioridad sobre las variables de entorno si se utilizan ambos.

| Argumento | Descripción | Ejemplo |
| ------------- | ------------- | ------------- |
| help | Muestra información de ayuda para todos los argumentos. | `--help` |
| version | Muestra la versión de MATLAB MCP Core Server. | `--version` |
| matlab-root | Ruta completa que especifica qué versión de MATLAB iniciar. No incluya `/bin` en la ruta. De forma predeterminada, el servidor intenta encontrar la primera versión de MATLAB en el PATH del sistema. | Windows: `--matlab-root=C:\\Program Files\\MATLAB\\R2026a` <br><br> Linux/macOS: `--matlab-root=/home/usr/MATLAB/R2026a`<br><br>Como variable de entorno: `MW_MCP_SERVER_MATLAB_ROOT=/home/usr/MATLAB/R2026a` |
| initialize-matlab-on-startup | Para inicializar MATLAB tan pronto como inicie el servidor, establezca este argumento en `true`. De forma predeterminada, MATLAB solo se inicia cuando se llama a la primera herramienta. | `--initialize-matlab-on-startup=true` |
| initial-working-folder | Especifique la carpeta donde MATLAB se inicia. Si no especifica un valor, MATLAB se inicia en la ruta del primer [Root (MCP)](https://modelcontextprotocol.io/specification/latest/client/roots) de su aplicación de IA. Si no ha definido un root, MATLAB se inicia en estas ubicaciones: <br> <ul><li>Linux: `/home/username` </li><li> Windows: `C:\Users\username\Documents`</li><li>Mac: `/Users/username/Documents`</li></ul> | Windows: `--initial-working-folder=C:\\Users\\username\\MyProject` <br><br> Linux/macOS: `--initial-working-folder=/Users/username/MyProject` |
| matlab-display-mode | Especifique si se muestra el escritorio de MATLAB. Use el modo `desktop` (predeterminado) para mostrar el escritorio de MATLAB. Use el modo `nodesktop` para utilizar MATLAB solo desde su aplicación de IA, sin el escritorio de MATLAB. Tenga en cuenta que en modo `nodesktop`, los comandos que requieren una interfaz gráfica (como `edit`, `open`, `open_system`, `uifigure` y `appdesigner`) seguirán abriendo ventanas de MATLAB en su escritorio. | `--matlab-display-mode=nodesktop` |
| matlab-session-mode | Especifique si el servidor MCP inicia una nueva sesión de MATLAB o se conecta a una sesión de MATLAB existente (compatible con MATLAB R2023a en adelante). El valor predeterminado es el modo **`auto`**.<br><br> **Modo `new`:** El servidor MCP inicia una nueva sesión de MATLAB. <br><br>**Modo `auto` (predeterminado):** El servidor intenta conectarse a una sesión de MATLAB existente, que debe haber configurado para el modo `existing` siguiendo las instrucciones siguientes. Si el servidor no puede encontrar una sesión de MATLAB existente, inicia una nueva. <br><br>**Modo `existing`:** El servidor intenta conectarse a una sesión de MATLAB existente. Debe haber configurado su sesión de MATLAB previamente para usar este modo, siguiendo estos pasos:<br><br><ol><li>Si está usando el modo `existing` por primera vez, ejecute `./matlab-mcp-core-server --setup-matlab`.<br><br>Este comando instala un complemento llamado MATLAB MCP Core Server Toolbox en MATLAB. Puede personalizar el comando con otros argumentos de esta tabla. Por ejemplo, para especificar qué versión de MATLAB usar para instalar la toolbox, puede usar `./matlab-mcp-core-server --setup-matlab --matlab-root=/home/usr/MATLAB/R2026a`.<br><br>Para Claude Desktop, debe descargar el binario de MATLAB MCP Core Server siguiendo las instrucciones en [Configuración](#configuración) antes de ejecutar `./matlab-mcp-core-server --setup-matlab`.<br><br></li><li>En la ventana de comandos de una sesión de MATLAB en ejecución, ejecute `shareMATLABSession()`. El servidor MCP se conectará a esta sesión de MATLAB cuando inicie el servidor con `--matlab-session-mode=existing`. Si está ejecutando múltiples sesiones de MATLAB, el servidor se conecta a la sesión de MATLAB donde ejecutó más recientemente el comando `shareMATLABSession()`.<br><br>Como alternativa a ejecutar `shareMATLABSession()` manualmente, puede agregar el comando a su [script de inicio de MATLAB (MathWorks)](https://www.mathworks.com/help/matlab/ref/startup.html).</li></ol> | `--matlab-session-mode=existing` |
| extension-file | Para usar herramientas MCP personalizadas, proporcione una ruta a un archivo JSON que defina sus herramientas. También puede usar múltiples archivos de extensión. Para obtener más información sobre el uso de herramientas personalizadas, consulte [Use Custom Tools with the MATLAB MCP Core Server](guides/custom-tools.es.md). | <br><br>Windows: `--extension-file=C:\\Users\\name\\my-tools.json` <br><br> Linux/macOS: `--extension-file=/path/to/my-tools.json` <br><br> **Uso de múltiples archivos de extensión:**<br><br>Windows:`--extension-file=C:\\path\\to\\tools-1.json --extension-file=C:\\path\\to\\tools-2.json`<br><br>Linux/macOS:`--extension-file=/path/to/tools1.json --extension-file=/path/to/tools2.json` <br><br> **Uso de variables de entorno:** <br><br> Windows: `MW_MCP_SERVER_EXTENSION_FILE=C:\Users\name\tools1.json;C:\Users\name\tools2.json` <br><br> Linux/macOS: `MW_MCP_SERVER_EXTENSION_FILE=/path/to/tools1.json:/path/to/tools2.json` |
| log-folder | Especifique la carpeta donde el servidor MCP almacena los archivos de registro. Si no se especifica, el servidor usa la carpeta temporal predeterminada de su sistema operativo. | Windows: `--log-folder=C:\\Users\\name\\AppData\\Local\\Temp` <br><br> Linux/macOS: `--log-folder=/tmp/my-logs`  |
| log-level | Los niveles de registro del servidor MCP. Los valores válidos, en orden decreciente de detalle, son `debug`, `info`, `warn` y `error`. | `--log-level=debug` |
| disable-telemetry | Para deshabilitar la recopilación de datos anónimos, establezca este argumento en `true`. Para obtener más información, consulte [Recopilación de datos](#recopilación-de-datos). | `--disable-telemetry=true` |

**Múltiples archivos de extensión**

Windows:
```
--extension-file=C:\\path\\to\\my-tools.json --extension-file=C:\\path\\to\\my-other-tools.json
```

Linux y macOS:
```
--extension-file=/path/to/my-tools.json --extension-file=/path/to/my-other-tools.json
```

**Variables de entorno**

Windows:
```
MW_MCP_SERVER_EXTENSION_FILE=C:\Users\name\my-tools.json;C:\Users\name\my-other-tools.json
```

Linux y macOS:
```
MW_MCP_SERVER_EXTENSION_FILE=/path/to/my-tools.json:/path/to/my-other-tools.json
```
## Herramientas

1. `detect_matlab_toolboxes`
    - Devuelve información sobre MATLAB y las toolboxes instaladas, incluyendo números de versión.  

1. `check_matlab_code`
    - Realiza un análisis de código estático en un script de MATLAB. Devuelve advertencias sobre estilo de programación, errores potenciales, funciones obsoletas, problemas de rendimiento y violaciones de buenas prácticas. Esta es una operación de solo lectura no destructiva que ayuda a identificar problemas de calidad del código sin ejecutar el script.
    - Entradas:
        - `script_path` (string): Ruta absoluta al archivo de script de MATLAB que se desea analizar. Debe ser un archivo `.m` válido. El archivo no se modifica durante el análisis. Ejemplo: `C:\Users\username\matlab\myFunction.m` o `/home/user/scripts/analysis.m`.

1. `evaluate_matlab_code`
    - Evalúa una cadena de código de MATLAB y devuelve la salida.
    - Entradas:
        - `code` (string): Código de MATLAB que se desea evaluar.
        - `project_path` (string): Ruta absoluta al directorio de su proyecto. MATLAB establece este directorio como la carpeta de trabajo actual. Ejemplo: `C:\Users\username\matlab-project` o `/home/user/research`.

1. `run_matlab_file`
    - Ejecuta un script de MATLAB y devuelve la salida. El script debe ser un archivo `.m` válido.
    - Entradas:
        - `script_path` (string): Ruta absoluta al archivo de script de MATLAB que se desea ejecutar. Debe ser un archivo `.m` válido. Ejemplo: `C:\Users\username\projects\analysis.m` o `/home/user/matlab/simulation.m`.

1. `run_matlab_test_file`
    - Ejecuta un script de pruebas de MATLAB y devuelve resultados completos de las pruebas. Diseñado específicamente para archivos de pruebas unitarias de MATLAB que siguen las convenciones del marco de pruebas de MATLAB.
    - Entradas:
        - `script_path` (string): Ruta absoluta al archivo de script de pruebas de MATLAB. Debe ser un archivo `.m` válido que contenga pruebas unitarias de MATLAB. Ejemplo: `C:\Users\username\tests\testMyFunction.m` o `/home/user/matlab/tests/test_analysis.m`.

## Recursos

El servidor MCP proporciona [Resources (MCP)](https://modelcontextprotocol.io/specification/latest/server/resources) para ayudar a su aplicación de IA a escribir código de MATLAB. Para ver instrucciones sobre el uso de este recurso, consulte la documentación de su aplicación de IA que explica cómo usar recursos.

1. `matlab_coding_guidelines`
    - Proporciona estándares de programación de MATLAB completos para mejorar la legibilidad, mantenibilidad y colaboración del código. Las directrices abarcan convenciones de nomenclatura, formato, comentarios, optimización del rendimiento y manejo de errores.
    - URI: `guidelines://coding`
    - MIME Type: `text/markdown`
    - Fuente: [MATLAB Coding Standards (GitHub)](https://github.com/matlab/rules/blob/main/matlab-coding-standards.md)

1. `plain_text_live_code_guidelines`
    - Proporciona reglas y directrices para generar live scripts usando el formato de archivo Live Code de texto plano `.m`, adecuado para control de versiones y desarrollo asistido por IA. Tenga en cuenta que para ejecutar live scripts de texto plano necesita MATLAB R2025a o posterior. Para obtener más información, consulte [Live Code File Format (MathWorks)](https://www.mathworks.com/help/matlab/matlab_prog/plain-text-file-format-for-live-scripts.html).
    - URI: `guidelines://plain-text-live-code`
    - MIME Type: `text/markdown`
    - Fuente: [Plain Text Live Code Generation (GitHub)](https://github.com/matlab/rules/blob/main/live-script-generation.md)

## Recopilación de datos

MATLAB MCP Core Server puede recopilar información completamente anónima sobre el uso del servidor y enviarla a MathWorks. Esta recopilación de datos ayuda a MathWorks a mejorar sus productos y está activada de forma predeterminada. Para desactivar la recopilación de datos, establezca el argumento `--disable-telemetry` en `true`.

## Consideraciones de seguridad

Al usar MATLAB MCP Core Server, debe revisar y validar exhaustivamente todas las llamadas a herramientas antes de ejecutarlas. Mantenga siempre a un humano en el proceso para acciones importantes y proceda solo cuando tenga certeza de que la llamada hará exactamente lo que espera. Para obtener más información, consulte [User Interaction Model (MCP)](https://modelcontextprotocol.io/specification/latest/server/tools#user-interaction-model) y [Security Considerations (MCP)](https://modelcontextprotocol.io/specification/latest/server/tools#security-considerations).

## Licencia y uso

La licencia está disponible en el archivo [LICENSE.md](../LICENSE.md) de este repositorio de GitHub.

Los servidores MCP solo están permitidos para su uso con MATLAB de acuerdo con el Contrato de Licencia de Software de MathWorks, y no deben ser compartidos por múltiples usuarios. Contacte a MathWorks si necesita soporte para uso compartido o centralizado del servidor.

## Contactar con soporte

MathWorks le anima a utilizar este repositorio y proporcionar comentarios. Para solicitar soporte técnico o enviar una solicitud de mejora, [cree un issue en GitHub](https://github.com/matlab/matlab-mcp-core-server/issues) o contacte con el [soporte técnico de MathWorks](https://www.mathworks.com/support/contact_us.html).

---

Copyright 2025-2026 The MathWorks, Inc.

---
