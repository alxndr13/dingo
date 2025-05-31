package main

import (
	"os"
	"path/filepath"
	"strings"
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

func TestTemplateFilesSprigIntegration(t *testing.T) {
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

	// Test various sprig functions
	testCases := []struct {
		name           string
		templateContent string
		data           Data
		expectedContent string
	}{
		{
			name:           "string_functions",
			templateContent: "{{.name | upper}} - {{.name | title}}",
			data:           Data{"name": "hello world"},
			expectedContent: "HELLO WORLD - Hello World",
		},
		{
			name:           "math_functions",
			templateContent: "Sum: {{add .a .b}}, Max: {{max .a .b .c}}",
			data:           Data{"a": 5, "b": 3, "c": 8},
			expectedContent: "Sum: 8, Max: 8",
		},
		{
			name:           "list_functions",
			templateContent: "First: {{first .items}}, Last: {{last .items}}, Join: {{join \", \" .items}}",
			data:           Data{"items": []string{"apple", "banana", "cherry"}},
			expectedContent: "First: apple, Last: cherry, Join: apple, banana, cherry",
		},
		{
			name:           "default_function",
			templateContent: "Value: {{.missing | default \"fallback\"}}",
			data:           Data{},
			expectedContent: "Value: fallback",
		},
		{
			name:           "quote_function",
			templateContent: "Quoted: {{.text | quote}}",
			data:           Data{"text": "hello world"},
			expectedContent: "Quoted: \"hello world\"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create template file
			fileName := tc.name + ".tmpl"
			templateFilePath := filepath.Join(templateDir, fileName)
			if err := os.WriteFile(templateFilePath, []byte(tc.templateContent), 0644); err != nil {
				t.Fatalf("failed to write template file: %v", err)
			}

			// Template the files
			err = templateFiles(templateDir, outputDir, tc.data)
			if err != nil {
				t.Fatalf("templateFiles returned error: %v", err)
			}

			// Check output
			outputFilePath := filepath.Join(outputDir, fileName)
			outputContentBytes, err := os.ReadFile(outputFilePath)
			if err != nil {
				t.Fatalf("failed to read output file: %v", err)
			}
			outputContent := string(outputContentBytes)
			if outputContent != tc.expectedContent {
				t.Errorf("expected file content %q, got %q", tc.expectedContent, outputContent)
			}

			// Clean up for next test
			os.Remove(templateFilePath)
			os.Remove(outputFilePath)
		})
	}
}

func TestTemplateFilesSprigDateFunctions(t *testing.T) {
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

	// Create a template that uses sprig date functions
	fileName := "date.tmpl"
	templateFilePath := filepath.Join(templateDir, fileName)
	templateContent := "{{now | date \"2006-01-02\"}} is today"
	if err := os.WriteFile(templateFilePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("failed to write template file: %v", err)
	}

	data := Data{}

	// Template the files
	err = templateFiles(templateDir, outputDir, data)
	if err != nil {
		t.Fatalf("templateFiles returned error: %v", err)
	}

	// Check that output file exists and contains date pattern
	outputFilePath := filepath.Join(outputDir, fileName)
	outputContentBytes, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	outputContent := string(outputContentBytes)
	
	// Check that output contains a date pattern (YYYY-MM-DD) and "is today"
	if len(outputContent) < 15 || !contains(outputContent, "is today") {
		t.Errorf("expected output to contain date and 'is today', got %q", outputContent)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
