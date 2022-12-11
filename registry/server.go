package registry

import (
	"bytes"
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
	regMutex      *sync.RWMutex
}

// 将注册好的服务添加到注册列表中
func (r *registry) add(reg Registration) error {
	r.regMutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.regMutex.Unlock()

	err := r.sendRequiredServices(reg)
	if nil != err {
		log.Printf("failed to send required services to service at url %s: %s", reg.UpdateURL, err.Error())
	}

	r.notify(patch{
		Added: []patchEntry{
			{
				Name: reg.ServiceName,
				URL:  reg.ServiceURL,
			},
		},
	})

	fmt.Printf("Service at URL %s registed\n", reg.ServiceURL)

	return nil
}

// 将服务的依赖发送给它
func (r registry) sendRequiredServices(reg Registration) error {
	r.regMutex.RLock()
	defer r.regMutex.RUnlock()

	var p patch
	// 循环已经注册的服务，如果找到当前服务所依赖的服务，就添加到 patch 里
	for _, r := range r.registrations {
		for _, req := range reg.RequiredServices {
			if r.ServiceName == req {
				p.Added = append(p.Added, patchEntry{
					Name: r.ServiceName,
					URL:  r.ServiceURL,
				})
			}
		}
	}

	err := r.sendPatch(p, reg.UpdateURL) // 将找到的依赖的信息发送给服务
	if nil != err {
		return err
	}
	return nil
}

// 发送依赖信息
func (r registry) sendPatch(p patch, url string) error {
	data, err := json.Marshal(p)
	if nil != err {
		return err
	}

	// 以 post 请求的形式发送依赖信息
	_, err = http.Post(url, "application/json", bytes.NewBuffer(data))
	if nil != err {
		return err
	}

	return nil
}

// 将服务移除注册列表
func (r *registry) remove(url string) error {
	r.regMutex.Lock()
	defer r.regMutex.Unlock()

	for i, registration := range r.registrations {
		if registration.ServiceURL == url {
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			fmt.Printf("Service at URL %s deregisted\n", url)
			return nil
		}
	}

	return fmt.Errorf("failed to find registed service with url %s", url)
}

// 建立一个包级的服务注册列表
var r = registry{
	registrations: make([]Registration, 0),
	regMutex:      new(sync.RWMutex),
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

// 将依赖的变化通知到需要该依赖的服务
func (r *registry) notify(fullPatch patch) {
	r.regMutex.RLock()
	defer r.regMutex.RUnlock()

	// 遍历已经注册的服务
	for _, reg := range r.registrations {
		go func(reg Registration) { // 将已经注册的服务传入到 go routine 里
			// 在 go routine 里，针对传入的服务，看看它的依赖是否需要更新
			for _, reqService := range reg.RequiredServices {
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				sendUpdate := false // 是否需要更新

				// 看看是否有新的依赖可用
				for _, added := range fullPatch.Added {
					if added.Name == reqService {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}
				}

				// 看看是否有依赖被删除了
				for _, removed := range fullPatch.Removed {
					if removed.Name == reqService {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}

				if sendUpdate {
					err := r.sendPatch(p, reg.UpdateURL)
					if err != nil {
						log.Fatal(err)
						return
					}
				}
			}
		}(reg)
	}
}
