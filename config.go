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
	TimeoutInSecs   int
}

const (
	DefaultNumBuildsToKeep = 10
	InvalidTimeoutInSecs   = -1
)

var DefaultBuildScriptArgs = []string{}

func ParseConfigFile(path string) (*Config, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("Could not read config file: %s", err)
	}

	config := Config{NumBuildsToKeep: DefaultNumBuildsToKeep, BuildScriptArgs: DefaultBuildScriptArgs, TimeoutInSecs: InvalidTimeoutInSecs}

	decoder := json.NewDecoder(file)

	if err = decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("Error parsing json: %s", err)
	}

	if err = checkRequiredConfig(config); err != nil {
		return nil, err
	}

	return &config, nil
}

func checkRequiredConfig(config Config) error {
	if config.BuildScript == "" {
		return fmt.Errorf("BuildScript is required in the config.")
	}

	if config.TimeoutInSecs == InvalidTimeoutInSecs {
		return fmt.Errorf("TimeoutInSecs is required in the config.")
	}

	return nil
}
