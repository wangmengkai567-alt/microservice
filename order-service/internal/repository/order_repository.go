package repository

import (
	"time"

	"gorm.io/gorm"
	"order-service/internal/model"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// DB 返回数据库连接，供 service 层使用
func (r *OrderRepository) DB() *gorm.DB {
	return r.db
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Where("id = ?", id).First(&order).Error
	return &order, err
}

func (r *OrderRepository) FindByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := r.db.Where("order_no = ?", orderNo).First(&order).Error
	return &order, err
}

func (r *OrderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

// UpdateStatus 专门用于更新订单状态
func (r *OrderRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) ListByUserID(userID uint, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error

	return orders, total, err
}

func (r *OrderRepository) DeleteByID(id uint) error {
	return r.db.Delete(&model.Order{}, id).Error
}

// FindExpiredPending 查找超时未支付的订单（实现 OrderRepo 接口）
func (r *OrderRepository) FindExpiredPending(before time.Time) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Where("status = ? AND created_at < ?", "pending", before).Find(&orders).Error
	return orders, err
}
