package v1

import (
	"io/ioutil" // 用于读取文件
	"net/http"  // 提供HTTP客户端和服务端的功能
	"strings"   // 用于字符串操作，如截取文件后缀名

	"chat-room/config"              // 引入项目的配置包，用于获取配置参数
	"chat-room/internal/service"    // 引入服务层，用于调用业务逻辑
	"chat-room/pkg/common/response" // 引入通用响应包，用于统一格式化HTTP响应
	"chat-room/pkg/global/log"      // 引入全局日志记录器，用于日志记录

	"github.com/gin-gonic/gin" // 引入Gin框架，用于处理HTTP请求
	"github.com/google/uuid"   // 引入UUID库，用于生成唯一标识符
)

// GetFile 函数通过文件名称从服务器获取文件流并返回给前端，通常用于显示图片或其他静态资源
func GetFile(c *gin.Context) {
	fileName := c.Param("fileName")                                               // 从请求路径中获取文件名参数
	log.Logger.Info(fileName)                                                     // 记录请求的文件名到日志中
	data, _ := ioutil.ReadFile(config.GetConfig().StaticPath.FilePath + fileName) // 读取文件内容，路径由配置文件提供
	c.Writer.Write(data)                                                          // 将文件内容写入HTTP响应，返回给前端
}

// SaveFile 函数用于处理文件上传，如头像等文件，保存到服务器并更新用户头像信息
func SaveFile(c *gin.Context) {
	namePreffix := uuid.New().String() // 生成一个新的UUID作为文件名前缀，确保文件名唯一

	userUuid := c.PostForm("uuid") // 从POST请求中获取用户UUID，用于关联用户信息

	file, _ := c.FormFile("file")             // 从请求中获取上传的文件
	fileName := file.Filename                 // 获取原始文件名
	index := strings.LastIndex(fileName, ".") // 找到文件名中最后一个'.'的位置，用于获取文件后缀
	suffix := fileName[index:]                // 提取文件后缀名

	newFileName := namePreffix + suffix // 将UUID和后缀拼接成新的文件名

	log.Logger.Info("file", log.Any("file name", config.GetConfig().StaticPath.FilePath+newFileName)) // 记录新文件名到日志中
	log.Logger.Info("userUuid", log.Any("userUuid name", userUuid))                                   // 记录关联的用户UUID到日志中

	c.SaveUploadedFile(file, config.GetConfig().StaticPath.FilePath+newFileName) // 将上传的文件保存到服务器指定路径
	err := service.UserService.ModifyUserAvatar(newFileName, userUuid)           // 调用业务逻辑层，更新用户的头像信息
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error())) // 如果更新失败，返回失败信息
	}
	c.JSON(http.StatusOK, response.SuccessMsg(newFileName)) // 更新成功，返回新文件名作为响应
}
