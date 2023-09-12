package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"internal/config"
	"internal/handlers"
	"internal/storage"
)

func main() {

	storage.Init()

	router := chi.NewRouter()
	config.ParseConfig()
	handlers.InitRouters(router)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.ServerConfig.NetAddress.String(), fmt.Sprint(config.ServerConfig.Port)), router)

	if err != nil {
		panic(err)
	}
}
