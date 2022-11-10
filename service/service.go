package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 使用 Gin 框架重写的服务启动函数
func Start(host, port, serviceName string, server *gin.Engine,
	registerHandlersFunc func()) error {
	// 调用该服务的请求处理函数
	registerHandlersFunc()
	// 启动该服务
	startService(host, port, serviceName, server)
	return nil
}

func startService(host, port, serviceName string, server *gin.Engine) {
	fmt.Println("Service", serviceName, "has been started")
	server.Run(fmt.Sprintf(":%s", port))
}
