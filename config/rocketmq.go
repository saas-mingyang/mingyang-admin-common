package config

import (
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"mingyang.com/admin-common/plugins/mq/rocketmq"
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

func (c RocketMqConf) MustNewProducer() *rocketmq.Sender {
	if err := c.validate(); err != nil {
		logx.Must(err)
	}
	return rocketmq.MustNewSender(c.NameServers, c.GroupName, c.Namespace, c.AccessKey, c.SecretKey, c.SendTimeout, c.Retry)
}

func (c RocketMqConf) MustNewConsumer(opts ...rocketmq.ConsumerSetup) *rocketmq.Consumer {
	if err := c.validate(); err != nil {
		logx.Must(err)
	}
	rmq, err := rocketmq.NewConsumer(rocketmq.ConsumerConf{
		NsResolver: c.NameServers,
		GroupName:  c.GroupName,
		Namespace:  c.Namespace,
		AccessKey:  c.AccessKey,
		SecretKey:  c.SecretKey,
	}, opts...)
	logx.Must(err)
	return rmq
}

func (c RocketMqConf) validate() error {
	if len(c.NameServers) == 0 {
		return errors.New("RocketMqConf.NameServers must not be empty")
	}
	if c.GroupName == "" {
		return errors.New("RocketMqConf.GroupName must not be empty")
	}
	return nil
}
