package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Megis82/shortener/internal/code"
	"github.com/Megis82/shortener/internal/config"
	"github.com/Megis82/shortener/internal/storage"

	chi "github.com/go-chi/chi/v5"
)

func ProcessPost(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	hashString := code.CodeString(string(body))
	storage.Add(hashString, string(body))

	retBody := fmt.Sprintf("%s%s", config.ServerConfig.BaseURL, hashString)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(201)
	w.Write([]byte(retBody))
}

func ProcessGet(w http.ResponseWriter, r *http.Request) {
	body := chi.URLParam(r, "shortURL")
	//fmt.Println(body, " body fron inner")
	//shortURL := strings.TrimPrefix(r.URL.Host, "/")
	//fmt.Println(r.URL., " url from ")
	w.Header().Set("Content-Type", "text/plain")
	redirectURL := storage.Find(body)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func InitRouters(router *chi.Mux) {
	router.Get("/{shortURL}", ProcessGet)
	router.Post("/", ProcessPost)
}
