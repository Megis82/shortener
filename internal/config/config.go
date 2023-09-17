package config

import (
	"flag"
)

type Config struct {
	//NetAddress net.IP
	NetAddress string
	//Port       int
	BaseURL string
}

var ServerConfig Config
var RunningConfig Config

func ParseConfig() {

	flag.StringVar(&ServerConfig.NetAddress, "a", ":8080", "network address")
	flag.StringVar(&ServerConfig.BaseURL, "b", "", "base URL")
	flag.Parse()
	//fmt.Println(ServerConfig.NetAddress, ServerConfig.BaseURL)
}
