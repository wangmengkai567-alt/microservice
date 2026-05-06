package processor

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"order-service/internal/logger"
	"order-service/internal/queue"
	"order-service/internal/repository"
)

type OrderProcessor struct {
	orderRepo *repository.OrderRepository
	mq        *queue.RabbitMQ
}

func NewOrderProcessor(orderRepo *repository.OrderRepository, mq *queue.RabbitMQ) *OrderProcessor {
	return &OrderProcessor{
		orderRepo: orderRepo,
		mq:        mq,
	}
}

func (p *OrderProcessor) Start() error {
	return p.mq.ConsumeOrders(p.processOrder)
}

func (p *OrderProcessor) processOrder(msg queue.OrderMessage) error {
	ctx := context.Background()
	
	logger.Info("Processing order asynchronously",
		zap.Uint("order_id", msg.OrderID),
		zap.Uint("user_id", msg.UserID),
		zap.Uint("product_id", msg.ProductID),
	)

	// 模拟订单处理逻辑（支付、库存确认、通知等）
	// 在真实场景中，这里会调用支付服务、发送邮件通知等
	time.Sleep(2 * time.Second) // 模拟处理时间

	// 更新订单状态为已完成
	order, err := p.orderRepo.FindByID(msg.OrderID)
	if err != nil {
		return fmt.Errorf("failed to find order: %w", err)
	}

	// 这里可以添加更多业务逻辑
	// 例如：调用支付服务、发送通知、更新库存等
	
	order.Status = "completed"
	if err := p.orderRepo.Update(order); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	logger.Info("Order processed and completed",
		zap.Uint("order_id", msg.OrderID),
		zap.String("status", order.Status),
	)

	// 这里可以发送订单完成通知
	// 例如：发送邮件、短信、推送通知等

	_ = ctx // 使用 context 进行超时控制等

	return nil
}
