package main

import (
	"fmt"
	stlog "log"

	"github.com/luoying_gh/distributed/log"
	"github.com/luoying_gh/distributed/registry"
	"github.com/luoying_gh/distributed/service"
)

func main() {
	server := log.Run("./distributed.log")
	host, port := "localhost", "4000"
	reg := registry.Registration{
		ServiceName: "LogService",
		ServiceURL:  fmt.Sprintf("http://%s:%s", host, port),
	}
	err := service.Start(host, port, reg, server, log.RegisterHandlers)
	if nil != err {
		stlog.Fatalln(err)
	}
}
