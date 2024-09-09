package router

import (
	"chat-room/api/v1"              // 引入API v1版本的路由处理函数
	"chat-room/pkg/common/response" // 引入通用响应包，用于统一格式化HTTP响应
	"chat-room/pkg/global/log"      // 引入全局日志记录器，用于日志记录
	"net/http"                      // 引入HTTP包，用于处理HTTP相关操作

	"github.com/gin-gonic/gin" // 引入Gin框架，用于处理HTTP请求
	"go.uber.org/zap"          // 引入Zap日志库，用于日志记录
)

// NewRouter 函数初始化Gin引擎，并配置路由规则和中间件
func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // 设置Gin为发布模式，减少日志输出

	server := gin.Default() // 使用默认配置初始化Gin引擎，包括日志和恢复中间件
	server.Use(Cors())      // 使用跨域请求中间件
	server.Use(Recovery)    // 使用自定义恢复中间件，捕获并处理panic
	// server.Use(gin.Recovery()) // 可以选择使用Gin自带的恢复中间件

	socket := RunSocekt // 定义WebSocket路由处理函数

	group := server.Group("") // 定义一个基础路径分组
	{
		// 用户相关路由
		group.GET("/user", v1.GetUserList)               // 获取用户列表
		group.GET("/user/:uuid", v1.GetUserDetails)      // 获取指定UUID用户的详细信息
		group.GET("/user/name", v1.GetUserOrGroupByName) // 通过用户名或群组名获取信息
		group.POST("/user/register", v1.Register)        // 用户注册
		group.POST("/user/login", v1.Login)              // 用户登录
		group.PUT("/user", v1.ModifyUserInfo)            // 修改用户信息

		// 好友相关路由
		group.POST("/friend", v1.AddFriend) // 添加好友

		// 消息相关路由
		group.GET("/message", v1.GetMessage) // 获取消息列表

		// 文件相关路由
		group.GET("/file/:fileName", v1.GetFile) // 获取文件
		group.POST("/file", v1.SaveFile)         // 上传文件

		// 群组相关路由
		group.GET("/group/:uuid", v1.GetGroup)                       // 获取群组信息
		group.POST("/group/:uuid", v1.SaveGroup)                     // 保存群组信息
		group.POST("/group/join/:userUuid/:groupUuid", v1.JoinGroup) // 加入群组
		group.GET("/group/user/:uuid", v1.GetGroupUsers)             // 获取群组用户列表

		// WebSocket相关路由
		group.GET("/socket.io", socket) // WebSocket连接
	}
	return server // 返回配置好的Gin引擎
}

// Cors 中间件用于处理跨域请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") // 获取请求的Origin头
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")                                                                                                                          // 允许所有来源的跨域请求
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")                                                                                   // 允许的HTTP方法
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")                                                             // 允许的请求头
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type") // 公开的响应头
			c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                  // 允许携带认证信息的跨域请求
		}
		// 处理预检请求
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		// 捕获可能的panic，并记录错误日志
		defer func() {
			if err := recover(); err != nil {
				log.Logger.Error("HttpError", zap.Any("HttpError", err))
			}
		}()

		c.Next() // 执行下一个中间件或处理器
	}
}

// Recovery 中间件用于捕获Gin处理过程中的panic，并返回标准化错误响应
func Recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Error("gin catch error: ", log.Any("gin catch error: ", r)) // 记录错误日志
			c.JSON(http.StatusOK, response.FailMsg("系统内部错误"))                      // 返回系统错误响应
		}
	}()
	c.Next() // 执行下一个中间件或处理器
}
