package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kimvlry/simple-order-service/internal/domain"
	mock_interfaces "github.com/kimvlry/simple-order-service/internal/mocks"
)

func TestOrderService_GetOrderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_interfaces.NewMockOrderRepository(ctrl)
	mockCache := mock_interfaces.NewMockCache(ctrl)

	service := NewOrderService(mockRepo, mockCache)

	tests := []struct {
		name          string
		orderID       string
		setupMocks    func()
		expectedOrder *domain.Order
		expectedError error
	}{
		{
			name:    "order found in cache",
			orderID: "test-order-1",
			setupMocks: func() {
				expectedOrder := &domain.Order{
					OrderUid:    "test-order-1",
					TrackNumber: "TRACK123",
				}
				mockCache.EXPECT().
					GetOrder(gomock.Any(), "test-order-1").
					Return(expectedOrder, nil)
			},
			expectedOrder: &domain.Order{
				OrderUid:    "test-order-1",
				TrackNumber: "TRACK123",
			},
			expectedError: nil,
		},
		{
			name:    "order not in cache, found in repo",
			orderID: "test-order-2",
			setupMocks: func() {
				expectedOrder := &domain.Order{
					OrderUid:    "test-order-2",
					TrackNumber: "TRACK456",
				}
				mockCache.EXPECT().
					GetOrder(gomock.Any(), "test-order-2").
					Return(nil, nil)
				mockRepo.EXPECT().
					GetById(gomock.Any(), "test-order-2").
					Return(expectedOrder, nil)
				mockCache.EXPECT().
					SaveOrder(gomock.Any(), expectedOrder).
					Return(nil)
			},
			expectedOrder: &domain.Order{
				OrderUid:    "test-order-2",
				TrackNumber: "TRACK456",
			},
			expectedError: nil,
		},
		{
			name:    "order not found anywhere",
			orderID: "test-order-3",
			setupMocks: func() {
				mockCache.EXPECT().
					GetOrder(gomock.Any(), "test-order-3").
					Return(nil, nil)
				mockRepo.EXPECT().
					GetById(gomock.Any(), "test-order-3").
					Return(nil, nil)
			},
			expectedOrder: nil,
			expectedError: nil,
		},
		{
			name:    "cache error",
			orderID: "test-order-4",
			setupMocks: func() {
				mockCache.EXPECT().
					GetOrder(gomock.Any(), "test-order-4").
					Return(nil, errors.New("cache error"))
			},
			expectedOrder: nil,
			expectedError: errors.New("cache error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			order, err := service.GetOrderByID(tt.orderID, context.Background())

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.expectedOrder == nil {
				if order != nil {
					t.Errorf("expected nil order, got %v", order)
				}
				return
			}

			if order == nil {
				t.Errorf("expected order %v, got nil", tt.expectedOrder)
				return
			}

			if order.OrderUid != tt.expectedOrder.OrderUid {
				t.Errorf("expected order ID %s, got %s", tt.expectedOrder.OrderUid, order.OrderUid)
			}
		})
	}
}
