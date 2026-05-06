package scheduler

import (
	"context"
	"time"

	"go.uber.org/zap"
	"order-service/internal/logger"
	"order-service/internal/service"
)

// OrderScheduler 订单定时任务调度器
type OrderScheduler struct {
	orderService *service.OrderService
	interval     time.Duration
	stopChan     chan struct{}
}

// NewOrderScheduler 创建订单调度器
func NewOrderScheduler(orderService *service.OrderService, interval time.Duration) *OrderScheduler {
	return &OrderScheduler{
		orderService: orderService,
		interval:     interval,
		stopChan:     make(chan struct{}),
	}
}

// Start 启动定时任务
func (s *OrderScheduler) Start() {
	ticker := time.NewTicker(s.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.runAutoCancel()
			case <-s.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
	logger.Info("Order scheduler started", zap.Duration("interval", s.interval))
}

// Stop 停止定时任务
func (s *OrderScheduler) Stop() {
	close(s.stopChan)
	logger.Info("Order scheduler stopped")
}

func (s *OrderScheduler) runAutoCancel() {
	ctx := context.Background()
	count, err := s.orderService.AutoCancelExpiredOrders(ctx)
	if err != nil {
		logger.Error("Auto cancel expired orders failed", zap.Error(err))
		return
	}
	if count > 0 {
		logger.Info("Auto cancelled expired orders", zap.Int("count", count))
	}
}
