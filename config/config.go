/*
Package config provides configuration management for the application.
It uses Viper for reading configuration settings from environment variables,
configuration files, and other sources.
*/
package config

import (
	"github.com/heyrovsky/disturbdb/utils"
	"github.com/spf13/viper"
)

// LOG_LEVEL defines the current log level for the application.
// It is set during the initialization of the configuration.
var (
	LOG_LEVEL string = "INFO"
)

// InitConfig initializes the configuration settings for the application.
// It imports environment variables and retrieves the LOG_LEVEL setting
// from the configuration sources defined in Viper.
func InitConfig() {
	utils.ImportEnv()                        // Import environment variables into the application
	LOG_LEVEL = viper.GetString("LOG_LEVEL") // Retrieve the LOG_LEVEL from the configuration
}
