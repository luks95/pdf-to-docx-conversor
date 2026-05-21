# Documentación Técnica: Word to PDF Desktop Converter

Esta aplicación permite convertir documentos Word (`.doc`, `.docx`) a PDF de manera local y eficiente. A continuación se detallan los aspectos clave de su arquitectura y funcionamiento.

## 1. Arquitectura General
El proyecto utiliza **Wails v2**, que combina:
- **Backend:** Go (Lógica de negocio, acceso al sistema de archivos y automatización de Word/LibreOffice).
- **Frontend:** HTML, CSS y JavaScript (Vanilla) ejecutándose en un componente WebView2 nativo de Windows.

## 2. Motor de Conversión (Lógica "Word-First")
La aplicación emplea una estrategia de conversión en dos niveles (ubicada en `converter.go`):

### Nivel 1: Microsoft Word (Vía OLE/COM)
Si el usuario tiene Microsoft Word instalado, la aplicación lo utiliza directamente.
- **Tecnología:** Usa la librería `go-ole` para comunicarse con la API de automatización de Word.
- **Ventaja:** Máxima fidelidad en la conversión, respetando formatos complejos, tablas e imágenes.
- **Detalle Técnico:** Se utiliza `runtime.LockOSThread()` y `ole.COINIT_APARTMENTTHREADED` para evitar errores de hilos (threading) comunes en la automatización de Windows.

### Nivel 2: LibreOffice (Fallback)
Si Word no está disponible, la aplicación busca automáticamente `soffice.exe` (LibreOffice).
- **Rutas de búsqueda:** Carpeta de programa por defecto o mediante la variable de entorno `LIBREOFFICE_PATH`.
- **Modo:** Ejecuta el comando en modo `--headless` para que el usuario no vea la interfaz de LibreOffice.

## 3. Sistema de Arrastre de Archivos (Drag & Drop)
Uno de los mayores retos técnicos fue obtener la **ruta completa** de los archivos arrastrados.
- **Problema:** Por seguridad, los navegadores web (y el WebView2) a menudo ocultan la ruta real del archivo.
- **Solución:** Se activó `EnableFileDrop: true` en la configuración de Wails (`main.go`). Esto permite que el backend capture la ruta real desde el sistema operativo y la envíe al frontend mediante el evento nativo `wails:file-drop`.

## 4. Mejoras de Interfaz (UX)
- **Multi-archivo:** Permite procesar una lista de archivos simultáneamente.
- **Loader Visual:** Un spinner indica cuando hay una conversión en curso.
- **Manejo de Errores:** Mensajes detallados si un archivo específico falla (por ejemplo, por estar bloqueado o dañado).
- **Botón de Respaldo:** Se incluyó un botón de selección manual por si el arrastre falla en configuraciones de Windows muy restrictivas.

## 5. Estructura de Archivos Clave
- `main.go`: Punto de entrada, configuración de la ventana y eventos nativos.
- `app.go`: Bindings entre JS y Go (selección de carpetas, llamada a conversión).
- `converter.go`: El "motor". Contiene toda la lógica de automatización de Word y LibreOffice.
- `frontend/src/main.js`: Gestiona la lista de archivos, el estado de la UI y las llamadas al backend.
- `frontend/src/style.css`: Estilos modernos con soporte para modo oscuro.

## 6. Solución de Problemas Comunes (Troubleshooting)
- **Error "Archivo no existe":** Suele ocurrir si se arrastra un archivo y el programa solo recibe el nombre. Asegúrate de que `wails dev` o el `.exe` final estén funcionando correctamente para capturar el evento nativo.
- **Error de "Subproceso":** Corregido mediante el bloqueo de hilos en Go para OLE. Si reaparece, cerrar cualquier instancia colgada de Word en el Administrador de Tareas.
- **LibreOffice:** Si no tienes Word, LibreOffice debe estar instalado en `C:\Program Files\LibreOffice`.

## 7. Compilación para Producción
Para generar el ejecutable final sin consola:
```bash
wails build
```
El archivo resultante se encontrará en `build/bin/`.
