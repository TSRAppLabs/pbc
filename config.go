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
	DataDir        string `json:"datadir"`
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
		config = defaultConfig
	}

	config = reconcileConfig(config)

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
	}

	return config, err
}

var _data_dir string

func getDataDir() string {
	if _data_dir == "" {
		_data_dir = ExpandPath(replaceTilde(GetConfig().DataDir))
	}
	return _data_dir
}

func reconcileConfig(config PBCConfig) PBCConfig {
	if config.DataDir == "" {
		config.DataDir = defaultConfig.DataDir
	}

	return config
}

var defaultConfig = PBCConfig{
	IgnorePatterns: []string{},
	DataDir:        "~/pbc/",
}
