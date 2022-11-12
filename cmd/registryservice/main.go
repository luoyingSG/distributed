package main

import "github.com/luoying_gh/distributed/registry"

func main() {
	registry.RegistryService.Run(registry.ServerPort)
}
