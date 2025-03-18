package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTemplateFilesBasic(t *testing.T) {
	// Create temporary directories for templates and output.
	templateDir, err := os.MkdirTemp("", "templateDir")
	if err != nil {
		t.Fatalf("failed to create temporary template dir: %v", err)
	}
	defer os.RemoveAll(templateDir)

	outputDir, err := os.MkdirTemp("", "outputDir")
	if err != nil {
		t.Fatalf("failed to create temporary output dir: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// Create a sample template file in the template directory.
	fileName := "greeting.tmpl"
	templateFilePath := filepath.Join(templateDir, fileName)
	templateContent := "Hello {{.Name}}"
	if err := os.WriteFile(templateFilePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("failed to write template file: %v", err)
	}

	data := Data{"Name": "Go Developer"}

	// Invoke the function under test.
	err = templateFiles(templateDir, outputDir, data)
	if err != nil {
		t.Fatalf("templateFiles returned error: %v", err)
	}

	// Check that the output file is created and contains the expected content.
	outputFilePath := filepath.Join(outputDir, fileName)
	outputContentBytes, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	outputContent := string(outputContentBytes)
	expectedContent := "Hello Go Developer"
	if outputContent != expectedContent {
		t.Errorf("expected file content %q, got %q", expectedContent, outputContent)
	}
}
