package main

import (
	"app/apps/im/ws/internal/config"
	"app/apps/im/ws/internal/handler"
	"app/apps/im/ws/internal/svc"
	server "app/apps/im/ws/websocket"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/local/im.yaml", "the config file")


// websocket 服务启动入口
func main() {
	// 加载配置文件
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	err := c.SetUp()
	if err != nil {
		panic(err)
	}

	// 装载服务上下文
	ctx := svc.NewServiceContext(c)

	// 创建一个 websocket 服务器实例
	srv := server.NewServer(c.ListenOn, server.WithAuthentication(handler.NewJwtAuto(ctx)))
	defer srv.Stop()

	// 注册路由
	handler.RegisterHandlers(srv, ctx)

	fmt.Printf("ws server starting at %v ...\n", c.ListenOn)
	srv.Start() // 启动服务器

}