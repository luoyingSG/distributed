package registry

// 注册服务类型
type Registration struct {
	ServiceName ServiceName
	ServiceURL  string
}

type ServiceName string

const (
	LogService = ServiceName("LogService")
)
