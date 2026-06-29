package bootstrap

import (
	"fmt"
	"github.com/saas-mingyang/mingyang-admin-common/enum/common"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
)

type BootstrapConf struct {
	Type string `json:",default=yaml,options=[yaml,json,toml]"`
	Etcd discov.EtcdConf
}

func Load[T any](localFile, bootstrapFile string) (T, configurator.Configurator[T]) {
	var c T

	// 获取环境
	appEnv := strings.TrimSpace(os.Getenv("APP_ENV"))
	if appEnv == common.EmptyString {
		appEnv = "dev"
	}

	logx.Infof("APP_ENV=%s", appEnv)

	// 自动选择本地配置
	if localFile == common.EmptyString || strings.HasSuffix(localFile, "dev.yaml") {
		localFile = filepath.Join("etc", appEnv+".yaml")
	}

	logx.Infof("local config file: %s", localFile)

	// 仅本地配置
	if bootstrapFile == common.EmptyString {
		conf.MustLoad(localFile, &c, conf.UseEnv())
		return c, nil
	}

	var bc BootstrapConf
	conf.MustLoad(bootstrapFile, &bc, conf.UseEnv())

	// 根据环境自动切换 Etcd Key（可选）
	//
	// bootstrap.yaml 原来：
	// Key: mingyang/gateway/api.conf
	//
	// APP_ENV=test 后：
	// Key: mingyang/test/gateway/api.conf
	//
	if bc.Etcd.Key != common.EmptyString {
		bc.Etcd.Key = buildEtcdKey(bc.Etcd.Key, appEnv)
	}

	logx.Infof("etcd hosts: %v", bc.Etcd.Hosts)
	logx.Infof("etcd key: %s", bc.Etcd.Key)

	ss := subscriber.MustNewEtcdSubscriber(bc.Etcd)

	cc := configurator.MustNewConfigCenter[T](configurator.Config{
		Type: bc.Type,
		Log:  false,
	}, ss)

	v, err := cc.GetConfig()
	logx.Must(err)

	return v, cc
}

func buildEtcdKey(key, env string) string {
	key = strings.Trim(key, "/")

	// 已经包含环境
	if strings.HasPrefix(key, env+"/") {
		return key
	}

	// mingyang/gateway/api.conf
	parts := strings.Split(key, "/")
	if len(parts) >= 2 {
		return fmt.Sprintf("%s/%s/%s",
			parts[0],
			env,
			strings.Join(parts[1:], "/"))
	}

	return fmt.Sprintf("%s/%s", env, key)
}
