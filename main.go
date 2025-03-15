package main

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed schema.cue
var schemaFile string

// Define a generic map to store YAML data
type Data map[string]any

func main() {
	// Load and merge YAML files
	mergedData, err := loadAndMergeYAMLFiles("data/base", "data/overlays/dev")
	if err != nil {
		fmt.Printf("Error loading YAML files: %v\n", err)
		os.Exit(1)
	}

	// Validate the merged data
	if err := validateData(mergedData); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data loaded and validated successfully")
	print("data:")
	fmt.Println(mergedData)
}
