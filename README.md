# 📄 Conversor de Documentos de Escritorio (Word ⇄ PDF) 

Este repositorio contiene la solución completa para la aplicación **Conversor Inteligente (Word to PDF Desktop Converter)** construida en **Go** y **Wails v2** para Windows.

El proyecto permite convertir documentos de manera bidireccional entre Word y PDF localmente y sin depender de servicios en la nube, optimizando la precisión de formato mediante una arquitectura de conversión de doble nivel (MS Word OLE/COM & LibreOffice Headless Fallback).

## 📂 Estructura del Repositorio

*   **[`word-to-pdf-desktop-go/`](file:///C:/Dev/ConversorWordAPDF/word-to-pdf-desktop-go)**: Carpeta principal del proyecto de software con el código de backend en Go, el frontend en JavaScript vanilla y la configuración de Wails.
*   **[`word-to-pdf-desktop-go.md`](file:///C:/Dev/ConversorWordAPDF/word-to-pdf-desktop-go.md)**: Documento de especificación de requisitos originales y reglas de negocio del proyecto.

## 📖 Documentación Principal

Para ver la guía de instalación detallada, configuración del entorno, arquitectura del motor de conversión, comandos de desarrollo y solución de problemas comunes, por favor lee el README oficial del proyecto:

👉 **[Ver README detallado del proyecto](file:///C:/Dev/ConversorWordAPDF/word-to-pdf-desktop-go/README.md)**

---

## 🛠️ Tecnologías Clave
*   **Backend:** Go 1.21+ con Wails v2
*   **Frontend:** HTML5, CSS3 y JavaScript Vanilla (sin frameworks pesados)
*   **Motores de Conversión Local:**
    *   Microsoft Word OLE/COM Automatización nativa (`go-ole`)
    *   LibreOffice Headless (`soffice.exe`) como fallback local de alta fidelidad
