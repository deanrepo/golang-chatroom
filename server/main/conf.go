package main

import (
	"encoding/json"
	"os"
)

// Config represents configuration files
type Config struct {
	ListenAddress string `json:"listenAddress"`
}

// loadConfig loads configurations from file
func loadConfig(filename string) (config *Config, err error) {
	configFile, err := os.Open(filename)
	if err != nil {
		return
	}
	defer configFile.Close()

	config = &Config{}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(config)
	return
}
