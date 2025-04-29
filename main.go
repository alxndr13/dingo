package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/alxndr13/dingo/decrypt"
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
	decryptor    string
	logger       *zap.Logger
)

func initDecryptor(decryptor string) (Decryptor, error) {
	switch decryptor {
	case "example":
		return decrypt.NewExampleDecryptor(), nil

	case "google":
		return decrypt.NewGoogleDecryptor(), nil

	}
	return nil, fmt.Errorf("no such decryptor")

}

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

			if len(decryptor) > 0 {
				// Decrypt secrets in mergedData
				decryptor, err := initDecryptor(decryptor)
				if err != nil {
					zap.Error(err)
				}
				if err := decryptSecrets(&mergedData, decryptor); err != nil {
					logger.Error("secret decryption failed",
						zap.Error(err),
						zap.Any("data", mergedData),
					)
					os.Exit(1)
				}

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
	rootCmd.PersistentFlags().StringVar(&decryptor, "decryptor", "", "Decryptor in case you're using secrets, leave empty if you do not want to use one. available values [example, google]")

	// bindFlags binds command line flags to viper configuration
	flags := []string{"basepath", "overlaypath", "templatepath", "logmode"}
	for _, flag := range flags {
		if err := viper.BindPFlag(flag, rootCmd.PersistentFlags().Lookup(flag)); err != nil {
			logger.Fatal("failed to bind flag",
				zap.Error(err),
				zap.String("flag", flag),
			)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("command execution failed",
			zap.Error(err),
		)
	}
}
