package service

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/luoying_gh/distributed/registry"
)

// 使用 Gin 框架重写的服务启动函数
func Start(host, port string, reg registry.Registration, server *gin.Engine,
	registerHandlersFunc func()) error {
	// 调用该服务的请求处理函数
	registerHandlersFunc()

	var wg sync.WaitGroup
	wg.Add(2)

	// 启动该服务
	go startService(host, port, reg.ServiceName, server)

	// 给用户一个停止服务的方式
	sig := make(chan string)
	go func() {
		fmt.Println("Press any key to stop this sevice.")

		var s string
		fmt.Scanln(&s) // 一旦用户按下任何按键，停止当前服务

		sig <- s
	}()

	// 在启动服务后，将该服务注册进 registry
	err := registry.RegistryEntrance(reg)
	if nil != err {
		return err
	}
	defer registry.RegistryExit(reg.ServiceURL)

	s := <-sig
	if s != "" {
		return nil
	}
	wg.Wait()

	return nil
}

func startService(host, port string, serviceName registry.ServiceName, server *gin.Engine) {
	fmt.Println("Service", serviceName, "has been started")
	server.Run(fmt.Sprintf(":%s", port))
}
