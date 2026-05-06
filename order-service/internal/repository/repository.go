package repository

import (
	"time"

	"order-service/internal/model"
)

// OrderRepo 定义订单仓储接口，方便测试时 mock
type OrderRepo interface {
	Create(order *model.Order) error
	FindByID(id uint) (*model.Order, error)
	FindByOrderNo(orderNo string) (*model.Order, error)
	Update(order *model.Order) error
	UpdateStatus(id uint, status string) error
	ListByUserID(userID uint, page, pageSize int) ([]model.Order, int64, error)
	DeleteByID(id uint) error
	// FindExpiredPending 查找超时未支付的订单
	FindExpiredPending(before time.Time) ([]model.Order, error)
}
