package logger

import (
	"net/http"

	"go.uber.org/zap"
)

type Mock struct {
	Log *zap.Logger
}

func InitMock() *Mock {
	lvl, _ := zap.ParseAtomicLevel("info")
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, _ := cfg.Build()

	return &Mock{Log: zl}
}

func (l *Mock) LogInfo(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	}
}
