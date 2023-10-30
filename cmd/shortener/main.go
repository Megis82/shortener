package main

import (
	"context"

	"github.com/Megis82/shortener/internal/config"
	"github.com/Megis82/shortener/internal/logger"
	"github.com/Megis82/shortener/internal/server"
	"github.com/Megis82/shortener/internal/storage"
)

func main() {
	ctx := context.Background()

	logger, err := logger.NewLogger()
	if err != nil {
		return
	}

	defer logger.Sync()

	config, err := config.Init()
	if err != nil {
		return
	}

	storage, err := storage.NewStorage(config)

	if err != nil {
		return
	}
	err = storage.Init()

	if err != nil {
		return
	}

	defer storage.Close()

	server, err := server.NewServer(config, storage, logger)
	if err != nil {
		return
	}

	server.Run(ctx)

}
