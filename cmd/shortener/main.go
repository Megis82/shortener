package main

import (
	"github.com/Megis82/shortener/internal/config"
	"github.com/Megis82/shortener/internal/logger"
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

	logger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	server, err := server.NewServer(config, storage, logger)
	if err != nil {
		panic(err)
	}

	server.Run()

}
