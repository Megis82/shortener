package server

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"go.uber.org/zap"
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

func printCallStack() {
	// Создаем буфер для хранения адресов вызовов функций
	callers := make([]uintptr, 10)

	// Получаем стек вызова функций
	numCallers := runtime.Callers(0, callers)

	// Преобразуем адреса в названия функций
	frames := runtime.CallersFrames(callers[:numCallers])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		fmt.Printf("- %s\n", frame.Function)
	}
}
