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
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"mingyang.com/admin-simple-admin-file/internal/config"
	"mingyang.com/admin-simple-admin-file/internal/handler"
	"mingyang.com/admin-simple-admin-file/internal/svc"
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

	server.Start()
}
