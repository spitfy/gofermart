package handler

import (
	"net/http"
	"testing"
)

func TestHandler_RegisterUser(t *testing.T) {
	type fields struct {
		s Service
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			h := &Handler{
				s: tt.fields.s,
			}
			h.RegisterUser(tt.args.w, tt.args.r)
		})
	}
}
