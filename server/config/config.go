package config

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

type BOTConfig struct {
	Bot BOTConfigs
}

type BOTConfigs struct {
	TokenID   string `toml:"tokenID"`
	BotID     string `toml:"botID"`
	ChannelID string `toml:"channelID"`
	SignUpCh  string `toml:"signupCh"`
}

const (
	configpath = "config/config.toml"
)

func ReadAPIConfig() (APIConfigs, error) {
	var config APIConfig
	_, err := toml.DecodeFile(configpath, &config)
	if err != nil {
		log.Println("filed:read APIconfig")
	}
	return config.API, err
}

func ReadDBConfig() (DBConfigs, error) {
	var config DBConfig
	_, err := toml.DecodeFile(configpath, &config)
	if err != nil {
		log.Println("filed:read DBconfig")
	}
	return config.DB, err
}

func ReadBOTConfig() (BOTConfigs, error) {
	var config BOTConfig
	_, err := toml.DecodeFile(configpath, &config)
	if err != nil {
		log.Println("filed:read Botconfig")
	}
	return config.Bot, err
}
