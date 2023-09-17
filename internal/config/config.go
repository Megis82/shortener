package config

import (
	"flag"
	"log"

	env "github.com/caarlos0/env/v6"
)

type Config struct {
	//NetAddress net.IP
	NetAddress string
	//Port       int
	BaseURL string
}

type configEnv struct {
	NetAddress string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

var ServerConfig Config

func ParseConfig() {

	flag.StringVar(&ServerConfig.NetAddress, "a", ":8080", "network address")
	flag.StringVar(&ServerConfig.BaseURL, "b", "", "base URL")
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

	//fmt.Println(ServerConfig.NetAddress, ServerConfig.BaseURL)
}
