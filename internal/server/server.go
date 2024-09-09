package server

import (
	"chat-room/config"              // 引入配置包，用于读取配置信息
	"chat-room/internal/service"    // 引入服务层，用于业务逻辑处理
	"chat-room/pkg/common/constant" // 引入常量包，用于定义全局常量
	"chat-room/pkg/common/util"     // 引入工具包，提供各种实用函数
	"chat-room/pkg/global/log"      // 引入全局日志记录器，用于日志记录
	"chat-room/pkg/protocol"        // 引入协议包，用于消息协议处理
	"encoding/base64"               // 引入base64编码解码库
	"io/ioutil"                     // 引入I/O实用函数库，用于文件读写
	"strings"                       // 引入字符串处理库
	"sync"                          // 引入同步包，用于并发控制

	"github.com/gogo/protobuf/proto" // 引入protobuf库，用于序列化和反序列化消息
	"github.com/google/uuid"         // 引入UUID库，用于生成唯一标识符
)

// MyServer 是全局的Server实例，用于管理WebSocket客户端
var MyServer = NewServer()

// Server 结构体用于管理连接的客户端和消息的处理
type Server struct {
	Clients   map[string]*Client // 存储连接的客户端，以客户端名称为键
	mutex     *sync.Mutex        // 互斥锁，用于保护共享资源
	Broadcast chan []byte        // 广播通道，用于发送消息给所有客户端
	Register  chan *Client       // 注册通道，用于注册新客户端
	Ungister  chan *Client       // 注销通道，用于注销客户端
}

// NewServer 初始化并返回一个Server实例
func NewServer() *Server {
	return &Server{
		mutex:     &sync.Mutex{},
		Clients:   make(map[string]*Client),
		Broadcast: make(chan []byte),
		Register:  make(chan *Client),
		Ungister:  make(chan *Client),
	}
}

// ConsumerKafkaMsg 函数用于消费Kafka中的消息，并将消息发送到Broadcast通道
func ConsumerKafkaMsg(data []byte) {
	MyServer.Broadcast <- data
}

// Start 函数启动服务器，处理客户端注册、注销和消息广播
func (s *Server) Start() {
	log.Logger.Info("start server", log.Any("start server", "start server..."))
	for {
		select {
		case conn := <-s.Register: // 处理新客户端注册
			log.Logger.Info("login", log.Any("login", "new user login in"+conn.Name))
			s.Clients[conn.Name] = conn
			msg := &protocol.Message{
				From:    "System",
				To:      conn.Name,
				Content: "welcome!",
			}
			protoMsg, _ := proto.Marshal(msg)
			conn.Send <- protoMsg

		case conn := <-s.Ungister: // 处理客户端注销
			log.Logger.Info("loginout", log.Any("loginout", conn.Name))
			if _, ok := s.Clients[conn.Name]; ok {
				close(conn.Send)
				delete(s.Clients, conn.Name)
			}

		case message := <-s.Broadcast: // 处理广播消息
			msg := &protocol.Message{}
			proto.Unmarshal(message, msg)

			if msg.To != "" {
				// 处理点对点消息或群组消息
				if msg.ContentType >= constant.TEXT && msg.ContentType <= constant.VIDEO {
					// 保存消息到数据库
					_, exists := s.Clients[msg.From]
					if exists {
						saveMessage(msg)
					}

					// 单聊消息
					if msg.MessageType == constant.MESSAGE_TYPE_USER {
						client, ok := s.Clients[msg.To]
						if ok {
							msgByte, err := proto.Marshal(msg)
							if err == nil {
								client.Send <- msgByte
							}
						}
					} else if msg.MessageType == constant.MESSAGE_TYPE_GROUP {
						// 群聊消息
						sendGroupMessage(msg, s)
					}
				} else {
					// 处理语音或视频聊天，直接转发消息
					client, ok := s.Clients[msg.To]
					if ok {
						client.Send <- message
					}
				}

			} else {
				// 广播消息，发送给所有客户端
				for id, conn := range s.Clients {
					log.Logger.Info("allUser", log.Any("allUser", id))

					select {
					case conn.Send <- message:
					default:
						close(conn.Send)
						delete(s.Clients, conn.Name)
					}
				}
			}
		}
	}
}

// sendGroupMessage 函数发送群组消息，遍历群组所有成员并逐个发送
func sendGroupMessage(msg *protocol.Message, s *Server) {
	// 获取群组成员列表
	users := service.GroupService.GetUserIdByGroupUuid(msg.To)
	for _, user := range users {
		if user.Uuid == msg.From {
			continue
		}

		client, ok := s.Clients[user.Uuid]
		if !ok {
			continue
		}

		// 获取发送者详情
		fromUserDetails := service.UserService.GetUserDetails(msg.From)
		// 修改消息的From字段，使其表示群组
		msgSend := protocol.Message{
			Avatar:       fromUserDetails.Avatar,
			FromUsername: msg.FromUsername,
			From:         msg.To,
			To:           msg.From,
			Content:      msg.Content,
			ContentType:  msg.ContentType,
			Type:         msg.Type,
			MessageType:  msg.MessageType,
			Url:          msg.Url,
		}

		// 将消息序列化并发送给群成员
		msgByte, err := proto.Marshal(&msgSend)
		if err == nil {
			client.Send <- msgByte
		}
	}
}

// saveMessage 函数保存消息，如果是文件消息则保存文件并更新消息内容
func saveMessage(message *protocol.Message) {
	// 处理base64编码的文件内容
	if message.ContentType == 2 {
		url := uuid.New().String() + ".png"
		index := strings.Index(message.Content, "base64")
		index += 7

		content := message.Content
		content = content[index:]

		dataBuffer, dataErr := base64.StdEncoding.DecodeString(content)
		if dataErr != nil {
			log.Logger.Error("transfer base64 to file error", log.String("transfer base64 to file error", dataErr.Error()))
			return
		}
		err := ioutil.WriteFile(config.GetConfig().StaticPath.FilePath+url, dataBuffer, 0666)
		if err != nil {
			log.Logger.Error("write file error", log.String("write file error", err.Error()))
			return
		}
		message.Url = url
		message.Content = ""
	} else if message.ContentType == 3 {
		// 处理普通文件的二进制数据
		fileSuffix := util.GetFileType(message.File)
		if fileSuffix == "" {
			fileSuffix = strings.ToLower(message.FileSuffix)
		}
		contentType := util.GetContentTypeBySuffix(fileSuffix)
		url := uuid.New().String() + "." + fileSuffix
		err := ioutil.WriteFile(config.GetConfig().StaticPath.FilePath+url, message.File, 0666)
		if err != nil {
			log.Logger.Error("write file error", log.String("write file error", err.Error()))
			return
		}
		message.Url = url
		message.File = nil
		message.ContentType = contentType
	}

	// 将消息保存到数据库
	service.MessageService.SaveMessage(*message)
}
