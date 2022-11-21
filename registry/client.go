package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// 注册服务的客户端的入口
// 其他服务调用 RegistryEntrance 函数，就能自动将自己注册进 registry
func RegistryEntrance(r Registration) error {
	// 将 UpdateURL 进行解析
	url, err := url.Parse(r.UpdateURL)
	if nil != err {
		return err
	}
	http.Handle(url.Path, &serviceUpdateHandler{})

	buf := new(bytes.Buffer)    // 实现了 io.Reader 和 io.Writer
	enc := json.NewEncoder(buf) // 新建一个编码器
	err = enc.Encode(r)
	if nil != err {
		return err
	}

	// 向注册服务发送一个注册请求
	resp, err := http.Post(ServerURL, "application/json", buf)
	if nil != err {
		return err
	}

	// 如果请求结果的状态码不等于 200
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service %s", r.ServiceName) // 返回注册失败
	}

	return nil
}

// 注册服务的客户端的出口
// 其他服务调用 RegistryExit 函数，就能从 registry 中取消注册
func RegistryExit(serviceURL string) error {
	// 向服务器发送一个取消注册请求
	req, _ := http.NewRequest("DELETE", ServerURL, bytes.NewBuffer([]byte(serviceURL)))
	req.Header.Add("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deregistry service, registry service responsed with code %v", resp.StatusCode)
	}
	return nil
}

type serviceUpdateHandler struct{}

// 服务更新的 Web 请求处理
func (suh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // 只允许 POST 请求
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	dec := json.NewDecoder(r.Body) // patch 被存放在请求体中
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	prov.Update(p)
}

// A 依赖于 B
// 则 B 对于 A 来说就是一个 provider
type providers struct {
	services map[ServiceName][]string // k：服务名，v：服务的 URL
	mutex    *sync.RWMutex
}

var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}

func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 处理“新增”的服务
	for _, addedProv := range pat.Added {
		if _, ok := p.services[addedProv.Name]; !ok {
			p.services[addedProv.Name] = []string{addedProv.URL}
		} else {
			p.services[addedProv.Name] = append(p.services[addedProv.Name], addedProv.URL)
		}
	}

	// 处理“移除”的服务
	for _, rmProv := range pat.Removed {
		if _, ok := p.services[rmProv.Name]; !ok {
			for i := 0; i < len(p.services[rmProv.Name]); i++ {
				if p.services[rmProv.Name][i] == rmProv.URL {
					p.services[rmProv.Name] = append(p.services[rmProv.Name][:i], p.services[rmProv.Name][i+1:]...)
				}
			}
		}
	}
}

// 通过服务的名称获取服务实例的 URL
func (p providers) get(serviceName ServiceName) (string, error) {
	// 获取该服务的所有实例的 URL
	urls, ok := p.services[serviceName]
	if !ok {
		return "", fmt.Errorf("Unable to find service named %s", serviceName)
	}

	// 随机返回其中的一个
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(urls))

	return urls[idx], nil
}

func GetProvider(serviceName ServiceName) (string, error) {
	return prov.get(serviceName)
}
