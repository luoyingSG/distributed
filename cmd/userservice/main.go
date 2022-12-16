package main

import (
	"fmt"
	stdlog "log"

	"github.com/mothbuzzing/distributed/log"
	"github.com/mothbuzzing/distributed/registry"
	"github.com/mothbuzzing/distributed/service"
	"github.com/mothbuzzing/distributed/user"
)

func main() {
	server := user.Run()
	host, port := "localhost", "5678"
	reg := registry.Registration{
		ServiceName:      "UserService",
		ServiceURL:       fmt.Sprintf("http://%s:%s", host, port),
		RequiredServices: []registry.ServiceName{registry.LogService},
		UpdateURL:        fmt.Sprintf("http://%s:%s/services", host, port),
	}

	err := service.Start(host, port, reg, server, user.UserHandlers)
	if nil != err {
		stdlog.Fatalln(err)
	}

	// 服务启动之后，开始使用 LogService
	if logProvider, err := registry.GetProvider(registry.LogService); nil != err {
		stdlog.Printf("failed to find required log service: %s", err.Error())
	} else {
		fmt.Printf("Log service found at %s", logProvider)
		log.SetClientLogger(logProvider, registry.LogService)
	}
}
