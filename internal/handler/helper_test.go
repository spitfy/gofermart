package handler

import (
	"net/http"
	"strings"
	"testing"

	"github.com/spitfy/gofermart/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

type customWriter struct {
	headers http.Header
	status  int
	body    []byte
}

func (w *customWriter) Header() http.Header  { return w.headers }
func (w *customWriter) WriteHeader(code int) { w.status = code }
func (w *customWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return len(b), nil
}

func newCustomWriter() *customWriter {
	return &customWriter{headers: make(http.Header)}
}

func Test_validateContentType(t *testing.T) {
	body := strings.NewReader(`{"foo":"bar"}`)
	req, err := http.NewRequest(http.MethodPost, "https://example.com/api", body)
	if err != nil {
		t.Errorf("error create request: %v", err)
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		ct   string
		want bool
	}{
		{"success", args{newCustomWriter(), req}, "application/json", true},
		{"fail", args{newCustomWriter(), req}, "application/octet-stream", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.r.Header.Set("Content-Type", tt.ct)
			if got := validateContentType(tt.args.w, tt.args.r); got != tt.want {
				t.Errorf("validateContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeJSONBody(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
		u    user.User
	}{
		{name: "success", body: `{"login":"user","password":"pass"}`, want: false, u: user.User{Login: "user", Password: "pass"}},
		{name: "fail", body: `{"Login":"user","PASSWORD":"pass"}`, want: false, u: user.User{Login: "", Password: ""}},
		{name: "empty", body: `{""}`, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			req, err := http.NewRequest(http.MethodPost, "https://example.com/api", body)
			if err != nil {
				t.Errorf("error create request: %v", err)
			}
			var u user.User
			if got := decodeJSONBody(newCustomWriter(), req, u); got != tt.want {
				t.Errorf("decodeJSONBody() = %v, want %v", got, tt.want)
				if !tt.want {
					assert.Equal(t, tt.u, u)
				}
			}
		})
	}
}
