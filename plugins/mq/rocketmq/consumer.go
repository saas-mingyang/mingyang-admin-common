package rocketmq

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/semaphore"
)

type MessageHandler func(ctx context.Context, msg *primitive.MessageExt) error

type subscribeEntry struct {
	topic     string
	tag       string
	handler   MessageHandler
	isOrderly bool
}

type Consumer struct {
	pushConsumer rocketmq.PushConsumer
	entries      []subscribeEntry
	concurrency  int64
	started      bool
	mu           sync.Mutex
}

type ConsumerOption func(*Consumer)

func WithConcurrency(n int64) ConsumerOption {
	return func(c *Consumer) {
		if n > 0 {
			c.concurrency = n
		}
	}
}

func NewConsumer(conf ConsumerConf, opts ...ConsumerOption) (*Consumer, error) {
	rlog.SetLogLevel("warn")

	model := consumer.Clustering
	if strings.EqualFold(conf.ConsumerModel, "BroadCasting") {
		model = consumer.BroadCasting
	}

	resolver := conf.NsResolver
	if len(resolver) == 0 {
		return nil, fmt.Errorf("NameServer address must not be empty")
	}

	groupName := conf.GroupName
	if groupName == "" {
		groupName = "DEFAULT_CONSUMER"
	}

	pushConsumer, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(groupName),
		consumer.WithNameServer(resolver),
		consumer.WithConsumerModel(model),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: conf.AccessKey,
			SecretKey: conf.SecretKey,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("create push consumer failed: %w", err)
	}

	c := &Consumer{
		pushConsumer: pushConsumer,
		concurrency:  20,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

func (c *Consumer) RegisterHandler(topic, tag string, handler MessageHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.started {
		logx.Errorf("cannot register handler after consumer started, topic=%s, tag=%s", topic, tag)
		return
	}

	isOrderly := strings.HasSuffix(topic, "_sort")

	c.entries = append(c.entries, subscribeEntry{
		topic:     topic,
		tag:       tag,
		handler:   handler,
		isOrderly: isOrderly,
	})
	logx.Infof("registered rocketmq handler: topic=%s, tag=%s, orderly=%v", topic, tag, isOrderly)
}

func (c *Consumer) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.started {
		return fmt.Errorf("consumer already started")
	}

	if len(c.entries) == 0 {
		logx.Info("no handlers registered, consumer will not subscribe to any topic")
		return nil
	}

	sem := semaphore.NewWeighted(c.concurrency)

	for _, entry := range c.entries {
		entry := entry
		var consumeFunc func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)

		if entry.isOrderly {
			consumeFunc = func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
				for _, msg := range msgs {
					if err := entry.handler(ctx, msg); err != nil {
						logx.Errorf("orderly message handle failed, topic=%s, tag=%s, msgId=%s, error=%v",
							entry.topic, entry.tag, msg.MsgId, err)
						return consumer.ConsumeRetryLater, nil
					}
				}
				return consumer.ConsumeSuccess, nil
			}
		} else {
			consumeFunc = func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
				for _, msg := range msgs {
					m := msg
					if err := sem.Acquire(ctx, 1); err != nil {
						logx.Error("acquire semaphore failed: %v", err)
						return consumer.ConsumeRetryLater, nil
					}
					go func() {
						defer sem.Release(1)
						if err := entry.handler(ctx, m); err != nil {
							logx.Errorf("concurrent message handle failed, topic=%s, tag=%s, msgId=%s, error=%v",
								entry.topic, entry.tag, m.MsgId, err)
						}
					}()
				}
				return consumer.ConsumeSuccess, nil
			}
		}

		if err := c.pushConsumer.Subscribe(
			entry.topic,
			consumer.MessageSelector{Type: consumer.TAG, Expression: entry.tag},
			consumeFunc,
		); err != nil {
			return fmt.Errorf("subscribe failed, topic=%s, tag=%s: %w", entry.topic, entry.tag, err)
		}
		logx.Infof("subscribed: topic=%s, tag=%s, orderly=%v", entry.topic, entry.tag, entry.isOrderly)
	}

	if err := c.pushConsumer.Start(); err != nil {
		return fmt.Errorf("start push consumer failed: %w", err)
	}

	c.started = true
	logx.Info("rocketmq consumer started successfully")
	return nil
}

func (c *Consumer) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.started {
		return nil
	}

	if err := c.pushConsumer.Shutdown(); err != nil {
		return fmt.Errorf("shutdown consumer failed: %w", err)
	}
	c.started = false
	logx.Info("rocketmq consumer stopped")
	return nil
}
