package v1

import (
	"chat-room/internal/model"      // 引入数据模型包，定义了数据库表结构
	"chat-room/internal/service"    // 引入服务层，用于调用业务逻辑
	"chat-room/pkg/common/response" // 引入通用响应包，用于统一格式化HTTP响应
	"net/http"                      // 提供HTTP客户端和服务端的功能

	"github.com/gin-gonic/gin" // 引入Gin框架，用于处理HTTP请求
)

// GetGroup 函数用于获取某个用户的分组列表
func GetGroup(c *gin.Context) {
	uuid := c.Param("uuid")                             // 从请求路径中获取用户的UUID
	groups, err := service.GroupService.GetGroups(uuid) // 调用服务层方法，获取用户的分组列表
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error())) // 如果出现错误，返回失败信息
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(groups)) // 返回分组列表信息，响应成功
}

// SaveGroup 函数用于保存用户的分组信息
func SaveGroup(c *gin.Context) {
	uuid := c.Param("uuid")  // 从请求路径中获取用户的UUID
	var group model.Group    // 声明一个Group类型的变量，用于接收客户端发送的分组信息
	c.ShouldBindJSON(&group) // 将请求中的JSON数据绑定到group变量上

	service.GroupService.SaveGroup(uuid, group)     // 调用服务层方法，保存分组信息
	c.JSON(http.StatusOK, response.SuccessMsg(nil)) // 返回成功响应
}

// JoinGroup 函数用于将用户加入某个组
func JoinGroup(c *gin.Context) {
	userUuid := c.Param("userUuid")                            // 从请求路径中获取用户的UUID
	groupUuid := c.Param("groupUuid")                          // 从请求路径中获取组的UUID
	err := service.GroupService.JoinGroup(groupUuid, userUuid) // 调用服务层方法，将用户加入指定的组
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error())) // 如果出现错误，返回失败信息
		return
	}
	c.JSON(http.StatusOK, response.SuccessMsg(nil)) // 返回成功响应
}

// GetGroupUsers 函数用于获取指定组内的所有用户信息
func GetGroupUsers(c *gin.Context) {
	groupUuid := c.Param("uuid")                                  // 从请求路径中获取组的UUID
	users := service.GroupService.GetUserIdByGroupUuid(groupUuid) // 调用服务层方法，获取组内的用户列表
	c.JSON(http.StatusOK, response.SuccessMsg(users))             // 返回用户列表，响应成功
}
