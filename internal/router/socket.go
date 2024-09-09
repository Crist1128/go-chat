package router

import (
	"chat-room/internal/server" // 引入服务器包，处理客户端连接和消息
	"chat-room/pkg/global/log"  // 引入全局日志记录器，用于日志记录
	"net/http"                  // 引入HTTP包，用于处理HTTP相关操作

	"github.com/gin-gonic/gin"     // 引入Gin框架，用于处理HTTP请求
	"github.com/gorilla/websocket" // 引入Gorilla WebSocket库，用于WebSocket连接
	"go.uber.org/zap"              // 引入Zap日志库，用于日志记录
)

// upGrader 用于升级HTTP连接为WebSocket连接
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有源的连接
	},
}

// RunSocekt 函数处理WebSocket连接
func RunSocekt(c *gin.Context) {
	user := c.Query("user") // 获取用户参数
	if user == "" {
		return // 如果没有用户参数，直接返回
	}
	log.Logger.Info("newUser", zap.String("newUser", user)) // 记录新用户连接的日志
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)   // 将HTTP连接升级为WebSocket连接
	if err != nil {
		return // 如果升级失败，直接返回
	}

	// 创建一个新的客户端实例
	client := &server.Client{
		Name: user,              // 设置客户端的用户名
		Conn: ws,                // WebSocket连接
		Send: make(chan []byte), // 创建一个发送消息的通道
	}

	// 将新客户端注册到服务器中
	server.MyServer.Register <- client
	// 启动客户端的读写协程
	go client.Read()
	go client.Write()
}
