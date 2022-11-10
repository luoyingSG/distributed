package registry

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

const ServerPort = ":3000"

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

	return nil
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
	RegistryService.GET("/registry", func(ctx *gin.Context) {
		body, _ := ctx.GetRawData()

		var data Registration
		err := json.Unmarshal(body, &data)
		if nil != err {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		}

		err = r.add(data)
		if nil != err {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		}
	})
}
