package main

import (
	"github.com/Megis82/shortener/internal/config"
	"github.com/Megis82/shortener/internal/server"
	"github.com/Megis82/shortener/internal/storage"
)

func main() {

	config, err := config.Init()

	if err != nil {
		panic(err)
	}

	storage, err := storage.NewMemoryStorage()

	if err != nil {
		panic(err)
	}

	server, err := server.NewServer(config, storage)

	if err != nil {
		panic(err)
	}

	server.Run()

}
