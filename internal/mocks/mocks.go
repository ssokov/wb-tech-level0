package mock_interfaces

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/kimvlry/simple-order-service/internal/domain"
)

type MockOrderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRepositoryMockRecorder
}

type MockOrderRepositoryMockRecorder struct {
	mock *MockOrderRepository
}

func NewMockOrderRepository(ctrl *gomock.Controller) *MockOrderRepository {
	mock := &MockOrderRepository{ctrl: ctrl}
	mock.recorder = &MockOrderRepositoryMockRecorder{mock}
	return mock
}

func (m *MockOrderRepository) EXPECT() *MockOrderRepositoryMockRecorder {
	return m.recorder
}

func (m *MockOrderRepository) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockOrderRepositoryMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockOrderRepository)(nil).Close))
}

func (m *MockOrderRepository) GetAll(ctx context.Context) ([]domain.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]domain.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockOrderRepositoryMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockOrderRepository)(nil).GetAll), ctx)
}

func (m *MockOrderRepository) GetById(ctx context.Context, id string) (*domain.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, id)
	ret0, _ := ret[0].(*domain.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockOrderRepositoryMockRecorder) GetById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockOrderRepository)(nil).GetById), ctx, id)
}

func (m *MockOrderRepository) Save(ctx context.Context, order *domain.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockOrderRepositoryMockRecorder) Save(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockOrderRepository)(nil).Save), ctx, order)
}

type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

type MockCacheMockRecorder struct {
	mock *MockCache
}

func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

func (m *MockCache) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockCacheMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCache)(nil).Close))
}

func (m *MockCache) GetOrder(ctx context.Context, orderUid string) (*domain.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrder", ctx, orderUid)
	ret0, _ := ret[0].(*domain.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockCacheMockRecorder) GetOrder(ctx, orderUid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrder", reflect.TypeOf((*MockCache)(nil).GetOrder), ctx, orderUid)
}

func (m *MockCache) RestoreCache(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestoreCache", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockCacheMockRecorder) RestoreCache(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreCache", reflect.TypeOf((*MockCache)(nil).RestoreCache), ctx)
}

func (m *MockCache) SaveOrder(ctx context.Context, order *domain.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockCacheMockRecorder) SaveOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveOrder", reflect.TypeOf((*MockCache)(nil).SaveOrder), ctx, order)
}

type MockMessageConsumer struct {
	ctrl     *gomock.Controller
	recorder *MockMessageConsumerMockRecorder
}

type MockMessageConsumerMockRecorder struct {
	mock *MockMessageConsumer
}

func NewMockMessageConsumer(ctrl *gomock.Controller) *MockMessageConsumer {
	mock := &MockMessageConsumer{ctrl: ctrl}
	mock.recorder = &MockMessageConsumerMockRecorder{mock}
	return mock
}

func (m *MockMessageConsumer) EXPECT() *MockMessageConsumerMockRecorder {
	return m.recorder
}

func (m *MockMessageConsumer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMessageConsumerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockMessageConsumer)(nil).Close))
}

func (m *MockMessageConsumer) Consume(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Consume", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMessageConsumerMockRecorder) Consume(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consume", reflect.TypeOf((*MockMessageConsumer)(nil).Consume), ctx)
}

type MockOrderService struct {
	ctrl     *gomock.Controller
	recorder *MockOrderServiceMockRecorder
}

type MockOrderServiceMockRecorder struct {
	mock *MockOrderService
}

func NewMockOrderService(ctrl *gomock.Controller) *MockOrderService {
	mock := &MockOrderService{ctrl: ctrl}
	mock.recorder = &MockOrderServiceMockRecorder{mock}
	return mock
}

func (m *MockOrderService) EXPECT() *MockOrderServiceMockRecorder {
	return m.recorder
}

func (m *MockOrderService) GetOrderByID(id string, ctx context.Context) (*domain.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByID", id, ctx)
	ret0, _ := ret[0].(*domain.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockOrderServiceMockRecorder) GetOrderByID(id, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByID", reflect.TypeOf((*MockOrderService)(nil).GetOrderByID), id, ctx)
}
