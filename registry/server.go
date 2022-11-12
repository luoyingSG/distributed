package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	ServerPort = ":3000"
	ServerURL  = "http://localhost" + ServerPort + "/registry"
)

// 服务注册列表
type registry struct {
	registrations []Registration
	regMutex      *sync.Mutex
}

// 将注册好的服务添加到注册列表中
func (r *registry) add(reg Registration) error {
	r.regMutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.regMutex.Unlock()

	fmt.Printf("Service at URL %s registed\n", reg.ServiceURL)

	return nil
}

// 将服务移除注册列表
func (r *registry) remove(url string) error {
	r.regMutex.Lock()
	for i, registration := range r.registrations {
		if registration.ServiceURL == url {
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			fmt.Printf("Service at URL %s deregisted\n", url)
			return nil
		}
	}
	r.regMutex.Unlock()

	return fmt.Errorf("failed to find registed service with url %s", url)
}

// 建立一个包级的服务注册列表
var r = registry{
	registrations: make([]Registration, 0),
	regMutex:      new(sync.Mutex),
}

// 服务注册模块的 WEB 服务
var RegistryService *gin.Engine

func init() {
	RegistryService = gin.Default()

	RegistryService.POST("/registry", func(ctx *gin.Context) {
		body, _ := ctx.GetRawData()

		var data Registration
		err := json.Unmarshal(body, &data)
		if nil != err {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		}
		fmt.Println(data)

		err = r.add(data)
		if nil != err {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		}
	})

	RegistryService.DELETE("/registry", func(ctx *gin.Context) {
		payload, err := ioutil.ReadAll(ctx.Request.Body)
		if nil != err {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		}

		url := string(payload)
		log.Printf("Removing service at URL %s", url)
		err = r.remove(url)
		if nil != err {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		}
	})
}
