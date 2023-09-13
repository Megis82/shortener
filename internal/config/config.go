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

func ParseConfig() {

	//var addr string
	flag.StringVar(&ServerConfig.NetAddress, "a", ":8080", "network address")
	flag.StringVar(&ServerConfig.BaseURL, "b", "http://localhost:8080/", "base URL")
	flag.Parse()

	// addrSlice := strings.Split(addr, ":")
	// fmt.Println(addrSlice)
	// fmt.Println(len(addrSlice))
	// if (len(addrSlice[0]) != 0) {
	// 	ServerConfig.NetAddress = net.ParseIP(addrSlice[0])
	// }
	// ServerConfig.NetAddress = net.ParseIP(addrSlice[0])
	// fmt.Println(ServerConfig.NetAddress.String())
	// ServerConfig.Port, _ = strconv.Atoi(addrSlice[1])
}
