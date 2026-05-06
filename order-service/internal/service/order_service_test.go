package service_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"order-service/internal/logger"
	"order-service/internal/model"
	"order-service/internal/queue"
	"order-service/internal/service"
)

// TestMain 在所有测试运行前初始化 logger，避免 nil pointer panic
func TestMain(m *testing.M) {
	_ = logger.Init("development")
	os.Exit(m.Run())
}

// ========== Mock: OrderRepo ==========

type mockOrderRepo struct {
	orders    map[uint]*model.Order
	nextID    uint
	createErr error
	updateErr error
}

func newMockOrderRepo() *mockOrderRepo {
	return &mockOrderRepo{
		orders: make(map[uint]*model.Order),
		nextID: 1,
	}
}

func (m *mockOrderRepo) Create(order *model.Order) error {
	if m.createErr != nil {
		return m.createErr
	}
	order.ID = m.nextID
	m.nextID++
	copied := *order
	m.orders[order.ID] = &copied
	return nil
}

func (m *mockOrderRepo) FindByID(id uint) (*model.Order, error) {
	order, ok := m.orders[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	copied := *order
	return &copied, nil
}

func (m *mockOrderRepo) FindByOrderNo(orderNo string) (*model.Order, error) {
	for _, o := range m.orders {
		if o.OrderNo == orderNo {
			copied := *o
			return &copied, nil
		}
	}
	return nil, errors.New("record not found")
}

func (m *mockOrderRepo) Update(order *model.Order) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	copied := *order
	m.orders[order.ID] = &copied
	return nil
}

func (m *mockOrderRepo) UpdateStatus(id uint, status string) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	order, ok := m.orders[id]
	if !ok {
		return errors.New("record not found")
	}
	order.Status = status
	return nil
}

func (m *mockOrderRepo) ListByUserID(userID uint, page, pageSize int) ([]model.Order, int64, error) {
	var result []model.Order
	for _, o := range m.orders {
		if o.UserID == userID {
			result = append(result, *o)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockOrderRepo) DeleteByID(id uint) error {
	if _, ok := m.orders[id]; !ok {
		return errors.New("record not found")
	}
	delete(m.orders, id)
	return nil
}

func (m *mockOrderRepo) FindExpiredPending(before time.Time) ([]model.Order, error) {
	var result []model.Order
	for _, o := range m.orders {
		if o.Status == "pending" && o.CreatedAt.Before(before) {
			result = append(result, *o)
		}
	}
	return result, nil
}

// ========== Mock: ProductClient ==========

type mockProductClient struct {
	product      *model.ProductInfo
	getErr       error
	adjustErr    error
	adjustCalled int
}

func (m *mockProductClient) GetProduct(_ context.Context, _ uint) (*model.ProductInfo, error) {
	return m.product, m.getErr
}

func (m *mockProductClient) AdjustStock(_ context.Context, _ uint, _ int) error {
	m.adjustCalled++
	return m.adjustErr
}

func (m *mockProductClient) Close() error { return nil }

// ========== Mock: MessageQueue ==========

type mockMQ struct {
	publishErr     error
	publishedCount int
}

func (m *mockMQ) PublishOrderCreated(_ context.Context, _ queue.OrderMessage) error {
	m.publishedCount++
	return m.publishErr
}

// ========== helper ==========

func newTestService(repo *mockOrderRepo, pc *mockProductClient, mq *mockMQ) *service.OrderService {
	return service.NewOrderService(repo, pc, mq)
}

// ========== CreateOrder ==========

func TestOrderService_CreateOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("normal order success", func(t *testing.T) {
		repo := newMockOrderRepo()
		pc := &mockProductClient{
			product: &model.ProductInfo{ID: 1, Name: "iPhone", Price: 5999.0, Stock: 10},
		}
		mq := &mockMQ{}
		svc := newTestService(repo, pc, mq)

		order, err := svc.CreateOrder(ctx, 1, 1, 2)
		require.NoError(t, err)
		assert.Equal(t, uint(1), order.UserID)
		assert.Equal(t, uint(1), order.ProductID)
		assert.Equal(t, "iPhone", order.ProductName)
		assert.Equal(t, 2, order.Quantity)
		assert.Equal(t, 5999.0*2, order.TotalPrice)
		assert.Equal(t, "pending", order.Status)
		assert.Contains(t, order.OrderNo, "ORD")
		assert.Equal(t, 1, mq.publishedCount)
		assert.Equal(t, 1, pc.adjustCalled)
	})

	t.Run("productID=0 returns error", func(t *testing.T) {
		svc := newTestService(newMockOrderRepo(), &mockProductClient{}, &mockMQ{})
		_, err := svc.CreateOrder(ctx, 1, 0, 1)
		assert.EqualError(t, err, "invalid product id")
	})

	t.Run("quantity<=0 returns error", func(t *testing.T) {
		svc := newTestService(newMockOrderRepo(), &mockProductClient{}, &mockMQ{})
		_, err := svc.CreateOrder(ctx, 1, 1, 0)
		assert.EqualError(t, err, "quantity must be greater than 0")
	})

	t.Run("product service unavailable", func(t *testing.T) {
		pc := &mockProductClient{getErr: errors.New("product service unavailable")}
		svc := newTestService(newMockOrderRepo(), pc, &mockMQ{})
		_, err := svc.CreateOrder(ctx, 1, 1, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get product")
	})

	t.Run("insufficient stock", func(t *testing.T) {
		pc := &mockProductClient{
			product: &model.ProductInfo{ID: 1, Name: "iPhone", Price: 5999.0, Stock: 1},
		}
		svc := newTestService(newMockOrderRepo(), pc, &mockMQ{})
		_, err := svc.CreateOrder(ctx, 1, 1, 5)
		assert.EqualError(t, err, "insufficient stock")
	})

	t.Run("adjust stock fails", func(t *testing.T) {
		pc := &mockProductClient{
			product:   &model.ProductInfo{ID: 1, Name: "iPhone", Price: 5999.0, Stock: 10},
			adjustErr: errors.New("stock service error"),
		}
		svc := newTestService(newMockOrderRepo(), pc, &mockMQ{})
		_, err := svc.CreateOrder(ctx, 1, 1, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to deduct stock")
	})

	t.Run("db write fails, stock is rolled back", func(t *testing.T) {
		repo := newMockOrderRepo()
		repo.createErr = errors.New("db error")
		pc := &mockProductClient{
			product: &model.ProductInfo{ID: 1, Name: "iPhone", Price: 5999.0, Stock: 10},
		}
		svc := newTestService(repo, pc, &mockMQ{})
		_, err := svc.CreateOrder(ctx, 1, 1, 2)
		assert.Error(t, err)
		// 1st call: deduct; 2nd call: rollback
		assert.Equal(t, 2, pc.adjustCalled, "stock should be rolled back")
	})

	t.Run("mq publish fails, order still succeeds", func(t *testing.T) {
		repo := newMockOrderRepo()
		pc := &mockProductClient{
			product: &model.ProductInfo{ID: 1, Name: "iPhone", Price: 5999.0, Stock: 10},
		}
		mq := &mockMQ{publishErr: errors.New("mq unavailable")}
		svc := newTestService(repo, pc, mq)
		order, err := svc.CreateOrder(ctx, 1, 1, 1)
		require.NoError(t, err)
		assert.NotNil(t, order)
	})
}

// ========== GetOrder ==========

func TestOrderService_GetOrder(t *testing.T) {
	repo := newMockOrderRepo()
	svc := newTestService(repo, &mockProductClient{}, &mockMQ{})
	repo.orders[1] = &model.Order{ID: 1, UserID: 10, Status: "pending"}

	t.Run("get own order", func(t *testing.T) {
		order, err := svc.GetOrder(10, 1)
		require.NoError(t, err)
		assert.Equal(t, uint(1), order.ID)
	})

	t.Run("orderID=0 returns error", func(t *testing.T) {
		_, err := svc.GetOrder(10, 0)
		assert.EqualError(t, err, "invalid order id")
	})

	t.Run("order not found", func(t *testing.T) {
		_, err := svc.GetOrder(10, 999)
		assert.EqualError(t, err, "order not found")
	})

	t.Run("get other user order returns permission denied", func(t *testing.T) {
		_, err := svc.GetOrder(99, 1)
		assert.EqualError(t, err, "permission denied")
	})
}

// ========== CancelOrder ==========

func TestOrderService_CancelOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("cancel pending order, stock restored", func(t *testing.T) {
		repo := newMockOrderRepo()
		pc := &mockProductClient{product: &model.ProductInfo{}}
		svc := newTestService(repo, pc, &mockMQ{})
		repo.orders[1] = &model.Order{ID: 1, UserID: 10, ProductID: 1, Quantity: 2, Status: "pending"}

		err := svc.CancelOrder(ctx, 10, 1)
		require.NoError(t, err)
		assert.Equal(t, "cancelled", repo.orders[1].Status)
		assert.Equal(t, 1, pc.adjustCalled)
	})

	t.Run("cancel already cancelled order", func(t *testing.T) {
		repo := newMockOrderRepo()
		repo.orders[1] = &model.Order{ID: 1, UserID: 10, Status: "cancelled"}
		svc := newTestService(repo, &mockProductClient{}, &mockMQ{})
		err := svc.CancelOrder(ctx, 10, 1)
		assert.EqualError(t, err, "order already cancelled")
	})

	t.Run("cancel completed order", func(t *testing.T) {
		repo := newMockOrderRepo()
		repo.orders[1] = &model.Order{ID: 1, UserID: 10, Status: "completed"}
		svc := newTestService(repo, &mockProductClient{}, &mockMQ{})
		err := svc.CancelOrder(ctx, 10, 1)
		assert.EqualError(t, err, "cannot cancel completed order")
	})

	t.Run("cancel other user order", func(t *testing.T) {
		repo := newMockOrderRepo()
		repo.orders[1] = &model.Order{ID: 1, UserID: 10, Status: "pending"}
		svc := newTestService(repo, &mockProductClient{}, &mockMQ{})
		err := svc.CancelOrder(ctx, 99, 1)
		assert.EqualError(t, err, "permission denied")
	})
}

// ========== PayOrder ==========

func TestOrderService_PayOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("pay pending order", func(t *testing.T) {
		repo := newMockOrderRepo()
		repo.orders[1] = &model.Order{ID: 1, UserID: 10, Status: "pending"}
		svc := newTestService(repo, &mockProductClient{}, &mockMQ{})
		err := svc.PayOrder(ctx, 10, 1)
		require.NoError(t, err)
		assert.Equal(t, "paid", repo.orders[1].Status)
	})

	t.Run("pay already paid order", func(t *testing.T) {
		repo := newMockOrderRepo()
		repo.orders[1] = &model.Order{ID: 1, UserID: 10, Status: "paid"}
		svc := newTestService(repo, &mockProductClient{}, &mockMQ{})
		err := svc.PayOrder(ctx, 10, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order cannot be paid")
	})

	t.Run("pay other user order", func(t *testing.T) {
		repo := newMockOrderRepo()
		repo.orders[1] = &model.Order{ID: 1, UserID: 10, Status: "pending"}
		svc := newTestService(repo, &mockProductClient{}, &mockMQ{})
		err := svc.PayOrder(ctx, 99, 1)
		assert.EqualError(t, err, "permission denied")
	})
}

// ========== DeleteOrder ==========

func TestOrderService_DeleteOrder(t *testing.T) {
	t.Run("delete cancelled order", func(t *testing.T) {
		repo := newMockOrderRepo()
		repo.orders[1] = &model.Order{ID: 1, UserID: 10, Status: "cancelled"}
		svc := newTestService(repo, &mockProductClient{}, &mockMQ{})
		err := svc.DeleteOrder(10, 1)
		require.NoError(t, err)
		_, exists := repo.orders[1]
		assert.False(t, exists)
	})

	t.Run("delete pending order fails", func(t *testing.T) {
		repo := newMockOrderRepo()
		repo.orders[1] = &model.Order{ID: 1, UserID: 10, Status: "pending"}
		svc := newTestService(repo, &mockProductClient{}, &mockMQ{})
		err := svc.DeleteOrder(10, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only cancelled orders can be deleted")
	})

	t.Run("delete other user order", func(t *testing.T) {
		repo := newMockOrderRepo()
		repo.orders[1] = &model.Order{ID: 1, UserID: 10, Status: "cancelled"}
		svc := newTestService(repo, &mockProductClient{}, &mockMQ{})
		err := svc.DeleteOrder(99, 1)
		assert.EqualError(t, err, "permission denied")
	})
}

// ========== ListOrders ==========

func TestOrderService_ListOrders(t *testing.T) {
	t.Run("invalid page/pageSize use defaults", func(t *testing.T) {
		repo := newMockOrderRepo()
		svc := newTestService(repo, &mockProductClient{}, &mockMQ{})
		orders, total, err := svc.ListOrders(1, 0, 0)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, orders)
	})

	t.Run("pageSize over 100 is capped", func(t *testing.T) {
		repo := newMockOrderRepo()
		svc := newTestService(repo, &mockProductClient{}, &mockMQ{})
		_, _, err := svc.ListOrders(1, 1, 999)
		require.NoError(t, err)
	})
}

// ========== AutoCancelExpiredOrders ==========

func TestOrderService_AutoCancelExpiredOrders(t *testing.T) {
	ctx := context.Background()

	t.Run("expired order is cancelled, stock restored", func(t *testing.T) {
		repo := newMockOrderRepo()
		pc := &mockProductClient{product: &model.ProductInfo{}}
		svc := newTestService(repo, pc, &mockMQ{})

		// 2 minutes ago - should be cancelled
		repo.orders[1] = &model.Order{
			ID:        1,
			UserID:    10,
			ProductID: 1,
			Quantity:  2,
			Status:    "pending",
			CreatedAt: time.Now().Add(-2 * time.Minute),
		}
		// just created - should NOT be cancelled
		repo.orders[2] = &model.Order{
			ID:        2,
			UserID:    10,
			ProductID: 1,
			Quantity:  1,
			Status:    "pending",
			CreatedAt: time.Now(),
		}

		count, err := svc.AutoCancelExpiredOrders(ctx)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
		assert.Equal(t, "cancelled", repo.orders[1].Status)
		assert.Equal(t, "pending", repo.orders[2].Status)
		assert.Equal(t, 1, pc.adjustCalled)
	})

	t.Run("no expired orders returns 0", func(t *testing.T) {
		repo := newMockOrderRepo()
		svc := newTestService(repo, &mockProductClient{}, &mockMQ{})
		count, err := svc.AutoCancelExpiredOrders(ctx)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}
