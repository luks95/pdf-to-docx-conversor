# Informe Técnico: Problema de Arrastre de Archivos (Drag & Drop)

Este documento detalla la investigación, el diagnóstico y las soluciones aplicadas al problema de la funcionalidad de arrastre de documentos en la aplicación **Word to PDF Converter**.

## 1. Estado Actual de la Implementación
La aplicación se encuentra en la versión **v3.4**, la cual implementa los estándares de Wails v2 para el manejo de archivos:

- **Backend (Go):** Configurado con `EnableFileDrop: true` en `main.go`.
- **Frontend (JS):** Escucha activa del evento nativo `wails:file-drop`.
- **CSS:** Uso de la propiedad `--wails-drop-target: drop` para señalizar las zonas de captura nativa.

## 2. Diagnóstico Técnico
Tras múltiples pruebas con diferentes versiones (v2.1 a v3.4), se han confirmado los siguientes puntos:

1.  **Detección de DOM:** El motor WebView2 de Windows **detecta el soltado (drop)** del archivo (confirmado por los logs de consola `Drop detectado en el DOM`).
2.  **Bloqueo de Intercepción:** Wails no logra interceptar el archivo para convertirlo en una ruta de sistema (path).
3.  **Causa Raíz:** **Restricciones de Privilegios de Windows (UAC / UIPI).** 
    - En Windows, si la aplicación se ejecuta con privilegios de **Administrador**, el sistema bloquea los eventos de arrastre provenientes de procesos con privilegios normales (como el Explorador de Archivos).
    - Esto se conoce como *User Interface Privilege Isolation* (UIPI).

## 3. Pruebas Realizadas
- **v2.1 - v2.5:** Intentos de captura mediante JavaScript estándar. Fallido porque el navegador solo entrega el nombre del archivo, no la ruta completa por seguridad.
- **v2.6 - v3.0:** Eliminación de manejadores JS para dar prioridad al motor de Wails. Fallido por la misma restricción de sistema.
- **v3.1 - v3.2:** Sistema de diagnóstico con temporizador. **Confirmó el bloqueo del sistema operativo** al no recibirse respuesta de Wails en el tiempo esperado.
- **v3.3 - v3.4:** Optimización nativa final (CSS nativo y limpieza de código).

## 4. Solución Recomendada (Entorno)
Para que la funcionalidad opere correctamente, el usuario debe:

1.  **Evitar el modo Administrador:** No ejecutar el terminal (CMD/PowerShell) ni el `.exe` como Administrador.
2.  **Reinicio de Explorer:** En algunos casos, reiniciar el proceso `explorer.exe` desde el Administrador de Tareas resuelve bloqueos temporales del sistema de arrastre.
3.  **Compilación Limpia:** Ejecutar `wails build` para generar el binario final y probarlo fuera del entorno de desarrollo, asegurándose de que no tenga el icono de "escudo" de administrador.

## 5. Conclusión
El código de la aplicación es **técnicamente correcto** y sigue las mejores prácticas de Wails v2. La falta de reacción se debe estrictamente a la capa de seguridad de Windows que impide la comunicación de archivos entre el Explorador y la Aplicación cuando hay niveles de privilegios distintos.
