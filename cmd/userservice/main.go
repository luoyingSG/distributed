package main

import (
	"fmt"
	"log"

	"github.com/luoying_gh/distributed/registry"
	"github.com/luoying_gh/distributed/service"
	"github.com/luoying_gh/distributed/user"
)

func main() {
	server := user.Run()
	host, port := "localhost", "5678"
	reg := registry.Registration{
		ServiceName: "UserService",
		ServiceURL:  fmt.Sprintf("http://%s:%s", host, port),
	}

	err := service.Start(host, port, reg, server, user.UserHandlers)
	if nil != err {
		log.Fatalln(err)
	}
}
