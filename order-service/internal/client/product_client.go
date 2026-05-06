package client

import (
	"context"
	"fmt"
	pb "order-service/proto"
	"order-service/internal/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type productClient struct {
	conn   *grpc.ClientConn
	client pb.ProductServiceClient
}

// NewProductClient 创建商品服务客户端
func NewProductClient(address string) (ProductClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	return &productClient{
		conn:   conn,
		client: pb.NewProductServiceClient(conn),
	}, nil
}

func (c *productClient) GetProduct(ctx context.Context, productID uint) (*model.ProductInfo, error) {
	resp, err := c.client.GetProduct(ctx, &pb.GetProductRequest{
		Id: uint32(productID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call product service: %w", err)
	}

	return &model.ProductInfo{
		ID:          uint(resp.Id),
		Name:        resp.Name,
		Description: resp.Description,
		Price:       resp.Price,
		Stock:       int(resp.Stock),
	}, nil
}

func (c *productClient) AdjustStock(ctx context.Context, productID uint, delta int) error {
	_, err := c.client.AdjustStock(ctx, &pb.AdjustStockRequest{
		Id:    uint32(productID),
		Delta: int32(delta),
	})
	if err != nil {
		return fmt.Errorf("failed to adjust product stock: %w", err)
	}
	return nil
}

func (c *productClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
