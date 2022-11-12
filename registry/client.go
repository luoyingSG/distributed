package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// 注册服务的客户端的入口
// 其他服务调用 RegistryEntrance 函数，就能自动将自己注册进 registry
func RegistryEntrance(r Registration) error {
	buf := new(bytes.Buffer)    // 实现了 io.Reader 和 io.Writer
	enc := json.NewEncoder(buf) // 新建一个编码器
	err := enc.Encode(r)
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
