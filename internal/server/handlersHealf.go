package server

import "net/http"

func (s Server) handleGetHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.storage.Ping(); err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}
