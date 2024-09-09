package v1

import (
	"net/http" // 提供HTTP客户端和服务端的功能

	"chat-room/internal/model"      // 引入数据模型包，定义了数据库表结构
	"chat-room/internal/service"    // 引入服务层，用于调用业务逻辑
	"chat-room/pkg/common/request"  // 引入通用请求包，定义了请求参数结构体
	"chat-room/pkg/common/response" // 引入通用响应包，用于统一格式化HTTP响应
	"chat-room/pkg/global/log"      // 引入全局日志记录器，用于日志记录

	"github.com/gin-gonic/gin" // 引入Gin框架，用于处理HTTP请求
)

// Login 函数用于处理用户登录请求
func Login(c *gin.Context) {
	var user model.User // 声明一个User类型的变量，用于接收客户端发送的登录信息
	// c.BindJSON(&user)  // 解析请求中的JSON数据，绑定到user变量
	c.ShouldBindJSON(&user)                         // 使用ShouldBindJSON方法绑定请求中的JSON数据到user变量
	log.Logger.Debug("user", log.Any("user", user)) // 记录用户登录信息到日志中

	if service.UserService.Login(&user) { // 调用服务层的Login方法验证用户信息
		c.JSON(http.StatusOK, response.SuccessMsg(user)) // 登录成功，返回用户信息
		return
	}

	c.JSON(http.StatusOK, response.FailMsg("Login failed")) // 登录失败，返回失败信息
}

// Register 函数用于处理用户注册请求
func Register(c *gin.Context) {
	var user model.User                        // 声明一个User类型的变量，用于接收客户端发送的注册信息
	c.ShouldBindJSON(&user)                    // 将请求中的JSON数据绑定到user变量
	err := service.UserService.Register(&user) // 调用服务层的Register方法注册新用户
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error())) // 如果注册失败，返回错误信息
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(user)) // 注册成功，返回用户信息
}

// ModifyUserInfo 函数用于处理用户信息修改请求
func ModifyUserInfo(c *gin.Context) {
	var user model.User                             // 声明一个User类型的变量，用于接收客户端发送的修改信息
	c.ShouldBindJSON(&user)                         // 将请求中的JSON数据绑定到user变量
	log.Logger.Debug("user", log.Any("user", user)) // 记录用户信息到日志中
	if err := service.UserService.ModifyUserInfo(&user); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error())) // 如果修改失败，返回错误信息
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(nil)) // 修改成功，返回成功信息
}

// GetUserDetails 函数用于获取某个用户的详细信息
func GetUserDetails(c *gin.Context) {
	uuid := c.Param("uuid") // 从请求路径中获取用户的UUID

	c.JSON(http.StatusOK, response.SuccessMsg(service.UserService.GetUserDetails(uuid))) // 返回用户详细信息
}

// GetUserOrGroupByName 函数通过用户名或组名获取用户或组信息
func GetUserOrGroupByName(c *gin.Context) {
	name := c.Query("name") // 从查询参数中获取名称

	c.JSON(http.StatusOK, response.SuccessMsg(service.UserService.GetUserOrGroupByName(name))) // 返回匹配的用户或组信息
}

// GetUserList 函数用于获取用户的好友列表
func GetUserList(c *gin.Context) {
	uuid := c.Query("uuid")                                                           // 从查询参数中获取用户的UUID
	c.JSON(http.StatusOK, response.SuccessMsg(service.UserService.GetUserList(uuid))) // 返回用户好友列表
}

// AddFriend 函数用于处理添加好友请求
func AddFriend(c *gin.Context) {
	var userFriendRequest request.FriendRequest // 声明一个FriendRequest类型的变量，用于接收客户端发送的好友请求信息
	c.ShouldBindJSON(&userFriendRequest)        // 将请求中的JSON数据绑定到userFriendRequest变量

	err := service.UserService.AddFriend(&userFriendRequest) // 调用服务层的AddFriend方法添加好友
	if nil != err {
		c.JSON(http.StatusOK, response.FailMsg(err.Error())) // 如果添加失败，返回错误信息
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(nil)) // 添加成功，返回成功信息
}
