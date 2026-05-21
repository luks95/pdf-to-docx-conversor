// Elementos del DOM
const dropZone = document.getElementById('drop-zone');
const fileListContainer = document.getElementById('file-list');
const status = document.getElementById('status');
const selectDirBtn = document.getElementById('select-dir-btn');
const convertBtn = document.getElementById('convert-btn');
const openFolderBtn = document.getElementById('open-folder-btn');
const outputDirLabel = document.getElementById('output-dir-label');
const loader = document.getElementById('loader');
const browseFilesBtn = document.getElementById('browse-files-btn');

// Elementos del Selector de Modo
const btnWordToPdf = document.getElementById('btn-word-to-pdf');
const btnPdfToWord = document.getElementById('btn-pdf-to-word');
const dropIcon = document.getElementById('drop-icon');
const dropText = document.getElementById('drop-text');

let selectedFiles = [];
let outputDir = "";
let currentMode = "word-to-pdf"; // "word-to-pdf" o "pdf-to-word"

// VERSIÓN 5.0 - SIMPLICIDAD ELEGANTE
console.log("Frontend cargado - v5.0");
showStatus("Listo. Arrastra archivos aquí para comenzar.", "info");

// 1. REGISTRO DE EVENTO WAILS (NATIVO Y PERSONALIZADO)
function initWailsEvents() {
    if (window.runtime) {
        console.log("Wails runtime detectado. Registrando eventos...");
        
        // Evento nativo de Wails
        window.runtime.EventsOn("wails:file-drop", (x, y, paths) => {
            console.log("¡WAILS NATIVO CAPTURÓ ARCHIVOS!", paths);
            if (paths && paths.length > 0) {
                handleFilesSelection(paths);
            }
        });

        // Receptor OnFileDrop explícito sin requerir CSS drop target (¡infalible!)
        window.runtime.OnFileDrop((x, y, paths) => {
            console.log("¡WAILS ONFILEDROP CAPTURÓ ARCHIVOS!", paths);
            if (paths && paths.length > 0) {
                handleFilesSelection(paths);
            }
        }, false);

        // Evento personalizado emitido desde Go (doble seguro)
        window.runtime.EventsOn("custom:file-drop", (paths) => {
            console.log("¡WAILS PERSONALIZADO (GO) CAPTURÓ ARCHIVOS!", paths);
            if (paths && paths.length > 0) {
                handleFilesSelection(paths);
            }
        });
    } else {
        console.warn("Wails runtime no detectado inmediatamente, reintentando en 100ms...");
        setTimeout(initWailsEvents, 100);
    }
}
initWailsEvents();

// 2. MANEJADORES DE EVENTOS
// Prevenimos dragover y drop a nivel de ventana para evitar el comportamiento predeterminado del navegador (abrir/descargar el archivo)
window.addEventListener('dragover', (e) => e.preventDefault());
window.addEventListener('drop', (e) => e.preventDefault());

if (dropZone) {
    dropZone.addEventListener('dragover', (e) => {
        e.preventDefault();
        dropZone.classList.add('hover');
    });

    dropZone.addEventListener('dragleave', () => {
        dropZone.classList.remove('hover');
    });

    // Limpiamos el estado hover al soltar y prevenimos comportamiento por defecto en la zona de drop
    dropZone.addEventListener('drop', (e) => {
        e.preventDefault();
        dropZone.classList.remove('hover');
    });
}

// 3. GESTIÓN DEL SELECTOR DE MODO
function setMode(mode) {
    if (currentMode === mode) return;
    
    currentMode = mode;
    selectedFiles = [];
    updateFileListUI();
    openFolderBtn.classList.add('hidden');
    
    if (mode === "word-to-pdf") {
        btnWordToPdf.classList.add('active');
        btnPdfToWord.classList.remove('active');
        
        dropIcon.textContent = "📄";
        dropText.textContent = "Arrastra aquí tus archivos Word o";
        convertBtn.textContent = "Convertir a PDF";
        
        showStatus("Modo Word a PDF activado. Selecciona archivos .doc o .docx", "info");
    } else {
        btnWordToPdf.classList.remove('active');
        btnPdfToWord.classList.add('active');
        
        dropIcon.textContent = "📕";
        dropText.textContent = "Arrastra aquí tus archivos PDF o";
        convertBtn.textContent = "Convertir a Word";
        
        showStatus("Modo PDF a Word activado. Selecciona archivos .pdf", "info");
    }
}

btnWordToPdf.addEventListener('click', () => setMode('word-to-pdf'));
btnPdfToWord.addEventListener('click', () => setMode('pdf-to-word'));

// 4. PROCESAMIENTO CENTRAL
function handleFilesSelection(filePaths) {
    console.log("Procesando:", filePaths);
    if (!filePaths || !Array.isArray(filePaths)) return;
    
    const validFiles = filePaths.filter(path => {
        if (typeof path !== 'string') return false;
        const ext = path.split('.').pop().toLowerCase();
        
        if (currentMode === "word-to-pdf") {
            return ext === 'doc' || ext === 'docx';
        } else {
            return ext === 'pdf';
        }
    });

    if (validFiles.length > 0) {
        validFiles.forEach(file => {
            if (!selectedFiles.includes(file)) selectedFiles.push(file);
        });
        updateFileListUI();
        convertBtn.disabled = false;
        showStatus(`${selectedFiles.length} archivo(s) listo(s) para convertir.`, 'success');
    } else {
        if (currentMode === "word-to-pdf") {
            showStatus("Error: Solo se admiten archivos Word (.doc, .docx)", "error");
        } else {
            showStatus("Error: Solo se admiten archivos PDF (.pdf)", "error");
        }
    }
}

function updateFileListUI() {
    fileListContainer.innerHTML = '';
    if (selectedFiles.length === 0) {
        dropText.classList.remove('hidden');
        browseFilesBtn.classList.remove('hidden');
        convertBtn.disabled = true;
        return;
    }
    dropText.classList.add('hidden');
    browseFilesBtn.classList.add('hidden');
    selectedFiles.forEach(path => {
        const div = document.createElement('div');
        div.className = 'file-item';
        const fileName = path.split(/[\\/]/).pop();
        div.innerHTML = `<span>${fileName}</span><span class="remove-btn" style="cursor:pointer;margin-left:10px">✕</span>`;
        div.querySelector('.remove-btn').onclick = (e) => {
            e.stopPropagation(); // Evitar abrir selector de archivos al hacer click en la cruz
            selectedFiles = selectedFiles.filter(f => f !== path);
            updateFileListUI();
            if (selectedFiles.length === 0) {
                showStatus("Listo. Arrastra archivos aquí para comenzar.", "info");
            } else {
                showStatus(`${selectedFiles.length} archivo(s) listo(s) para convertir.`, 'success');
            }
        };
        fileListContainer.appendChild(div);
    });
}

browseFilesBtn.addEventListener('click', async () => {
    try {
        if (window.go && window.go.main && window.go.main.App) {
            let result;
            if (currentMode === "word-to-pdf") {
                result = await window.go.main.App.SelectFiles();
            } else {
                result = await window.go.main.App.SelectPDFFiles();
            }
            if (result) handleFilesSelection(result);
        }
    } catch (e) { console.error(e); }
});

selectDirBtn.addEventListener('click', async () => {
    try {
        const result = await window.go.main.App.SelectOutputFolder();
        if (result) {
            outputDir = result;
            outputDirLabel.textContent = outputDir;
        }
    } catch (e) { showStatus("Error al seleccionar carpeta", "error"); }
});

convertBtn.addEventListener('click', async () => {
    if (selectedFiles.length === 0) return;
    loader.classList.remove('hidden');
    convertBtn.disabled = true;
    showStatus("Convirtiendo archivos...", "info");
    try {
        let results;
        if (currentMode === "word-to-pdf") {
            results = await window.go.main.App.ConvertToPDF(selectedFiles, outputDir);
        } else {
            results = await window.go.main.App.ConvertToWord(selectedFiles, outputDir);
        }
        let count = results.filter(r => r.success).length;
        if (count > 0) {
            showStatus(`¡Conversión exitosa! ${count} de ${results.length} archivos convertidos.`, "success");
            openFolderBtn.classList.remove('hidden');
        } else {
            const firstErr = results[0]?.error || "Error desconocido";
            showStatus(`Error en la conversión: ${firstErr}`, "error");
        }
    } catch (e) { 
        showStatus(`Error fatal en conversión: ${e.message || e}`, "error"); 
    } finally { 
        loader.classList.add('hidden'); 
        convertBtn.disabled = false; 
    }
});

openFolderBtn.addEventListener('click', async () => {
    try {
        // Usar la carpeta de salida seleccionada, o el directorio del primer archivo si no hay una carpeta elegida
        let folderToOpen = outputDir;
        if (!folderToOpen && selectedFiles.length > 0) {
            // Obtener el directorio del primer archivo.
            const file = selectedFiles[0];
            const lastSlash = Math.max(file.lastIndexOf('/'), file.lastIndexOf('\\'));
            if (lastSlash !== -1) {
                folderToOpen = file.substring(0, lastSlash);
            }
        }
        
        if (folderToOpen) {
            await window.go.main.App.OpenFolder(folderToOpen);
        } else {
            showStatus("No se pudo determinar la carpeta de destino.", "error");
        }
    } catch (e) { 
        showStatus("Error al abrir la carpeta de salida", "error"); 
    }
});

function showStatus(message, type) {
    status.textContent = message;
    status.className = `status ${type}`;
}
