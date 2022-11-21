package log

import (
	"bytes"
	"fmt"
	stdlog "log"
	"net/http"

	"github.com/luoying_gh/distributed/registry"
)

func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stdlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService)) // 设置客户端日志的前缀
	stdlog.SetFlags(0)                                      // 设置客户端日志的标记（无）
	stdlog.SetOutput(clientLogger{url: serviceURL})         // 设置客户端日志的输出
}

// 实现了 io.Writer
type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer(data)
	resp, err := http.Post(cl.url+"/log", "text/plain", b) // 向服务端发送 Post 请求
	if nil != err {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed to write log message, Service responsed with cdoe %v", resp.StatusCode)
	}

	return len(data), nil
}
