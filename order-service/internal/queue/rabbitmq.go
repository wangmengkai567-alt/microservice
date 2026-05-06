package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"order-service/internal/logger"
)

const (
	OrderQueue    = "order_processing"
	OrderExchange = "orders"
	OrderRoutingKey = "order.created"
)

type OrderMessage struct {
	OrderID   uint    `json:"order_id"`
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	CreatedAt string  `json:"created_at"`
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// 声明交换机
	err = channel.ExchangeDeclare(
		OrderExchange, // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// 声明队列
	_, err = channel.QueueDeclare(
		OrderQueue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// 绑定队列到交换机
	err = channel.QueueBind(
		OrderQueue,      // queue name
		OrderRoutingKey, // routing key
		OrderExchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	logger.Info("RabbitMQ connected successfully")

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		url:     url,
	}, nil
}

func (r *RabbitMQ) PublishOrderCreated(ctx context.Context, msg OrderMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = r.channel.PublishWithContext(
		ctx,
		OrderExchange,   // exchange
		OrderRoutingKey, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	logger.Info("Order message published",
		zap.Uint("order_id", msg.OrderID),
		zap.Uint("user_id", msg.UserID),
	)

	return nil
}

func (r *RabbitMQ) ConsumeOrders(handler func(OrderMessage) error) error {
	msgs, err := r.channel.Consume(
		OrderQueue, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	logger.Info("Starting to consume order messages")

	go func() {
		for msg := range msgs {
			var orderMsg OrderMessage
			if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
				logger.Error("Failed to unmarshal message", zap.Error(err))
				msg.Nack(false, false) // 拒绝消息，不重新入队
				continue
			}

			logger.Info("Processing order message",
				zap.Uint("order_id", orderMsg.OrderID),
			)

			if err := handler(orderMsg); err != nil {
				logger.Error("Failed to process order",
					zap.Uint("order_id", orderMsg.OrderID),
					zap.Error(err),
				)
				// 处理失败，重新入队（最多重试3次）
				msg.Nack(false, true)
			} else {
				logger.Info("Order processed successfully",
					zap.Uint("order_id", orderMsg.OrderID),
				)
				msg.Ack(false)
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
