package server_test

import (
	"testing"
	//"github.com/Megis82/shortener/internal/handlers"
)

func TestMainPage(t *testing.T) {
	// storage.Init()
	// config.ParseConfig()
	// type want struct {
	// 	code        int
	// 	contentType string
	// }
	// tests := []struct {
	// 	name     string
	// 	body     string
	// 	respAddr string
	// 	wantPost want
	// 	wantGet  want
	// }{
	// 	{
	// 		name: "positive test #1",
	// 		body: "http://mail.ru",
	// 		wantPost: want{
	// 			code:        201,
	// 			contentType: "text/plain",
	// 		},
	// 		wantGet: want{
	// 			code:        307,
	// 			contentType: "text/plain",
	// 		},
	// 	},
	// 	{
	// 		name: "positive test #2",
	// 		body: "http://yandex.ru/defaultPage",
	// 		wantPost: want{
	// 			code:        201,
	// 			contentType: "text/plain",
	// 		},
	// 		wantGet: want{
	// 			code:        307,
	// 			contentType: "text/plain",
	// 		},
	// 	},
	// }

	// for _, test := range tests {
	// 	t.Run(test.name, func(t *testing.T) {
	// 		requestPost := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))

	// 		wPost := httptest.NewRecorder()
	// 		handlers.ProcessPost(wPost, requestPost)

	// 		resPost := wPost.Result()

	// 		assert.Equal(t, resPost.StatusCode, test.wantPost.code)
	// 		defer resPost.Body.Close()
	// 		responseBody, err := io.ReadAll(resPost.Body)

	// 		require.NoError(t, err)
	// 		assert.Equal(t, resPost.Header.Get("Content-Type"), test.wantPost.contentType)

	// 		//requestGet := httptest.NewRequest(http.MethodGet, string(responseBody), nil)
	// 		shortURL := string(responseBody)

	// 		u, err := url.Parse(shortURL)
	// 		if err != nil {
	// 			fmt.Println("Ошибка при разборе URL:", err)
	// 			return
	// 		}

	// 		u.Host = ""
	// 		u.Scheme = ""
	// 		newURL := u.String()[1:]

	// 		requestGet := httptest.NewRequest(http.MethodGet, "/{shortURL}", nil)
	// 		rctx := chi.NewRouteContext()
	// 		rctx.URLParams.Add("shortURL", newURL)
	// 		wGet := httptest.NewRecorder()

	// 		requestGet = requestGet.WithContext(context.WithValue(requestGet.Context(), chi.RouteCtxKey, rctx))

	// 		handlers.ProcessGet(wGet, requestGet)

	// 		resGet := wGet.Result()

	// 		// fmt.Println("Заголовки:")
	// 		// for key, value := range resGet.Header {
	// 		// 	fmt.Printf("%s: %s\n", key, value)
	// 		// }

	// 		defer resGet.Body.Close()

	// 		assert.Equal(t, resGet.StatusCode, test.wantGet.code)

	// 		require.NoError(t, err)
	// 		assert.Equal(t, resGet.Header.Get("Content-Type"), test.wantGet.contentType)
	// 		assert.Equal(t, resGet.Header.Get("Location"), test.body)

	// 		fmt.Println(string(test.body), " test body")
	// 		fmt.Println(resGet.Header.Get("Location"), " answer body")

	// 		fmt.Println(storage.Find(code.CodeString(test.body)))

	// 	})
	// }
}
