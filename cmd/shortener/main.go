package main

import (
	b64 "encoding/base64"
	"fmt"
	"io"
	"net/http"
)

var memStor map[string]string

func codeString(input string) string {
	retStr := b64.StdEncoding.EncodeToString([]byte(input))
	return retStr[:10]
}

func processPost(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	hashString := codeString(string(body))
	memStor[hashString] = string(body)

	retBody := fmt.Sprintf("http://localhost:8080/%s", hashString)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(201)
	w.Write([]byte(retBody))
}

func processGet(w http.ResponseWriter, r *http.Request) {
	body := r.URL.Path[1:]
	w.Header().Set("Content-Type", "text/plain")
	redirectURL := memStor[body]
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		processPost(w, r)
	} else if r.Method == http.MethodGet {
		processGet(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

func main() {

	memStor = make(map[string]string, 100)

	mux := http.NewServeMux()

	mux.HandleFunc(`/`, mainPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
