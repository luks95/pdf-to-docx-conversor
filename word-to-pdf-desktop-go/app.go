package main

import (
	"context"
	"os/exec"
	"path/filepath"
	"runtime"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	wruntime.OnFileDrop(ctx, func(x, y int, paths []string) {
		wruntime.EventsEmit(ctx, "custom:file-drop", paths)
	})
}

// SelectFiles opens a file dialog to select multiple Word files
func (a *App) SelectFiles() ([]string, error) {
	files, err := wruntime.OpenMultipleFilesDialog(a.ctx, wruntime.OpenDialogOptions{
		Title: "Seleccionar archivos Word",
		Filters: []wruntime.FileFilter{
			{
				DisplayName: "Documentos de Word (*.doc;*.docx)",
				Pattern:     "*.doc;*.docx",
			},
		},
	})
	return files, err
}

// SelectPDFFiles opens a file dialog to select multiple PDF files
func (a *App) SelectPDFFiles() ([]string, error) {
	files, err := wruntime.OpenMultipleFilesDialog(a.ctx, wruntime.OpenDialogOptions{
		Title: "Seleccionar archivos PDF",
		Filters: []wruntime.FileFilter{
			{
				DisplayName: "Documentos PDF (*.pdf)",
				Pattern:     "*.pdf",
			},
		},
	})
	return files, err
}

// SelectOutputFolder opens a directory dialog
func (a *App) SelectOutputFolder() (string, error) {
	folder, err := wruntime.OpenDirectoryDialog(a.ctx, wruntime.OpenDialogOptions{
		Title: "Seleccionar carpeta de salida",
	})
	if err != nil {
		return "", err
	}
	return folder, nil
}

// ConversionResult represents the status of a single file conversion
type ConversionResult struct {
	FileName string `json:"fileName"`
	Success  bool   `json:"success"`
	Path     string `json:"path"`
	Error    string `json:"error"`
}

// ConvertToPDF converts multiple Word files to PDF
func (a *App) ConvertToPDF(inputFiles []string, outputDir string) []ConversionResult {
	results := make([]ConversionResult, 0, len(inputFiles))
	for _, file := range inputFiles {
		pdfPath, err := ConvertWordToPDF(file, outputDir)
		res := ConversionResult{
			FileName: filepath.Base(file),
			Success:  err == nil,
			Path:     pdfPath,
		}
		if err != nil {
			res.Error = err.Error()
		}
		results = append(results, res)
	}
	return results
}

// ConvertToWord converts multiple PDF files to Word (.docx)
func (a *App) ConvertToWord(inputFiles []string, outputDir string) []ConversionResult {
	results := make([]ConversionResult, 0, len(inputFiles))
	for _, file := range inputFiles {
		docxPath, err := ConvertPDFToWord(file, outputDir)
		res := ConversionResult{
			FileName: filepath.Base(file),
			Success:  err == nil,
			Path:     docxPath,
		}
		if err != nil {
			res.Error = err.Error()
		}
		results = append(results, res)
	}
	return results
}

// OpenFolder opens the specified path in the default file explorer
func (a *App) OpenFolder(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", filepath.Clean(path))
	case "darwin":
		cmd = exec.Command("open", path)
	default: // Linux and others
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}
