package pbc

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
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
		configIsCached = true
	}

	return cacheConfig
}

func obtainConfig() PBCConfig {
	config, err := configFromFile(getHomeConfigPath())

	if err != nil {
		fmt.Printf("Failing to get config %v", err)
		config = defaultConfig()
	}

	return config
}

func getHomeConfigPath() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".pbcconfig")
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
	} else {
		fmt.Printf("Using config file: %v\n", path)
	}

	return config, err
}

func defaultConfig() PBCConfig {
	return PBCConfig{
		IgnorePatterns: []string{},
	}
}
