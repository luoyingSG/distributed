package registry

// 注册服务类型
type Registration struct {
	ServiceName      ServiceName   // 服务名
	ServiceURL       string        // 服务 URL
	RequiredServices []ServiceName // 服务依赖
	UpdateURL        string        // registryservice 在发现其依赖的服务后，将该依赖发送到这个 url
}

type ServiceName string

const (
	LogService  = ServiceName("LogService")
	UserService = ServiceName("UserService")
)

// 注册情况的更新，这是每一条更新
type patchEntry struct {
	Name ServiceName
	URL  string
}

// 注册情况的更新，这是所有的更新
type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
