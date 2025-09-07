package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInitialize(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name    string
		args    args
		want    zapcore.Level
		wantErr bool
	}{
		{"default", args{"info"}, zap.InfoLevel, false},
		{"wrong", args{"info22"}, zap.InfoLevel, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Initialize(tt.args.level)
			if (err != nil) && tt.wantErr {
				assert.Contains(t, err.Error(), "unrecognized level")
				return
			}
			if got.Log.Level() != tt.want {
				t.Errorf("Initialize() level = %v, want %v", got.Log.Level(), tt.want)
			}
		})
	}
}
