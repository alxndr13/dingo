package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"

	"github.com/goccy/go-yaml"
)

func mergeData(base, overlay Data) Data {
	for key, value := range overlay {
		if baseValue, exists := base[key]; exists {
			// If both values are maps, merge them recursively
			if baseMap, ok := baseValue.(map[string]any); ok {
				if overlayMap, ok := value.(map[string]any); ok {
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

func loadAndMergeYAMLFiles(baseDirPath string, overlayDirPath string) (Data, error) {
	mergedData := make(Data)
	var errs error
	// Walk through the baseDir data directory
	errWalkBase := filepath.Walk(baseDirPath, func(path string, info os.FileInfo, err error) error {
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

	if errWalkBase != nil {
		errs = errors.Join(errWalkBase)
	}

	errWalkOverlays := filepath.Walk(overlayDirPath, func(path string, info os.FileInfo, err error) error {
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

	if errWalkOverlays != nil {
		errs = errors.Join(errWalkOverlays)
	}

	if errs != nil {
		return nil, errs
	}

	return mergedData, nil
}
