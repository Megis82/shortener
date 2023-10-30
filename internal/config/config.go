package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	env "github.com/caarlos0/env/v6"
)

type Config struct {
	NetAddress  string
	BaseURL     string
	FileStorage string
	DatabaseDSN string
	//DatabaseType string
}

type configEnv struct {
	NetAddress  string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	FileStorage string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN string `env:"DATABASE_DSN"`

	//DatabaseType string `env:"DBTYPE"`
}

func Init() (Config, error) {

	var ServerConfig Config
	tmpFile := filepath.Join(os.TempDir(), "short-url-db.json")
	flag.StringVar(&ServerConfig.NetAddress, "a", ":8080", "network address")
	flag.StringVar(&ServerConfig.BaseURL, "b", "", "base URL")
	flag.StringVar(&ServerConfig.FileStorage, "f", tmpFile, "file storage")
	flag.StringVar(&ServerConfig.DatabaseDSN, "d", "", "database storage")
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
		ServerConfig.FileStorage = cfg.FileStorage
	}

	if cfg.DatabaseDSN != "" {
		ServerConfig.DatabaseDSN = cfg.DatabaseDSN
	}

	return ServerConfig, nil
}
