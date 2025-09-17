package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/spitfy/gofermart/internal/domain/withdraw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupWithdrawalTest(t *testing.T) (userServicer *withdraw.MockServicer, handler *Handler) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wd := withdraw.NewMockServicer(ctrl)
	s := Service{WithdrawService: wd}
	return wd, newHandler(s)
}

func TestHandler_WithdrawBalance(t *testing.T) {
	validContentType := "application/json"
	invalidContentType := "text/plain"

	validBody := `{"order":"12345", "sum":100.5}`
	invalidJSONBody := `{"order":"12345", "sum":`

	tests := []struct {
		name           string
		userID         interface{}
		contentType    string
		body           string
		mockSetup      func(wd *withdraw.MockServicer)
		expectedStatus int
	}{
		{
			name:           "no userID unauthorized",
			userID:         nil,
			contentType:    validContentType,
			body:           validBody,
			mockSetup:      nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid content-type",
			userID:         1,
			contentType:    invalidContentType,
			body:           validBody,
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest, // validateContentType возвращает false и не устанавливает статус - проверим отсутствие вызова WriteHeader
		},
		{
			name:           "invalid JSON body",
			userID:         1,
			contentType:    validContentType,
			body:           invalidJSONBody,
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest, // decodeJSONBody возвращает true и функция завершается - не вызывается WriteHeader явно
		},
		{
			name:        "low balance error",
			userID:      1,
			contentType: validContentType,
			body:        validBody,
			mockSetup: func(wd *withdraw.MockServicer) {
				wd.EXPECT().Add(gomock.Any(), 1, gomock.Any()).Return(withdraw.ErrLowBalance).Times(1)
			},
			expectedStatus: http.StatusPaymentRequired,
		},
		{
			name:        "internal server error",
			userID:      1,
			contentType: validContentType,
			body:        validBody,
			mockSetup: func(wd *withdraw.MockServicer) {
				wd.EXPECT().Add(gomock.Any(), 1, gomock.Any()).Return(errors.New("db error")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "success",
			userID:      1,
			contentType: validContentType,
			body:        validBody,
			mockSetup: func(wd *withdraw.MockServicer) {
				wd.EXPECT().Add(gomock.Any(), 1, gomock.Any()).Return(nil).Times(1)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd, h := setupWithdrawalTest(t)

			handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				if tt.userID != nil {
					ctx = context.WithValue(ctx, userIDKey, tt.userID)
				}
				r = r.WithContext(ctx)
				h.WithdrawBalance(w, r)
			})

			srv := httptest.NewServer(handlerFunc)
			defer srv.Close()

			if tt.mockSetup != nil {
				tt.mockSetup(wd)
			}

			req := resty.New().R().SetBody(tt.body).SetHeader("Content-Type", tt.contentType)

			resp, err := req.Post(srv.URL)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode())
		})
	}
}
