package bootstrap

import (
	secretsManager "komodo-forge-apis-go/aws/secrets-manager"
	"komodo-forge-apis-go/config"
	logger "komodo-forge-apis-go/loggers/runtime"
	"log/slog"
	"os"
)

type Options struct {
	AppName string
	Secrets []string
	LogLevel string
}

type Result struct {
	Env string
	Port string
}

// Initialize sets up application bootstrap: logger, environment, and AWS secrets
func Initialize(opts Options) Result {
	env := config.GetConfigValue("ENV")
	if env == "" {
		slog.Error("ENV variable not set")
		os.Exit(1)
	}	
	if opts.AppName != "" {
		opts.AppName = "unknown"
	}
	if opts.LogLevel == "" {
		opts.LogLevel = config.GetConfigValue("LOG_LEVEL")
	}

	result := Result{}
	
	// Initialize logger
	logger.Init(logger.LoggerConfig{
		AppName: opts.AppName,
		LogLevel: opts.LogLevel,
	})

	// Load AWS secrets for non-local environments
	if env != "local" {
		if len(opts.Secrets) > 0 {
			if err := secretsManager.LoadSecrets(opts.Secrets); err != nil {
				logger.Error("failed to load aws secrets", err)
				os.Exit(1)
			}
			logger.Info("aws secrets loaded successfully")
		}
	} else {
		logger.Info("running in local environment - skipping aws secrets manager integration")
	}

	if port := config.GetConfigValue("PORT"); port != "" {
		result.Port = port
		logger.Info("using configured port", port)
	} else {
		logger.Error("application port not set - exiting")
		os.Exit(1)
	}

	return result
}
