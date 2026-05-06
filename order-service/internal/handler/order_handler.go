package handler

import (
	"context"
	"order-service/internal/logger"
	"order-service/internal/middleware"
	"order-service/internal/service"
	pb "order-service/proto"
	"time"

	"go.uber.org/zap"
)

// beijingLoc 北京时区 UTC+8
var beijingLoc = time.FixedZone("CST", 8*3600)

// formatBeijingTime 格式化为北京时间字符串
func formatBeijingTime(t time.Time) string {
	return t.In(beijingLoc).Format("2006-01-02 15:04:05")
}

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	traceID := middleware.GetTraceID(ctx)
	
	// 从 context 中获取 userID
	userID, ok := ctx.Value("userID").(uint)
	if !ok {
		logger.Error("Failed to get userID from context", zap.String("trace_id", traceID))
		return nil, ErrUnauthorized
	}

	logger.Info("Creating order",
		zap.String("trace_id", traceID),
		zap.Uint("user_id", userID),
		zap.Uint32("product_id", req.ProductId),
		zap.Int32("quantity", req.Quantity),
	)

	order, err := h.service.CreateOrder(ctx, userID, uint(req.ProductId), int(req.Quantity))
	if err != nil {
		logger.Error("Failed to create order",
			zap.String("trace_id", traceID),
			zap.Uint("user_id", userID),
			zap.Uint32("product_id", req.ProductId),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("Order created successfully",
		zap.String("trace_id", traceID),
		zap.Uint("user_id", userID),
		zap.Uint("order_id", order.ID),
		zap.String("order_no", order.OrderNo),
	)

	return &pb.CreateOrderResponse{
		Id:      uint32(order.ID),
		OrderNo: order.OrderNo,
		Message: "order created successfully",
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	userID, ok := ctx.Value("userID").(uint)
	if !ok {
		return nil, ErrUnauthorized
	}

	order, err := h.service.GetOrder(userID, uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.GetOrderResponse{
		Id:          uint32(order.ID),
		OrderNo:     order.OrderNo,
		UserId:      uint32(order.UserID),
		ProductId:   uint32(order.ProductID),
		ProductName: order.ProductName,
		Quantity:    int32(order.Quantity),
		TotalPrice:  order.TotalPrice,
		Status:      order.Status,
		CreatedAt:   formatBeijingTime(order.CreatedAt),
		UpdatedAt:   formatBeijingTime(order.UpdatedAt),
	}, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	userID, ok := ctx.Value("userID").(uint)
	if !ok {
		return nil, ErrUnauthorized
	}

	orders, total, err := h.service.ListOrders(userID, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.Order
	for _, o := range orders {
		pbOrders = append(pbOrders, &pb.Order{
			Id:          uint32(o.ID),
			OrderNo:     o.OrderNo,
			UserId:      uint32(o.UserID),
			ProductId:   uint32(o.ProductID),
			ProductName: o.ProductName,
			Quantity:    int32(o.Quantity),
			TotalPrice:  o.TotalPrice,
			Status:      o.Status,
			CreatedAt:   formatBeijingTime(o.CreatedAt),
			UpdatedAt:   formatBeijingTime(o.UpdatedAt),
		})
	}

	return &pb.ListOrdersResponse{
		Orders: pbOrders,
		Total:  int32(total),
	}, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	traceID := middleware.GetTraceID(ctx)
	
	userID, ok := ctx.Value("userID").(uint)
	if !ok {
		logger.Error("Failed to get userID from context", zap.String("trace_id", traceID))
		return nil, ErrUnauthorized
	}

	logger.Info("Cancelling order",
		zap.String("trace_id", traceID),
		zap.Uint("user_id", userID),
		zap.Uint32("order_id", req.Id),
	)

	err := h.service.CancelOrder(ctx, userID, uint(req.Id))
	if err != nil {
		logger.Error("Failed to cancel order",
			zap.String("trace_id", traceID),
			zap.Uint("user_id", userID),
			zap.Uint32("order_id", req.Id),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("Order cancelled successfully",
		zap.String("trace_id", traceID),
		zap.Uint("user_id", userID),
		zap.Uint32("order_id", req.Id),
	)

	return &pb.CancelOrderResponse{
		Message: "order cancelled successfully",
	}, nil
}

func (h *OrderHandler) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	userID, ok := ctx.Value("userID").(uint)
	if !ok {
		return nil, ErrUnauthorized
	}

	if err := h.service.DeleteOrder(userID, uint(req.Id)); err != nil {
		return nil, err
	}

	return &pb.DeleteOrderResponse{
		Message: "order deleted successfully",
	}, nil
}

// PayOrder 支付订单
func (h *OrderHandler) PayOrder(ctx context.Context, req *pb.PayOrderRequest) (*pb.PayOrderResponse, error) {
	traceID := middleware.GetTraceID(ctx)

	userID, ok := ctx.Value("userID").(uint)
	if !ok {
		logger.Error("Failed to get userID from context", zap.String("trace_id", traceID))
		return nil, ErrUnauthorized
	}

	logger.Info("Paying order",
		zap.String("trace_id", traceID),
		zap.Uint("user_id", userID),
		zap.Uint32("order_id", req.Id),
	)

	err := h.service.PayOrder(ctx, userID, uint(req.Id))
	if err != nil {
		logger.Error("Failed to pay order",
			zap.String("trace_id", traceID),
			zap.Uint("user_id", userID),
			zap.Uint32("order_id", req.Id),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("Order paid successfully",
		zap.String("trace_id", traceID),
		zap.Uint("user_id", userID),
		zap.Uint32("order_id", req.Id),
	)

	return &pb.PayOrderResponse{
		Message: "order paid successfully",
	}, nil
}

var ErrUnauthorized = &UnauthorizedError{}

type UnauthorizedError struct{}

func (e *UnauthorizedError) Error() string {
	return "unauthorized"
}
