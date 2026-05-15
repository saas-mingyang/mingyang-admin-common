package rocketmq

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/zeromicro/go-zero/core/logx"
)

type MessageOption func(*primitive.Message)

type Sender struct {
	producer rocketmq.Producer
}

func NewSender(p rocketmq.Producer) *Sender {
	return &Sender{producer: p}
}

func MustNewSender(nameServers []string, groupName, namespace, accessKey, secretKey string, sendTimeout, retry int) *Sender {
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(nameServers)),
		producer.WithGroupName(groupName),
		producer.WithNamespace(namespace),
		producer.WithSendMsgTimeout(time.Duration(sendTimeout)*time.Second),
		producer.WithRetry(retry),
		producer.WithCredentials(primitive.Credentials{AccessKey: accessKey, SecretKey: secretKey}),
	)
	logx.Must(err)
	logx.Must(p.Start())
	return &Sender{producer: p}
}

func WithTag(tag string) MessageOption {
	return func(msg *primitive.Message) {
		msg.WithTag(tag)
	}
}

func WithKeys(keys []string) MessageOption {
	return func(msg *primitive.Message) {
		msg.WithKeys(keys)
	}
}

func WithShardingKey(key string) MessageOption {
	return func(msg *primitive.Message) {
		msg.WithShardingKey(key)
	}
}

func WithDelayTimeLevel(level int) MessageOption {
	return func(msg *primitive.Message) {
		msg.WithDelayTimeLevel(level)
	}
}

func WithProperties(properties map[string]string) MessageOption {
	return func(msg *primitive.Message) {
		for k, v := range properties {
			msg.WithProperty(k, v)
		}
	}
}

func (s *Sender) SendSync(ctx context.Context, topic string, body []byte, opts ...MessageOption) error {
	if s.producer == nil {
		return fmt.Errorf("producer is nil")
	}
	msg := &primitive.Message{Topic: topic, Body: body}
	for _, opt := range opts {
		opt(msg)
	}
	res, err := s.producer.SendSync(ctx, msg)
	if err != nil {
		return fmt.Errorf("send message failed: %w", err)
	}
	if res.Status != primitive.SendOK {
		return fmt.Errorf("send status abnormal: %v", res.Status)
	}
	return nil
}

func (s *Sender) SendAsync(ctx context.Context, topic string, body []byte, callback func(context.Context, *primitive.SendResult, error), opts ...MessageOption) error {
	if s.producer == nil {
		return fmt.Errorf("producer is nil")
	}
	msg := &primitive.Message{Topic: topic, Body: body}
	for _, opt := range opts {
		opt(msg)
	}
	return s.producer.SendAsync(ctx, callback, msg)
}

func (s *Sender) SendOneWay(ctx context.Context, topic string, body []byte, opts ...MessageOption) error {
	if s.producer == nil {
		return fmt.Errorf("producer is nil")
	}
	msg := &primitive.Message{Topic: topic, Body: body}
	for _, opt := range opts {
		opt(msg)
	}
	return s.producer.SendOneWay(ctx, msg)
}

func (s *Sender) SendWithDelay(ctx context.Context, topic string, body []byte, delayLevel int, opts ...MessageOption) error {
	if delayLevel < 1 || delayLevel > 18 {
		return fmt.Errorf("delay level must be 1-18")
	}
	opts = append(opts, WithDelayTimeLevel(delayLevel))
	return s.SendSync(ctx, topic, body, opts...)
}

func (s *Sender) Close() error {
	if s.producer != nil {
		return s.producer.Shutdown()
	}
	return nil
}

func SendSyncMessage(producer rocketmq.Producer, ctx context.Context, topic, tag string, body []byte) error {
	msg := &primitive.Message{
		Topic: topic,
		Body:  body,
	}
	msg.WithTag(tag)

	res, err := producer.SendSync(ctx, msg)
	if err != nil {
		logx.Errorf("send rocketmq message failed, topic=%s, tag=%s, error=%v", topic, tag, err)
		return fmt.Errorf("send rocketmq message failed: %w", err)
	}
	if res.Status != primitive.SendOK {
		logx.Errorf("send rocketmq message status abnormal, topic=%s, tag=%s, status=%v", topic, tag, res.Status)
		return fmt.Errorf("send rocketmq message status abnormal: %v", res.Status)
	}
	logx.Infof("send rocketmq message success, topic=%s, tag=%s, msgId=%s", topic, tag, res.MsgID)
	return nil
}

func SendDelayMessage(producer rocketmq.Producer, ctx context.Context, topic, tag string, body []byte, delaySecond int) error {
	var delayLevel int
	switch {
	case delaySecond <= 0:
		return SendSyncMessage(producer, ctx, topic, tag, body)
	case delaySecond <= 5:
		delayLevel = 1
	case delaySecond <= 10:
		delayLevel = 2
	case delaySecond <= 30:
		delayLevel = 3
	case delaySecond <= 60:
		delayLevel = 4
	case delaySecond <= 120:
		delayLevel = 5
	case delaySecond <= 180:
		delayLevel = 6
	case delaySecond <= 240:
		delayLevel = 7
	case delaySecond <= 300:
		delayLevel = 8
	case delaySecond <= 600:
		delayLevel = 9
	case delaySecond <= 1800:
		delayLevel = 10
	case delaySecond <= 3600:
		delayLevel = 11
	default:
		delayLevel = 12
	}

	return sendWithDelayLevel(producer, ctx, topic, tag, body, delayLevel)
}

func sendWithDelayLevel(producer rocketmq.Producer, ctx context.Context, topic, tag string, body []byte, delayLevel int) error {
	msg := &primitive.Message{
		Topic: topic,
		Body:  body,
	}
	msg.WithTag(tag)
	msg.WithDelayTimeLevel(delayLevel)

	res, err := producer.SendSync(ctx, msg)
	if err != nil {
		logx.Errorf("send delay rocketmq message failed, topic=%s, tag=%s, delayLevel=%d, error=%v", topic, tag, delayLevel, err)
		return fmt.Errorf("send rocketmq message failed: %w", err)
	}
	if res.Status != primitive.SendOK {
		logx.Errorf("send delay rocketmq message status abnormal, topic=%s, tag=%s, status=%v", topic, tag, res.Status)
		return fmt.Errorf("send rocketmq message status abnormal: %v", res.Status)
	}
	logx.Infof("send delay rocketmq message success, topic=%s, tag=%s, delayLevel=%d, msgId=%s", topic, tag, delayLevel, res.MsgID)
	return nil
}
