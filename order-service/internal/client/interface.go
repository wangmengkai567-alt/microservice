package client

import (
	"context"

	"order-service/internal/model"
)

// ProductClient 商品服务客户端接口（不依赖 proto，可被测试安全引用）
type ProductClient interface {
	GetProduct(ctx context.Context, productID uint) (*model.ProductInfo, error)
	AdjustStock(ctx context.Context, productID uint, delta int) error
	Close() error
}
