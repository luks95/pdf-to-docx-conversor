# Proyecto: Word to PDF Desktop Converter en Go

## Objetivo

Crear una aplicación de escritorio para Windows que permita al usuario arrastrar un archivo Word (`.docx` o `.doc`) y convertirlo a PDF de forma local.

La aplicación debe ser simple, rápida y funcional. El usuario debe poder:

1. Abrir la aplicación.
2. Arrastrar un archivo Word a una zona visual.
3. Validar que el archivo sea `.doc` o `.docx`.
4. Elegir una carpeta de salida o usar la misma carpeta del archivo original.
5. Presionar un botón para convertir.
6. Generar el PDF usando LibreOffice en modo headless.
7. Mostrar mensajes claros de éxito o error.
8. Opcionalmente abrir la carpeta donde se generó el PDF.

## Stack solicitado

Usar:

- Go
- Wails
- HTML / CSS / JavaScript para la interfaz
- LibreOffice Headless como motor de conversión
- Windows como sistema principal

No usar:

- Java
- JavaFX
- Spring Boot
- Electron
- Base de datos
- Librerías pagas

## Nombre del proyecto

```text
word-to-pdf-desktop-go
```

## Motor de conversión

La conversión debe ejecutarse usando `soffice.exe` de LibreOffice.

Ruta principal esperada en Windows:

```text
C:\Program Files\LibreOffice\program\soffice.exe
```

Ruta alternativa:

```text
C:\Program Files (x86)\LibreOffice\program\soffice.exe
```

Comando base:

```bash
"C:\Program Files\LibreOffice\program\soffice.exe" --headless --convert-to pdf "archivo.docx" --outdir "carpeta_salida"
```

## Arquitectura esperada

```text
Frontend Wails
    ↓
Usuario arrastra archivo Word
    ↓
Frontend envía ruta del archivo al backend Go
    ↓
Backend Go valida archivo y carpeta
    ↓
Backend Go ejecuta LibreOffice Headless con os/exec
    ↓
Se genera el PDF
    ↓
Frontend muestra resultado
```

## Funcionalidades requeridas

### 1. Interfaz gráfica

Crear una ventana con:

- Título: `Word to PDF Converter`
- Zona grande para arrastrar archivos.
- Texto visible: `Arrastra aquí tu archivo Word`
- Mostrar el nombre del archivo seleccionado.
- Botón: `Elegir carpeta de salida`
- Botón: `Convertir a PDF`
- Botón opcional: `Abrir carpeta`
- Label o bloque de estado para mostrar:
  - Archivo seleccionado.
  - Carpeta de salida seleccionada.
  - Conversión en proceso.
  - Conversión exitosa.
  - Error de conversión.
  - LibreOffice no encontrado.
  - Archivo inválido.

### 2. Drag and drop

Permitir arrastrar archivos sobre la ventana.

Validar que el archivo tenga una de estas extensiones:

```text
.doc
.docx
```

Si el archivo no es válido, mostrar:

```text
Solo se permiten archivos .doc o .docx
```

### 3. Backend en Go

Crear una estructura principal:

```go
type App struct {
    ctx context.Context
}
```

Crear métodos exportados para Wails:

```go
func (a *App) SelectOutputFolder() (string, error)
func (a *App) ConvertToPDF(inputFile string, outputDir string) (string, error)
func (a *App) OpenFolder(path string) error
```

### 4. Conversión a PDF

Crear la lógica de conversión usando `os/exec`.

Ejemplo conceptual:

```go
cmd := exec.Command(
    libreOfficePath,
    "--headless",
    "--convert-to",
    "pdf",
    inputFile,
    "--outdir",
    outputDir,
)
```

La función debe:

1. Verificar que el archivo exista.
2. Verificar que la extensión sea `.doc` o `.docx`.
3. Buscar LibreOffice en las rutas comunes.
4. Si no encuentra LibreOffice, devolver un error claro.
5. Si `outputDir` está vacío, usar la carpeta del archivo original.
6. Ejecutar `soffice.exe`.
7. Capturar `stdout` y `stderr`.
8. Esperar a que termine el proceso.
9. Validar que el PDF fue creado.
10. Devolver la ruta del PDF generado.

## Nombre esperado del PDF

Si el archivo es:

```text
contrato.docx
```

El PDF generado debe ser:

```text
contrato.pdf
```

En la carpeta de salida seleccionada.

## Manejo de errores

Controlar estos casos:

- LibreOffice no está instalado.
- Archivo inválido.
- Archivo no existe.
- La extensión no es `.doc` ni `.docx`.
- Error durante la conversión.
- PDF no generado.
- Carpeta de salida inexistente.
- Permiso denegado en carpeta de salida.
- Usuario cancela selección de carpeta.

Los errores deben mostrarse al usuario de manera clara.

Ejemplos:

```text
No se encontró LibreOffice. Instale LibreOffice o verifique la ruta.
```

```text
No se pudo convertir el archivo. Revise que el documento no esté dañado.
```

```text
El PDF no fue generado correctamente.
```

## Estructura del proyecto esperada

Crear una estructura similar a esta:

```text
word-to-pdf-desktop-go/
 ├─ go.mod
 ├─ go.sum
 ├─ main.go
 ├─ app.go
 ├─ converter.go
 ├─ README.md
 ├─ wails.json
 └─ frontend/
    ├─ index.html
    ├─ package.json
    ├─ src/
    │  ├─ main.js
    │  ├─ App.js
    │  └─ style.css
    └─ dist/
```

La estructura puede ajustarse al template actual de Wails, pero debe mantener separación clara entre:

- lógica principal;
- conversión;
- interfaz;
- estilos.

## Archivos esperados

### main.go

Debe inicializar la aplicación Wails.

Debe configurar:

- título de ventana;
- ancho y alto inicial;
- binding del struct `App`.

### app.go

Debe contener:

- struct `App`;
- método `startup`;
- método `SelectOutputFolder`;
- método `ConvertToPDF`;
- método `OpenFolder`.

### converter.go

Debe contener:

- función para buscar LibreOffice;
- función para validar extensión;
- función para obtener nombre de PDF;
- función para ejecutar la conversión.

Funciones sugeridas:

```go
func FindLibreOffice() (string, error)
func IsWordFile(path string) bool
func BuildOutputPDFPath(inputFile string, outputDir string) string
func ConvertWordToPDF(inputFile string, outputDir string) (string, error)
```

### Frontend

Debe tener:

- zona de drag and drop;
- selección de carpeta;
- botón convertir;
- estado visible;
- diseño simple y moderno.

## UX deseada

La ventana debe ser sencilla y moderna.

Sugerencia visual:

```text
+------------------------------------------------+
|              Word to PDF Converter             |
|                                                |
|   +----------------------------------------+   |
|   |                                        |   |
|   |     Arrastra aquí tu archivo Word      |   |
|   |                                        |   |
|   +----------------------------------------+   |
|                                                |
|   Archivo: contrato.docx                       |
|   Salida: C:\Users\Lucas\Documents           |
|                                                |
|   [Elegir carpeta] [Convertir a PDF]           |
|                                                |
|   Estado: PDF generado correctamente           |
|                                                |
|   [Abrir carpeta]                              |
+------------------------------------------------+
```

## Estilos CSS

Crear un diseño limpio:

- fondo claro;
- tarjeta central;
- zona de arrastre con borde punteado;
- efecto visual cuando el usuario arrastra un archivo encima;
- botones modernos;
- mensaje de estado visible;
- mensaje de error diferenciado;
- diseño responsive dentro de la ventana.

## Validaciones importantes

Antes de convertir:

1. Verificar que haya archivo seleccionado.
2. Verificar que sea `.doc` o `.docx`.
3. Verificar que exista el archivo.
4. Verificar que exista LibreOffice.
5. Verificar la carpeta de salida.
6. Si no se eligió carpeta de salida, usar la carpeta original del archivo.

## Comandos esperados

Para crear el proyecto con Wails:

```bash
wails init -n word-to-pdf-desktop-go -t vanilla
```

Para ejecutar en desarrollo:

```bash
wails dev
```

Para compilar:

```bash
wails build
```

El resultado debe ser un `.exe` para Windows.

## Requisitos técnicos

- Usar Go estándar para la conversión.
- Usar `os/exec`.
- Usar `filepath` para rutas.
- Usar `os.Stat` para validar archivos y carpetas.
- No bloquear la interfaz mientras convierte.
- Mostrar estado de proceso.
- Código claro y mantenible.
- Priorizar que funcione rápido.

## Ejemplo de lógica de conversión en Go

```go
func ConvertWordToPDF(inputFile string, outputDir string) (string, error) {
    if inputFile == "" {
        return "", fmt.Errorf("no se seleccionó ningún archivo")
    }

    if !IsWordFile(inputFile) {
        return "", fmt.Errorf("solo se permiten archivos .doc o .docx")
    }

    if _, err := os.Stat(inputFile); err != nil {
        return "", fmt.Errorf("el archivo no existe: %w", err)
    }

    if outputDir == "" {
        outputDir = filepath.Dir(inputFile)
    }

    if info, err := os.Stat(outputDir); err != nil || !info.IsDir() {
        return "", fmt.Errorf("la carpeta de salida no existe")
    }

    libreOfficePath, err := FindLibreOffice()
    if err != nil {
        return "", err
    }

    cmd := exec.Command(
        libreOfficePath,
        "--headless",
        "--convert-to",
        "pdf",
        inputFile,
        "--outdir",
        outputDir,
    )

    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("error al convertir: %v - %s", err, string(output))
    }

    pdfPath := BuildOutputPDFPath(inputFile, outputDir)

    if _, err := os.Stat(pdfPath); err != nil {
        return "", fmt.Errorf("el PDF no fue generado correctamente")
    }

    return pdfPath, nil
}
```

## Extra opcional

Agregar detección de LibreOffice usando variable de entorno:

```text
LIBREOFFICE_PATH
```

Orden de búsqueda:

1. Variable de entorno `LIBREOFFICE_PATH`
2. `C:\Program Files\LibreOffice\program\soffice.exe`
3. `C:\Program Files (x86)\LibreOffice\program\soffice.exe`
4. Buscar `soffice` en el `PATH`

## Instrucción para Codex CLI

Genera el proyecto completo en Go usando Wails y LibreOffice Headless siguiendo esta especificación.

Crea todos los archivos necesarios para que el proyecto pueda ejecutarse con:

```bash
wails dev
```

Y pueda compilarse con:

```bash
wails build
```

Priorizar que funcione rápido antes que agregar demasiadas funciones.

El objetivo principal es tener una aplicación desktop donde el usuario arrastre un archivo Word y lo convierta a PDF localmente usando LibreOffice.
