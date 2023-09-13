package config

import (
	"flag"
	"net"
	"strconv"
	"strings"
)

type Config struct {
	NetAddress net.IP
	Port       int
	BaseURL    string
}

var ServerConfig Config

func ParseConfig() {

	var addr string
	flag.StringVar(&addr, "a", "localhost:8080", "network address")
	flag.StringVar(&ServerConfig.BaseURL, "b", "http://localhost:8080/", "base URL")
	flag.Parse()

	addrSlice := strings.Split(addr, ":")
	ServerConfig.NetAddress = net.ParseIP(addrSlice[0])
	ServerConfig.Port, _ = strconv.Atoi(addrSlice[1])
}
