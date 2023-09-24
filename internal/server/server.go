package server

import (
	"net/http"

	"github.com/Megis82/shortener/internal/config"
	"github.com/Megis82/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	config  config.Config
	storage storage.DataStorage
	routers *chi.Mux
}

func NewServer(conf config.Config, stor storage.DataStorage) (*Server, error) {
	server := &Server{config: conf, storage: stor, routers: chi.NewRouter()}
	server.routes()
	return server, nil
}

func (s *Server) routes() {
	s.routers.Get("/{shortURL}", s.ProcessGet)
	s.routers.Post("/", s.ProcessPost)
	//router.NotFound()
	//return router, nil
}

func (s *Server) Run() {
	err1 := http.ListenAndServe(s.config.NetAddress, s.routers)

	if err1 != nil {
		panic(err1)
	}
}
