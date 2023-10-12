package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Megis82/shortener/internal/config"
	"github.com/Megis82/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	config  config.Config
	storage storage.DataStorage
	routers *chi.Mux
	logger  *zap.Logger
}

func NewServer(conf config.Config, stor storage.DataStorage, log *zap.Logger) (*Server, error) {
	server := &Server{config: conf, storage: stor, routers: chi.NewRouter(), logger: log}
	server.routes()
	return server, nil
}

func (s *Server) routes() {
	s.routers.Use(s.WithLogging)
	s.routers.Use(GzipHandle)
	s.routers.Get("/{shortURL}", s.ProcessGet)
	s.routers.Post("/", s.ProcessPost)
	s.routers.Post("/api/shorten", s.ProcessPostApi)
}

func (s *Server) Run() {
	//err1 := http.ListenAndServe(s.config.NetAddress, s.routers)

	var srv http.Server
	srv.Addr = s.config.NetAddress
	srv.Handler = s.routers

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed

}
