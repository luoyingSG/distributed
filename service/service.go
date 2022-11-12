package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/luoying_gh/distributed/registry"
)

// 使用 Gin 框架重写的服务启动函数
func Start(host, port string, reg registry.Registration, server *gin.Engine,
	registerHandlersFunc func()) error {
	// 调用该服务的请求处理函数
	registerHandlersFunc()
	// 启动该服务
	startService(host, port, reg.ServiceName, server)

	// 在启动服务后，将该服务注册进 registry
	err := registry.RegistryEntrance(reg)
	if nil != err {
		return nil
	}

	return nil
}

func startService(host, port string, serviceName registry.ServiceName, server *gin.Engine) {
	fmt.Println("Service", serviceName, "has been started")
	server.Run(fmt.Sprintf(":%s", port))
}
