package client

import (
	"api-gateway/config"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GRPCClients struct {
	UserConn    *grpc.ClientConn
	ProductConn *grpc.ClientConn
	OrderConn   *grpc.ClientConn
}

var Clients *GRPCClients

func InitClients() error {
	userConn, err := grpc.Dial(
		config.Global.Services.UserService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to user service: %w", err)
	}

	productConn, err := grpc.Dial(
		config.Global.Services.ProductService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to product service: %w", err)
	}

	orderConn, err := grpc.Dial(
		config.Global.Services.OrderService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to order service: %w", err)
	}

	Clients = &GRPCClients{
		UserConn:    userConn,
		ProductConn: productConn,
		OrderConn:   orderConn,
	}

	return nil
}

func (c *GRPCClients) Close() {
	if c.UserConn != nil {
		c.UserConn.Close()
	}
	if c.ProductConn != nil {
		c.ProductConn.Close()
	}
	if c.OrderConn != nil {
		c.OrderConn.Close()
	}
}

// CreateContextWithToken 创建带 Token 的 context
func CreateContextWithToken(token string) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if token != "" {
		md := metadata.Pairs("authorization", token)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	
	return ctx
}
