package logger

import (
	"go.uber.org/zap"
	"net/http"
)

type LoggerMock struct {
	Log *zap.Logger
}

func InitMock() *LoggerMock {
	lvl, _ := zap.ParseAtomicLevel("info")
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, _ := cfg.Build()

	return &LoggerMock{Log: zl}
}

func (l *LoggerMock) LogInfo(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	}
}
