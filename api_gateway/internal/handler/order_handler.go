package handler

import (
	"api-gateway/internal/client"
	"context"
	"net/http"
	"strconv"
	"time"

	orderPb "api-gateway/proto/order"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

// CreateOrder 创建订单
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, _ := c.Get("token")
	tokenStr, _ := token.(string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if tokenStr != "" {
		md := metadata.Pairs("authorization", tokenStr)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	orderClient := orderPb.NewOrderServiceClient(client.Clients.OrderConn)
	resp, err := orderClient.CreateOrder(ctx, &orderPb.CreateOrderRequest{
		ProductId: uint32(req.ProductID),
		Quantity:  int32(req.Quantity),
	})

	if err != nil {
		// 添加详细的错误日志
		println("CreateOrder error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       resp.Id,
		"order_no": resp.OrderNo,
	})
}

// GetOrder 获取订单详情
func (h *OrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	token, _ := c.Get("token")
	tokenStr, _ := token.(string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if tokenStr != "" {
		md := metadata.Pairs("authorization", tokenStr)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	orderClient := orderPb.NewOrderServiceClient(client.Clients.OrderConn)
	resp, err := orderClient.GetOrder(ctx, &orderPb.GetOrderRequest{
		Id: uint32(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           resp.Id,
		"order_no":     resp.OrderNo,
		"user_id":      resp.UserId,
		"product_id":   resp.ProductId,
		"product_name": resp.ProductName,
		"quantity":     resp.Quantity,
		"total_price":  resp.TotalPrice,
		"status":       resp.Status,
		"created_at":   resp.CreatedAt,
		"updated_at":   resp.UpdatedAt,
	})
}

// ListOrders 订单列表
func (h *OrderHandler) ListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	token, _ := c.Get("token")
	tokenStr, _ := token.(string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if tokenStr != "" {
		md := metadata.Pairs("authorization", tokenStr)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	orderClient := orderPb.NewOrderServiceClient(client.Clients.OrderConn)
	resp, err := orderClient.ListOrders(ctx, &orderPb.ListOrdersRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orders := make([]gin.H, 0, len(resp.Orders))
	for _, o := range resp.Orders {
		orders = append(orders, gin.H{
			"id":           o.Id,
			"order_no":     o.OrderNo,
			"product_name": o.ProductName,
			"quantity":     o.Quantity,
			"total_price":  o.TotalPrice,
			"status":       o.Status,
			"created_at":   o.CreatedAt,
			"updated_at":   o.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  resp.Total,
	})
}

// CancelOrder 取消订单
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	token, _ := c.Get("token")
	tokenStr, _ := token.(string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if tokenStr != "" {
		md := metadata.Pairs("authorization", tokenStr)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	orderClient := orderPb.NewOrderServiceClient(client.Clients.OrderConn)
	_, err = orderClient.CancelOrder(ctx, &orderPb.CancelOrderRequest{
		Id: uint32(id),
	})

	if err != nil {
		// 添加详细的错误日志
		println("CancelOrder error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order cancelled successfully",
	})
}

// DeleteOrder 删除订单
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	token, _ := c.Get("token")
	tokenStr, _ := token.(string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if tokenStr != "" {
		md := metadata.Pairs("authorization", tokenStr)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	orderClient := orderPb.NewOrderServiceClient(client.Clients.OrderConn)
	_, err = orderClient.DeleteOrder(ctx, &orderPb.DeleteOrderRequest{
		Id: uint32(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order deleted successfully",
	})
}

// PayOrder 支付订单
func (h *OrderHandler) PayOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	token, _ := c.Get("token")
	tokenStr, _ := token.(string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if tokenStr != "" {
		md := metadata.Pairs("authorization", tokenStr)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	orderClient := orderPb.NewOrderServiceClient(client.Clients.OrderConn)
	_, err = orderClient.PayOrder(ctx, &orderPb.PayOrderRequest{
		Id: uint32(id),
	})

	if err != nil {
		println("PayOrder error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order paid successfully",
	})
}
