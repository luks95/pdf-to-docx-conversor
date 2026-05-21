package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// ConvertWithWord uses Microsoft Word COM automation to convert to PDF
func ConvertWithWord(inputFile string, outputDir string) (string, error) {
	// Importante: Bloquear el hilo actual para operaciones OLE en Windows
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	absInputPath, err := filepath.Abs(inputFile)
	if err != nil {
		return "", fmt.Errorf("error al obtener ruta absoluta de entrada: %v", err)
	}

	pdfPath := BuildOutputPDFPath(absInputPath, outputDir)
	absOutputPath, err := filepath.Abs(pdfPath)
	if err != nil {
		return "", fmt.Errorf("error al obtener ruta absoluta de salida: %v", err)
	}

	fmt.Printf("[DEBUG] Intentando convertir con Word: %s\n", absInputPath)

	// Initialize OLE
	err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)
	if err != nil {
		fmt.Printf("[DEBUG] CoInitializeEx nota: %v\n", err)
	}
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("Word.Application")
	if err != nil {
		return "", fmt.Errorf("Microsoft Word no parece estar instalado o accesible: %v", err)
	}
	word, _ := unknown.QueryInterface(ole.IID_IDispatch)
	defer word.Release()

	// Configurar propiedades de Word
	oleutil.PutProperty(word, "Visible", false)
	oleutil.PutProperty(word, "DisplayAlerts", 0) // wdAlertsNone = 0

	// Open Document
	documents := oleutil.MustGetProperty(word, "Documents").ToIDispatch()
	
	fmt.Printf("[DEBUG] Abriendo documento...\n")
	document, err := oleutil.CallMethod(documents, "Open", absInputPath)
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el documento (puede estar bloqueado): %v", err)
	}
	doc := document.ToIDispatch()
	defer doc.Release()

	fmt.Printf("[DEBUG] Exportando a PDF...\n")
	// wdExportFormatPDF = 17
	_, err = oleutil.CallMethod(doc, "ExportAsFixedFormat", absOutputPath, 17)
	if err != nil {
		fmt.Printf("[DEBUG] ExportAsFixedFormat falló, reintentando con SaveAs2: %v\n", err)
		// wdFormatPDF = 17
		_, err = oleutil.CallMethod(doc, "SaveAs2", absOutputPath, 17)
		if err != nil {
			return "", fmt.Errorf("error al guardar como PDF: %v", err)
		}
	}

	// Cerrar sin guardar cambios
	oleutil.CallMethod(doc, "Close", 0) // wdDoNotSaveChanges = 0
	
	fmt.Printf("[DEBUG] Conversión completada con éxito\n")
	return absOutputPath, nil
}

// FindLibreOffice tries to find the soffice.exe executable
func FindLibreOffice() (string, error) {
	if path := os.Getenv("LIBREOFFICE_PATH"); path != "" {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	commonPaths := []string{
		`C:\Program Files\LibreOffice\program\soffice.exe`,
		`C:\Program Files (x86)\LibreOffice\program\soffice.exe`,
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	path, err := exec.LookPath("soffice")
	if err == nil {
		return path, nil
	}

	return "", fmt.Errorf("LibreOffice no encontrado")
}

// IsWordFile checks if the file has a .doc or .docx extension
func IsWordFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".doc" || ext == ".docx"
}

// BuildOutputPDFPath constructs the expected PDF path
func BuildOutputPDFPath(inputFile string, outputDir string) string {
	fileName := filepath.Base(inputFile)
	ext := filepath.Ext(fileName)
	nameWithoutExt := fileName[:len(fileName)-len(ext)]
	if outputDir == "" {
		outputDir = filepath.Dir(inputFile)
	}
	return filepath.Join(outputDir, nameWithoutExt+".pdf")
}

// ConvertWordToPDF executes the conversion, trying Word first, then LibreOffice
func ConvertWordToPDF(inputFile string, outputDir string) (string, error) {
	fmt.Printf("[DEBUG] ConvertWordToPDF recibió: %s\n", inputFile)
	
	// Verificar si el archivo existe antes de hacer nada
	if _, err := os.Stat(inputFile); err != nil {
		// Si falla, intentamos obtener la ruta absoluta por si acaso
		abs, _ := filepath.Abs(inputFile)
		if _, errAbs := os.Stat(abs); errAbs == nil {
			inputFile = abs
		} else {
			return "", fmt.Errorf("el archivo no existe. Ruta intentada: %s", inputFile)
		}
	}

	if !IsWordFile(inputFile) {
		return "", fmt.Errorf("solo se permiten archivos .doc o .docx")
	}

	var absOutputDir string
	if outputDir == "" {
		absOutputDir = filepath.Dir(inputFile)
	} else {
		var err error
		absOutputDir, err = filepath.Abs(outputDir)
		if err != nil {
			absOutputDir = outputDir
		}
	}

	fmt.Printf("[DEBUG] Usando ruta entrada: %s\n", inputFile)
	fmt.Printf("[DEBUG] Usando ruta salida: %s\n", absOutputDir)

	// 1. Try with Microsoft Word (COM Automation)
	pdfPath, err := ConvertWithWord(inputFile, absOutputDir)
	if err == nil {
		return pdfPath, nil
	}

	fmt.Printf("[DEBUG] Word falló: %v. Intentando LibreOffice...\n", err)

	// 2. Try with LibreOffice (as fallback)
	libreOfficePath, errLO := FindLibreOffice()
	if errLO != nil {
		return "", fmt.Errorf("no se pudo convertir: Word falló (%v) y LibreOffice no está instalado", err)
	}

	cmd := exec.Command(
		libreOfficePath,
		"--headless",
		"--convert-to",
		"pdf",
		inputFile,
		"--outdir",
		absOutputDir,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error al convertir con LibreOffice: %v - %s", err, string(output))
	}

	finalPdfPath := BuildOutputPDFPath(inputFile, absOutputDir)
	return finalPdfPath, nil
}

// IsPDFFile checks if the file has a .pdf extension
func IsPDFFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".pdf"
}

// BuildOutputWordPath constructs the expected .docx path
func BuildOutputWordPath(inputFile string, outputDir string) string {
	fileName := filepath.Base(inputFile)
	ext := filepath.Ext(fileName)
	nameWithoutExt := fileName[:len(fileName)-len(ext)]
	if outputDir == "" {
		outputDir = filepath.Dir(inputFile)
	}
	return filepath.Join(outputDir, nameWithoutExt+".docx")
}

// ConvertPDFWithWord uses Microsoft Word COM automation to convert PDF to DOCX
func ConvertPDFWithWord(inputFile string, outputDir string) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	absInputPath, err := filepath.Abs(inputFile)
	if err != nil {
		return "", fmt.Errorf("error al obtener ruta absoluta de entrada: %v", err)
	}

	docxPath := BuildOutputWordPath(absInputPath, outputDir)
	absOutputPath, err := filepath.Abs(docxPath)
	if err != nil {
		return "", fmt.Errorf("error al obtener ruta absoluta de salida: %v", err)
	}

	fmt.Printf("[DEBUG] Intentando convertir PDF a Word con MS Word: %s\n", absInputPath)

	// Initialize OLE
	err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)
	if err != nil {
		fmt.Printf("[DEBUG] CoInitializeEx nota: %v\n", err)
	}
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("Word.Application")
	if err != nil {
		return "", fmt.Errorf("Microsoft Word no parece estar instalado o accesible: %v", err)
	}
	word, _ := unknown.QueryInterface(ole.IID_IDispatch)
	defer word.Release()

	// Configurar propiedades de Word
	oleutil.PutProperty(word, "Visible", false)
	oleutil.PutProperty(word, "DisplayAlerts", 0) // wdAlertsNone = 0

	// Open Document
	documents := oleutil.MustGetProperty(word, "Documents").ToIDispatch()
	
	fmt.Printf("[DEBUG] Abriendo PDF...\n")
	document, err := oleutil.CallMethod(documents, "Open", absInputPath)
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el PDF (puede estar bloqueado): %v", err)
	}
	doc := document.ToIDispatch()
	defer doc.Release()

	fmt.Printf("[DEBUG] Guardando como DOCX...\n")
	// wdFormatXMLDocument = 16 (formato .docx nativo)
	_, err = oleutil.CallMethod(doc, "SaveAs2", absOutputPath, 16)
	if err != nil {
		return "", fmt.Errorf("error al guardar como DOCX: %v", err)
	}

	// Cerrar sin guardar cambios
	oleutil.CallMethod(doc, "Close", 0) // wdDoNotSaveChanges = 0
	
	fmt.Printf("[DEBUG] Conversión completada con éxito\n")
	return absOutputPath, nil
}

// ConvertPDFToWord executes the conversion, trying Word first, then LibreOffice
func ConvertPDFToWord(inputFile string, outputDir string) (string, error) {
	fmt.Printf("[DEBUG] ConvertPDFToWord recibió: %s\n", inputFile)
	
	if _, err := os.Stat(inputFile); err != nil {
		abs, _ := filepath.Abs(inputFile)
		if _, errAbs := os.Stat(abs); errAbs == nil {
			inputFile = abs
		} else {
			return "", fmt.Errorf("el archivo no existe. Ruta intentada: %s", inputFile)
		}
	}

	if !IsPDFFile(inputFile) {
		return "", fmt.Errorf("solo se permiten archivos .pdf")
	}

	var absOutputDir string
	if outputDir == "" {
		absOutputDir = filepath.Dir(inputFile)
	} else {
		var err error
		absOutputDir, err = filepath.Abs(outputDir)
		if err != nil {
			absOutputDir = outputDir
		}
	}

	fmt.Printf("[DEBUG] Usando ruta entrada: %s\n", inputFile)
	fmt.Printf("[DEBUG] Usando ruta salida: %s\n", absOutputDir)

	// 1. Try with Microsoft Word (COM Automation)
	docxPath, err := ConvertPDFWithWord(inputFile, absOutputDir)
	if err == nil {
		return docxPath, nil
	}

	fmt.Printf("[DEBUG] Word falló al importar PDF: %v. Intentando LibreOffice...\n", err)

	// 2. Try with LibreOffice (as fallback)
	libreOfficePath, errLO := FindLibreOffice()
	if errLO != nil {
		return "", fmt.Errorf("no se pudo convertir: Word falló (%v) y LibreOffice no está instalado", err)
	}

	cmd := exec.Command(
		libreOfficePath,
		"--headless",
		"--infilter=writer_pdf_import",
		"--convert-to",
		"docx",
		inputFile,
		"--outdir",
		absOutputDir,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error al convertir con LibreOffice: %v - %s", err, string(output))
	}

	finalDocxPath := BuildOutputWordPath(inputFile, absOutputDir)
	return finalDocxPath, nil
}
