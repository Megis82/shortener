package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Megis82/shortener/internal/code"
	"go.uber.org/zap"

	chi "github.com/go-chi/chi/v5"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func (s *Server) ProcessPost(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	hashString := code.CodeString(string(body))
	s.storage.Add(hashString, string(body))

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

func (s *Server) ProcessGet(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) ProcessPostApi(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)
		duration := time.Since(start)
		s.logger.Info("Request started",
			zap.String("method", method),
			zap.String("path", uri),
			zap.String("status", fmt.Sprint(responseData.status)),
			zap.String("size", fmt.Sprint(responseData.size)),
			zap.String("duration", duration.String()),
		)
	}
	return http.HandlerFunc(logFn)
}
