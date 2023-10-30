package server

import (
	"fmt"
	"net/http"
)

func (s Server) handleGetHealth(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = s.storage.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(err)
	w.WriteHeader(http.StatusOK)
}
