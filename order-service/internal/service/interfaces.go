package service

import (
	"context"

	"order-service/internal/model"
)

// ProductClient 商品服务客户端接口（service 包内部使用，不依赖 client 包）
// client.productClient 通过鸭子类型自动满足此接口
type ProductClient interface {
	GetProduct(ctx context.Context, productID uint) (*model.ProductInfo, error)
	AdjustStock(ctx context.Context, productID uint, delta int) error
	Close() error
}
