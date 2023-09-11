package handlers

import (
	"fmt"
	"io"
	"net/http"

	"internal/code"
	"internal/storage"
)

func processPost(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	hashString := code.CodeString(string(body))
	storage.Add(hashString, string(body))

	retBody := fmt.Sprintf("http://localhost:8080/%s", hashString)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(201)
	w.Write([]byte(retBody))
}

func processGet(w http.ResponseWriter, r *http.Request) {
	body := r.URL.Path[1:]
	w.Header().Set("Content-Type", "text/plain")
	redirectURL := storage.Find(body)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		processPost(w, r)
	} else if r.Method == http.MethodGet {
		processGet(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}
