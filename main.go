package main

import (
	"fmt"
	"os"
  "io"
  _ "embed"
	"path/filepath"
  "cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"

	"github.com/goccy/go-yaml"
)

//go:embed schema.cue
var schemaFile string

// Define a generic map to store YAML data
type Data map[string]interface{}

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

func loadAndMergeYAMLFiles(baseDirPath string, overlayDirPath string) (Data, error) {
	mergedData := make(Data)

	// Walk through the baseDir data directory
	err := filepath.Walk(baseDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a YAML file
		if !info.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			// Read the YAML file
      f, err := os.Open(path)
      if err != nil {
				return fmt.Errorf("error opening file %s: %v", path, err)
      }
			content, err := io.ReadAll(f)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", path, err)
			}

			// Parse YAML content
			var data Data
			if err := yaml.Unmarshal(content, &data); err != nil {
				return fmt.Errorf("error parsing YAML from %s: %v", path, err)
			}

			// Merge with existing data
			mergedData = mergeData(mergedData, data)
		}
		return nil
	})
	err = filepath.Walk(overlayDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a YAML file
		if !info.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			// Read the YAML file
      f, err := os.Open(path)
      if err != nil {
				return fmt.Errorf("error opening file %s: %v", path, err)
      }
			content, err := io.ReadAll(f)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", path, err)
			}

			// Parse YAML content
			var data Data
			if err := yaml.Unmarshal(content, &data); err != nil {
				return fmt.Errorf("error parsing YAML from %s: %v", path, err)
			}

			// Merge with existing data
			mergedData = mergeData(mergedData, data)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return mergedData, nil
}

func mergeData(base, overlay Data) Data {
	for key, value := range overlay {
		if baseValue, exists := base[key]; exists {
			// If both values are maps, merge them recursively
			if baseMap, ok := baseValue.(map[string]interface{}); ok {
				if overlayMap, ok := value.(map[string]interface{}); ok {
					base[key] = mergeData(baseMap, overlayMap)
					continue
				}
			}
		}
		// Otherwise, overlay value takes precedence
		base[key] = value
	}
	return base
}

func validateData(data Data) error {
  ctx := cuecontext.New()
	schema := ctx.CompileString(schemaFile).LookupPath(cue.ParsePath("#Schema"))

  dataAsCue := ctx.Encode(data)

  unified := schema.Unify(dataAsCue)
  err := unified.Validate()
  if err != nil {
    return err
  }
  return nil
}
