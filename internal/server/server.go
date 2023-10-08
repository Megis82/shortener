package server

import (
	"net/http"

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
	s.routers.Get("/{shortURL}", s.ProcessGet)
	s.routers.Post("/", s.ProcessPost)
	s.routers.Post("/api/shorten", s.ProcessPostApi)
}

func (s *Server) Run() {
	err1 := http.ListenAndServe(s.config.NetAddress, s.routers)

	if err1 != nil {
		panic(err1)
	}
}
