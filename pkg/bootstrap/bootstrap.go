package bootstrap

import (
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
	if appEnv == "" {
		appEnv = "dev"
	}

	logx.Infof("APP_ENV=%s", appEnv)

	// 自动选择本地配置
	if localFile == "" || strings.HasSuffix(localFile, "dev.yaml") {
		localFile = filepath.Join("etc", appEnv+".yaml")
	}

	logx.Infof("local config file: %s", localFile)

	// 仅本地配置
	if bootstrapFile == "" {
		conf.MustLoad(localFile, &c, conf.UseEnv())
		return c, nil
	}

	var bc BootstrapConf
	conf.MustLoad(bootstrapFile, &bc, conf.UseEnv())

	// ==========================
	// 使用环境变量覆盖 Etcd Hosts
	// ==========================
	if hosts := strings.TrimSpace(os.Getenv("ETCD_HOSTS")); hosts != "" {
		var etcdHosts []string

		for _, host := range strings.Split(hosts, ",") {
			host = strings.TrimSpace(host)
			if host != "" {
				etcdHosts = append(etcdHosts, host)
			}
		}

		if len(etcdHosts) > 0 {
			bc.Etcd.Hosts = etcdHosts
		}
	}

	// 如果需要，也允许覆盖 Key
	if key := strings.TrimSpace(os.Getenv("ETCD_KEY")); key != "" {
		bc.Etcd.Key = key
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
