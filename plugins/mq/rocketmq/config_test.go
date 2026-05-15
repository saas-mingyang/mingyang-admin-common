package rocketmq

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestProducerAndConsumer(t *testing.T) {

	producerConf := &ProducerConf{
		NsResolver: []string{"192.168.201.58:9876"},
		GroupName:  "test_producer_group",
		MsgTimeOut: 3,
		Retry:      2,
	}

	producer := producerConf.MustNewProducer()

	consumerConf := ConsumerConf{
		NsResolver:       []string{"192.168.201.58:9876"},
		GroupName:        fmt.Sprintf("test_consumer_group_%d", time.Now().UnixMilli()),
		ConsumerModel:    "Clustering",
		ConsumeFromWhere: "last",
	}

	consumer, err := NewConsumer(consumerConf)
	require.Nil(t, err)

	msgCh := make(chan string, 10)

	consumer.RegisterHandler("mingyang_admin_camera_topic_order", "test_tag", func(ctx context.Context, msg *primitive.MessageExt) error {
		body := string(msg.Body)
		if body == "hello rocketmq" {
			msgCh <- body
		}
		return nil
	})

	err = consumer.Start()
	require.Nil(t, err)
	defer consumer.Stop()

	time.Sleep(2 * time.Second)

	sender := NewSender(producer)
	err = sender.SendSync(context.Background(), "mingyang_admin_camera_topic_order", []byte("hello rocketmq"),
		WithTag("test_tag"),
	)
	assert.Nil(t, err)

	select {
	case got := <-msgCh:
		assert.Equal(t, "hello rocketmq", got)
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestConf_Validate(t *testing.T) {
	c := &ProducerConf{
		NsResolver:                 []string{"192.168.201.58:9876"},
		GroupName:                  "",
		Namespace:                  "",
		InstanceName:               "",
		MsgTimeOut:                 0,
		DefaultTopicQueueNums:      0,
		CreateTopicKey:             "",
		CompressMsgBodyOverHowMuch: 0,
		CompressLevel:              0,
		Retry:                      0,
	}

	err := c.Validate()
	assert.Nil(t, err)
}

func TestSendSyncMessage(t *testing.T) {

	producerConf := &ProducerConf{
		NsResolver: []string{"192.168.201.58:9876"},
		GroupName:  "test_group",
	}
	producer := producerConf.MustNewProducer()

	err := SendSyncMessage(producer, context.Background(), "mingyang_admin_camera_topic_order", "test_tag", []byte("test123213313321313132313213212131"))
	assert.Nil(t, err)
}

func TestSendDelayMessage(t *testing.T) {
	t.Skip("RocketMQ server required")

	producerConf := &ProducerConf{
		NsResolver: []string{"192.168.201.58:9876"},
		GroupName:  "test_group",
	}
	producer := producerConf.MustNewProducer()

	err := SendDelayMessage(producer, context.Background(), "mingyang_admin_camera_topic_order", "test_tag", []byte("delay test"), 5)
	assert.Nil(t, err)
}

func TestConsumerWithOptions(t *testing.T) {

	consumerConf := ConsumerConf{
		NsResolver:    []string{"192.168.201.58:9876"},
		GroupName:     "test_group",
		ConsumerModel: "Clustering",
	}

	consumer, err := NewConsumer(consumerConf, WithConcurrency(10))
	require.Nil(t, err)

	consumer.RegisterHandler("mingyang_admin_camera_topic_order", "test_tag", func(ctx context.Context, msg *primitive.MessageExt) error {
		fmt.Printf("received: %s\n", string(msg.Body))
		return nil
	})

	err = consumer.Start()
	require.Nil(t, err)
	defer consumer.Stop()

	time.Sleep(1000 * time.Second)
}
