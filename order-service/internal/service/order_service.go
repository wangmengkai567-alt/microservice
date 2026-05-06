package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"order-service/internal/logger"
	"order-service/internal/model"
	"order-service/internal/pkg"
	"order-service/internal/queue"
	"order-service/internal/repository"
)

type OrderService struct {
	repo          repository.OrderRepo
	productClient ProductClient
	mq            queue.MessageQueue
}

func NewOrderService(repo repository.OrderRepo, productClient ProductClient, mq queue.MessageQueue) *OrderService {
	return &OrderService{
		repo:          repo,
		productClient: productClient,
		mq:            mq,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID, productID uint, quantity int) (*model.Order, error) {
	// 1. 验证输入
	if productID == 0 {
		return nil, errors.New("invalid product id")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	// 2. 调用商品服务获取商品信息
	product, err := s.productClient.GetProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// 3. 检查库存
	if product.Stock < quantity {
		return nil, errors.New("insufficient stock")
	}

	// 4. 计算总价
	totalPrice := product.Price * float64(quantity)

	// 5. 生成订单号
	orderNo := pkg.GenerateOrderNo()

	// 6. 创建订单
	order := &model.Order{
		OrderNo:     orderNo,
		UserID:      userID,
		ProductID:   productID,
		ProductName: product.Name,
		Quantity:    quantity,
		TotalPrice:  totalPrice,
		Status:      "pending",
	}

	// 先扣减库存，确保下单后库存实时变化
	if err := s.productClient.AdjustStock(ctx, productID, -quantity); err != nil {
		return nil, fmt.Errorf("failed to deduct stock: %w", err)
	}

	if err := s.repo.Create(order); err != nil {
		// 订单落库失败时尝试回滚库存（尽力而为）
		_ = s.productClient.AdjustStock(ctx, productID, quantity)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// 发送订单创建消息到队列，异步处理
	orderMsg := queue.OrderMessage{
		OrderID:    order.ID,
		UserID:     order.UserID,
		ProductID:  order.ProductID,
		Quantity:   order.Quantity,
		TotalPrice: order.TotalPrice,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}
	
	if err := s.mq.PublishOrderCreated(ctx, orderMsg); err != nil {
		logger.Error("Failed to publish order message",
			zap.Uint("order_id", order.ID),
			zap.Error(err),
		)
		// 消息发送失败不影响订单创建，只记录日志
	}

	return order, nil
}

func (s *OrderService) GetOrder(userID, orderID uint) (*model.Order, error) {
	if orderID == 0 {
		return nil, errors.New("invalid order id")
	}

	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// 验证订单所属用户
	if order.UserID != userID {
		return nil, errors.New("permission denied")
	}

	return order, nil
}

func (s *OrderService) ListOrders(userID uint, page, pageSize int) ([]model.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.ListByUserID(userID, page, pageSize)
}

func (s *OrderService) CancelOrder(ctx context.Context, userID, orderID uint) error {
	if orderID == 0 {
		return errors.New("invalid order id")
	}

	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	// 验证订单所属用户
	if order.UserID != userID {
		return errors.New("permission denied")
	}

	// 检查订单状态
	if order.Status == "cancelled" {
		return errors.New("order already cancelled")
	}
	if order.Status == "completed" {
		return errors.New("cannot cancel completed order")
	}

	// 恢复库存
	if err := s.productClient.AdjustStock(ctx, order.ProductID, order.Quantity); err != nil {
		return fmt.Errorf("failed to restore stock: %w", err)
	}

	// 使用专门的 UpdateStatus 方法更新订单状态
	if err := s.repo.UpdateStatus(orderID, "cancelled"); err != nil {
		// 如果状态更新失败，尝试回滚库存（尽力而为）
		_ = s.productClient.AdjustStock(ctx, order.ProductID, -order.Quantity)
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// TODO: 发送订单取消事件/消息

	return nil
}

func (s *OrderService) DeleteOrder(userID, orderID uint) error {
	if orderID == 0 {
		return errors.New("invalid order id")
	}

	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.UserID != userID {
		return errors.New("permission denied")
	}

	// 仅允许删除已取消订单，避免误删有效订单记录
	// 使用 strings.TrimSpace 去除可能的空格
	if order.Status != "cancelled" {
		return fmt.Errorf("only cancelled orders can be deleted, current status: %s", order.Status)
	}

	if err := s.repo.DeleteByID(orderID); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}

// PayOrder 支付订单
func (s *OrderService) PayOrder(ctx context.Context, userID, orderID uint) error {
	if orderID == 0 {
		return errors.New("invalid order id")
	}

	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	// 验证订单所属用户
	if order.UserID != userID {
		return errors.New("permission denied")
	}

	// 检查订单状态
	if order.Status != "pending" {
		return fmt.Errorf("order cannot be paid, current status: %s", order.Status)
	}

	// 更新订单状态为已支付
	if err := s.repo.UpdateStatus(orderID, "paid"); err != nil {
		return fmt.Errorf("failed to pay order: %w", err)
	}

	logger.Info("Order paid successfully",
		zap.Uint("order_id", orderID),
		zap.Uint("user_id", userID),
	)

	return nil
}

// AutoCancelExpiredOrders 自动取消超时未支付的订单
func (s *OrderService) AutoCancelExpiredOrders(ctx context.Context) (int, error) {
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)

	orders, err := s.repo.FindExpiredPending(oneMinuteAgo)
	if err != nil {
		return 0, fmt.Errorf("failed to find expired orders: %w", err)
	}

	cancelledCount := 0
	for _, order := range orders {
		// 恢复库存
		if err := s.productClient.AdjustStock(ctx, order.ProductID, order.Quantity); err != nil {
			logger.Error("Failed to restore stock for expired order",
				zap.Uint("order_id", order.ID),
				zap.Uint("product_id", order.ProductID),
				zap.Error(err),
			)
			continue
		}

		// 更新订单状态为已取消
		if err := s.repo.UpdateStatus(order.ID, "cancelled"); err != nil {
			logger.Error("Failed to cancel expired order",
				zap.Uint("order_id", order.ID),
				zap.Error(err),
			)
			// 回滚库存
			_ = s.productClient.AdjustStock(ctx, order.ProductID, -order.Quantity)
			continue
		}

		cancelledCount++
		logger.Info("Auto cancelled expired order",
			zap.Uint("order_id", order.ID),
			zap.String("order_no", order.OrderNo),
		)
	}

	return cancelledCount, nil
}
