package main

import "github.com/mothbuzzing/distributed/registry"

func main() {
	registry.RegistryService.Run(registry.ServerPort)
}
