# Word to PDF Desktop Converter (Go + Wails)

Esta es una aplicación de escritorio para Windows que permite convertir archivos Word (`.doc`, `.docx`) a PDF de forma local utilizando LibreOffice en modo headless.

## Requisitos

- [Go](https://go.dev/dl/) (v1.21+)
- [Wails](https://wails.io/docs/gettingstarted/installation) (v2)
- [Node.js y npm](https://nodejs.org/) (para el frontend)
- [LibreOffice](https://es.libreoffice.org/descarga/libreoffice/) instalado en la ruta por defecto o configurado en `LIBREOFFICE_PATH`.

## Cómo ejecutar en desarrollo

1. Abre una terminal en la carpeta raíz del proyecto.
2. Ejecuta:
   ```bash
   wails dev
   ```

## Cómo compilar para producción

1. Ejecuta:
   ```bash
   wails build
   ```
2. El ejecutable `.exe` se generará en la carpeta `build/bin/`.

## Características

- Interfaz moderna y limpia.
- Soporte para Drag & Drop (arrastrar y soltar archivos).
- Selección de carpeta de salida personalizada.
- Conversión rápida usando LibreOffice Headless.
- Validación de tipos de archivos.
- Opción para abrir la carpeta contenedora tras la conversión.
