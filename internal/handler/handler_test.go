package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/spitfy/gofermart/internal/domain/user"
	"github.com/spitfy/gofermart/internal/middleware/auth"
	"github.com/stretchr/testify/assert"
)

func TestHandler_RegisterUser(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		expectedCode int
		body         string
		user         user.User
		times        int
	}{
		{
			name:         "success",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			body:         `{"login":"test", "password": "test"}`,
			user: user.User{
				Login:    "test",
				Password: "test",
			},
			times: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := user.NewMockServicer(ctrl)
			a := auth.NewMockServicer(ctrl)
			s := Service{UserService: m, Auth: a}
			h := newHandler(s)

			srv := httptest.NewServer(http.HandlerFunc(h.RegisterUser))

			a.EXPECT().CreateToken(gomock.Any(), gomock.Any()).
				Return("1", nil).Times(tt.times)
			m.EXPECT().RegisterUser(gomock.Any(), tt.user).
				Return(1, nil).Times(tt.times)

			req := resty.New().R()
			req.SetHeader("Content-Type", "application/json")

			if tt.body != "" {
				req.SetBody(tt.body)
			}

			resp, err := req.Execute(tt.method, srv.URL)

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code mismatch")
		})
	}
}
