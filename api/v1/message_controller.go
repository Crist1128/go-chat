package v1

import (
	"net/http" // 提供HTTP客户端和服务端的功能

	"chat-room/internal/service"    // 引入服务层，用于调用业务逻辑
	"chat-room/pkg/common/request"  // 引入通用请求包，定义了请求参数结构体
	"chat-room/pkg/common/response" // 引入通用响应包，用于统一格式化HTTP响应
	"chat-room/pkg/global/log"      // 引入全局日志记录器，用于日志记录

	"github.com/gin-gonic/gin" // 引入Gin框架，用于处理HTTP请求
)

// GetMessage 函数用于获取某个用户的消息列表
func GetMessage(c *gin.Context) {
	log.Logger.Info(c.Query("uuid"))          // 记录请求中的UUID参数到日志中
	var messageRequest request.MessageRequest // 声明一个MessageRequest类型的变量，用于接收请求参数
	err := c.BindQuery(&messageRequest)       // 将查询参数绑定到messageRequest变量
	if nil != err {
		log.Logger.Error("bindQueryError", log.Any("bindQueryError", err)) // 如果绑定失败，记录错误日志
	}
	log.Logger.Info("messageRequest params: ", log.Any("messageRequest", messageRequest)) // 记录绑定后的请求参数到日志中

	messages, err := service.MessageService.GetMessages(messageRequest) // 调用服务层方法，获取用户的消息列表
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error())) // 如果出现错误，返回失败信息
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(messages)) // 返回消息列表，响应成功
}
