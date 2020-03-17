package configs

import (
	"log"

	"github.com/BurntSushi/toml"
)

type APIConfig struct {
	API APIConfigs
}

type APIConfigs struct {
	Port    string `toml:"port"`
	View    string `toml:"view"`
	Version string `toml:"version"`
}

type DBConfig struct {
	DB DBConfigs
}

type DBConfigs struct {
	Dbms     string `toml:"dbms"`
	User     string `toml:"user"`
	Pass     string `toml:"pass"`
	Protocol string `toml:"protocol"`
	Dbname   string `toml:"dbname"`
}

func ReadAPIConfig(path string) (APIConfigs, error) {
	var config APIConfig
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		log.Println("filed:read APIconfig")
	}
	return config.API, err
}

func ReadDBConfig(path string) (DBConfigs, error) {
	var config DBConfig
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		log.Println("filed:read DBconfig")
	}
	return config.DB, err
}
