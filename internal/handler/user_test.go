package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/spitfy/gofermart/internal/domain/user"
	"github.com/spitfy/gofermart/internal/middleware/auth"
	"github.com/stretchr/testify/assert"
)

func setupUserTest(t *testing.T) (userServicer *user.MockServicer, authServicer *auth.MockServicer, handler *Handler) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	u := user.NewMockServicer(ctrl)
	a := auth.NewMockServicer(ctrl)
	s := Service{UserService: u, Auth: a}
	return u, a, newHandler(s)
}

func TestHandler_RegisterUser(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		expectedCode int
		body         string
		user         user.User
		mockSetup    func(u *user.MockServicer, a *auth.MockServicer)
	}{
		{
			name:         "success",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			body:         `{"login":"test", "password":"test"}`,
			user:         user.User{Login: "test", Password: "test"},
			mockSetup: func(u *user.MockServicer, a *auth.MockServicer) {
				u.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(1, nil).Times(1)
				a.EXPECT().CreateToken(gomock.Any(), 1).Return("token", nil).Times(1)
			},
		},
		{
			name:         "conflict user exists",
			method:       http.MethodPost,
			expectedCode: http.StatusConflict,
			body:         `{"login":"test", "password":"test"}`,
			user:         user.User{Login: "test", Password: "test"},
			mockSetup: func(u *user.MockServicer, a *auth.MockServicer) {
				u.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(0, user.ErrExistsUser).Times(1)
			},
		},
		{
			name:         "register user internal error",
			method:       http.MethodPost,
			expectedCode: http.StatusInternalServerError,
			body:         `{"login":"test", "password":"test"}`,
			user:         user.User{Login: "test", Password: "test"},
			mockSetup: func(u *user.MockServicer, a *auth.MockServicer) {
				u.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(0, errors.New("fail")).Times(1)
			},
		},
		{
			name:         "token creation error",
			method:       http.MethodPost,
			expectedCode: http.StatusInternalServerError,
			body:         `{"login":"test", "password":"test"}`,
			user:         user.User{Login: "test", Password: "test"},
			mockSetup: func(u *user.MockServicer, a *auth.MockServicer) {
				u.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(1, nil).Times(1)
				a.EXPECT().CreateToken(gomock.Any(), 1).Return("", errors.New("token error")).Times(1)
			},
		},
		{
			name:         "bad method returns bad request",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			body:         "",
			user:         user.User{},
			mockSetup:    func(u *user.MockServicer, a *auth.MockServicer) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, a, h := setupUserTest(t)

			srv := httptest.NewServer(http.HandlerFunc(h.RegisterUser))
			defer srv.Close()

			if tt.mockSetup != nil {
				tt.mockSetup(u, a)
			}

			req := resty.New().R().SetHeader("Content-Type", "application/json").SetBody(tt.body)

			resp, err := req.Execute(tt.method, srv.URL)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode())
		})
	}
}

func TestHandler_LoginUser(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		expectedCode int
		body         string
		user         user.User
		mockSetup    func(u *user.MockServicer, a *auth.MockServicer)
	}{
		{
			name:         "success login",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			body:         `{"login":"test", "password":"test"}`,
			user:         user.User{Login: "test", Password: "test"},
			mockSetup: func(u *user.MockServicer, a *auth.MockServicer) {
				u.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(1, nil).Times(1)
				a.EXPECT().CreateToken(gomock.Any(), 1).Return("token", nil).Times(1)
			},
		},
		{
			name:         "unauthorized error",
			method:       http.MethodPost,
			expectedCode: http.StatusUnauthorized,
			body:         `{"login":"wrong", "password":"wrong"}`,
			user:         user.User{Login: "wrong", Password: "wrong"},
			mockSetup: func(u *user.MockServicer, a *auth.MockServicer) {
				u.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(0, auth.ErrUnAuth).Times(1)
			},
		},
		{
			name:         "internal server error on login",
			method:       http.MethodPost,
			expectedCode: http.StatusInternalServerError,
			body:         `{"login":"test", "password":"test"}`,
			user:         user.User{Login: "test", Password: "test"},
			mockSetup: func(u *user.MockServicer, a *auth.MockServicer) {
				u.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(0, errors.New("db error")).Times(1)
			},
		},
		{
			name:         "internal server error on token creation",
			method:       http.MethodPost,
			expectedCode: http.StatusInternalServerError,
			body:         `{"login":"test", "password":"test"}`,
			user:         user.User{Login: "test", Password: "test"},
			mockSetup: func(u *user.MockServicer, a *auth.MockServicer) {
				u.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(1, nil).Times(1)
				a.EXPECT().CreateToken(gomock.Any(), 1).Return("", errors.New("token error")).Times(1)
			},
		},
		{
			name:         "bad method",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			body:         "",
			user:         user.User{},
			mockSetup: func(u *user.MockServicer, a *auth.MockServicer) {
			},
		},
		{
			name:         "invalid JSON body",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			body:         `{"login": "test", "password":`, // invalid JSON
			user:         user.User{},
			mockSetup: func(u *user.MockServicer, a *auth.MockServicer) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, a, h := setupUserTest(t)

			srv := httptest.NewServer(http.HandlerFunc(h.LoginUser))
			defer srv.Close()

			if tt.mockSetup != nil {
				tt.mockSetup(u, a)
			}

			req := resty.New().R().SetHeader("Content-Type", "application/json").SetBody(tt.body)

			resp, err := req.Execute(tt.method, srv.URL)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode())
		})
	}
}

func TestHandler_GetUserBalance(t *testing.T) {
	tests := []struct {
		name         string
		userID       interface{}
		expectedCode int
		balance      user.Balance
		mockSetup    func(u *user.MockServicer)
	}{
		{
			name:         "unauthorized no userID",
			userID:       nil,
			expectedCode: http.StatusUnauthorized,
			mockSetup:    func(u *user.MockServicer) {},
		},
		{
			name:         "internal error from balance",
			userID:       1,
			expectedCode: http.StatusInternalServerError,
			balance: user.Balance{
				Current:   0,
				Withdrawn: 0,
			},
			mockSetup: func(u *user.MockServicer) {
				u.EXPECT().Balance(gomock.Any(), 1).Return(user.Balance{
					Current:   0,
					Withdrawn: 0,
				}, errors.New("fail")).Times(1)
			},
		},
		{
			name:         "success get balance",
			userID:       1,
			expectedCode: http.StatusOK,
			balance: user.Balance{
				Current:   100.5,
				Withdrawn: 17.4,
			},
			mockSetup: func(u *user.MockServicer) {
				u.EXPECT().Balance(gomock.Any(), 1).Return(user.Balance{
					Current:   100.5,
					Withdrawn: 17.4,
				}, nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, _, h := setupUserTest(t)

			uMock := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				if tt.userID != nil {
					ctx = context.WithValue(ctx, userIDKey, tt.userID)
				}
				r = r.WithContext(ctx)
				h.GetUserBalance(w, r)
			})

			srv := httptest.NewServer(uMock)
			defer srv.Close()

			tt.mockSetup(u)

			resp, err := resty.New().R().Get(srv.URL)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode())

			if tt.expectedCode == http.StatusOK {
				var gotBalance user.Balance
				err = json.Unmarshal(resp.Body(), &gotBalance)
				require.NoError(t, err)
				assert.Equal(t, tt.balance, gotBalance)
			}
		})
	}
}
