package main

import (
	"github.com/Megis82/shortener/internal/config"
	"github.com/Megis82/shortener/internal/logger"
	"github.com/Megis82/shortener/internal/server"
	"github.com/Megis82/shortener/internal/storage"
)

func main() {

	logger, err := logger.NewLogger()
	if err != nil {
		return
	}

	defer logger.Sync()

	config, err := config.Init()
	if err != nil {
		return
	}

	storage, err := storage.NewMemoryStorage(config)
	if err != nil {
		return
	}
	defer storage.Close()

	server, err := server.NewServer(config, storage, logger)
	if err != nil {
		return
	}

	server.Run()

}
