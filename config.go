package pbc

import (
	"encoding/json"
	"os"
)

type PBCConfig struct {
	IgnorePatterns []string
}

var cacheConfig PBCConfig
var configIsCached bool

/*
  Gets the configuration for the passbook compiler
*/

func GetConfig() PBCConfig {
	if !configIsCached {
		cacheConfig = obtainConfig()
	}

	return cacheConfig
}

func obtainConfig() PBCConfig {
	config, err := configFromFile("~/.pbconfig")

	if err != nil {
		config = defaultConfig()
	}

	return config
}

func configFromFile(path string) (PBCConfig, error) {
	var config PBCConfig
	file, err := os.Open(path)

	if err != nil {
		return config, err
	}

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return config, err
	}

	return config, err
}

func defaultConfig() PBCConfig {
	return PBCConfig{
		IgnorePatterns: []string{},
	}
}
