package rocketmq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/primitive"
)

// MessageConsumer 消费者接口,消费者需要实现
type MessageConsumer interface {
	Topic() string
	Tag() string
	Consume(ctx context.Context, msg *primitive.MessageExt) error
}
