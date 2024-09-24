/*
Package utils provides utility functions for the application.
This includes functions for importing environment variables from
a configuration file using Viper.
*/
package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

// ImportEnv initializes the Viper configuration by loading settings
// from an environment file (.env) and setting default values.
// It sets the LOG_LEVEL to "INFO" if not provided in the environment.
func ImportEnv() {
	viper.SetConfigName(".env") // Set the name of the configuration file
	viper.SetConfigType("env")  // Set the type of the configuration file
	viper.AddConfigPath(".")    // Add the current directory as a config path

	viper.SetDefault("LOG_LEVEL", "INFO") // Set default log level to "INFO"

	viper.AutomaticEnv() // Automatically override config values with environment variables

	// Read in the configuration file
	if err := viper.ReadInConfig(); err != nil {
		// Handle the case where the configuration file is not found
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Print("Configuration file not loading error: ", err)
			panic(1) // Panic if there is an error loading the config
		}
	}
}
