package queue

import "context"

// MessageQueue 定义消息队列接口，方便测试时 mock
type MessageQueue interface {
	PublishOrderCreated(ctx context.Context, msg OrderMessage) error
}
