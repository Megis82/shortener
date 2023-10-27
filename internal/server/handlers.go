package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Megis82/shortener/internal/code"
	"github.com/Megis82/shortener/internal/storage"

	chi "github.com/go-chi/chi/v5"
)

func (s *Server) PostLinkAdd(w http.ResponseWriter, r *http.Request) {
	CodeStatus := 201

	body, _ := io.ReadAll(r.Body)
	reqURL := string(body)

	// fmt.Println("URL start")
	// fmt.Println([]byte(reqURL))
	// fmt.Println("URL end")

	hashString := code.CodeString(reqURL)

	err := s.storage.Add(r.Context(), hashString, reqURL)

	if err != nil {

		if errors.Is(err, storage.ErrConflict) {
			// ошибка специфична для случая конфликта имён пользователей
			// w.WriteHeader(409)
			CodeStatus = 409

			hashString, _ = s.storage.FindShortByFullPath(r.Context(), reqURL)
		} else {

			//fmt.Println("Ошибка работы с хранилищем (((((", err)
			return
		}
	}

	retURL := ""
	if s.config.BaseURL == "" {
		requestURL := r.Host
		retURL = fmt.Sprintf("%s%s/", "http://", requestURL)
	} else {
		retURL = fmt.Sprintf("%s/", s.config.BaseURL)
	}

	retBody := fmt.Sprintf("%s%s", retURL, hashString)

	w.WriteHeader(CodeStatus)
	w.Write([]byte(retBody))
}

func (s *Server) GetLinkAdd(w http.ResponseWriter, r *http.Request) {
	body := chi.URLParam(r, "shortURL")
	if redirectURL, _, err := s.storage.Find(r.Context(), body); err != nil {

		htmlFile, err := os.Open("static/error404.html")
		if err != nil {
			//error
			//fmt.Println("Ошибка открытия файла:", err)
			return
		}
		defer htmlFile.Close()

		htmlData, err := io.ReadAll(htmlFile)
		if err != nil {
			//error
			// fmt.Println("Ошибка чтения файла:", err)
			return
		}
		w.WriteHeader(404)
		w.Write(htmlData)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}
}

func (s *Server) PostAPILinkAddBatch(w http.ResponseWriter, r *http.Request) {

	type req struct {
		CorrID  string `json:"correlation_id"`
		OrigURL string `json:"original_url"`
	}

	type answ struct {
		CorrID   string `json:"correlation_id"`
		ShortURL string `json:"short_url"`
	}

	reqBody := make([]req, 0)
	answBody := make([]answ, 0)
	batchArray := make(map[string]string, 0)
	jsonReq, _ := io.ReadAll(r.Body)

	err := json.Unmarshal(jsonReq, &reqBody)

	if err != nil {
		// fmt.Println("Ошибка при парсинге JSON:", err)
		return
	}

	for _, reqStr := range reqBody {
		hashString := code.CodeString(reqStr.OrigURL)

		batchArray[hashString] = reqStr.OrigURL

		retURL := ""
		if s.config.BaseURL == "" {
			requestURL := r.Host
			retURL = fmt.Sprintf("%s%s/", "http://", requestURL)
		} else {
			retURL = fmt.Sprintf("%s/", s.config.BaseURL)
		}

		retBody := fmt.Sprintf("%s%s", retURL, hashString)

		answBody = append(answBody, answ{CorrID: reqStr.CorrID, ShortURL: retBody})
	}

	s.storage.AddBatch(r.Context(), batchArray)

	jsonAns, err := json.Marshal(answBody)

	if err != nil {
		// fmt.Println("Ошибка при преобразовании структуры в JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(jsonAns))
}

// func (s *Server) PostAPILinkAdd(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")
// 	CodeStatus := 201

// 	var req struct {
// 		URL string `json:"url"`
// 	}

// 	jsonReq, _ := io.ReadAll(r.Body)

// 	err := json.Unmarshal(jsonReq, &req)

// 	if err != nil {
// 		fmt.Println("Ошибка при парсинге JSON:", err)
// 		return
// 	}

// 	hashString := code.CodeString(req.URL)

// 	// var shortPath string

// 	err = s.storage.Add(r.Context(), hashString, req.URL)

// 	if err != nil {

// 		if errors.Is(err, storage.ErrConflict) {
// 			// ошибка специфична для случая конфликта имён пользователей
// 			// w.WriteHeader(409)
// 			CodeStatus = 409
// 			fmt.Println("!!!!!!!!!!!!Ошибка работы с хранилищем", err)

// 			hashString, _ = s.storage.FindShortByFullPath(r.Context(), req.URL)
// 		} else {

// 			fmt.Println("Ошибка работы с хранилищем (((((", err)
// 			return
// 		}
// 	}

// 	retURL := ""
// 	if s.config.BaseURL == "" {
// 		requestURL := r.Host
// 		retURL = fmt.Sprintf("%s%s/", "http://", requestURL)
// 	} else {
// 		retURL = fmt.Sprintf("%s/", s.config.BaseURL)
// 	}

// 	retBody := fmt.Sprintf("%s%s", retURL, hashString)

// 	var ans struct {
// 		Result string
// 	}

// 	ans.Result = retBody

// 	jsonAns, err := json.Marshal(ans)

// 	if err != nil {
// 		fmt.Println("Ошибка при преобразовании структуры в JSON:", err)
// 		return
// 	}

// 	w.WriteHeader(CodeStatus)
// 	w.Write([]byte(jsonAns))
// }
