package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kimvlry/simple-order-service/internal/domain"
	mock_interfaces "github.com/kimvlry/simple-order-service/internal/mocks"
)

func TestOrderHandler_GetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_interfaces.NewMockOrderService(ctrl)
	handler := NewOrderHandler(mockService)

	tests := []struct {
		name           string
		url            string
		setupMocks     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful order retrieval",
			url:  "/order?id=test-order-1",
			setupMocks: func() {
				expectedOrder := &domain.Order{
					OrderUid:    "test-order-1",
					TrackNumber: "TRACK123",
				}
				mockService.EXPECT().
					GetOrderByID("test-order-1", gomock.Any()).
					Return(expectedOrder, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"order_uid":"test-order-1","track_number":"TRACK123"`,
		},
		{
			name: "missing id parameter",
			url:  "/order",
			setupMocks: func() {
				// No mock expectations
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing id parameter",
		},
		{
			name: "order not found",
			url:  "/order?id=non-existent",
			setupMocks: func() {
				mockService.EXPECT().
					GetOrderByID("non-existent", gomock.Any()).
					Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			name: "internal server error",
			url:  "/order?id=error-order",
			setupMocks: func() {
				mockService.EXPECT().
					GetOrderByID("error-order", gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal error: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()

			handler.GetOrder(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" {
				body := w.Body.String()
				if tt.name == "successful order retrieval" {
					// For JSON response, just check if it contains the expected part
					if !contains(body, tt.expectedBody) {
						t.Errorf("expected body to contain %s, got %s", tt.expectedBody, body)
					}
				} else {
					// Remove trailing newline for comparison
					expectedBody := tt.expectedBody
					if body == expectedBody+"\n" {
						body = expectedBody
					}
					if body != expectedBody {
						t.Errorf("expected body %s, got %s", expectedBody, body)
					}
				}
			}
		})
	}
}

func TestOrderHandler_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_interfaces.NewMockOrderService(ctrl)
	handler := NewOrderHandler(mockService)

	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	handler.Ping(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	expectedBody := "hello from simple order-service\n"
	if w.Body.String() != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
