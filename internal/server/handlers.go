package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Megis82/shortener/internal/code"

	chi "github.com/go-chi/chi/v5"
)

func (s *Server) PostLinkAdd(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	hashString := code.CodeString(string(body))
	s.storage.Add(hashString, string(body))
	s.storage.Close()

	retURL := ""
	if s.config.BaseURL == "" {
		requestURL := r.Host
		retURL = fmt.Sprintf("%s%s/", "http://", requestURL)
	} else {
		retURL = fmt.Sprintf("%s/", s.config.BaseURL)
	}
	//fmt.Println(retURL, "req url")
	retBody := fmt.Sprintf("%s%s", retURL, hashString)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(201)
	w.Write([]byte(retBody))
}

func (s *Server) GetLinkAdd(w http.ResponseWriter, r *http.Request) {
	body := chi.URLParam(r, "shortURL")
	if redirectURL, _, err := s.storage.Find(body); err != nil {

		htmlFile, err := os.Open("static/error404.html")
		if err != nil {
			//error
			fmt.Println("Ошибка открытия файла:", err)
			return
		}
		defer htmlFile.Close()

		htmlData, err := io.ReadAll(htmlFile)
		if err != nil {
			//error
			fmt.Println("Ошибка чтения файла:", err)
			return
		}
		w.WriteHeader(404)
		w.Write(htmlData)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}
}

func (s *Server) PostApiLinkAdd(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Url string
	}

	jsonReq, _ := io.ReadAll(r.Body)

	err := json.Unmarshal(jsonReq, &req)

	if err != nil {
		fmt.Println("Ошибка при парсинге JSON:", err)
		return
	}

	hashString := code.CodeString(req.Url)
	s.storage.Add(hashString, req.Url)

	retURL := ""
	if s.config.BaseURL == "" {
		requestURL := r.Host
		retURL = fmt.Sprintf("%s%s/", "http://", requestURL)
	} else {
		retURL = fmt.Sprintf("%s/", s.config.BaseURL)
	}

	retBody := fmt.Sprintf("%s%s", retURL, hashString)

	var ans struct {
		Result string
	}

	ans.Result = retBody

	jsonAns, err := json.Marshal(ans)

	if err != nil {
		fmt.Println("Ошибка при преобразовании структуры в JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(jsonAns))
}
