package repository

import (
	"product-service/internal/model"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	return &product, err
}

func (r *ProductRepository) Update(product *model.Product) error {
	// 使用 Save 方法，它会更新所有字段包括零值
	return r.db.Save(product).Error
}

// UpdateStock 专门用于更新库存，避免零值问题
func (r *ProductRepository) UpdateStock(id uint, stock int) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).Update("stock", stock).Error
}

// AdjustStockAtomic 使用原子操作调整库存，避免并发问题
func (r *ProductRepository) AdjustStockAtomic(id uint, delta int) error {
	// 使用 SQL 原子操作：UPDATE products SET stock = stock + delta WHERE id = ?
	result := r.db.Model(&model.Product{}).
		Where("id = ? AND stock + ? >= 0", id, delta).
		UpdateColumn("stock", gorm.Expr("stock + ?", delta))
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	return nil
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *ProductRepository) List(page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(pageSize).Find(&products).Error
	return products, total, err
}
