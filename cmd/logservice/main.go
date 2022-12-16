package main

import (
	"fmt"
	stlog "log"

	"github.com/mothbuzzing/distributed/log"
	"github.com/mothbuzzing/distributed/registry"
	"github.com/mothbuzzing/distributed/service"
)

func main() {
	server := log.Run("./distributed.log")
	host, port := "localhost", "4000"
	reg := registry.Registration{
		ServiceName:      "LogService",
		ServiceURL:       fmt.Sprintf("http://%s:%s", host, port),
		RequiredServices: make([]registry.ServiceName, 0),
		UpdateURL:        fmt.Sprintf("http://%s:%s/services", host, port),
	}
	err := service.Start(host, port, reg, server, log.RegisterHandlers)
	if nil != err {
		stlog.Fatalln(err)
	}
}
