package bootstrap

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"mingyang.com/admin-common/enum/common"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/logx"
)

type BootstrapConf struct {
	Type string `json:",default=yaml,options=[yaml,json,toml]"`
	Etcd discov.EtcdConf
}

type Bootstrap[T any] struct {
	Config       T
	Configurator configurator.Configurator[T]
}

func (b *Bootstrap[T]) OnChange(fn func()) {
	if b.Configurator != nil {
		b.Configurator.AddListener(fn)
	}
}
func Load[T any](configDir string) *Bootstrap[T] {
	var cfg T

	bc := BootstrapConf{}

	bootstrapFile := filepath.Join(configDir, "bootstrap.yaml")

	env := strings.TrimSpace(os.Getenv(common.APP_ENV))

	if env == common.EmptyString {
		env = service.DevMode
	}

	localFile := filepath.Join(configDir, env+".yaml")

	// 读取 bootstrap
	if _, err := os.Stat(bootstrapFile); err == nil {
		conf.MustLoad(bootstrapFile, &bc, conf.UseEnv())
	}

	// 第一优先：环境变量
	if hosts := strings.TrimSpace(os.Getenv(common.ETCD_HOSTS)); hosts != common.EmptyString {

		bc.Etcd.Hosts = splitHosts(hosts)

		if key := strings.TrimSpace(os.Getenv(common.ETCD_CONF_KEY)); key != common.EmptyString {
			bc.Etcd.Key = key
		}

		cfg, cc := loadFromEtcd[T](bc)

		return &Bootstrap[T]{
			Config:       cfg,
			Configurator: cc,
		}
	}

	// 第二优先：bootstrap
	if len(bc.Etcd.Hosts) > common.Zero {

		cfg, cc := loadFromEtcd[T](bc)

		return &Bootstrap[T]{
			Config:       cfg,
			Configurator: cc,
		}
	}

	// 本地配置
	logx.Infof("Config Source : Local")
	logx.Infof("Config File   : %s", localFile)

	conf.MustLoad(localFile, &cfg, conf.UseEnv())

	return &Bootstrap[T]{
		Config: cfg,
	}
}

func splitHosts(hosts string) []string {
	var result []string
	for _, host := range strings.Split(hosts, common.Comma) {
		host = strings.TrimSpace(host)
		if host != common.EmptyString {
			result = append(result, host)
		}
	}
	return result
}

func loadFromEtcd[T any](bc BootstrapConf) (T, configurator.Configurator[T]) {
	// 默认配置类型
	if strings.TrimSpace(bc.Type) == common.EmptyString {
		bc.Type = "yaml"
	}

	// 参数校验
	if len(bc.Etcd.Hosts) == common.Zero {
		logx.Must(fmt.Errorf("etcd hosts is empty"))
	}

	if strings.TrimSpace(bc.Etcd.Key) == common.EmptyString {
		logx.Must(fmt.Errorf("etcd key is empty"))
	}

	logx.Infof("==================================================")
	logx.Infof("Config Source : ETCD")
	logx.Infof("Config Type   : %s", bc.Type)
	logx.Infof("Etcd Hosts    : %v", bc.Etcd.Hosts)
	logx.Infof("Etcd Key      : %s", bc.Etcd.Key)
	logx.Infof("==================================================")

	ss := subscriber.MustNewEtcdSubscriber(bc.Etcd)

	cc := configurator.MustNewConfigCenter[T](
		configurator.Config{
			Type: bc.Type,
			Log:  false,
		},
		ss,
	)

	cfg, err := cc.GetConfig()
	logx.Must(err)

	return cfg, cc
}

func Health() rest.Route {
	return rest.Route{
		Method: http.MethodGet,
		Path:   "/health",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			httpx.OkJson(w, map[string]string{"status": "ok"})
		},
	}
}
