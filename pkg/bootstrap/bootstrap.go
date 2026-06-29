package bootstrap

import (
	"github.com/saas-mingyang/mingyang-admin-common/enum/common"
	"os"
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

	// 有 ETCD_HOSTS，使用配置中心
	if strings.TrimSpace(os.Getenv("ETCD_HOSTS")) != common.EmptyString {

		var bc BootstrapConf
		conf.MustLoad(bootstrapFile, &bc, conf.UseEnv())

		// 使用环境变量覆盖 Hosts
		bc.Etcd.Hosts = splitHosts(os.Getenv("ETCD_HOSTS"))

		// 可选：覆盖 Key
		if key := strings.TrimSpace(os.Getenv("ETCD_KEY")); key != common.EmptyString {
			bc.Etcd.Key = key
		}

		logx.Infof("use etcd config center")
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

	// 没有 ETCD_HOSTS，直接读取本地配置
	if localFile == common.EmptyString {
		localFile = "etc/dev.yaml"
	}

	logx.Infof("use local config: %s", localFile)

	conf.MustLoad(localFile, &c, conf.UseEnv())
	return c, nil
}

func splitHosts(hosts string) []string {
	var result []string
	for _, host := range strings.Split(hosts, ",") {
		host = strings.TrimSpace(host)
		if host != "" {
			result = append(result, host)
		}
	}
	return result
}
