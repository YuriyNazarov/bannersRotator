package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	Database DBConf     `json:"database"`
	Logger   LoggerCfg  `json:"logs"`
	Queue    QueueCfg   `json:"queue"`
	Server   Connection `json:"server"`
}

type DBConf struct {
	DSN string `json:"dsn"`
}

type LoggerCfg struct {
	Level       string `json:"level"`
	Destination string `json:"destination"`
}

type QueueCfg struct {
	DSN      string `json:"dsn"`
	Exchange string `json:"exchange"`
	Queue    string `json:"queue"`
}

type Connection struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func LoadConfig() (Config, error) {
	config := Config{}
	pwd, err := os.Getwd()
	if err != nil {
		return config, fmt.Errorf("could not open config file: %w", err)
	}
	path := filepath.Join(pwd, "../../configs/config.json")
	file, err := os.Open(path)
	if err != nil {
		return config, fmt.Errorf("could not open config file: %w", err)
	}
	rawCfg, err := io.ReadAll(file)
	if err != nil {
		return config, fmt.Errorf("could not parse config, %w", err)
	}
	err = json.Unmarshal(rawCfg, &config)
	return config, err
}
