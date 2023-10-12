package config

import (
	"flag"
	"log"

	env "github.com/caarlos0/env/v6"
)

type Config struct {
	NetAddress  string
	BaseURL     string
	FileStorage string
	//DatabaseType string
}

type configEnv struct {
	NetAddress  string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	FileStorage string `env:"FILE_STORAGE_PATH"`
	//DatabaseType string `env:"DBTYPE"`
}

func Init() (Config, error) {

	var ServerConfig Config

	flag.StringVar(&ServerConfig.NetAddress, "a", ":8080", "network address")
	flag.StringVar(&ServerConfig.BaseURL, "b", "", "base URL")
	flag.StringVar(&ServerConfig.FileStorage, "f", "/tmp/short-url-db.json", "file storage")
	//flag.StringVar(&ServerConfig.DatabaseType, "dbtype", "memory", "Database type")
	flag.Parse()

	var cfg configEnv

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.NetAddress != "" {
		ServerConfig.NetAddress = cfg.NetAddress
	}

	if cfg.BaseURL != "" {
		ServerConfig.BaseURL = cfg.BaseURL
	}

	if cfg.FileStorage != "" {
		ServerConfig.FileStorage = cfg.BaseURL
	}

	// if cfg.DatabaseType != "" {
	// 	ServerConfig.DatabaseType = cfg.DatabaseType
	// }

	return ServerConfig, nil
}
