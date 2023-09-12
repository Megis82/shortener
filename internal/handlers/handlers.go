package handlers

import (
	"fmt"
	"io"
	"net/http"

	"internal/code"
	"internal/storage"

	"github.com/go-chi/chi/v5"
)

func ProcessPost(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	hashString := code.CodeString(string(body))
	storage.Add(hashString, string(body))

	retBody := fmt.Sprintf("http://localhost:8080/%s", hashString)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(201)
	w.Write([]byte(retBody))
}

func ProcessGet(w http.ResponseWriter, r *http.Request) {
	body := chi.URLParam(r, "shortURL")
	fmt.Println(body)
	w.Header().Set("Content-Type", "text/plain")
	redirectURL := storage.Find(body)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func InitRouters(router *chi.Mux) {
	router.Get("/{shortURL}", ProcessGet)
	router.Post("/", ProcessPost)
}
