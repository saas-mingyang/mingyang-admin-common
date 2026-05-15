package config

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type RocketMqConf struct {
	NameServers []string `json:",env=ROCKETMQ_NAME_SERVERS"`
	GroupName   string   `json:",optional,env=ROCKETMQ_GROUP"`
	Namespace   string   `json:",optional,env=ROCKETMQ_NAMESPACE"`
	AccessKey   string   `json:",optional,env=ROCKETMQ_ACCESS_KEY"`
	SecretKey   string   `json:",optional,env=ROCKETMQ_SECRET_KEY"`
	SendTimeout int      `json:",optional,default=3,env=ROCKETMQ_SEND_TIMEOUT"`
	Retry       int      `json:",optional,default=2,env=ROCKETMQ_RETRY"`
}

func (c RocketMqConf) MustNewProducer() rocketmq.Producer {
	err := c.validate()
	logx.Must(err)

	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(c.NameServers)),
		producer.WithGroupName(c.GroupName),
		producer.WithNamespace(c.Namespace),
		producer.WithSendMsgTimeout(time.Duration(c.SendTimeout)*time.Second),
		producer.WithRetry(c.Retry),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		}),
	)
	logx.Must(err)

	err = p.Start()
	logx.Must(err)

	return p
}

func (c RocketMqConf) MustNewPushConsumer() rocketmq.PushConsumer {
	err := c.validate()
	logx.Must(err)

	cs, err := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver(c.NameServers)),
		consumer.WithGroupName(c.GroupName),
		consumer.WithNamespace(c.Namespace),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		}),
	)
	logx.Must(err)

	return cs
}

func (c RocketMqConf) validate() error {
	if len(c.NameServers) == 0 {
		logx.Error("RocketMqConf.NameServers must not be empty")
	}
	if c.GroupName == "" {
		logx.Error("RocketMqConf.GroupName must not be empty")
	}
	return nil
}
