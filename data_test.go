package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadAndMergeYAMLFiles(t *testing.T) {
	// Create temporary directories for base and overlay YAML files.
	baseDir, err := os.MkdirTemp("", "baseDir")
	if err != nil {
		t.Fatalf("failed to create temporary base directory: %v", err)
	}
	defer os.RemoveAll(baseDir)

	overlayDir, err := os.MkdirTemp("", "overlayDir")
	if err != nil {
		t.Fatalf("failed to create temporary overlay directory: %v", err)
	}
	defer os.RemoveAll(overlayDir)

	// Create a YAML file in the base directory.
	baseYAMLContent := `
    base: value1
    overlay: mustbeOverWritten
    this: "should be in the final file"`
	baseYAMLPath := filepath.Join(baseDir, "base.yaml")
	if err := os.WriteFile(baseYAMLPath, []byte(baseYAMLContent), 0644); err != nil {
		t.Fatalf("failed to write base YAML file: %v", err)
	}

	// Create another YAML file in the overlay directory.
	overlayYAMLContent := "overlay: value2"
	overlayYAMLPath := filepath.Join(overlayDir, "overlay.yaml")
	if err := os.WriteFile(overlayYAMLPath, []byte(overlayYAMLContent), 0644); err != nil {
		t.Fatalf("failed to write overlay YAML file: %v", err)
	}

	// Call the function under test.
	merged, err := loadAndMergeYAMLFiles(baseDir, overlayDir)
	if err != nil {
		t.Fatalf("loadAndMergeYAMLFiles returned error: %v", err)
	}

	// Define the expected merged data.
	expected := Data{
		"base":    "value1",
		"overlay": "value2",
    "this": "should be in the final file",
	}

	if !reflect.DeepEqual(merged, expected) {
		t.Errorf("expected merged data %v, got %v", expected, merged)
	}
}


func TestLoadAndMergeYAMLFiles_InvalidYAML(t *testing.T) {
	// Create temporary directories for base and overlay YAML files.
	baseDir, err := os.MkdirTemp("", "baseDir")
	if err != nil {
		t.Fatalf("failed to create temporary base directory: %v", err)
	}
	defer os.RemoveAll(baseDir)

	overlayDir, err := os.MkdirTemp("", "overlayDir")
	if err != nil {
		t.Fatalf("failed to create temporary overlay directory: %v", err)
	}
	defer os.RemoveAll(overlayDir)

	// Create an invalid YAML file in the base directory.
	invalidYAMLContent := "invalid: yaml: : content" // This is intentionally malformed.
	invalidYAMLPath := filepath.Join(baseDir, "invalid.yaml")
	if err := os.WriteFile(invalidYAMLPath, []byte(invalidYAMLContent), 0644); err != nil {
		t.Fatalf("failed to write invalid YAML file: %v", err)
	}

	// Optionally, you can create a valid YAML file in the overlay directory.
	overlayYAMLContent := "overlay: valid"
	overlayYAMLPath := filepath.Join(overlayDir, "overlay.yaml")
	if err := os.WriteFile(overlayYAMLPath, []byte(overlayYAMLContent), 0644); err != nil {
		t.Fatalf("failed to write overlay YAML file: %v", err)
	}

	// Call the function and expect an error due to invalid YAML.
	_, err = loadAndMergeYAMLFiles(baseDir, overlayDir)
	if err == nil {
		t.Fatalf("expected an error due to invalid YAML, but got nil")
	}

	// Optionally, check that the error message contains specific text.
	expectedErrSubstring := "error parsing YAML"
	if err != nil && !strings.Contains(err.Error(), expectedErrSubstring) {
		t.Errorf("expected error to contain %q, but got %q", expectedErrSubstring, err.Error())
	}
}
