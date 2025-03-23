package main

import (
	_ "embed"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//go:embed schema.cue
var schemaFile string

type Data map[string]any

var (
	basePath     string
	overlayPath  string
	templatePath string
	logMode      string
	logger       *zap.Logger
)

func initLogger() error {
	var err error

	if len(logMode) == 0 {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return err
		}

	}

	switch logMode {
	case "human":
		logger, err = zap.NewDevelopment()
		if err != nil {
			return err
		}
	case "json":
		logger, err = zap.NewProduction()
		if err != nil {
			return err
		}
	default:
		logger, err = zap.NewProduction()
		if err != nil {
			return err
		}
	}
	defer logger.Sync()
	return nil
}

// Decryptor interface for decrypting secrets
type Decryptor interface {
	Decrypt(secretName string) (string, error)
}

// ExampleDecryptor is a simple implementation of Decryptor
type ExampleDecryptor struct{}

func (d *ExampleDecryptor) Decrypt(secretName string) (string, error) {
	// Implement your decryption logic here
	// For example, you might look up the secret in a secure store
	// Here, we'll just return a dummy value for demonstration
	return "decryptedValue", nil
}

// decryptSecrets recursively decrypts values in a Data map.
func decryptSecrets(data *Data, decryptor Decryptor) error {
	secretPattern := regexp.MustCompile(`\$\$(.*?)\$\$`)
	for k, v := range *data {
		switch value := v.(type) {
		case string:
			// Handle string values.
			matches := secretPattern.FindAllStringSubmatch(value, -1)
			newStr := value
			for _, match := range matches {
				secretName := match[1]
				decryptedValue, err := decryptor.Decrypt(secretName)
				if err != nil {
					return err
				}
				newStr = strings.ReplaceAll(newStr, match[0], decryptedValue)
			}
			(*data)[k] = newStr
		case map[string]interface{}:
			// Recursively process nested maps.
			nestedData := Data(value)
			if err := decryptSecrets(&nestedData, decryptor); err != nil {
				return err
			}
			(*data)[k] = nestedData
		case []interface{}:
			// Process slices using a helper function.
			newList, err := decryptSecretsInSlice(value, decryptor)
			if err != nil {
				return err
			}
			(*data)[k] = newList
		}
	}
	return nil
}

// decryptSecretsInSlice recursively processes entries in a slice.
func decryptSecretsInSlice(list []interface{}, decryptor Decryptor) ([]interface{}, error) {
	secretPattern := regexp.MustCompile(`\$\$(.*?)\$\$`)
	for i, v := range list {
		switch value := v.(type) {
		case string:
			matches := secretPattern.FindAllStringSubmatch(value, -1)
			newStr := value
			for _, match := range matches {
				secretName := match[1]
				decryptedValue, err := decryptor.Decrypt(secretName)
				if err != nil {
					return nil, err
				}
				newStr = strings.ReplaceAll(newStr, match[0], decryptedValue)
			}
			list[i] = newStr
		case map[string]interface{}:
			nestedData := Data(value)
			if err := decryptSecrets(&nestedData, decryptor); err != nil {
				return nil, err
			}
			list[i] = nestedData
		case []interface{}:
			// Call recursively for nested slices.
			newList, err := decryptSecretsInSlice(value, decryptor)
			if err != nil {
				return nil, err
			}
			list[i] = newList
		}
	}
	return list, nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "dingo",
		Short: "Merges and validates data to template $stuff",
		Run: func(cmd *cobra.Command, args []string) {
			err := initLogger()
			if err != nil {
				zap.Error(err)
				os.Exit(1)
			}
			mergedData, err := loadAndMergeYAMLFiles(basePath, overlayPath)
			if err != nil {
				logger.Error("failed to load YAML files",
					zap.Error(err),
					zap.String("basePath", basePath),
					zap.String("overlayPath", overlayPath),
				)
				os.Exit(1)
			}

			if err := validateData(mergedData); err != nil {
				logger.Error("validation failed",
					zap.Error(err),
					zap.Any("data", mergedData),
				)
				os.Exit(1)
			}

			// Decrypt secrets in mergedData
			decryptor := &ExampleDecryptor{}
			if err := decryptSecrets(&mergedData, decryptor); err != nil {
				logger.Error("secret decryption failed",
					zap.Error(err),
					zap.Any("data", mergedData),
				)
				os.Exit(1)
			}

			logger.Info("data loaded and validated successfully",
				zap.Any("data", mergedData),
			)

			if err := templateFiles("./templates", "./output", mergedData); err != nil {
				logger.Error("templating failed",
					zap.Error(err),
					zap.Any("data", mergedData),
				)
				os.Exit(1)
			}
		},
	}

	rootCmd.PersistentFlags().StringVar(&basePath, "basepath", "data/base", "Base directory for YAML files")
	rootCmd.PersistentFlags().StringVar(&overlayPath, "overlaypath", "data/overlays/dev", "Overlay directory for YAML files")
	rootCmd.PersistentFlags().StringVar(&templatePath, "templatepath", "templates", "Template files to template")
	rootCmd.PersistentFlags().StringVar(&logMode, "logmode", "human", "Log Mode, available values [human, json]")

	err := viper.BindPFlag("basepath", rootCmd.PersistentFlags().Lookup("basepath"))
	if err != nil {
		logger.Fatal("failed to bind flag",
			zap.Error(err),
			zap.String("flag", "basepath"),
		)
	}

	err = viper.BindPFlag("overlaypath", rootCmd.PersistentFlags().Lookup("overlaypath"))
	if err != nil {
		logger.Fatal("failed to bind flag",
			zap.Error(err),
			zap.String("flag", "overlaypath"),
		)
	}
	err = viper.BindPFlag("templatepath", rootCmd.PersistentFlags().Lookup("templatepath"))
	if err != nil {
		logger.Fatal("failed to bind flag",
			zap.Error(err),
			zap.String("flag", "templatepath"),
		)
	}
	err = viper.BindPFlag("logmode", rootCmd.PersistentFlags().Lookup("logmode"))
	if err != nil {
		logger.Fatal("failed to bind flag",
			zap.Error(err),
			zap.String("flag", "logmode"),
		)
	}

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("command execution failed",
			zap.Error(err),
		)
	}
}
