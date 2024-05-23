package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config struct defines the application's configuration structure.
type Config struct {
	Logger     LoggerCfg
	HTTPServer HTTPServerCfg
}

// LoggerCfg struct defines the logger configuration.
type LoggerCfg struct {
	Level string
}

// HTTPServerCfg struct defines the HTTP server configuration.
type HTTPServerCfg struct {
	Port int
	Host string
}

func loadConfig() (*Config, error) {
	// Set the name of the configurations file
	viper.SetConfigName("config")
	// Set the file type of the configurations file
	viper.SetConfigType("yml")
	// Set the path to look for the configurations file
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Warnf("config file not found: %v", err)
	}

	viper.AutomaticEnv()

	// set defaults
	viper.SetDefault("Logger.Level", "info")
	viper.SetDefault("HTTPServer.Port", 8080)
	viper.SetDefault("HTTPServer.Host", "0.0.0.0")

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("viper unmarshal to config failed: %v", err)
	}

	return &config, nil
}

func configureLogger(cfg *LoggerCfg) error {
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("can not parse logrus level: %v", err)
	}

	logrus.SetLevel(level)

	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
		PrettyPrint:      true,
	})

	return nil
}
