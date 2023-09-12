package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"internal/handlers"
	"internal/storage"
)

func main() {

	storage.Init()

	router := chi.NewRouter()

	handlers.InitRouters(router)

	err := http.ListenAndServe(`:8080`, router)

	if err != nil {
		panic(err)
	}
}
