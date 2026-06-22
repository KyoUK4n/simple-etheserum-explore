// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/KyoUK4n/etherscan/internal/config"
	"github.com/KyoUK4n/etherscan/internal/handler"
	"github.com/KyoUK4n/etherscan/internal/svc"
	_ "github.com/joho/godotenv/autoload" // 必须最先导入，以加载配置项

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/etherscan-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	server := rest.MustNewServer(c.RestConf, rest.WithCorsHeaders("Access-Control-Allow-Origin"))

	svcCtx := svc.NewServiceContext(c)
	// api routers
	handler.RegisterHandlers(server, svcCtx)
	// web routers
	handler.RegisterWebRouter(server)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	go server.Start()

	// 优雅退出
	gracefulStop(server, svcCtx)
}

func gracefulStop(server *rest.Server, svcCtx *svc.ServiceContext) {
	sysSig := make(chan os.Signal, 1)
	signal.Notify(sysSig, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-sysSig
		log.Printf("get a signal: %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Println("terminating...")
			svcCtx.Close()
			server.Stop()
			return
		case syscall.SIGHUP:
			log.Println("terminal unconnected")
			return
		default:
			log.Println("no match signal")
			return
		}
	}
}
