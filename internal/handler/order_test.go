package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/spitfy/gofermart/internal/domain/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupOrderTest(t *testing.T) (orderServicer *order.MockServicer, handler *Handler) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	o := order.NewMockServicer(ctrl)
	s := Service{OrderService: o}
	return o, newHandler(s)
}

func TestHandler_CreateOrder(t *testing.T) {
	tests := []struct {
		name         string
		contentType  string
		body         string
		userID       interface{}
		expectedCode int
		mockSetup    func(o *order.MockServicer)
	}{
		{
			name:         "bad content type",
			contentType:  "application/json",
			body:         "12345678903",
			userID:       1,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "empty body",
			contentType:  "text/plain",
			body:         "",
			userID:       1,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "no userID in context",
			contentType:  "text/plain",
			body:         "12345678903",
			userID:       nil,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid Luhn",
			contentType:  "text/plain",
			body:         "12345678900",
			userID:       1,
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "conflict order another user",
			contentType:  "text/plain",
			body:         "79927398713", // valid luhn
			userID:       1,
			expectedCode: http.StatusConflict,
			mockSetup: func(o *order.MockServicer) {
				o.EXPECT().AddOrder(gomock.Any(), 1, "79927398713").Return(order.ErrOrderAnotherUser).Times(1)
			},
		},
		{
			name:         "order already exists",
			contentType:  "text/plain",
			body:         "79927398713",
			userID:       1,
			expectedCode: http.StatusOK,
			mockSetup: func(o *order.MockServicer) {
				o.EXPECT().AddOrder(gomock.Any(), 1, "79927398713").Return(order.ErrExistsOrder).Times(1)
			},
		},
		{
			name:         "internal server error",
			contentType:  "text/plain",
			body:         "79927398713",
			userID:       1,
			expectedCode: http.StatusInternalServerError,
			mockSetup: func(o *order.MockServicer) {
				o.EXPECT().AddOrder(gomock.Any(), 1, "79927398713").Return(errors.New("db error")).Times(1)
			},
		},
		{
			name:         "success add order",
			contentType:  "text/plain",
			body:         "79927398713",
			userID:       1,
			expectedCode: http.StatusAccepted,
			mockSetup: func(o *order.MockServicer) {
				o.EXPECT().AddOrder(gomock.Any(), 1, "79927398713").Return(nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, h := setupOrderTest(t)

			handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				if tt.userID != nil {
					ctx = context.WithValue(ctx, userIDKey, tt.userID)
				}
				r = r.WithContext(ctx)
				h.CreateOrder(w, r)
			})

			srv := httptest.NewServer(handlerFunc)
			defer srv.Close()

			if tt.mockSetup != nil {
				tt.mockSetup(o)
			}

			req := resty.New().R().SetHeader("Content-Type", tt.contentType).SetBody(tt.body)

			resp, err := req.Post(srv.URL)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode())
		})
	}
}

func TestHandler_ListOrders(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		userID       interface{}
		expectedCode int
		orders       []order.Order
		mockSetup    func(o *order.MockServicer)
	}{
		{
			name:         "unauthorized no userID",
			userID:       nil,
			expectedCode: http.StatusUnauthorized,
			mockSetup:    func(o *order.MockServicer) {},
		},
		{
			name:         "internal error from order service",
			userID:       1,
			expectedCode: http.StatusInternalServerError,
			mockSetup: func(o *order.MockServicer) {
				o.EXPECT().ListOrders(gomock.Any(), 1).Return(nil, errors.New("db error")).Times(1)
			},
		},
		{
			name:         "empty orders list",
			userID:       1,
			expectedCode: http.StatusNoContent,
			mockSetup: func(o *order.MockServicer) {
				o.EXPECT().ListOrders(gomock.Any(), 1).Return([]order.Order{}, nil).Times(1)
			},
		},
		{
			name:         "success list orders",
			userID:       1,
			expectedCode: http.StatusOK,
			orders: []order.Order{
				{
					Number:    "12345678903",
					Status:    order.StatusInvalid, // пример значения, подставьте ваше
					Accrual:   15.5,
					CreatedAt: now,
				},
				{
					Number:    "98765432109",
					Status:    order.StatusNew,
					Accrual:   0,
					CreatedAt: now.Add(-time.Hour),
				},
			},
			mockSetup: func(o *order.MockServicer) {
				o.EXPECT().ListOrders(gomock.Any(), 1).Return([]order.Order{
					{
						Number:    "12345678903",
						Status:    order.StatusInvalid,
						Accrual:   15.5,
						CreatedAt: now,
					},
					{
						Number:    "98765432109",
						Status:    order.StatusNew,
						Accrual:   0,
						CreatedAt: now.Add(-time.Hour),
					},
				}, nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, h := setupOrderTest(t)

			handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				if tt.userID != nil {
					ctx = context.WithValue(ctx, userIDKey, tt.userID)
				}
				r = r.WithContext(ctx)
				h.ListOrders(w, r)
			})

			srv := httptest.NewServer(handlerFunc)
			defer srv.Close()

			tt.mockSetup(o)

			resp, err := resty.New().R().Get(srv.URL)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode())

			if tt.expectedCode == http.StatusOK {
				var got []struct {
					Number     string  `json:"number"`
					Status     string  `json:"status"`
					Accrual    float64 `json:"accrual"`
					UploadedAt string  `json:"uploaded_at"`
				}
				err = json.Unmarshal(resp.Body(), &got)
				require.NoError(t, err)

				assert.Len(t, got, len(tt.orders))
				for i, o := range tt.orders {
					assert.Equal(t, o.Number, got[i].Number)
					assert.Equal(t, string(o.Status), got[i].Status)
					assert.Equal(t, o.Accrual, got[i].Accrual)
					assert.WithinDuration(t, o.CreatedAt, parseTime(got[i].UploadedAt), time.Second)
				}
			}
		})
	}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
