package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	BuildScript     string
	BuildScriptArgs []string
	NumBuildsToKeep int
}

const DefaultNumBuildsToKeep = 10

var DefaultBuildScriptArgs = []string{}

func ParseConfigFile(path string) (*Config, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("Could not read config file: %s", err)
	}

	config := Config{NumBuildsToKeep: DefaultNumBuildsToKeep, BuildScriptArgs: DefaultBuildScriptArgs}

	decoder := json.NewDecoder(file)

	if err = decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("Error parsing json: %s", err)
	}

	return &config, nil
}
