package configs

import (
	"github.com/spf13/viper"
)

/*
InitializeViper Function initializes viper to read config.yml file and environment variables in the application.
*/
func InitializeViper() error {
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	return viper.ReadInConfig()
}
