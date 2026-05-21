package main

import (
	"testing"
)

func TestIsWordFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"document.docx", true},
		{"document.doc", true},
		{"document.DOCX", true},
		{"document.DOC", true},
		{"document.pdf", false},
		{"document.txt", false},
		{"no_extension", false},
		{"path/to/doc.docx", true},
		{"path/to/image.png", false},
	}

	for _, test := range tests {
		result := IsWordFile(test.path)
		if result != test.expected {
			t.Errorf("IsWordFile(%q) = %v; want %v", test.path, result, test.expected)
		}
	}
}

func TestBuildOutputPDFPath(t *testing.T) {
	tests := []struct {
		inputFile string
		outputDir string
		expected  string
	}{
		{"contrato.docx", "", "contrato.pdf"},
		{"contrato.docx", "C:\\Salida", "C:\\Salida\\contrato.pdf"},
		{"C:\\Documentos\\contrato.docx", "", "C:\\Documentos\\contrato.pdf"},
		{"C:\\Documentos\\contrato.docx", "D:\\PDFs", "D:\\PDFs\\contrato.pdf"},
	}

	for _, test := range tests {
		result := BuildOutputPDFPath(test.inputFile, test.outputDir)
		if result != test.expected {
			t.Errorf("BuildOutputPDFPath(%q, %q) = %q; want %q", test.inputFile, test.outputDir, result, test.expected)
		}
	}
}

func TestIsPDFFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"document.pdf", true},
		{"document.PDF", true},
		{"document.docx", false},
		{"document.txt", false},
		{"no_extension", false},
		{"path/to/doc.pdf", true},
		{"path/to/image.png", false},
	}

	for _, test := range tests {
		result := IsPDFFile(test.path)
		if result != test.expected {
			t.Errorf("IsPDFFile(%q) = %v; want %v", test.path, result, test.expected)
		}
	}
}

func TestBuildOutputWordPath(t *testing.T) {
	tests := []struct {
		inputFile string
		outputDir string
		expected  string
	}{
		{"contrato.pdf", "", "contrato.docx"},
		{"contrato.pdf", "C:\\Salida", "C:\\Salida\\contrato.docx"},
		{"C:\\Documentos\\contrato.pdf", "", "C:\\Documentos\\contrato.docx"},
		{"C:\\Documentos\\contrato.pdf", "D:\\PDFs", "D:\\PDFs\\contrato.docx"},
	}

	for _, test := range tests {
		result := BuildOutputWordPath(test.inputFile, test.outputDir)
		if result != test.expected {
			t.Errorf("BuildOutputWordPath(%q, %q) = %q; want %q", test.inputFile, test.outputDir, result, test.expected)
		}
	}
}

