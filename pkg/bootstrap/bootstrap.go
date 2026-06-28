package bootstrap

import (
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
)

// BootstrapConf is the minimal local configuration that points to the etcd
// config center. The etcd address itself cannot live inside etcd, so this is
// the only piece that must stay on local disk (or be supplied via env).
type BootstrapConf struct {
	// Type is the format of the config stored in etcd: yaml, json or toml.
	Type string `json:",default=yaml,options=[yaml,json,toml]"`
	// Etcd holds the config-center coordinates (Hosts + Key).
	Etcd discov.EtcdConf
}

// Load loads the application config of type T.
//
// When bootstrapFile is empty it falls back to loading localFile directly via
// conf.MustLoad — i.e. the original behaviour, so existing deployments that do
// not pass -cc keep working unchanged.
//
// When bootstrapFile is provided it reads the etcd coordinates from it and
// pulls the full config from the etcd config center. The returned Configurator
// is non-nil in that case so the caller can register hot-reload listeners; it
// is nil in the local-file path.
func Load[T any](localFile, bootstrapFile string) (T, configurator.Configurator[T]) {
	var c T

	if bootstrapFile == "" {
		conf.MustLoad(localFile, &c, conf.UseEnv())
		return c, nil
	}

	var bc BootstrapConf
	conf.MustLoad(bootstrapFile, &bc, conf.UseEnv())

	ss := subscriber.MustNewEtcdSubscriber(bc.Etcd)
	cc := configurator.MustNewConfigCenter[T](configurator.Config{
		Type: bc.Type,
		Log:  false,
	}, ss)

	v, err := cc.GetConfig()
	logx.Must(err)

	return v, cc
}
