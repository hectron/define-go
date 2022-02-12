package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var (
	homeDir, _ = os.UserHomeDir()
	ConfigPath = filepath.Join(homeDir, ".config", AppName)
)

func LoadConfig(c *cli.Context) {
	configFilePath := c.String("config")

	if configFilePath == "" {
		return
	}

	configPath, filename := filepath.Split(configFilePath)

	if configPath == "" {
		configPath = "."
	}

	viper.SetConfigName(filename)
	viper.SetConfigType("yml")
	viper.AddConfigPath(configPath)

	viper.ReadInConfig()
}
