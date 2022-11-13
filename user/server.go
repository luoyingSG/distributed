package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var server *gin.Engine

// 用户服务的请求处理函数
func UserHandlers() {
	// 用户注册
	server.POST("/user", func(ctx *gin.Context) {
		// request.Body
		body, _ := ctx.GetRawData()

		// 解析序列化
		var request map[string]interface{}
		json.Unmarshal(body, &request)

		fmt.Println("request:", request)

		// 将用户数据添加到数据库
		id := addUser(request)

		if id < 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "邮箱地址和用户名不能为空！"})
		} else if id == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "邮箱地址或用户名非法！"})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"msg": "请求成功！", "ID": id})
		}
	})
}

// 注册用户
func addUser(request map[string]interface{}) int {
	u := new(User)
	if email, ok := request["Email"]; ok {
		u.Email = fmt.Sprintf("%v", email)
	} else {
		return -1
	}
	if username, ok := request["Username"]; ok {
		u.Username = fmt.Sprintf("%v", username)
	} else {
		return -1
	}

	fmt.Println("signin user:", *u)

	err := u.signin()
	if nil != err {
		fmt.Printf("Failed to sign in: %s", err.Error())
		return 0
	}

	return int(u.ID)
}

// 初始化路由、中间件
func Run() *gin.Engine {
	server = gin.Default()
	return server
}
