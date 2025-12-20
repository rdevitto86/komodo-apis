package bootstrap

import (
	secretsManager "komodo-forge-apis-go/aws/secrets-manager"
	"komodo-forge-apis-go/config"
	logger "komodo-forge-apis-go/logging/runtime"
	"os"
)

type Options struct {
	AppName   string
	Secrets   []string
	LoggerCfg *logger.LoggerConfig
}

type Result struct {
	Env  string
	Port string
}

// Initialize sets up application bootstrap: logger, environment, and AWS secrets
func Initialize(opts Options) Result {
	appName := "application"
	if opts.AppName != "" { appName = opts.AppName }

	// Initialize logger with AppName from options
	loggerCfg := logger.LoggerConfig{AppName: appName}
	if opts.LoggerCfg != nil {
		loggerCfg = *opts.LoggerCfg
		if loggerCfg.AppName == "" { loggerCfg.AppName = appName }
	}
	logger.InitLogger(loggerCfg)

	env := config.GetConfigValue("ENV")
	if env == "" {
		logger.Fatal("ENV variable not set")
		os.Exit(1)
	}	

	logger.Info("starting " + appName + " in " + env + " environment")

	// Load AWS secrets for non-local environments
	if env != "local" {
		if !secretsManager.IsUsingAWS() {
			logger.Fatal("aws secrets manager disabled in non-local environment")
			os.Exit(1)
		}

		if len(opts.Secrets) > 0 {
			if err := secretsManager.LoadSecrets(opts.Secrets); err != nil {
				logger.Fatal("failed to load aws secrets", err)
				os.Exit(1)
			}
			
			logger.Info("aws secrets loaded successfully")
		}
	} else {
		logger.Info("running in local environment - skipping aws secrets manager integration")
	}

	port := config.GetConfigValue("PORT")
	if port == "" {
		port = "8080" // Default port
		logger.Warn("PORT not set - defaulting to " + port)
	}

	return Result{ Env: env, Port: port }
}
