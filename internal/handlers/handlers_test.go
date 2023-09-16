package handlers_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	chi "github.com/go-chi/chi/v5"

	"github.com/Megis82/shortener/internal/code"
	"github.com/Megis82/shortener/internal/config"
	"github.com/Megis82/shortener/internal/handlers"
	"github.com/Megis82/shortener/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainPage(t *testing.T) {
	storage.Init()
	config.ParseConfig()
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name     string
		body     string
		respAddr string
		wantPost want
		wantGet  want
	}{
		{
			name: "positive test #1",
			body: "http://mail.ru",
			wantPost: want{
				code:        201,
				contentType: "text/plain",
			},
			wantGet: want{
				code:        307,
				contentType: "text/plain",
			},
		},
		{
			name: "positive test #2",
			body: "http://yandex.ru/defaultPage",
			wantPost: want{
				code:        201,
				contentType: "text/plain",
			},
			wantGet: want{
				code:        307,
				contentType: "text/plain",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestPost := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))

			wPost := httptest.NewRecorder()
			handlers.ProcessPost(wPost, requestPost)

			resPost := wPost.Result()

			assert.Equal(t, resPost.StatusCode, test.wantPost.code)
			defer resPost.Body.Close()
			responseBody, err := io.ReadAll(resPost.Body)

			require.NoError(t, err)
			assert.Equal(t, resPost.Header.Get("Content-Type"), test.wantPost.contentType)

			//requestGet := httptest.NewRequest(http.MethodGet, string(responseBody), nil)
			shortURL := string(responseBody)

			u, err := url.Parse(shortURL)
			if err != nil {
				fmt.Println("Ошибка при разборе URL:", err)
				return
			}

			// Удаление домена
			u.Host = ""

			// Удаление порта
			//u.Port = ""

			// Удаление протокола
			u.Scheme = ""

			// Получение обновленной URL
			newURL := u.String()[1:]

			requestGet := httptest.NewRequest(http.MethodGet, "/{shortURL}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("shortURL", newURL)
			wGet := httptest.NewRecorder()

			requestGet = requestGet.WithContext(context.WithValue(requestGet.Context(), chi.RouteCtxKey, rctx))

			handlers.ProcessGet(wGet, requestGet)

			resGet := wGet.Result()

			fmt.Println("Заголовки:")
			for key, value := range resGet.Header {
				fmt.Printf("%s: %s\n", key, value)
			}

			defer resGet.Body.Close()

			assert.Equal(t, resGet.StatusCode, test.wantGet.code)

			require.NoError(t, err)
			assert.Equal(t, resGet.Header.Get("Content-Type"), test.wantGet.contentType)
			assert.Equal(t, resGet.Header.Get("Location"), test.body)

			fmt.Println(string(test.body), " test body")
			fmt.Println(resGet.Header.Get("Location"), " answer body")

			fmt.Println(storage.Find(code.CodeString(test.body)))

		})
	}
}

// server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "https://example.com", http.StatusMovedPermanently)
//    }))

//    // Создаем клиент с настройкой запрета редиректа
//    client := &http.Client{
// 	CheckRedirect: func(req *http.Request, via []*http.Request) error {
// 	 return http.ErrUseLastResponse
// 	},
//    }

//    // Создаем запрос к серверу
//    req, _ := http.NewRequest("GET", server.URL, nil)

//    // Отправляем запрос
//    resp, err := client.Do(req)
//    if err != nil {
// 	fmt.Println("Error:", err)
// 	return
//    }
