package main

import (
	_ "embed"
	"os"

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
