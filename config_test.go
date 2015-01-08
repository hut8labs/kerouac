package main

import (
	"reflect"
	"testing"
)

func TestConfigBarfsOnNonExistingFile(t *testing.T) {
	config, err := ParseConfigFile("testfiles/not_a_json_file.json")

	if config != nil {
		t.Errorf("Non existent file returned a config %+v\n", config)
	}

	if err == nil {
		t.Errorf("Err was nil on non existing config file")
	}
}

func TestConfigBarfsOnBadJson(t *testing.T) {
	config, err := ParseConfigFile("testfiles/bad_json_config.json")

	if config != nil {
		t.Errorf("Bad json returned a config %+v\n", config)
	}

	if err == nil {
		t.Errorf("Err was nil on bad json file")
	}
}

func TestConfigParsesGoodJson(t *testing.T) {
	config, err := ParseConfigFile("testfiles/good_json_config.json")

	assertExpectedConfig(config, t, "good config")

	if err != nil {
		t.Errorf("Err was non-nil on good json file")
	}
}

func TestConfigParsesGoodJsonWithExtraFlags(t *testing.T) {
	config, err := ParseConfigFile("testfiles/json_config_with_unused.json")

	assertExpectedConfig(config, t, "config with extra flags")

	if err != nil {
		t.Errorf("Err was non-nil on good json file")
	}
}

func TestDefaultsOnNonRequired(t *testing.T) {
	config, err := ParseConfigFile("testfiles/good_json_only_required.json")

	if err != nil {
		t.Errorf("Err was non-nil on good json file with only required")
	}

	if config == nil {
		t.Errorf("Config was nil on good json file with only required")
	}

	if config.NumBuildsToKeep != DefaultNumBuildsToKeep {
		t.Errorf("Did not use default NumBuildsToKeep: %+v", config)
	}

	if !reflect.DeepEqual(config.BuildScriptArgs, DefaultBuildScriptArgs) {
		t.Errorf("Did not use default BuildScriptArgs: %+v", config)
	}
}

func assertExpectedConfig(config *Config, t *testing.T, context string) {
	if config == nil {
		t.Errorf("%s returned a nil config", context)
	}

	if config.NumBuildsToKeep != 22 {
		t.Errorf("%s wrong number of builds %+v", context, config)
	}

	if config.BuildScript != "build.sh" {
		t.Errorf("%s wrong build script %+v", context, config)
	}

	if !reflect.DeepEqual(config.BuildScriptArgs, []string{"arg1", "arg 2"}) {
		t.Errorf("%s wrong build script args %+v", context, config)
	}
}
