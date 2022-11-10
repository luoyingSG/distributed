package main

import (
	stlog "log"

	"github.com/luoying_gh/distributed/log"
	"github.com/luoying_gh/distributed/service"
)

func main() {
	server := log.Run("./distributed.log")
	host, port := "localhost", "4000"
	err := service.Start(host, port, "LogService", server, log.RegisterHandlers)
	if nil != err {
		stlog.Fatalln(err)
	}
}
