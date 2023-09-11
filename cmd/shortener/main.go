package main

import (
	"net/http"

	"internal/handlers"
	"internal/storage"
)

func main() {

	storage.Init()

	mux := http.NewServeMux()

	mux.HandleFunc(`/`, handlers.MainPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
