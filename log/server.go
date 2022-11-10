package log

import (
	"fmt"
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var log *stlog.Logger

var server *gin.Engine

type fileLog string

// 使得 fileLog 类型实现 io.Writer 接口定义的 Write 方法
func (fl fileLog) Write(data []byte) (int, error) {
	// 打开日志文件         文件路径     文件权限
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if nil != err {
		fmt.Printf("Failed to open log file: %s\n", string(fl))
		return 0, err
	}
	// 关闭前关闭日志文件
	defer f.Close()
	// 写入日志
	return f.Write(data)
}

// 初始化日志系统
func Run(destination string) *gin.Engine {
	// 初始化传入：io.Writer 类型的值，日志的前缀，日志的标记（长日期｜长时间）
	log = stlog.New(fileLog(destination), "go ", stlog.LstdFlags)
	// 初始化服务
	server = gin.Default()
	return server
}

// 使用 Gin 框架重写请求处理函数
func RegisterHandlers() {
	server.POST("/log", func(ctx *gin.Context) {
		msg, err := ioutil.ReadAll(ctx.Request.Body)
		if nil != err || len(msg) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Could not read anything from the request body!"})
			return
		}
		write(string(msg))
	})
}

// 将网络日志写入文件
func write(message string) {
	log.Printf("%v\n", message)
}
