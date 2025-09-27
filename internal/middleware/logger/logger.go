package logger

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Logger struct {
	Log *zap.Logger
}

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
		body   []byte
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func Initialize(level string) (*Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Log: zl}, nil
}

// Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
// Сведения об ответах должны содержать код статуса и размер содержимого ответа.
func (l *Logger) LogInfo(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   &responseData{status: 0, size: 0, body: []byte{}},
		}

		bodyBuf := new(bytes.Buffer)
		tee := io.TeeReader(r.Body, bodyBuf)
		body, _ := io.ReadAll(tee)
		r.Body = io.NopCloser(bodyBuf)

		start := time.Now()
		h(&lw, r)
		duration := time.Since(start)

		l.Log.Info("request log",
			zap.String("url", r.URL.String()),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("status", lw.responseData.status),
			zap.Int("size", lw.responseData.size),
			zap.ByteString("request_body", body),
			zap.ByteString("response_body", lw.responseData.body),
		)
	}
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	r.responseData.body = b
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
