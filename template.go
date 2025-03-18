package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

func templateFiles(templateDir, outputDir string, data Data) error {

  // clean the output path before regenerating it
  err := os.RemoveAll(outputDir)
  if err != nil {
    return err
  }

	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Generate the relative path from the template directory.
		relativePath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return fmt.Errorf("failed to determine relative path for %s: %w", path, err)
		}

		// Compute the corresponding output path.
		outputPath := filepath.Join(outputDir, relativePath)

		if info.IsDir() {
			// If it's a directory, create it in the output directory.
			if err := os.MkdirAll(outputPath, 0755); err != nil {
				return fmt.Errorf("failed to create output directory %s: %w", outputPath, err)
			}
			return nil
		}

		// Read the file content.
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", path, err)
		}

		// Parse and execute the template.
		tmpl, err := template.New(info.Name()).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", path, err)
		}

		// Write the template output to the output path.
		if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write templated file %s: %w", outputPath, err)
		}

		return nil
	})
}

