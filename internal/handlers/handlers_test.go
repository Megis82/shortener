package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Megis82/shortener/internal/handlers"
	"github.com/Megis82/shortener/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainPage(t *testing.T) {
	storage.Init()
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

			requestGet := httptest.NewRequest(http.MethodGet, string(responseBody), strings.NewReader(test.body))

			wGet := httptest.NewRecorder()
			handlers.ProcessGet(wGet, requestGet)

			resGet := wGet.Result()
			assert.Equal(t, resGet.StatusCode, test.wantGet.code)

			require.NoError(t, err)
			assert.Equal(t, resGet.Header.Get("Content-Type"), test.wantGet.contentType)
			assert.Equal(t, resGet.Header.Get("Location"), test.body)

		})
	}
}
