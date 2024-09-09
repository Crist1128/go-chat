package server

import (
	"chat-room/config"              // 引入配置包，用于获取配置信息
	"chat-room/internal/kafka"      // 引入Kafka包，用于处理Kafka消息队列
	"chat-room/pkg/common/constant" // 引入常量包，定义项目中的常量
	"chat-room/pkg/global/log"      // 引入全局日志记录器，用于日志记录
	"chat-room/pkg/protocol"        // 引入协议包，用于处理消息的协议格式

	"github.com/gogo/protobuf/proto" // 引入Protobuf库，用于序列化和反序列化消息
	"github.com/gorilla/websocket"   // 引入Gorilla WebSocket库，用于处理WebSocket连接
)

// Client 结构体表示一个WebSocket连接客户端
type Client struct {
	Conn *websocket.Conn // WebSocket连接实例
	Name string          // 客户端的名称（通常是用户名）
	Send chan []byte     // 发送消息的通道
}

// Read 方法用于从WebSocket连接读取消息
func (c *Client) Read() {
	defer func() {
		MyServer.Ungister <- c // 当连接关闭时，将客户端从服务器的客户端列表中注销
		c.Conn.Close()         // 关闭连接
	}()

	for {
		c.Conn.PongHandler()                    // 处理WebSocket的Pong消息，保持连接活跃
		_, message, err := c.Conn.ReadMessage() // 读取消息
		if err != nil {
			// 如果读取消息失败，记录错误日志并关闭连接
			log.Logger.Error("client read message error", log.Any("client read message error", err.Error()))
			MyServer.Ungister <- c // 将客户端从服务器的客户端列表中注销
			c.Conn.Close()         // 关闭连接
			break
		}

		msg := &protocol.Message{}    // 创建一个空的Message对象
		proto.Unmarshal(message, msg) // 反序列化从客户端接收到的消息

		// 处理心跳消息（pong响应）
		if msg.Type == constant.HEAT_BEAT {
			pong := &protocol.Message{
				Content: constant.PONG,      // 响应内容设置为PONG
				Type:    constant.HEAT_BEAT, // 消息类型为心跳
			}
			pongByte, err2 := proto.Marshal(pong) // 将响应消息序列化为字节数组
			if nil != err2 {
				log.Logger.Error("client marshal message error", log.Any("client marshal message error", err2.Error()))
			}
			c.Conn.WriteMessage(websocket.BinaryMessage, pongByte) // 将响应消息写回客户端
		} else {
			// 如果消息不是心跳消息，则根据配置决定使用Kafka还是直接广播
			if config.GetConfig().MsgChannelType.ChannelType == constant.KAFKA {
				kafka.Send(message) // 通过Kafka发送消息
			} else {
				MyServer.Broadcast <- message // 通过服务器的广播通道广播消息
			}
		}
	}
}

// Write 方法用于向WebSocket连接发送消息
func (c *Client) Write() {
	defer func() {
		c.Conn.Close() // 当写操作完成后，关闭连接
	}()

	for message := range c.Send { // 读取发送通道中的消息
		c.Conn.WriteMessage(websocket.BinaryMessage, message) // 将消息发送给客户端
	}
}
