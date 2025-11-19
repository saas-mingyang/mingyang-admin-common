//		Simple Admin File
//
//		This is simple admin file manager api doc
//
//		Schemes: http, https
//		Host: localhost:9102
//		BasePath: /
//		Version: 1.7.0
//		Contact: yuansu.china.work@gmail.com
//		securityDefinitions:
//		  Token:
//		    type: apiKey
//		    name: Authorization
//		    in: header
//		security:
//		  - Token: []
//	    Consumes:
//		  - application/json
//
//		Produces:
//		  - application/json
//
// swagger:meta
package main

import (
	"flag"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"time"

	"github.com/suyuan32/simple-admin-file-tenant/internal/config"
	"github.com/suyuan32/simple-admin-file-tenant/internal/handler"
	"github.com/suyuan32/simple-admin-file-tenant/internal/svc"

	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/fms.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors(c.CROSConf.Address))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	GetPrivateURL("t5y897q8i.hn-bkt.clouddn.com",
		"2025-11-19/image/019a9a1b-c0f5-7e95-afcc-10f6072c56ce.png",
		"u1xd-w1-ezqdgHp14f1DcXoRsD0kMO7pbi84pvPy",
		"RVr3jzN-rpAzls_qwdQmP9U9U_I_IBKklp1ZqCsw")

	server.Start()
}

func GetPrivateURL(domain, key, accessKey, secretKey string) string {
	// 创建 MAC
	mac := auth.New(accessKey, secretKey)

	// 过期时间（1小时）
	deadline := time.Now().Unix() + 3600

	// 生成私有下载链接
	privateURL := storage.MakePrivateURL(mac, domain, key, deadline)
	fmt.Printf("Private URL: %s", privateURL)
	return privateURL
}
